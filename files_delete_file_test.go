package webpubsub

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	h "github.com/webpubsub/sdk-go/v7/tests/helpers"
)

func AssertDeleteFile(t *testing.T, checkQueryParam, testContext bool) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	if !checkQueryParam {
		queryParam = nil
	}

	o := newDeleteFileBuilder(pn)
	if testContext {
		o = newDeleteFileBuilderWithContext(pn, backgroundContext)
	}

	channel := "chan"
	id := "fileid"
	name := "filename"
	o.Channel(channel)
	o.QueryParam(queryParam)
	o.ID(id)
	o.Name(name)

	path, err := o.opts.buildPath()
	assert.Nil(err)

	h.AssertPathsEqual(t,
		fmt.Sprintf(deleteFilePath, pn.Config.SubscribeKey, channel, id, name),
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

func TestDeleteFile(t *testing.T) {
	AssertDeleteFile(t, true, false)
}

func TestDeleteFileContext(t *testing.T) {
	AssertDeleteFile(t, true, true)
}

func TestDeleteFileResponseValueError(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	opts := &deleteFileOpts{
		webpubsub: pn,
	}
	jsonBytes := []byte(`s`)

	_, _, err := newWPSDeleteFileResponse(jsonBytes, opts, StatusResponse{})
	assert.Equal("webpubsub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestDeleteFileResponseValuePass(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	opts := &deleteFileOpts{
		webpubsub: pn,
	}
	jsonBytes := []byte(`{"status":200}`)

	_, s, err := newWPSDeleteFileResponse(jsonBytes, opts, StatusResponse{StatusCode: 200})
	assert.Equal(200, s.StatusCode)

	assert.Nil(err)
}
