package heartbeat

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/logp"
)

// Config is a structure for heartbeat
type Config struct {
	Interval time.Duration
	Disabled bool
}

// IntervalValue is a default value for heartbeat interval
var IntervalValue = 5 * time.Minute

// NewHeartbeatConfig is a constructor to return the object of heartbeatConfig structure
func NewHeartbeatConfig(interval time.Duration, disabled bool) *Config {
	return &Config{
		Interval: interval,
		Disabled: disabled,
	}
}

// CreateEnabled will create all miscellaneous components
func (config *Config) CreateEnabled(doneChan chan struct{}, serviceName string) (*StatusBeater, error) {
	if config == nil {
		return nil, fmt.Errorf("no heartbeat specified. To disable, specify 'disabled: true' in the heartbeat configuration")
	}

	if config.Disabled {
		// Customer has explicitly disabled heart beating
		return nil, nil
	}

	if config.Interval <= 0 {
		// Shouldn't happen in regular code path because of our defaults / validation
		logp.Warn("Heartbeat interval can not be less than zero. Setting to default 5 minute")
		config.Interval = IntervalValue
	}

	return NewStatusBeater(serviceName, config.Interval, doneChan), nil
}
