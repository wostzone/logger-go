package internal_test

import (
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wostzone/gateway/pkg/lib"
	"github.com/wostzone/gateway/pkg/messaging"
	"github.com/wostzone/gateway/pkg/messaging/smbserver"
	"github.com/wostzone/logger/internal"
)

var homeFolder string

const pluginID = "logger-test"

var recConfig *internal.WostLoggerConfig = &internal.WostLoggerConfig{} // use defaults
var gwConfig *lib.GatewayConfig
var setupOnce = false

// Use the project app folder during testing
func setup() {
	if setupOnce {
		return
	}
	setupOnce = true
	cwd, _ := os.Getwd()
	homeFolder = path.Join(cwd, "../dist")
	recConfig = &internal.WostLoggerConfig{}
	os.Args = append(os.Args[0:1], strings.Split("", " ")...)
	gwConfig, _ = lib.SetupConfig(homeFolder, pluginID, recConfig)
}
func teardown() {
}

func TestStartStopRecorder(t *testing.T) {
	setup()
	// recConfig := &internal.WostLoggerConfig{} // use defaults
	// gwConfig, err := lib.SetupConfig(homeFolder, pluginID, recConfig)
	// assert.NoError(t, err)
	server, err := smbserver.StartSmbServer(gwConfig)
	require.NoError(t, err)

	svc := internal.NewWostLoggerService()
	err = svc.Start(gwConfig, recConfig)
	assert.NoError(t, err)
	svc.Stop()
	server.Stop()
	teardown()
}

func TestRecordMessage(t *testing.T) {
	setup()

	// recConfig := &internal.WostLoggerConfig{} // use defaults
	// os.Args = append(os.Args[0:1], strings.Split("", " ")...)
	// gwConfig, err := lib.SetupConfig(homeFolder, pluginID, recConfig)
	// lib.LoadConfig(gwConfig.ConfigFolder+"/gateway.yaml", gwConfig)
	// gwConfig.Messenger.HostPort = "localhost:9999"

	// assert.NoError(t, err)
	server, err := smbserver.StartSmbServer(gwConfig)
	require.NoError(t, err)

	svc := internal.NewWostLoggerService()
	err = svc.Start(gwConfig, recConfig)
	client, err := messaging.StartGatewayMessenger("test1", gwConfig)
	assert.NoError(t, err)

	client.Publish(messaging.EventsChannelID, []byte("Hello world"))
	time.Sleep(1 * time.Second)
	client.Disconnect()

	assert.NoError(t, err)
	svc.Stop()
	server.Stop()
	teardown()
}
