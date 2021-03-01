package internal_test

import (
	"os"
	"path"
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

// Use the project app folder during testing
func init() {
	cwd, _ := os.Getwd()
	homeFolder = path.Join(cwd, "../dist")
}

func TestStartStopRecorder(t *testing.T) {

	recConfig := &internal.WostLoggerConfig{} // use defaults
	gwConfig, err := lib.SetupConfig(homeFolder, pluginID, recConfig)
	assert.NoError(t, err)
	server, err := smbserver.StartSmbServer(gwConfig)
	require.NoError(t, err)

	svc := internal.NewWostLoggerService()
	err = svc.Start(gwConfig, recConfig)
	assert.NoError(t, err)
	svc.Stop()
	server.Stop()
}

func TestRecordMessage(t *testing.T) {

	recConfig := &internal.WostLoggerConfig{} // use defaults
	gwConfig, err := lib.SetupConfig(homeFolder, pluginID, recConfig)
	assert.NoError(t, err)
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
}
