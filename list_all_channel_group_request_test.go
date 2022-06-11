package webpubsub

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	h "github.com/webpubsub/go/v7/tests/helpers"
)

func TestListAllChannelGroupRequestBasic(t *testing.T) {
	assert := assert.New(t)

	opts := &allChannelGroupOpts{
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

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestListAllChannelGroupRequestBasicQueryParam(t *testing.T) {
	assert := assert.New(t)

	opts := &allChannelGroupOpts{
		ChannelGroup: "cg",
		webpubsub:    webpubsub,
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

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestNewAllChannelGroupBuilder(t *testing.T) {
	assert := assert.New(t)
	o := newAllChannelGroupBuilder(webpubsub)
	o.ChannelGroup("cg")

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

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestNewAllChannelGroupBuilderContext(t *testing.T) {
	assert := assert.New(t)
	o := newAllChannelGroupBuilderWithContext(webpubsub, backgroundContext)
	o.ChannelGroup("cg")

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

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestListAllChannelsNewAllChannelGroupResponseErrorUnmarshalling(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`s`)

	_, _, err := newAllChannelGroupResponse(jsonBytes, StatusResponse{})
	assert.Equal("webpubsub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestListAllChannelsValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &allChannelGroupOpts{
		ChannelGroup: "cg",
		webpubsub:    pn,
	}

	assert.Equal("webpubsub/validation: webpubsub: List Channels In Channel Group: Missing Subscribe Key", opts.validate().Error())
}

func TestListAllChannelsValidateChannelGrp(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	opts := &allChannelGroupOpts{
		webpubsub: pn,
	}

	assert.Equal("webpubsub/validation: webpubsub: List Channels In Channel Group: Missing Channel Group", opts.validate().Error())
}
