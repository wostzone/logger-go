package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/wostzone/gateway/pkg/lib"
	"github.com/wostzone/logger/internal"
)

var loggerConfig = &internal.WostLoggerConfig{}

func main() {
	gatewayConfig, err := lib.SetupConfig("", internal.LoggerPluginID, loggerConfig)

	svc := internal.NewWostLoggerService()
	err = svc.Start(gatewayConfig, loggerConfig)
	if err != nil {
		logrus.Errorf("Logger: Failed to start")
		os.Exit(1)
	}
	lib.WaitForSignal()
	svc.Stop()
	os.Exit(0)
}
