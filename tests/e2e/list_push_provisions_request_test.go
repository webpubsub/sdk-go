package e2e

import (
	"testing"

	webpubsub "github.com/webpubsub/go/v7"

	"github.com/stretchr/testify/assert"
)

func TestListPushProvisionsNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := webpubsub.NewWebPubSub(configCopy())
	ch1 := randomized("testChannel_sub_")
	cg1 := randomized("testCG_sub_")

	_, _, err := pn.AddPushNotificationsOnChannels().
		Channels([]string{ch1}).
		DeviceIDForPush(cg1).
		PushType(webpubsub.WPSPushTypeGCM).
		Execute()

	assert.Nil(err)

	resp, _, err := pn.ListPushProvisions().
		DeviceIDForPush(cg1).
		PushType(webpubsub.WPSPushTypeGCM).
		Execute()
	assert.Contains(resp.Channels, ch1)
	assert.Nil(err)
}

func TestListPushProvisionsNotStubbedContext(t *testing.T) {
	assert := assert.New(t)

	pn := webpubsub.NewWebPubSub(configCopy())
	ch1 := randomized("testChannel_sub_")
	cg1 := randomized("testCG_sub_")

	_, _, err := pn.AddPushNotificationsOnChannelsWithContext(backgroundContext).
		Channels([]string{ch1}).
		DeviceIDForPush(cg1).
		PushType(webpubsub.WPSPushTypeGCM).
		Execute()

	assert.Nil(err)

	resp, _, err := pn.ListPushProvisionsWithContext(backgroundContext).
		DeviceIDForPush(cg1).
		PushType(webpubsub.WPSPushTypeGCM).
		Execute()
	assert.Contains(resp.Channels, ch1)
	assert.Nil(err)
}

func TestListPushProvisionsTopicAndEnvNotStubbed(t *testing.T) {
	assert := assert.New(t)

	pn := webpubsub.NewWebPubSub(configCopy())
	ch1 := randomized("testChannel_sub_")
	cg1 := randomized("testCG_sub_")

	_, _, err := pn.AddPushNotificationsOnChannels().
		Channels([]string{ch1}).
		DeviceIDForPush(cg1).
		PushType(webpubsub.WPSPushTypeGCM).
		Execute()

	assert.Nil(err)

	resp, _, err := pn.ListPushProvisions().
		DeviceIDForPush(cg1).
		PushType(webpubsub.WPSPushTypeGCM).
		Topic("a").
		Environment(webpubsub.WPSPushEnvironmentProduction).
		Execute()
	assert.Contains(resp.Channels, ch1)
	assert.Nil(err)
}

func TestListPushProvisionsTopicAndEnvNotStubbedContext(t *testing.T) {
	assert := assert.New(t)

	pn := webpubsub.NewWebPubSub(configCopy())
	ch1 := randomized("testChannel_sub_")
	cg1 := randomized("testCG_sub_")

	_, _, err := pn.AddPushNotificationsOnChannelsWithContext(backgroundContext).
		Channels([]string{ch1}).
		DeviceIDForPush(cg1).
		PushType(webpubsub.WPSPushTypeGCM).
		Topic("a").
		Environment(webpubsub.WPSPushEnvironmentProduction).
		Execute()

	assert.Nil(err)

	resp, _, err := pn.ListPushProvisionsWithContext(backgroundContext).
		DeviceIDForPush(cg1).
		PushType(webpubsub.WPSPushTypeGCM).
		Topic("a").
		Environment(webpubsub.WPSPushEnvironmentProduction).
		Execute()
	assert.Contains(resp.Channels, ch1)
	assert.Nil(err)
}
