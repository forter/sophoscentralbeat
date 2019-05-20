// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import "time"

const (
	Filler = "FILLER"
)

type Config struct {
	Period        time.Duration `config:"period"`
	APIKey        string        `config:"api_key"`
	Authorization string        `config:"authorization"`
}

var DefaultConfig = Config{
	Period:        300 * time.Second,
	APIKey:        Filler,
	Authorization: Filler,
}
