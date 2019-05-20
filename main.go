package main

import (
	"os"

	"github.com/forter/sophoscentralbeat/cmd"

	_ "github.com/forter/sophoscentralbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
