package webpubsub

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	h "github.com/webpubsub/go/v7/tests/helpers"
	"github.com/webpubsub/go/v7/utils"
)

func AssertSetChannelMetadata(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	incl := []WPSChannelMetadataInclude{
		WPSChannelMetadataIncludeCustom,
	}
	custom := make(map[string]interface{})
	custom["a"] = "b"
	custom["c"] = "d"

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	inclStr := EnumArrayToStringArray(incl)

	o := newSetChannelMetadataBuilder(pn)
	if testContext {
		o = newSetChannelMetadataBuilderWithContext(pn, backgroundContext)
	}

	o.Include(incl)
	o.Channel("id0")
	o.Name("name")
	o.Description("exturl")
	o.Custom(custom)
	o.QueryParam(queryParam)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/channels/%s", pn.Config.SubscribeKey, "id0"),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	expectedBody := "{\"name\":\"name\",\"description\":\"exturl\",\"custom\":{\"a\":\"b\",\"c\":\"d\"}}"

	assert.Equal(expectedBody, string(body))

	if checkQueryParam {
		u, _ := o.opts.buildQuery()
		assert.Equal("v1", u.Get("q1"))
		assert.Equal("v2", u.Get("q2"))
		assert.Equal(string(utils.JoinChannels(inclStr)), u.Get("include"))
	}

}

func TestExcludeInChannelMetadataBodyNotSetFields(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())

	o := newSetChannelMetadataBuilder(pn)

	o.Channel("id0")
	o.Name("name")

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf("/v2/objects/%s/channels/%s", pn.Config.SubscribeKey, "id0"),
		path, []int{})

	body, err := o.opts.buildBody()
	assert.Nil(err)

	expectedBody := "{\"name\":\"name\"}"

	assert.Equal(expectedBody, string(body))
}

func TestSetChannelMetadata(t *testing.T) {
	AssertSetChannelMetadata(t, true, false)
}

func TestSetChannelMetadataContext(t *testing.T) {
	AssertSetChannelMetadata(t, true, true)
}

func TestSetChannelMetadataResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	opts := &setChannelMetadataOpts{
		webpubsub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newWPSSetChannelMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("webpubsub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestSetChannelMetadataResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	opts := &setChannelMetadataOpts{
		webpubsub: pn,
	}
	jsonBytes := []byte(`{"status":200,"data":{"id":"id0","name":"name","description":"desc","custom":{"a":"b","c":"d"},"created":"2019-08-20T13:26:08.341297Z","updated":"2019-08-20T14:48:11.675743Z","eTag":"AYKH2s7ZlYKoJA"}}`)

	r, _, err := newWPSSetChannelMetadataResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("id0", r.Data.ID)
	assert.Equal("name", r.Data.Name)
	assert.Equal("desc", r.Data.Description)
	// assert.Equal("2019-08-20T13:26:08.341297Z", r.Data.Created)
	assert.Equal("2019-08-20T14:48:11.675743Z", r.Data.Updated)
	assert.Equal("AYKH2s7ZlYKoJA", r.Data.ETag)
	assert.Equal("b", r.Data.Custom["a"])

	assert.Nil(err)
}
