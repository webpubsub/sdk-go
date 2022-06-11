package webpubsub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddChannelsToPushOptsValidate(t *testing.T) {
	assert := assert.New(t)

	opts := &addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
	}

	err := opts.validate()
	assert.Nil(err)

	opts1 := &addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeNone,
		webpubsub:       webpubsub,
	}

	err1 := opts1.validate()
	assert.Contains(err1.Error(), "Missing Push Type")

	opts2 := &addChannelsToPushOpts{
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
	}

	err2 := opts2.validate()
	assert.Contains(err2.Error(), "Missing Channel")

	opts3 := &addChannelsToPushOpts{
		Channels:  []string{"ch1", "ch2", "ch3"},
		PushType:  WPSPushTypeAPNS,
		webpubsub: webpubsub,
	}

	err3 := opts3.validate()
	assert.Contains(err3.Error(), "Missing Device ID")

}

func TestAddChannelsToPushOptsBuildPath(t *testing.T) {
	assert := assert.New(t)

	opts := &addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
	}

	str, err := opts.buildPath()
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceId", str)
	assert.Nil(err)

}

func TestAddChannelsToPushOptsBuildQuery(t *testing.T) {
	assert := assert.New(t)

	opts := &addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
	}

	u, err := opts.buildQuery()
	assert.Equal("ch1,ch2,ch3", u.Get("add"))
	assert.Equal("apns", u.Get("type"))
	assert.Nil(err)
}

func TestAddChannelsToPushOptsBuildQueryParams(t *testing.T) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := &addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
		QueryParam:      queryParam,
	}

	u, err := opts.buildQuery()
	assert.Equal("ch1,ch2,ch3", u.Get("add"))
	assert.Equal("apns", u.Get("type"))
	assert.Equal("v1", u.Get("q1"))
	assert.Equal("v2", u.Get("q2"))
	assert.Nil(err)
}

func TestAddChannelsToPushOptsBuildQueryParamsTopicAndEnv(t *testing.T) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := &addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
		QueryParam:      queryParam,
		Topic:           "a",
		Environment:     WPSPushEnvironmentProduction,
	}

	u, err := opts.buildQuery()
	assert.Equal("ch1,ch2,ch3", u.Get("add"))
	assert.Equal("apns", u.Get("type"))
	assert.Equal("v1", u.Get("q1"))
	assert.Equal("v2", u.Get("q2"))
	assert.Equal("production", u.Get("environment"))
	assert.Equal("a", u.Get("topic"))
	assert.Nil(err)
}

func TestAddChannelsToPushOptsBuildBody(t *testing.T) {
	assert := assert.New(t)

	opts := &addChannelsToPushOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
	}

	_, err := opts.buildBody()

	assert.Nil(err)

}

func TestNewAddPushNotificationsOnChannelsBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newAddPushNotificationsOnChannelsBuilder(webpubsub)
	o.Channels([]string{"ch1", "ch2", "ch3"})
	o.DeviceIDForPush("deviceID")
	o.PushType(WPSPushTypeAPNS)

	path, err := o.opts.buildPath()
	assert.Nil(err)
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceID", path)
}

func TestNewAddPushNotificationsOnChannelsBuilderWithContext(t *testing.T) {
	assert := assert.New(t)

	o := newAddPushNotificationsOnChannelsBuilderWithContext(webpubsub, backgroundContext)
	o.Channels([]string{"ch1", "ch2", "ch3"})
	o.DeviceIDForPush("deviceID")
	o.PushType(WPSPushTypeAPNS)

	path, err := o.opts.buildPath()
	assert.Nil(err)
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceID", path)
}

func TestAddChannelsToPushValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &addChannelsToPushOpts{
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       pn,
	}

	assert.Equal("webpubsub/validation: webpubsub: Remove Channel Group: Missing Subscribe Key", opts.validate().Error())
}
