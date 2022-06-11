package webpubsub

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	h "github.com/webpubsub/go/v7/tests/helpers"
)

func TestHistoryDeleteRequestAllParams(t *testing.T) {
	assert := assert.New(t)

	opts := &historyDeleteOpts{
		Channel:   "ch",
		SetStart:  true,
		SetEnd:    true,
		Start:     int64(123),
		End:       int64(456),
		webpubsub: webpubsub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v3/history/sub-key/sub_key/channel/%s", opts.Channel),
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("start", "123")
	expected.Set("end", "456")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestHistoryDeleteRequestQueryParams(t *testing.T) {
	assert := assert.New(t)

	opts := &historyDeleteOpts{
		Channel:   "ch",
		SetStart:  true,
		SetEnd:    true,
		Start:     int64(123),
		End:       int64(456),
		webpubsub: webpubsub,
	}

	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts.QueryParam = queryParam

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("start", "123")
	expected.Set("end", "456")
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestNewHistoryDeleteBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newHistoryDeleteBuilder(webpubsub)
	o.Channel("ch")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v3/history/sub-key/sub_key/channel/%s", o.opts.Channel),
		u.EscapedPath(), []int{})

	_, err1 := o.opts.buildQuery()
	assert.Nil(err1)

}

func TestNewHistoryDeleteBuilderContext(t *testing.T) {
	assert := assert.New(t)

	o := newHistoryDeleteBuilderWithContext(webpubsub, backgroundContext)
	o.Channel("ch")

	path, err := o.opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		fmt.Sprintf("/v3/history/sub-key/sub_key/channel/%s", o.opts.Channel),
		u.EscapedPath(), []int{})

	_, err1 := o.opts.buildQuery()
	assert.Nil(err1)

}

func TestHistoryDeleteOptsValidateSub(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &historyDeleteOpts{
		Channel:   "ch",
		SetStart:  true,
		SetEnd:    true,
		Start:     int64(123),
		End:       int64(456),
		webpubsub: pn,
	}

	assert.Equal("webpubsub/validation: webpubsub: Delete messages: Missing Subscribe Key", opts.validate().Error())
}

func TestHistoryDeleteOptsValidateSec(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	pn.Config.SecretKey = ""
	opts := &historyDeleteOpts{
		Channel:   "ch",
		SetStart:  true,
		SetEnd:    true,
		Start:     int64(123),
		End:       int64(456),
		webpubsub: pn,
	}

	assert.Equal("webpubsub/validation: webpubsub: Delete messages: Missing Secret Key", opts.validate().Error())
}
