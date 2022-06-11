package webpubsub

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	h "github.com/webpubsub/sdk-go/v7/tests/helpers"
)

func TestAddChannelRequestBasic(t *testing.T) {
	assert := assert.New(t)

	opts := &addChannelOpts{
		Channels:     []string{"ch1", "ch2", "ch3"},
		ChannelGroup: "cg",
		webpubsub:    webpubsub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/channel-registration/sub-key/sub_key/channel-group/cg"),
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("add", "ch1,ch2,ch3")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)

	assert.Equal([]byte{}, body)
}

func TestAddChannelRequestBasicQueryParam(t *testing.T) {
	assert := assert.New(t)

	opts := &addChannelOpts{
		Channels:     []string{"ch1", "ch2", "ch3"},
		ChannelGroup: "cg",
		webpubsub:    webpubsub,
	}
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}
	opts.QueryParam = queryParam

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/channel-registration/sub-key/sub_key/channel-group/cg"),
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")
	expected.Set("add", "ch1,ch2,ch3")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)

	assert.Equal([]byte{}, body)
}

func TestNewAddChannelToChannelGroupBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newAddChannelToChannelGroupBuilder(webpubsub)
	o.ChannelGroup("cg")
	o.Channels([]string{"ch1", "ch2", "ch3"})
	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/channel-registration/sub-key/sub_key/channel-group/cg"),
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("add", "ch1,ch2,ch3")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	assert.Equal([]byte{}, body)
}

func TestNewAddChannelToChannelGroupBuilderWithContext(t *testing.T) {
	assert := assert.New(t)

	o := newAddChannelToChannelGroupBuilderWithContext(webpubsub, backgroundContext)
	o.ChannelGroup("cg")
	o.Channels([]string{"ch1", "ch2", "ch3"})
	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/channel-registration/sub-key/sub_key/channel-group/cg"),
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("add", "ch1,ch2,ch3")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	assert.Equal([]byte{}, body)
}

func TestAddChannelOptsValidateSub(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &addChannelOpts{
		Channels:     []string{"ch1", "ch2", "ch3"},
		ChannelGroup: "cg",
		webpubsub:    pn,
	}
	assert.Equal("webpubsub/validation: webpubsub: Add Channel To Channel Group: Missing Subscribe Key", opts.validate().Error())
}
