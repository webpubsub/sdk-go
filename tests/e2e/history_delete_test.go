package e2e

import (
	"testing"

	"github.com/stretchr/testify/assert"
	webpubsub "github.com/webpubsub/sdk-go/v7"
)

func TestHistoryDeleteNotStubbed(t *testing.T) {
	assert := assert.New(t)

	ch := randomized("h-ch")
	pn := webpubsub.NewWebPubSub(pamConfigCopy())

	_, _, err := pn.DeleteMessages().
		Channel(ch).
		Execute()

	assert.Nil(err)
}

func TestHistoryDeleteNotStubbedContext(t *testing.T) {
	assert := assert.New(t)

	ch := randomized("h-ch")
	pn := webpubsub.NewWebPubSub(pamConfigCopy())

	_, _, err := pn.DeleteMessagesWithContext(backgroundContext).
		Channel(ch).
		Execute()

	assert.Nil(err)
}

func TestHistoryDeleteMissingChannelError(t *testing.T) {
	assert := assert.New(t)

	config2 := pamConfigCopy()

	pn := webpubsub.NewWebPubSub(config2)

	res, _, err := pn.DeleteMessages().
		Channel("").
		Execute()

	assert.Nil(res)
	assert.Contains(err.Error(), "Missing Channel")
}

func TestHistoryDeleteSuperCall(t *testing.T) {
	assert := assert.New(t)

	config := pamConfigCopy()

	// Not allowed characters: /?#,
	validCharacters := "-._~:[]@!$&'()*+;=`|"

	config.UUID = validCharacters

	pn := webpubsub.NewWebPubSub(config)

	_, _, err := pn.DeleteMessages().
		Channel(validCharacters).
		Start(int64(123)).
		End(int64(456)).
		Execute()

	assert.Nil(err)
}
