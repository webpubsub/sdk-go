package webpubsub

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	h "github.com/webpubsub/go/v7/tests/helpers"
)

func TestTimeRequestHTTP2(t *testing.T) {
	assert := assert.New(t)

	config := NewConfig(GenerateUUID())
	config.Origin = "ssp.webpubsub.com"
	config.UseHTTP2 = true

	pn := NewWebPubSub(config)

	_, s, err := pn.Time().Execute()

	assert.Nil(err)
	assert.Equal(200, s.StatusCode)
}

func TestNewTimeResponseUnmarshalling(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`s`)

	_, _, err := newTimeResponse(jsonBytes, fakeResponseState)
	assert.Equal("webpubsub/parsing: Error unmarshalling response: {s}", err.Error())

	opts := &timeOpts{}
	a, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal(a, []byte{})
}

func TestNewTimeResponseQueryParam(t *testing.T) {
	assert := assert.New(t)

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}
	config := NewConfig(GenerateUUID())
	pn := NewWebPubSub(config)

	opts := &timeOpts{}
	opts.webpubsub = pn
	opts.QueryParam = queryParam

	expected := &url.Values{}
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")

	path, err := opts.buildPath()
	u := &url.URL{
		Path: path,
	}
	assert.Nil(err)

	query, err := opts.buildQuery()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		"/time/0",
		u.EscapedPath(), []int{})

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	a, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal(a, []byte{})
}

func TestNewTimeBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newTimeBuilder(webpubsub)
	_, err := o.opts.buildBody()
	assert.Nil(err)
}

func TestNewTimeBuilderContext(t *testing.T) {
	assert := assert.New(t)

	o := newTimeBuilderWithContext(webpubsub, backgroundContext)
	_, err := o.opts.buildBody()
	assert.Nil(err)
}
