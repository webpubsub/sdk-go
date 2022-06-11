package e2e

import (
	//"fmt"

	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	webpubsub "github.com/webpubsub/sdk-go/v7"
)

func TestGrantToken(t *testing.T) {

	assert := assert.New(t)

	pn := webpubsub.NewWebPubSub(pamConfigCopy())
	if enableDebuggingInTests {
		pn.Config.Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
	ch1 := randomized("channel1")
	ch := map[string]webpubsub.ChannelPermissions{
		ch1: {
			Read:   true,
			Write:  true,
			Delete: false,
		},
	}

	cg1 := randomized("group1")
	cg2 := randomized("group2")
	cg := map[string]webpubsub.GroupPermissions{
		cg1: {
			Read:   true,
			Manage: true,
		},
		cg2: {
			Read:   true,
			Manage: false,
		},
	}

	res, _, err := pn.GrantToken().TTL(10).
		Channels(ch).
		ChannelGroups(cg).
		Execute()

	assert.Nil(err)

	assert.NotNil(res)
	if res != nil {
		token := res.Data.Token
		cborObject, err := webpubsub.GetPermissions(token)
		if err == nil {
			chResources := webpubsub.ParseGrantResources(cborObject.Resources, token, cborObject.Timestamp, cborObject.TTL)

			assert.Equal(ch[ch1], chResources.Channels[ch1].Permissions)
			assert.Equal(cg[cg1], chResources.Groups[cg1].Permissions)
			assert.Equal(cg[cg2], chResources.Groups[cg2].Permissions)
		}

	}

}
