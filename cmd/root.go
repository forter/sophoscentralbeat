package cmd

import (
	"github.com/forter/sophoscentralbeat/beater"

	cmd "github.com/elastic/beats/libbeat/cmd"
	"github.com/elastic/beats/libbeat/cmd/instance"
)

// Name of this beat
var Name = "sophoscentralbeat"

// RootCmd to handle beats cli
var RootCmd = cmd.GenRootCmdWithSettings(beater.New, instance.Settings{Name: Name})
