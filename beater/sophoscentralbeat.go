package beater

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/antihax/optional"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/paths"
	"github.com/logrhythm/sophoscentralbeat/config"
	"github.com/logrhythm/sophoscentralbeat/handlers"
	"github.com/logrhythm/sophoscentralbeat/sophoscentral"
	"github.com/mitchellh/mapstructure"
)

// Sophoscentralbeat configuration.
type Sophoscentralbeat struct {
	done            chan struct{}
	config          config.Config
	sophos          *sophoscentral.APIClient
	sophosAuth      context.Context
	client          beat.Client
	logger          logp.Logger
	basepath        string
	currentPosition *scbPosition
	posHandler      *handlers.PositionHandler
}

//Positionfile : position file data format
type scbPosition struct {
	EventsTimestamp int64 `json:"timestamp_events"`
	AlertsTimestamp int64 `json:"timestamp_alerts"`
}

// New creates an instance of sophoscentralbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	logger := logp.NewLogger("sophoscentralbeat-internal")
	c := config.DefaultConfig
	sophoscentralConfig := sophoscentral.NewConfiguration()
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	logger.Info("Period set to ", c.Period.String())
	sophos := sophoscentral.NewAPIClient(sophoscentralConfig)
	auth := context.WithValue(context.Background(), sophoscentral.ContextAPIKey, sophoscentral.APIKey{
		Key: c.APIKey,
	})

	pos, err := handlers.NewPostionHandler("")

	if err != nil {
		logp.Err("Unable to get position Handler %v", err)
		return nil, err
	}

	currentPos := new(scbPosition)
	poserr := pos.ReadPositionfromFile(currentPos)
	yesterdayTime := GenerateYesterdayTimeStamp()
	if poserr != nil {

		currentPos.EventsTimestamp = yesterdayTime
		currentPos.AlertsTimestamp = yesterdayTime
	}

	if currentPos.EventsTimestamp < yesterdayTime {
		currentPos.EventsTimestamp = yesterdayTime
	}

	if currentPos.AlertsTimestamp < yesterdayTime {
		currentPos.AlertsTimestamp = yesterdayTime
	}

	bt := &Sophoscentralbeat{
		done:            make(chan struct{}),
		sophos:          sophos,
		sophosAuth:      auth,
		config:          c,
		logger:          *logger,
		basepath:        c.Basepath,
		currentPosition: currentPos,
		posHandler:      pos,
	}
	return bt, nil
}

//GetSophosEvents : calls Sophos Events Api
func GetSophosEvents(scb Sophoscentralbeat) error {
	scb.logger.Info("Making sophos event call")
	var isDataReceived bool

	options := &sophoscentral.GetEventsUsingGET1Opts{
		Limit:    optional.NewInt32(1000),
		FromDate: optional.NewInt64(scb.currentPosition.EventsTimestamp),
	}

	value, _, err := scb.sophos.EventControllerV1ImplApi.GetEventsUsingGET1(scb.sophosAuth, scb.config.APIKey, scb.config.Authorization, scb.basepath, options)
	if err != nil {
		scb.logger.Error(err)
		return err
	}

	for _, item := range value.Items {
		scb.client.Publish(GetEvent(item))
		eventCreationTime, _ := time.Parse(time.RFC3339, item.CreatedAt)
		UpdateEventTime(&scb, eventCreationTime.Unix())
	}

	isDataReceived = len(value.Items) > 0

	for value.HasMore == true {

		options = &sophoscentral.GetEventsUsingGET1Opts{
			Limit:  optional.NewInt32(1000),
			Cursor: optional.NewString(value.NextCursor),
		}
		nestedVal, _, err := scb.sophos.EventControllerV1ImplApi.GetEventsUsingGET1(scb.sophosAuth, scb.config.APIKey, scb.config.Authorization, scb.basepath, options)
		if err != nil {
			scb.logger.Error(err)
			return err
		}
		for _, item := range nestedVal.Items {
			scb.client.Publish(GetEvent(item))
			eventCreationTime, _ := time.Parse(time.RFC3339, item.CreatedAt)
			UpdateEventTime(&scb, eventCreationTime.Unix())
		}
		value.HasMore = nestedVal.HasMore
		value.NextCursor = nestedVal.NextCursor
	}

	if isDataReceived {
		scb.currentPosition.EventsTimestamp = scb.currentPosition.EventsTimestamp + 1
	}

	scb.posHandler.WritePostiontoFile(scb.currentPosition)
	return nil
}

func LegacyEventEntityToCommonMap(entity sophoscentral.LegacyEventEntity) (common.MapStr, error) {
	var result common.MapStr
	mConfig := &mapstructure.DecoderConfig{
		TagName: "json",
		Result:  &result,
	}
	decoder, _ := mapstructure.NewDecoder(mConfig)
	err := decoder.Decode(entity)
	if err != nil {
		logp.L().Error("Error decoding Okta LogEvent record", err)
		return nil, err
	}
	return result, nil
}

