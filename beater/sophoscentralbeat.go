package beater

import (
	"context"
	"fmt"
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

func GetSophosEvents(scb Sophoscentralbeat) ([]sophoscentral.LegacyEventEntity, error) {
	scb.logger.Info("Making sophos event call")
	var items []sophoscentral.LegacyEventEntity
	now := time.Now().UTC()
	from := now.Add(scb.config.Period * -1)
	options := &sophoscentral.GetEventsUsingGET1Opts{
		Limit:    optional.NewInt32(1000),
		FromDate: optional.NewInt64(from.Unix()),
	}
	value, _, err := scb.sophos.EventControllerV1ImplApi.GetEventsUsingGET1(scb.sophosAuth, scb.config.APIKey, scb.config.Authorization, scb.basepath, options)
	if err != nil {
		scb.logger.Error(err)
		return nil, err
	}
	for _, item := range value.Items {
		// fmt.Println(item)
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

func GetSophosAlerts(scb Sophoscentralbeat) ([]sophoscentral.AlertEntity, error) {
	scb.logger.Info("Making sophos alert call")
	var items []sophoscentral.AlertEntity
	now := time.Now().UTC()
	from := now.Add(scb.config.Period * -1)
	options := &sophoscentral.GetAlertsUsingGET1Opts{
		Limit:    optional.NewInt32(1000),
		FromDate: optional.NewInt64(from.Unix()),
	}
	value, _, err := scb.sophos.AlertControllerV1ImplApi.GetAlertsUsingGET1(scb.sophosAuth, scb.config.APIKey, scb.config.Authorization, scb.basepath, options)
	if err != nil {
		scb.logger.Error(err)
		return nil, err
	}
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
func (bt *Sophoscentralbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
