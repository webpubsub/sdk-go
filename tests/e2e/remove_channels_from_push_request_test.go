package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	webpubsub "github.com/webpubsub/sdk-go/v7"
)

func TestRemovePushNotificationsFromChannels(t *testing.T) {
	assert := assert.New(t)

	pn := webpubsub.NewWebPubSub(configCopy())

	_, _, err := pn.RemovePushNotificationsFromChannels().
		Channels([]string{"ch"}).
		DeviceIDForPush("cg").
		PushType(webpubsub.WPSPushTypeGCM).
		Execute()
	assert.Nil(err)
}

func TestRemovePushNotificationsFromChannelsContext(t *testing.T) {
	assert := assert.New(t)

	pn := webpubsub.NewWebPubSub(configCopy())

	_, _, err := pn.RemovePushNotificationsFromChannelsWithContext(backgroundContext).
		Channels([]string{"ch"}).
		DeviceIDForPush("cg").
		PushType(webpubsub.WPSPushTypeGCM).
		Execute()
	assert.Nil(err)
}
func TestRemovePushNotificationsFromChannelsTopicAndEnv(t *testing.T) {
	assert := assert.New(t)

	pn := webpubsub.NewWebPubSub(configCopy())

	_, _, err := pn.RemovePushNotificationsFromChannels().
		Channels([]string{"ch"}).
		DeviceIDForPush("cg").
		PushType(webpubsub.WPSPushTypeGCM).
		Topic("a").
		Environment(webpubsub.WPSPushEnvironmentProduction).
		Execute()
	assert.Nil(err)
}

func TestRemovePushNotificationsFromChannelsTopicAndEnvContext(t *testing.T) {
	assert := assert.New(t)

	pn := webpubsub.NewWebPubSub(configCopy())

	_, _, err := pn.RemovePushNotificationsFromChannelsWithContext(backgroundContext).
		Channels([]string{"ch"}).
		DeviceIDForPush("cg").
		PushType(webpubsub.WPSPushTypeGCM).
		Topic("a").
		Environment(webpubsub.WPSPushEnvironmentProduction).
		Execute()
	assert.Nil(err)
}
