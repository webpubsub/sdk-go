package webpubsub

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	h "github.com/webpubsub/go/v7/tests/helpers"
)

func TestLeaveRequestSingleChannel(t *testing.T) {
	assert := assert.New(t)

	opts := &leaveOpts{
		Channels:  []string{"ch"},
		webpubsub: webpubsub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/presence/sub-key/sub_key/channel/%s/leave", opts.Channels[0]),
		u.EscapedPath(), []int{})
}

func TestLeaveRequestMultipleChannels(t *testing.T) {
	assert := assert.New(t)

	opts := &leaveOpts{
		Channels:  []string{"ch1", "ch2", "ch3"},
		webpubsub: webpubsub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/channel/ch1,ch2,ch3/leave",
		u.EscapedPath(), []int{})
}

func TestLeaveRequestSingleChannelGroup(t *testing.T) {
	assert := assert.New(t)

	opts := &leaveOpts{
		ChannelGroups: []string{"cg"},
		webpubsub:     webpubsub,
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})
}

func TestLeaveRequestSingleChannelGroupQueryParam(t *testing.T) {
	assert := assert.New(t)

	opts := &leaveOpts{
		ChannelGroups: []string{"cg"},
		webpubsub:     webpubsub,
	}
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts.QueryParam = queryParam

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")

	expected.Set("channel-group", "cg")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})
}

func TestLeaveRequestMultipleChannelGroups(t *testing.T) {
	assert := assert.New(t)

	opts := &leaveOpts{
		ChannelGroups: []string{"cg1", "cg2", "cg3"},
		webpubsub:     webpubsub,
	}

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})
}

func TestLeaveRequestChannelsAndGroups(t *testing.T) {
	assert := assert.New(t)

	opts := &leaveOpts{
		Channels:      []string{"ch1", "ch2", "ch3"},
		ChannelGroups: []string{"cg1", "cg2", "cg3"},
		webpubsub:     webpubsub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/channel/ch1,ch2,ch3/leave",
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})
}

func TestLeaveRequestBuildQuery(t *testing.T) {
	assert := assert.New(t)
	opts := &leaveOpts{
		Channels:      []string{"ch1", "ch2", "ch3"},
		ChannelGroups: []string{"cg1", "cg2", "cg3"},
		webpubsub:     webpubsub,
	}
	query, err := opts.buildQuery()
	assert.NotNil(query)
	assert.Nil(err)

}

func TestLeaveRequestBuildPath(t *testing.T) {
	assert := assert.New(t)
	opts := &leaveOpts{
		webpubsub: webpubsub,
	}
	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/channel/,/leave",
		u.EscapedPath(), []int{})

}

func TestNewLeaveBuilder(t *testing.T) {
	assert := assert.New(t)
	o := newLeaveBuilder(webpubsub)

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/channel/,/leave",
		u.EscapedPath(), []int{})

}

func TestNewLeaveBuilderContext(t *testing.T) {
	assert := assert.New(t)
	o := newLeaveBuilderWithContext(webpubsub, backgroundContext)

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/channel/,/leave",
		u.EscapedPath(), []int{})

}

func TestLeaveOptsValidateSub(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &leaveOpts{
		webpubsub: pn,
	}

	assert.Equal("webpubsub/validation: webpubsub: Unsubscribe: Missing Subscribe Key", opts.validate().Error())
}

func TestLeaveOptsValidateCH(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	opts := &leaveOpts{
		webpubsub: pn,
	}

	assert.Equal("webpubsub/validation: webpubsub: Unsubscribe: Missing Channel or Channel Group", opts.validate().Error())
}
