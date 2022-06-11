package webpubsub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPushString(t *testing.T) {
	assert := assert.New(t)

	pushAPNS := WPSPushTypeAPNS
	pushAPNS2 := WPSPushTypeAPNS2
	pushMPNS := WPSPushTypeMPNS
	pushGCM := WPSPushTypeGCM
	pushNONE := WPSPushTypeNone

	assert.Equal("apns", pushAPNS.String())
	assert.Equal("apns2", pushAPNS2.String())
	assert.Equal("mpns", pushMPNS.String())
	assert.Equal("gcm", pushGCM.String())
	assert.Equal("none", pushNONE.String())
}

func TestStatusCategoryString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("Unknown", WPSUnknownCategory.String())
	assert.Equal("Timeout", WPSTimeoutCategory.String())
	assert.Equal("Connected", WPSConnectedCategory.String())
	assert.Equal("Disconnected", WPSDisconnectedCategory.String())
	assert.Equal("Cancelled", WPSCancelledCategory.String())
	assert.Equal("Loop Stop", WPSLoopStopCategory.String())
	assert.Equal("Acknowledgment", WPSAcknowledgmentCategory.String())
	assert.Equal("Bad Request", WPSBadRequestCategory.String())
	assert.Equal("Access Denied", WPSAccessDeniedCategory.String())
	assert.Equal("Reconnected", WPSReconnectedCategory.String())
	assert.Equal("Reconnection Attempts Exhausted", WPSReconnectionAttemptsExhausted.String())
	assert.Equal("No Stub Matched", WPSNoStubMatchedCategory.String())
}

func TestOperationTypeString(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("Subscribe", WPSSubscribeOperation.String())
	assert.Equal("Unsubscribe", WPSUnsubscribeOperation.String())
	assert.Equal("Publish", WPSPublishOperation.String())
	assert.Equal("Fire", WPSFireOperation.String())
	assert.Equal("History", WPSHistoryOperation.String())
	assert.Equal("Fetch Messages", WPSFetchMessagesOperation.String())
	assert.Equal("Where Now", WPSWhereNowOperation.String())
	assert.Equal("Here Now", WPSHereNowOperation.String())
	assert.Equal("Heartbeat", WPSHeartBeatOperation.String())
	assert.Equal("Set State", WPSSetStateOperation.String())
	assert.Equal("Get State", WPSGetStateOperation.String())
	assert.Equal("Add Channel To Channel Group", WPSAddChannelsToChannelGroupOperation.String())
	assert.Equal("Remove Channel From Channel Group", WPSRemoveChannelFromChannelGroupOperation.String())
	assert.Equal("Remove Channel Group", WPSRemoveGroupOperation.String())
	assert.Equal("List Channels In Channel Group", WPSChannelsForGroupOperation.String())
	assert.Equal("List Push Enabled Channels", WPSPushNotificationsEnabledChannelsOperation.String())
	assert.Equal("Add Push From Channel", WPSAddPushNotificationsOnChannelsOperation.String())
	assert.Equal("Remove Push From Channel", WPSRemovePushNotificationsFromChannelsOperation.String())
	assert.Equal("Remove All Push Notifications", WPSRemoveAllPushNotificationsOperation.String())
	assert.Equal("Time", WPSTimeOperation.String())
	assert.Equal("Grant", WPSAccessManagerGrant.String())
	assert.Equal("Revoke", WPSAccessManagerRevoke.String())
	assert.Equal("Delete messages", WPSDeleteMessagesOperation.String())
}
