package heartbeat

import (
	"encoding/json"
	"os"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/golang/protobuf/ptypes/timestamp"
)

// Heartbeat is a structure for heartbeat
type Heartbeat struct {
	// Service name
	ServiceName string `json:"service_name"`
	// Current version of the service
	ServiceVersion string `json:"service_version"`

	Time timestamp.Timestamp `json:"time"`

	Status Status `json:"status"`
}

const (
	//ServiceStarted is a code for starting a particular service
	ServiceStarted = 1
	//ServiceRunning is a code for running instance a particular service
	ServiceRunning = 2
	//ServiceStopped is a code for stopping a particular service
	ServiceStopped = 3

	// FQBeatName variable name for fully qualified beat name
	FQBeatName = "FullyQualifiedBeatName"
)

// fqBeatName is the fully qualified beat name
var fqBeatName string

// Status is used for status of heartbeat1
type Status struct {
	Code        int64  `json:"code"`
	Description string `json:"description"`
}

// IntervalFunc is a function that can trigger a timing event based on a duration
type IntervalFunc func() <-chan time.Time

// StatusBeater reports simple service information
type StatusBeater struct {
	Name    string
	Version string

	IntervalFunc IntervalFunc
	doneChan     chan struct{}
}

// Start will begin reporting heartbeats through the beats
func (sb *StatusBeater) Start(stopChan chan struct{}, publish func(event beat.Event)) {
	go func() {
		fqBeatName = os.Getenv(FQBeatName)
		sb.Beat(ServiceStarted, "Service started", publish)
		for {
			select {
			case <-sb.IntervalFunc():
				sb.Beat(ServiceRunning, "Service is Running", publish)
			case <-stopChan:
				sb.Beat(ServiceStopped, "Service is Stopped", publish)
				sb.doneChan <- struct{}{}
				return
			}
		}
	}()
}

// Beat will send a beat containing simple service status information
func (sb *StatusBeater) Beat(status int64, description string, publish func(event beat.Event)) {
	now := time.Now().UnixNano()
	msg := Heartbeat{
		ServiceName:    sb.Name,
		ServiceVersion: sb.Version,
		Time: timestamp.Timestamp{
			Seconds: now / time.Nanosecond.Nanoseconds(),
		},
		Status: Status{
			Code:        status,
			Description: description,
		},
	}
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		logp.Warn("internal heartbeat message json conversion failed %s", err)
		return
	}
	sb.PublishEvent(msgJSON, publish)

}

// PublishEvent will publish passed Log
func (sb *StatusBeater) PublishEvent(logData []byte, publish func(event beat.Event)) {
	event := beat.Event{
		Timestamp: time.Now(),
		Fields: common.MapStr{
			"heartbeat":              string(logData),
			"fullyqualifiedbeatname": fqBeatName,
		},
	}
	publish(event)
	logp.Info("heartbeat sent")
	logp.Debug("Fully Qualified Beatname: %s", fqBeatName)
}

// NewStatusBeater will return a new StatusBeater with the provided base information
func NewStatusBeater(serviceName string, interval time.Duration, doneChan chan struct{}) *StatusBeater {
	return NewStatusBeaterWithFunc(
		serviceName,
		func() <-chan time.Time {
			return time.After(interval)
		},
		doneChan,
	)
}

// NewStatusBeaterWithFunc returns a new StatusBeater that uses the provided func as a trigger for sending beats
func NewStatusBeaterWithFunc(serviceName string, intervalFunc IntervalFunc, doneChan chan struct{}) *StatusBeater {
	return &StatusBeater{
		Name:         serviceName,
		IntervalFunc: intervalFunc,
		doneChan:     doneChan,
	}
}
