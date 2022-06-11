package webpubsub

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	h "github.com/webpubsub/go/v7/tests/helpers"
)

func AssertRemoveChannelMetadata(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newRemoveChannelMetadataBuilder(pn)
	if testContext {
		o = newRemoveChannelMetadataBuilderWithContext(pn, backgroundContext)
	}

	o.Channel("id0")
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/channels/%s", pn.Config.SubscribeKey, "id0"),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)
	assert.Empty(body)

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
	}

}

func TestRemoveChannelMetadata(t *testing.T) {
	AssertRemoveChannelMetadata(t, true, false)
}

func TestRemoveChannelMetadataContext(t *testing.T) {
	AssertRemoveChannelMetadata(t, true, true)
}

func TestRemoveChannelMetadataResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	opts := &removeChannelMetadataOpts{
		webpubsub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newWPSRemoveChannelMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("webpubsub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestRemoveChannelMetadataResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	opts := &removeChannelMetadataOpts{
		webpubsub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":null}`)

	r, _, err := newWPSRemoveChannelMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal(nil, r.Data)

	assert.Nil(err)
}
