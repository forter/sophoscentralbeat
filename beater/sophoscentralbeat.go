package beater

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/antihax/optional"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/mitchellh/mapstructure"

	"github.com/logrhythm/sophoscentralbeat/config"
	"github.com/logrhythm/sophoscentralbeat/sophoscentral"
)

// Sophoscentralbeat configuration.
type Sophoscentralbeat struct {
	done       chan struct{}
	config     config.Config
	sophos     *sophoscentral.APIClient
	sophosAuth context.Context
	client     beat.Client
	logger     logp.Logger
	basepath   string
}

//Positionfile : position file data format
type Positionfile struct {
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
	logger.Info("Period set to %s", c.Period.String())
	sophos := sophoscentral.NewAPIClient(sophoscentralConfig)
	auth := context.WithValue(context.Background(), sophoscentral.ContextAPIKey, sophoscentral.APIKey{
		Key: c.APIKey,
	})
	bt := &Sophoscentralbeat{
		done:       make(chan struct{}),
		sophos:     sophos,
		sophosAuth: auth,
		config:     c,
		logger:     *logger,
		basepath:   c.Basepath,
	}
	return bt, nil
}

//GetSophosEvents : calls Sophos Events Api
func GetSophosEvents(scb Sophoscentralbeat) ([]sophoscentral.LegacyEventEntity, error) {
	scb.logger.Info("Making sophos event call")
	var items []sophoscentral.LegacyEventEntity

	timeStamp, posFileStatus := ReadTimeStamp()

	from := GenerateYesterdayTimeStamp()

	if posFileStatus != false && timeStamp.EventsTimestamp > from && timeStamp.EventsTimestamp != 0 {
		from = timeStamp.EventsTimestamp
	}

	options := &sophoscentral.GetEventsUsingGET1Opts{
		Limit:    optional.NewInt32(1000),
		FromDate: optional.NewInt64(from),
	}

	value, _, err := scb.sophos.EventControllerV1ImplApi.GetEventsUsingGET1(scb.sophosAuth, scb.config.APIKey, scb.config.Authorization, scb.basepath, options)
	if err != nil {
		scb.logger.Error(err)
		return nil, err
	}

	//update timestamp
	WriteTimeStamp(time.Now().Unix(), 0)

	for _, item := range value.Items {
		items = append(items, item)
	}
	for value.HasMore == true {
		options.Cursor = optional.NewString(value.NextCursor)
		value, _, err := scb.sophos.EventControllerV1ImplApi.GetEventsUsingGET1(scb.sophosAuth, scb.config.APIKey, scb.config.Authorization, scb.basepath, options)
		if err != nil {
			scb.logger.Error(err)
			return nil, err
		}
		for _, item := range value.Items {
			items = append(items, item)
		}
	}
	return value.Items, nil
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

//GetSophosAlerts : call alerts API
func GetSophosAlerts(scb Sophoscentralbeat) ([]sophoscentral.AlertEntity, error) {
	scb.logger.Info("Making sophos alert call")
	var items []sophoscentral.AlertEntity

	timeStamp, posFileStatus := ReadTimeStamp()

	from := GenerateYesterdayTimeStamp()

	if posFileStatus != false && timeStamp.AlertsTimestamp > from && timeStamp.AlertsTimestamp != 0 {
		from = timeStamp.AlertsTimestamp
	}

	options := &sophoscentral.GetAlertsUsingGET1Opts{
		Limit:    optional.NewInt32(1000),
		FromDate: optional.NewInt64(from),
	}
	value, _, err := scb.sophos.AlertControllerV1ImplApi.GetAlertsUsingGET1(scb.sophosAuth, scb.config.APIKey, scb.config.Authorization, scb.basepath, options)
	if err != nil {
		scb.logger.Error(err)
		return nil, err
	}

	WriteTimeStamp(0, time.Now().Unix())
	for _, item := range value.Items {
		items = append(items, item)
	}
	for value.HasMore == true {
		options.Cursor = optional.NewString(value.NextCursor)
		value, _, err := scb.sophos.AlertControllerV1ImplApi.GetAlertsUsingGET1(scb.sophosAuth, scb.config.APIKey, scb.config.Authorization, scb.basepath, options)
		if err != nil {
			scb.logger.Error(err)
			return nil, err
		}
		for _, item := range value.Items {
			items = append(items, item)
		}
	}
	return items, nil
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
		events, err := GetSophosEvents(*scb)
		if err != nil {
			scb.logger.Error(err)
		}
		var toSend []beat.Event
		for _, event := range events {
			beatEvent := beat.Event{
				Timestamp: time.Now(),
				Fields: common.MapStr{
					"response": event,
				},
			}
			toSend = append(toSend, beatEvent)
		}

		scb.logger.Info("Attempting to fetch Sophos Alerts")
		alerts, err := GetSophosAlerts(*scb)
		if err != nil {
			scb.logger.Error(err)
		}
		for _, alert := range alerts {
			event := beat.Event{
				Timestamp: time.Now(),
				Fields: common.MapStr{
					"response": alert,
				},
			}
			toSend = append(toSend, event)
		}
		scb.client.PublishAll(toSend)
		scb.logger.Info("Events sent")
	}
}

// Stop stops sophoscentralbeat.
func (scb *Sophoscentralbeat) Stop() {
	scb.client.Close()
	close(scb.done)
}

//WriteTimeStamp : writes timestamp to file
func WriteTimeStamp(eventTimeStamp int64, alertTimeStamp int64) {

	filePath := "data/pos.json"
	var position Positionfile

	//position file unavailable
	if _, err := os.Stat(filePath); err != nil {

		position = Positionfile{
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
func ReadTimeStamp() (Positionfile, bool) {
	filePath := "data/pos.json"

	var pos Positionfile
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
	now := time.Now().UTC()

	//1 day = 1440 minutes (for exact 24 hours)
	count := 1440

	//generate 24 hours prior timestamp to compare with stored timestamp
	return now.Add(time.Duration(-count) * time.Minute).Unix()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
