package webpubsub

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	h "github.com/webpubsub/go/v7/tests/helpers"
)

func TestSubscribeSingleChannel(t *testing.T) {
	assert := assert.New(t)
	opts := &subscribeOpts{
		Channels:  []string{"ch"},
		webpubsub: webpubsub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/subscribe/sub_key/ch/0", u.EscapedPath(), []int{})
}

func TestSubscribeMultipleChannels(t *testing.T) {
	assert := assert.New(t)
	opts := &subscribeOpts{
		Channels:  []string{"ch-1", "ch-2", "ch-3"},
		webpubsub: webpubsub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		"/v2/subscribe/sub_key/ch-1,ch-2,ch-3/0", u.EscapedPath(), []int{})
}

func TestSubscribeChannelGroups(t *testing.T) {
	assert := assert.New(t)
	opts := &subscribeOpts{
		ChannelGroups: []string{"cg-1", "cg-2", "cg-3"},
		webpubsub:     webpubsub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/subscribe/sub_key/,/0", u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg-1,cg-2,cg-3")

	h.AssertQueriesEqual(t, expected, query,
		[]string{"pnsdk", "uuid"}, []string{})
}

func TestSubscribeMixedParams(t *testing.T) {
	assert := assert.New(t)

	opts := &subscribeOpts{
		Channels:         []string{"ch"},
		ChannelGroups:    []string{"cg"},
		Region:           "us-east-1",
		Timetoken:        123,
		FilterExpression: "abc",
		webpubsub:        webpubsub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)

	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		"/v2/subscribe/sub_key/ch/0", u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg")
	expected.Set("tr", "us-east-1")
	expected.Set("filter-expr", "abc")
	expected.Set("tt", "123")

	h.AssertQueriesEqual(t, expected, query,
		[]string{"pnsdk", "uuid"}, []string{})
}

func TestSubscribeMixedQueryParams(t *testing.T) {
	assert := assert.New(t)

	opts := &subscribeOpts{
		Channels:         []string{"ch"},
		ChannelGroups:    []string{"cg"},
		Region:           "us-east-1",
		Timetoken:        123,
		FilterExpression: "abc",
		webpubsub:        webpubsub,
	}
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts.QueryParam = queryParam

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg")
	expected.Set("tr", "us-east-1")
	expected.Set("filter-expr", "abc")
	expected.Set("tt", "123")
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")

	h.AssertQueriesEqual(t, expected, query,
		[]string{"pnsdk", "uuid"}, []string{})
}

func TestSubscribeValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &subscribeOpts{
		Channels:         []string{"ch"},
		ChannelGroups:    []string{"cg"},
		Region:           "us-east-1",
		Timetoken:        123,
		FilterExpression: "abc",
		webpubsub:        pn,
	}

	assert.Equal("webpubsub/validation: webpubsub: Subscribe: Missing Subscribe Key", opts.validate().Error())
}

func TestSubscribeValidatePublishKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	pn.Config.PublishKey = ""
	opts := &subscribeOpts{
		Channels:         []string{"ch"},
		ChannelGroups:    []string{"cg"},
		Region:           "us-east-1",
		Timetoken:        123,
		FilterExpression: "abc",
		webpubsub:        pn,
	}

	assert.Nil(opts.validate())
}

func TestSubscribeValidateCHAndCG(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	opts := &subscribeOpts{
		Region:           "us-east-1",
		Timetoken:        123,
		FilterExpression: "abc",
		webpubsub:        pn,
	}

	assert.Equal("webpubsub/validation: webpubsub: Subscribe: Missing Channel", opts.validate().Error())
}

func TestSubscribeValidateState(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	opts := &subscribeOpts{
		Channels:         []string{"ch"},
		ChannelGroups:    []string{"cg"},
		Region:           "us-east-1",
		Timetoken:        123,
		FilterExpression: "abc",
		webpubsub:        pn,
	}
	opts.State = map[string]interface{}{"a": "a"}

	assert.Nil(opts.validate())
}
