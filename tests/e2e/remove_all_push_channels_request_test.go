package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	webpubsub "github.com/webpubsub/sdk-go/v7"
)

func TestRemoveAllPushNotifications(t *testing.T) {
	assert := assert.New(t)

	pn := webpubsub.NewWebPubSub(configCopy())

	_, _, err := pn.RemoveAllPushNotifications().
		DeviceIDForPush("cg").
		PushType(webpubsub.WPSPushTypeGCM).
		Execute()
	assert.Nil(err)
}

func TestRemoveAllPushNotificationsContext(t *testing.T) {
	assert := assert.New(t)

	pn := webpubsub.NewWebPubSub(configCopy())

	_, _, err := pn.RemoveAllPushNotificationsWithContext(backgroundContext).
		DeviceIDForPush("cg").
		PushType(webpubsub.WPSPushTypeGCM).
		Execute()
	assert.Nil(err)
}
