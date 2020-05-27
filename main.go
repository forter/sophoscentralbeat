package main

import (
	"os"

	"github.com/logrhythm/sophoscentralbeat/cmd"

	_ "github.com/logrhythm/sophoscentralbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
