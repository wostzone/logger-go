package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/wostzone/hub/pkg/config"
	"github.com/wostzone/hub/pkg/hub"
	"github.com/wostzone/logger/internal"
)

var pluginConfig = &internal.WostLoggerConfig{}

func main() {
	hubConfig, err := config.SetupConfig("", internal.PluginID, pluginConfig)

	svc := internal.WostLogger{}
	err = svc.Start(hubConfig, pluginConfig)
	if err != nil {
		logrus.Errorf("Logger: Failed to start")
		os.Exit(1)
	}
	hub.WaitForSignal()
	svc.Stop()
	os.Exit(0)
}
