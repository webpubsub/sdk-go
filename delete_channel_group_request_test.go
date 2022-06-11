package webpubsub

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	h "github.com/webpubsub/go/v7/tests/helpers"
)

func TestDeleteChannelGroupRequestBasic(t *testing.T) {
	assert := assert.New(t)

	opts := &deleteChannelGroupOpts{
		ChannelGroup: "cg",
		webpubsub:    webpubsub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/channel-registration/sub-key/sub_key/channel-group/cg/remove"),
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)

	assert.Equal([]byte{}, body)
}

func TestDeleteChannelGroupRequestBasicQueryParam(t *testing.T) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts := &deleteChannelGroupOpts{
		ChannelGroup: "cg",
		webpubsub:    webpubsub,
		QueryParam:   queryParam,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/channel-registration/sub-key/sub_key/channel-group/cg/remove"),
		u.EscapedPath(), []int{})

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

func TestNewDeleteChannelGroupBuilder(t *testing.T) {
	assert := assert.New(t)
	o := newDeleteChannelGroupBuilder(webpubsub)
	o.ChannelGroup("cg")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/channel-registration/sub-key/sub_key/channel-group/cg/remove"),
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	assert.Equal([]byte{}, body)
}

func TestNewDeleteChannelGroupBuilderContext(t *testing.T) {
	assert := assert.New(t)
	o := newDeleteChannelGroupBuilderWithContext(webpubsub, backgroundContext)
	o.ChannelGroup("cg")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v1/channel-registration/sub-key/sub_key/channel-group/cg/remove"),
		u.EscapedPath(), []int{})

	query, err := o.opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	assert.Equal([]byte{}, body)
}
func TestDeleteChannelGroupOptsValidateSub(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &deleteChannelGroupOpts{
		ChannelGroup: "cg",
		webpubsub:    pn,
	}

	assert.Equal("webpubsub/validation: webpubsub: Remove Channel Group: Missing Subscribe Key", opts.validate().Error())
}
