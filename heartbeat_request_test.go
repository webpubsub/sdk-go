package webpubsub

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	h "github.com/webpubsub/sdk-go/v7/tests/helpers"
)

func TestHeartbeatRequestBasic(t *testing.T) {
	assert := assert.New(t)

	state := make(map[string]interface{})
	state["one"] = []string{"qwerty"}
	state["two"] = 2

	opts := &heartbeatOpts{
		webpubsub:     webpubsub,
		State:         state,
		Channels:      []string{"ch"},
		ChannelGroups: []string{"cg"},
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/%s/heartbeat",
			strings.Join(opts.Channels, ",")),
		u.EscapedPath(), []int{})

	u2, err := opts.buildQuery()
	assert.Nil(err)

	assert.Equal("cg", u2.Get("channel-group"))
	assert.Equal(`{"one":["qwerty"],"two":2}`, u2.Get("state"))
}

func TestNewHeartbeatBuilder(t *testing.T) {
	assert := assert.New(t)

	state := make(map[string]interface{})
	state["one"] = []string{"qwerty"}
	state["two"] = 2

	o := newHeartbeatBuilder(webpubsub)
	o.State(state)
	o.Channels([]string{"ch"})
	o.ChannelGroups([]string{"cg"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/%s/heartbeat",
			strings.Join(o.opts.Channels, ",")),
		u.EscapedPath(), []int{})

	u2, err := o.opts.buildQuery()
	assert.Nil(err)

	assert.Equal("cg", u2.Get("channel-group"))
	assert.Equal(`{"one":["qwerty"],"two":2}`, u2.Get("state"))
}

func TestNewHeartbeatBuilderContext(t *testing.T) {
	assert := assert.New(t)

	state := make(map[string]interface{})
	state["one"] = []string{"qwerty"}
	state["two"] = 2

	o := newHeartbeatBuilderWithContext(webpubsub, backgroundContext)
	o.State(state)
	o.Channels([]string{"ch"})
	o.ChannelGroups([]string{"cg"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/%s/heartbeat",
			strings.Join(o.opts.Channels, ",")),
		u.EscapedPath(), []int{})

	u2, err := o.opts.buildQuery()
	assert.Nil(err)

	assert.Equal("cg", u2.Get("channel-group"))
	assert.Equal(`{"one":["qwerty"],"two":2}`, u2.Get("state"))
}

func TestHeartbeatValidateChAndCg(t *testing.T) {
	assert := assert.New(t)

	opts := &heartbeatOpts{
		webpubsub: webpubsub,
	}
	err := opts.validate()
	assert.Equal("webpubsub/validation: webpubsub: Heartbeat: Missing Channel or Channel Group", err.Error())
}

func TestHeartbeatValidateSubKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &heartbeatOpts{
		webpubsub: pn,
	}
	err := opts.validate()
	assert.Equal("webpubsub/validation: webpubsub: Heartbeat: Missing Subscribe Key", err.Error())
}
