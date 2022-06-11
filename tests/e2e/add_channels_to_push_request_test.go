package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	webpubsub "github.com/webpubsub/sdk-go/v7"
)

func TestAddChannelToPushNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := webpubsub.NewWebPubSub(configCopy())

	_, _, err := pn.AddPushNotificationsOnChannels().
		Channels([]string{"ch"}).
		DeviceIDForPush("cg").
		PushType(webpubsub.WPSPushTypeGCM).
		Execute()
	assert.Nil(err)
}

func TestAddChannelToPushNotStubbedContext(t *testing.T) {
	assert := assert.New(t)

	pn := webpubsub.NewWebPubSub(configCopy())

	_, _, err := pn.AddPushNotificationsOnChannelsWithContext(backgroundContext).
		Channels([]string{"ch1"}).
		DeviceIDForPush("cg1").
		PushType(webpubsub.WPSPushTypeGCM).
		Execute()
	assert.Nil(err)
}

func TestAddChannelToPushNotStubbedContextWithQueryParam(t *testing.T) {
	assert := assert.New(t)

	pn := webpubsub.NewWebPubSub(configCopy())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	_, _, err := pn.AddPushNotificationsOnChannelsWithContext(backgroundContext).
		Channels([]string{"ch1"}).
		DeviceIDForPush("cg1").
		PushType(webpubsub.WPSPushTypeGCM).
		QueryParam(queryParam).
		Execute()
	assert.Nil(err)
}

func TestAddChannelToPushTopicAndEnvNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := webpubsub.NewWebPubSub(configCopy())

	_, _, err := pn.AddPushNotificationsOnChannels().
		Channels([]string{"ch"}).
		DeviceIDForPush("cg").
		PushType(webpubsub.WPSPushTypeGCM).
		Topic("a").
		Environment(webpubsub.WPSPushEnvironmentDevelopment).
		Execute()
	assert.Nil(err)
}

func TestAddChannelToPushTopicAndEnvNotStubbedContext(t *testing.T) {
	assert := assert.New(t)

	pn := webpubsub.NewWebPubSub(configCopy())

	_, _, err := pn.AddPushNotificationsOnChannelsWithContext(backgroundContext).
		Channels([]string{"ch1"}).
		DeviceIDForPush("cg1").
		PushType(webpubsub.WPSPushTypeGCM).
		Topic("a").
		Environment(webpubsub.WPSPushEnvironmentDevelopment).
		Execute()
	assert.Nil(err)
}

func TestAddChannelToPushTopicAndEnvNotStubbedContextWithQueryParam(t *testing.T) {
	assert := assert.New(t)

	pn := webpubsub.NewWebPubSub(configCopy())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	_, _, err := pn.AddPushNotificationsOnChannelsWithContext(backgroundContext).
		Channels([]string{"ch1"}).
		DeviceIDForPush("cg1").
		PushType(webpubsub.WPSPushTypeGCM).
		QueryParam(queryParam).
		Topic("a").
		Environment(webpubsub.WPSPushEnvironmentProduction).
		Execute()
	assert.Nil(err)
}
