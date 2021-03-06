package internal_test

import (
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/wostzone/hub/pkg/config"
	"github.com/wostzone/hub/pkg/messaging"
	"github.com/wostzone/hub/pkg/smbserver"
	"github.com/wostzone/logger/internal"
)

var homeFolder string

const pluginID = "logger-test"
const loremIpsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor " +
	"incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco " +
	"laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate " +
	"velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, " +
	"sunt in culpa qui officia deserunt mollit anim id est laborum."

var recConfig *internal.WostLoggerConfig = &internal.WostLoggerConfig{} // use defaults
var gwConfig *config.HubConfig
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
	gwConfig, _ = config.SetupConfig(homeFolder, pluginID, recConfig)
}
func teardown() {
}

func TestStartStop(t *testing.T) {
	setup()
	// recConfig := &internal.WostLoggerConfig{} // use defaults
	// gwConfig, err := lib.SetupConfig(homeFolder, pluginID, recConfig)
	// assert.NoError(t, err)
	server, err := smbserver.StartSmbServer(gwConfig)
	require.NoError(t, err)

	svc := internal.WostLogger{}
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
	// lib.LoadConfig(gwConfig.ConfigFolder+"/hub.yaml", gwConfig)
	// gwConfig.Messenger.HostPort = "localhost:9999"

	// assert.NoError(t, err)
	server, err := smbserver.StartSmbServer(gwConfig)
	require.NoError(t, err)

	svc := internal.WostLogger{}
	err = svc.Start(gwConfig, recConfig)
	client, err := messaging.StartHubMessenger("test1", gwConfig)
	assert.NoError(t, err)

	client.Publish(messaging.EventsChannelID, []byte("Hello world"))
	client.Publish(messaging.EventsChannelID, []byte(loremIpsum))
	time.Sleep(1 * time.Second)
	client.Disconnect()

	assert.NoError(t, err)
	svc.Stop()
	server.Stop()
	teardown()
}

func TestBadLoggingFolder(t *testing.T) {
	setup()

	// recConfig := &internal.WostLoggerConfig{} // use defaults
	// os.Args = append(os.Args[0:1], strings.Split("", " ")...)
	// gwConfig, err := lib.SetupConfig(homeFolder, pluginID, recConfig)
	// lib.LoadConfig(gwConfig.ConfigFolder+"/hub.yaml", gwConfig)
	// gwConfig.Messenger.HostPort = "localhost:9999"

	// assert.NoError(t, err)
	server, err := smbserver.StartSmbServer(gwConfig)
	require.NoError(t, err)

	svc := internal.WostLogger{}
	gwConfig.Logging.LogFile = "/notafolder"
	err = svc.Start(gwConfig, recConfig)
	assert.Error(t, err)

	client, err := messaging.StartHubMessenger("test1", gwConfig)
	assert.NoError(t, err)

	client.Publish(messaging.EventsChannelID, []byte("Hello world"))
	client.Publish(messaging.EventsChannelID, []byte(loremIpsum))
	time.Sleep(1 * time.Second)
	client.Disconnect()

	assert.NoError(t, err)
	svc.Stop()
	server.Stop()
	teardown()
}
