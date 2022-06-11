package webpubsub

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	h "github.com/webpubsub/sdk-go/v7/tests/helpers"
)

func TestHereNowChannelsGroups(t *testing.T) {
	assert := assert.New(t)

	opts := &hereNowOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		ChannelGroups:   []string{"cg1", "cg2", "cg3"},
		webpubsub:       webpubsub,
		IncludeUUIDs:    true,
		SetIncludeUUIDs: true,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub_key/sub_key/channel/ch1,ch2,ch3",
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")
	expected.Set("disable-uuids", "0")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()

	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestHereNowNoChannel(t *testing.T) {
	assert := assert.New(t)

	opts := &hereNowOpts{
		ChannelGroups: []string{"cg1", "cg2", "cg3"},
		webpubsub:     webpubsub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Equal("/v2/presence/sub_key/sub_key/channel/,", path)
}

func TestNewHereNowBuilder(t *testing.T) {
	assert := assert.New(t)

	o := newHereNowBuilder(webpubsub)
	o.ChannelGroups([]string{"cg1", "cg2", "cg3"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	assert.Equal("/v2/presence/sub_key/sub_key/channel/,", path)
}

func TestNewHereNowBuilderContext(t *testing.T) {
	assert := assert.New(t)

	o := newHereNowBuilderWithContext(webpubsub, backgroundContext)
	o.ChannelGroups([]string{"cg1", "cg2", "cg3"})

	path, err := o.opts.buildPath()
	assert.Nil(err)
	assert.Equal("/v2/presence/sub_key/sub_key/channel/,", path)
}

func TestHereNowMultipleWithOpts(t *testing.T) {
	assert := assert.New(t)

	opts := &hereNowOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		ChannelGroups:   []string{"cg1", "cg2", "cg3"},
		IncludeUUIDs:    false,
		IncludeState:    true,
		SetIncludeState: true,
		SetIncludeUUIDs: true,
		webpubsub:       webpubsub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub_key/sub_key/channel/ch1,ch2,ch3",
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")
	expected.Set("state", "1")
	expected.Set("disable-uuids", "1")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestHereNowMultipleWithOptsQueryParam(t *testing.T) {
	assert := assert.New(t)

	opts := &hereNowOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		ChannelGroups:   []string{"cg1", "cg2", "cg3"},
		IncludeUUIDs:    false,
		IncludeState:    true,
		SetIncludeState: true,
		SetIncludeUUIDs: true,
		webpubsub:       webpubsub,
	}
	queryParam := map[string]string{
		"q1": "v1",
		"q2": "v2",
	}

	opts.QueryParam = queryParam

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")
	expected.Set("state", "1")
	expected.Set("disable-uuids", "1")
	expected.Set("q1", "v1")
	expected.Set("q2", "v2")

	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestHereNowGlobal(t *testing.T) {
	assert := assert.New(t)

	opts := &hereNowOpts{
		webpubsub: webpubsub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	u := &url.URL{
		Path: path,
	}
	h.AssertPathsEqual(t,
		"/v2/presence/sub_key/sub_key",
		u.EscapedPath(), []int{})

	query, err := opts.buildQuery()
	assert.Nil(err)

	expected := &url.Values{}
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

	body, err := opts.buildBody()
	assert.Nil(err)
	assert.Equal([]byte{}, body)
}

func TestHereNowValidateSubscribeKey(t *testing.T) {
	assert := assert.New(t)
	pn := NewWebPubSub(NewDemoConfig())
	pn.Config.SubscribeKey = ""
	opts := &hereNowOpts{
		webpubsub: pn,
	}

	assert.Equal("webpubsub/validation: webpubsub: Here Now: Missing Subscribe Key", opts.validate().Error())
}

func TestHereNowBuildPath(t *testing.T) {
	assert := assert.New(t)
	opts := &hereNowOpts{
		webpubsub: webpubsub,
	}

	path, err := opts.buildPath()
	assert.Nil(err)
	assert.Equal("/v2/presence/sub_key/sub_key", path)

}

func TestHereNowBuildQuery(t *testing.T) {
	assert := assert.New(t)
	opts := &hereNowOpts{
		Channels:        []string{"ch1", "ch2", "ch3"},
		ChannelGroups:   []string{"cg1", "cg2", "cg3"},
		IncludeUUIDs:    false,
		IncludeState:    true,
		SetIncludeState: true,
		SetIncludeUUIDs: false,
		webpubsub:       webpubsub,
	}
	query, err := opts.buildQuery()
	assert.Nil(err)
	expected := &url.Values{}
	expected.Set("channel-group", "cg1,cg2,cg3")
	expected.Set("state", "1")
	h.AssertQueriesEqual(t, expected, query, []string{"pnsdk", "uuid"}, []string{})

}

func TestNewHereNowResponseErrorUnmarshalling(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte(`s`)

	_, _, err := newHereNowResponse(jsonBytes, nil, StatusResponse{})
	assert.Equal("webpubsub/parsing: Error unmarshalling response: {s}", err.Error())
}

func TestNewHereNowResponseOneChannel(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte("{\"status\":200,\"message\":\"OK\",\"service\":\"Presence\",\"uuids\":[{\"uuid\":\"a3ffd012-a3b9-478c-8705-64089f24d71e\",\"state\":{\"age\":10}}],\"occupancy\":1}")

	_, _, err := newHereNowResponse(jsonBytes, []string{"a"}, StatusResponse{})
	assert.Nil(err)
}

func TestNewHereNowResponseOccupancyZero(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte("{\"status\":200,\"message\":\"OK\",\"service\":\"Presence\",\"occupancy\":0,\"total_channels\":1,\"total_occupancy\":1}")

	r, _, err := newHereNowResponse(jsonBytes, []string{"a"}, StatusResponse{})
	assert.Nil(err)
	assert.Equal(1, r.TotalChannels)
	assert.Equal(0, r.TotalOccupancy)

}

func TestNewHereNowResponseOccupancyZeroPayload(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte("{\"status\":200,\"message\":\"OK\",\"service\":\"Presence\",\"occupancy\":\"0\",\"total_channels\":1,\"total_occupancy\":1}")

	r, _, err := newHereNowResponse(jsonBytes, []string{"a"}, StatusResponse{})
	assert.Nil(err)
	assert.Equal(1, r.TotalChannels)
	assert.Equal(0, r.TotalOccupancy)
}

func TestNewHereNowResponseOccupancyZeroPayloadWithCh(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte("{\"status\":200,\"message\":\"OK\",\"payload\":{\"total_occupancy\":3,\"total_channels\":1,\"channels\":{\"ch1\":{\"occupancy\":1,\"uuids\":[{\"uuid\":\"user1\",\"state\":{\"age\":10}}]}}},\"service\":\"Presence\"}")

	r, _, err := newHereNowResponse(jsonBytes, []string{"a"}, StatusResponse{})
	assert.Nil(err)
	assert.Equal(1, r.TotalChannels)
	assert.Equal(3, r.TotalOccupancy)
}

func TestNewHereNowResponseOccupancyZeroPayloadWithoutCh(t *testing.T) {
	assert := assert.New(t)
	jsonBytes := []byte("{\"status\":200,\"message\":\"OK\",\"payload\":{\"total_occupancy\":3,\"total_channels\":2},\"service\":\"Presence\"}")

	r, _, err := newHereNowResponse(jsonBytes, []string{"a"}, StatusResponse{})
	assert.Nil(err)
	assert.Equal(1, r.TotalChannels)
	assert.Equal(0, r.TotalOccupancy)

}
