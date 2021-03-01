package internal

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wostzone/gateway/pkg/lib"
	"github.com/wostzone/gateway/pkg/messaging"
)

// WostLoggerConfig with logger plugin configuration
type WostLoggerConfig struct {
	Channels []string `yaml:"channels"`
}

// WostLogger is a gateway plugin for recording channels
type WostLogger struct {
	config      WostLoggerConfig
	gwConfig    *lib.GatewayConfig
	messenger   messaging.IGatewayMessenger
	fileHandles map[string]*os.File
}

// LoggerPluginID is the default ID of the WoST Logger plugin
const LoggerPluginID = "logger"

// handleChannelMessage receives and records a channel message
func (wlog *WostLogger) handleChannelMessage(channel string, message []byte) {
	logrus.Infof("handleChannelMessage: Received message on channel %s: %s", channel, message)
	fileHandle := wlog.fileHandles[channel]
	if fileHandle != nil {
		sender := ""
		timeStamp := time.Now().Format("2006-01-02T15:04:05.000Z07:00")
		// timeStamp := time.Now().Format(time.RFC3339Nano)
		maxLen := len(message)
		if maxLen > 40 {
			maxLen = 40
		}
		line := fmt.Sprintf("[%s] %s %s: %s", timeStamp, sender, channel, message[:maxLen])
		n, err := fileHandle.WriteString(line + "\n")
		_ = n
		if err != nil {
			logrus.Errorf("handleChannelMessage: Unable to record channel '%s': %s", channel, err)
		}
	}
}

// StartRecordChannel setup recording of a channel
func (wlog *WostLogger) StartRecordChannel(channel string, messenger messaging.IGatewayMessenger) {
	logsFolder := path.Dir(wlog.gwConfig.Logging.LogFile)
	filename := path.Join(logsFolder, channel+".txt")
	fileHandle, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0640)

	if err != nil {
		logrus.Errorf("StartRecordChannel: Unable to open file '%s' for writing: %s. Channel '%s' ignored", filename, err, channel)
		return
	}
	wlog.fileHandles[channel] = fileHandle
	messenger.Subscribe(channel, wlog.handleChannelMessage)
}

// Start connects, subscribe and start the recording
func (wlog *WostLogger) Start(gwConfig *lib.GatewayConfig, recConfig *WostLoggerConfig) error {
	var err error
	wlog.config = *recConfig
	wlog.gwConfig = gwConfig
	wlog.messenger, err = messaging.StartGatewayMessenger(LoggerPluginID, gwConfig)

	// messaging.NewGatewayConnection(gwConfig.Messenger.Protocol)
	// load the default channels if config doesn't have any
	if wlog.config.Channels == nil || len(wlog.config.Channels) == 0 {
		wlog.config.Channels = []string{messaging.TDChannelID, messaging.EventsChannelID, messaging.ActionChannelID}
	}
	for _, channel := range wlog.config.Channels {
		wlog.StartRecordChannel(channel, wlog.messenger)
	}

	logrus.Infof("Start: channels: %s", wlog.config.Channels)
	return err
}

// Stop the logging
func (wlog *WostLogger) Stop() {
	logrus.Info("Recorder Stop: Stopping recorder service")
	for _, channel := range wlog.config.Channels {
		wlog.messenger.Unsubscribe(channel)
	}
	for _, fileHandle := range wlog.fileHandles {
		fileHandle.Close()
	}
	wlog.fileHandles = make(map[string]*os.File)

}

// NewWostLoggerService creates a new logger of WoST gateway messages
func NewWostLoggerService() *WostLogger {
	wlog := &WostLogger{
		fileHandles: make(map[string]*os.File),
	}
	return wlog
}
