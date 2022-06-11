package webpubsub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRemoveAllPushChannelsForDeviceBuilder(t *testing.T) {
	assert := assert.New(t)
	o := newRemoveAllPushChannelsForDeviceBuilder(webpubsub)
	o.DeviceIDForPush("deviceId")
	o.PushType(WPSPushTypeAPNS)

	str, err := o.opts.buildPath()
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceId/remove", str)
	assert.Nil(err)

}

func TestNewRemoveAllPushChannelsForDeviceBuilderContext(t *testing.T) {
	assert := assert.New(t)
	o := newRemoveAllPushChannelsForDeviceBuilderWithContext(webpubsub, backgroundContext)
	o.DeviceIDForPush("deviceId")
	o.PushType(WPSPushTypeAPNS)

	str, err := o.opts.buildPath()
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceId/remove", str)
	assert.Nil(err)

}

func TestRemoveAllPushNotificationsValidate(t *testing.T) {
	assert := assert.New(t)

	opts := &removeAllPushChannelsForDeviceOpts{
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
	}

	err := opts.validate()
	assert.Nil(err)

	opts1 := &removeAllPushChannelsForDeviceOpts{
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeNone,
		webpubsub:       webpubsub,
	}

	err1 := opts1.validate()
	assert.Contains(err1.Error(), "Missing Push Type")

	opts3 := &removeAllPushChannelsForDeviceOpts{
		PushType:  WPSPushTypeAPNS,
		webpubsub: webpubsub,
	}

	err3 := opts3.validate()

	assert.Contains(err3.Error(), "Missing Device ID")

}

func TestRemoveAllPushNotificationsBuildPath(t *testing.T) {
	assert := assert.New(t)

	opts := &removeAllPushChannelsForDeviceOpts{
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
	}

	str, err := opts.buildPath()
	assert.Equal("/v1/push/sub-key/sub_key/devices/deviceId/remove", str)
	assert.Nil(err)

}

func TestRemoveAllPushNotificationsBuildQueryParam(t *testing.T) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := &removeAllPushChannelsForDeviceOpts{
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
		QueryParam:      queryParam,
	}

	u, err := opts.buildQuery()
	assert.Equal("apns", u.Get("type"))
	assert.Equal("v1", u.Get("q1"))
	assert.Equal("v2", u.Get("q2"))

	assert.Nil(err)
}

func TestRemoveAllPushNotificationsBuildQuery(t *testing.T) {
	assert := assert.New(t)

	opts := &removeAllPushChannelsForDeviceOpts{
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
	}

	u, err := opts.buildQuery()
	assert.Equal("apns", u.Get("type"))

	assert.Nil(err)
}

func TestRemoveAllPushNotificationsBuildBody(t *testing.T) {
	assert := assert.New(t)

	opts := &removeAllPushChannelsForDeviceOpts{
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       webpubsub,
	}

	_, err := opts.buildBody()
	assert.Nil(err)

}

func TestRemoveAllPushChannelsForDeviceOptsValidateSub(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &removeAllPushChannelsForDeviceOpts{
		DeviceIDForPush: "deviceId",
		PushType:        WPSPushTypeAPNS,
		webpubsub:       pn,
	}

	assert.Equal("webpubsub/validation: webpubsub: Remove Channel Group: Missing Subscribe Key", opts.validate().Error())
}
