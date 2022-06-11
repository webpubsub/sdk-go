package webpubsub

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	h "github.com/webpubsub/sdk-go/v7/tests/helpers"
)

func init() {
	pnconfig = NewConfig(GenerateUUID())

	pnconfig.PublishKey = "pub_key"
	pnconfig.SubscribeKey = "sub_key"
	pnconfig.SecretKey = "secret_key"

	webpubsub = NewWebPubSub(pnconfig)
}

func TestWhereNowBasicRequest(t *testing.T) {
	assert := assert.New(t)

	opts := &whereNowOpts{
		UUID:      "my-custom-uuid",
		webpubsub: webpubsub,
	}
	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/uuid/my-custom-uuid",
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestWhereNowBasicRequestQueryParam(t *testing.T) {
	assert := assert.New(t)
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}
	opts := &whereNowOpts{
		UUID:      "my-custom-uuid",
		webpubsub: webpubsub,
	}
	opts.QueryParam = queryParam
	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}

	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/uuid/my-custom-uuid",
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

func TestNewWhereNowBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newWhereNowBuilder(webpubsub)
	o.UUID("my-custom-uuid")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/uuid/my-custom-uuid",
		u.EscapedPath(), []int{})
}

func TestNewWhereNowBuilderContext(t *testing.T) {
	assert := assert.New(t)

	o := newWhereNowBuilderWithContext(webpubsub, backgroundContext)
	o.UUID("my-custom-uuid")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub-key/sub_key/uuid/my-custom-uuid",
		u.EscapedPath(), []int{})
}

func TestNewWhereNowResponserrorUnmarshalling(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`s`)

	_, _, err := newWhereNowResponse(jsonBytes, StatusResponse{})
	assert.Equal("webpubsub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestWhereNowValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &whereNowOpts{
		UUID:      "my-custom-uuid",
		webpubsub: pn,
	}

	assert.Equal("webpubsub/validation: webpubsub: Where Now: Missing Subscribe Key", opts.validate().Error())
}
