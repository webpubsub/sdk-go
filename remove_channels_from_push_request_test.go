package webpubsub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveChannelsFromPushRequestValidate(t *testing.T) {
	assert := assert.New(t)

	opts := &removeChannelsFromPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
	}

	err := opts.validate()
	assert.Nil(err)

	opts1 := &removeChannelsFromPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeNone,
		webpubsub:       webpubsub,
	}

	err1 := opts1.validate()
	assert.Contains(err1.Error(), "Missing Push Type")

	opts2 := &removeChannelsFromPushOpts{
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
	}

	err2 := opts2.validate()
	assert.Contains(err2.Error(), "Missing Channel")

	opts3 := &removeChannelsFromPushOpts{
		Channels:  []string{"ch1", "ch2", "ch3"},
		PushType:  WPSPushTypeAPNS,
		webpubsub: webpubsub,
	}

	err3 := opts3.validate()
	assert.Contains(err3.Error(), "Missing Device ID")
}

func TestRemoveChannelsFromPushRequestBuildPath(t *testing.T) {
	assert := assert.New(t)

	opts := &removeChannelsFromPushOpts{
		DeviceIDForPush: "deviceId",
		Channels:        []string{"ch1", "ch2", "ch3"},
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
	}

	str, err := opts.buildPath()
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceId", str)
	assert.Nil(err)

}

func TestRemoveChannelsFromPushRequestBuildQueryParam(t *testing.T) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := &removeChannelsFromPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
		QueryParam:      queryParam,
	}

	u, err := opts.buildQuery()
	assert.Equal("ch1,ch2,ch3", u.Get("remove"))
	assert.Equal("apns", u.Get("type"))
	assert.Equal("v1", u.Get("q1"))
	assert.Equal("v2", u.Get("q2"))

	assert.Nil(err)
}

func TestRemoveChannelsFromPushRequestBuildQueryParamTopicAndEnv(t *testing.T) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := &removeChannelsFromPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
		QueryParam:      queryParam,
		Topic:           "a",
		Environment:     WPSPushEnvironmentProduction,
	}

	u, err := opts.buildQuery()
	assert.Equal("ch1,ch2,ch3", u.Get("remove"))
	assert.Equal("apns", u.Get("type"))
	assert.Equal("v1", u.Get("q1"))
	assert.Equal("v2", u.Get("q2"))
	assert.Equal("production", u.Get("environment"))
	assert.Equal("a", u.Get("topic"))

	assert.Nil(err)
}

func TestRemoveChannelsFromPushRequestBuildQuery(t *testing.T) {
	assert := assert.New(t)

	opts := &removeChannelsFromPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
	}

	u, err := opts.buildQuery()
	assert.Equal("ch1,ch2,ch3", u.Get("remove"))
	assert.Equal("apns", u.Get("type"))

	assert.Nil(err)
}

func TestRemoveChannelsFromPushRequestBuildBody(t *testing.T) {
	assert := assert.New(t)

	opts := &removeChannelsFromPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
	}

	_, err := opts.buildBody()
	assert.Nil(err)

}

func TestNewRemoveChannelsFromPushBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newRemoveChannelsFromPushBuilder(webpubsub)
	o.Channels([]string{"ch1", "ch2", "ch3"})
	o.DeviceIDForPush("deviceId")
	o.PushType(WPSPushTypeAPNS)
	u, err := o.opts.buildQuery()
	assert.Equal("ch1,ch2,ch3", u.Get("remove"))
	assert.Equal("apns", u.Get("type"))
	assert.Nil(err)
}

func TestNewRemoveChannelsFromPushBuilderWithContext(t *testing.T) {
	assert := assert.New(t)

	o := newRemoveChannelsFromPushBuilderWithContext(webpubsub, backgroundContext)
	o.Channels([]string{"ch1", "ch2", "ch3"})
	o.DeviceIDForPush("deviceId")
	o.PushType(WPSPushTypeAPNS)
	u, err := o.opts.buildQuery()
	assert.Equal("ch1,ch2,ch3", u.Get("remove"))
	assert.Equal("apns", u.Get("type"))
	assert.Nil(err)

}

func TestRemChannelsFromPushValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &removeChannelsFromPushOpts{
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       pn,
	}

	assert.Equal("webpubsub/validation: webpubsub: Remove Channel Group: Missing Subscribe Key", opts.validate().Error())
}