//GetSophosAlertsOld : call alerts API
func GetSophosAlerts(scb Sophoscentralbeat) error {
	scb.logger.Info("Making sophos alert call")

	var isDataReceived bool

	options := &sophoscentral.GetAlertsUsingGET1Opts{
		Limit:    optional.NewInt32(1000),
		FromDate: optional.NewInt64(scb.currentPosition.AlertsTimestamp),
	}

	value, _, err := scb.sophos.AlertControllerV1ImplApi.GetAlertsUsingGET1(scb.sophosAuth, scb.config.APIKey, scb.config.Authorization, scb.basepath, options)
	if err != nil {
		scb.logger.Error(err)
		return err
	}

	for _, item := range value.Items {
		scb.client.Publish(GetEvent(item))
		alertCreationTime, _ := time.Parse(time.RFC3339, item.CreatedAt)
		UpdateAlertTime(&scb, alertCreationTime.Unix())

	}

	isDataReceived = len(value.Items) > 0

	for value.HasMore == true {

		options = &sophoscentral.GetAlertsUsingGET1Opts{
			Limit:  optional.NewInt32(1000),
			Cursor: optional.NewString(value.NextCursor),
		}

		nestedVal, _, err := scb.sophos.AlertControllerV1ImplApi.GetAlertsUsingGET1(scb.sophosAuth, scb.config.APIKey, scb.config.Authorization, scb.basepath, options)
		if err != nil {
			scb.logger.Error(err)
			return err
		}
		for _, item := range nestedVal.Items {
			scb.client.Publish(GetEvent(item))
			alertCreationTime, _ := time.Parse(time.RFC3339, item.CreatedAt)
			UpdateAlertTime(&scb, alertCreationTime.Unix())
		}
		value.HasMore = nestedVal.HasMore
		value.NextCursor = nestedVal.NextCursor
	}

	if isDataReceived {
		scb.currentPosition.AlertsTimestamp = scb.currentPosition.AlertsTimestamp + 1
	}

	scb.posHandler.WritePostiontoFile(scb.currentPosition)
	return nil
}

func AlertEntityToCommonMap(entity sophoscentral.AlertEntity) (common.MapStr, error) {
	var result common.MapStr
	mConfig := &mapstructure.DecoderConfig{
		TagName: "json",
		Result:  &result,
	}
	decoder, _ := mapstructure.NewDecoder(mConfig)
	err := decoder.Decode(entity)
	if err != nil {
		logp.L().Error("Error decoding Okta LogEvent record", err)
		return nil, err
	}
	return result, nil
}

// Run starts sophoscentralbeat.
func (scb *Sophoscentralbeat) Run(b *beat.Beat) error {

	scb.logger.Info("sophoscentralbeat is running! Hit CTRL-C to stop it.")

	var err error
	scb.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(scb.config.Period)
	for {
		select {
		case <-scb.done:
			return nil
		case <-ticker.C:
			scb.logger.Info("Tick")
		}
		scb.logger.Info("Attempting to fetch Sophos Central Events")
		err := GetSophosEvents(*scb)
		if err != nil {
			scb.logger.Error(err)
		}

		scb.logger.Info("Attempting to fetch Sophos Alerts")
		err = GetSophosAlerts(*scb)
		if err != nil {
			scb.logger.Error(err)
		}

		scb.logger.Info("Events sent")
	}
}

// Stop stops sophoscentralbeat.
func (scb *Sophoscentralbeat) Stop() {
	scb.client.Close()
	scb.posHandler.WritePostiontoFile(scb.currentPosition)
	close(scb.done)
}

func UpdateEventTime(scb *Sophoscentralbeat, eventTimeStamp int64) {
	scb.currentPosition.EventsTimestamp = eventTimeStamp
}

func UpdateAlertTime(scb *Sophoscentralbeat, alertTimeStamp int64) {
	scb.currentPosition.AlertsTimestamp = alertTimeStamp
}

//WriteTimeStamp : writes timestamp to file
func WriteTimeStamp(eventTimeStamp int64, alertTimeStamp int64) {

	filePath := filepath.Join(paths.Paths.Home, "logs/pos.json")
	var position scbPosition

	//position file unavailable
	if _, err := os.Stat(filePath); err != nil {

		position = scbPosition{
			EventsTimestamp: eventTimeStamp,
			AlertsTimestamp: alertTimeStamp,
		}
	} else {
		//position file available

		position, _ = ReadTimeStamp()

		if eventTimeStamp != 0 && alertTimeStamp == 0 {
			position.EventsTimestamp = eventTimeStamp
		} else if eventTimeStamp == 0 && alertTimeStamp != 0 {
			position.AlertsTimestamp = alertTimeStamp
		} else if eventTimeStamp != 0 && alertTimeStamp != 0 {
			position.EventsTimestamp = eventTimeStamp
			position.AlertsTimestamp = alertTimeStamp
		}
	}

	jsonFile, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		return
	}
	jsonWriter := io.Writer(jsonFile)
	encoder := json.NewEncoder(jsonWriter)
	err = encoder.Encode(&position)
	if err != nil {
		fmt.Println("Error encoding JSON to file:", err)
		return
	}
}

//ReadTimeStamp : read tiemstamp from file
func ReadTimeStamp() (scbPosition, bool) {
	filePath := filepath.Join(paths.Paths.Home, "logs/pos.json")
	var pos scbPosition
	status := false

	if _, err := os.Stat(filePath); err == nil {
		// path/to/whatever exists

		jsonFile, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error opening JSON file:", err)
			return pos, false
		}
		defer jsonFile.Close()
		decoder := json.NewDecoder(jsonFile)

		err = decoder.Decode(&pos)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return pos, false
		}

		status = true
	}

	return pos, status
}

//GenerateYesterdayTimeStamp : generate 24 hour prior timestamp
func GenerateYesterdayTimeStamp() int64 {
	return time.Now().AddDate(0, 0, -1).UTC().Unix()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//GetEvent converts json data to beats json response
func GetEvent(data interface{}) beat.Event {

	event := beat.Event{

		Timestamp: time.Now(),

		Fields: common.MapStr{

			"response": data,
		},
	}

	return event

}
