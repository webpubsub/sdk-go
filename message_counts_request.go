package webpubsub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"

	"reflect"
	"strconv"
	"strings"

	"github.com/webpubsub/sdk-go/v7/pnerr"
	"github.com/webpubsub/sdk-go/v7/utils"

	"net/http"
	"net/url"
)

var emptyMessageCountsResp *MessageCountsResponse

const messageCountsPath = "/v3/history/sub-key/%s/message-counts/%s"

type messageCountsBuilder struct {
	opts *messageCountsOpts
}

func newMessageCountsBuilder(webpubsub *WebPubSub) *messageCountsBuilder {
	builder := messageCountsBuilder{
		opts: &messageCountsOpts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newMessageCountsBuilderWithContext(webpubsub *WebPubSub,
	context Context) *messageCountsBuilder {
	builder := messageCountsBuilder{
		opts: &messageCountsOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

// Channels sets the Channels for the MessageCounts request.
func (b *messageCountsBuilder) Channels(channels []string) *messageCountsBuilder {
	b.opts.Channels = channels
	return b
}

// Deprecated: Use ChannelsTimetoken instead, pass one value in ChannelsTimetoken to achieve the same results.
// TODO: Remove in next major version bump
func (b *messageCountsBuilder) Timetoken(timetoken int64) *messageCountsBuilder {
	b.opts.Timetoken = timetoken
	return b
}

// ChannelsTimetoken Array of timetokens, in order of the channels list..
func (b *messageCountsBuilder) ChannelsTimetoken(channelsTimetoken []int64) *messageCountsBuilder {
	b.opts.ChannelsTimetoken = channelsTimetoken
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *messageCountsBuilder) QueryParam(queryParam map[string]string) *messageCountsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the MessageCounts request.
func (b *messageCountsBuilder) Transport(tr http.RoundTripper) *messageCountsBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the MessageCounts request.
func (b *messageCountsBuilder) Execute() (*MessageCountsResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyMessageCountsResp, status, err
	}

	return newMessageCountsResponse(rawJSON, b.opts, status)
}

type messageCountsOpts struct {
	webpubsub *WebPubSub

	Channels          []string
	Timetoken         int64
	ChannelsTimetoken []int64

	QueryParam map[string]string

	// nil hacks
	Transport http.RoundTripper

	ctx Context
}

func (o *messageCountsOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *messageCountsOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *messageCountsOpts) context() Context {
	return o.ctx
}

func (o *messageCountsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if len(o.Channels) <= 0 {
		return newValidationError(o, StrMissingChannel)
	}

	if (len(o.ChannelsTimetoken) <= 0) && (o.Timetoken == 0) {
		return newValidationError(o, StrChannelsTimetoken)
	}

	if (len(o.ChannelsTimetoken) > 1) && (len(o.Channels) != len(o.ChannelsTimetoken)) {
		return newValidationError(o, StrChannelsTimetokenLength)
	}

	return nil
}

func (o *messageCountsOpts) buildPath() (string, error) {
	channels := utils.JoinChannels(o.Channels)

	return fmt.Sprintf(messageCountsPath,
		o.webpubsub.Config.SubscribeKey,
		channels), nil
}

func (o *messageCountsOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)

	if (o.ChannelsTimetoken != nil) && (len(o.ChannelsTimetoken) == 1) {
		q.Set("timetoken", strconv.FormatInt(o.ChannelsTimetoken[0], 10))
		q.Set("channelsTimetoken", "")
	} else if o.ChannelsTimetoken != nil {
		q.Set("timetoken", "")
		q.Set("channelsTimetoken", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(o.ChannelsTimetoken)), ","), "[]"))
	} else {
		// TODO: Remove in next major version bump
		q.Set("timetoken", strconv.FormatInt(o.Timetoken, 10))
		q.Set("channelsTimetoken", "")
	}

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *messageCountsOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *messageCountsOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *messageCountsOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *messageCountsOpts) httpMethod() string {
	return "GET"
}

func (o *messageCountsOpts) isAuthRequired() bool {
	return true
}

func (o *messageCountsOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *messageCountsOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *messageCountsOpts) operationType() OperationType {
	return WPSMessageCountsOperation
}

func (o *messageCountsOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *messageCountsOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// MessageCountsResponse is the response to MessageCounts request. It contains a map of type MessageCountsResponseItem
type MessageCountsResponse struct {
	Channels map[string]int
}

//http://ps.pndsn.com/v3/history/sub-key/demo/message-counts/my-channel,my-channel1?timestamp=1549982652&pnsdk=WebPubSub-Go/4.1.6&uuid=pn-82f145ea-adc3-4917-a11d-76a957347a82&timetoken=15499825804610610&channelsTimetoken=15499825804610610,15499925804610615&auth=akey&signature=pVDVge_suepcOlSMllpsXg_jpOjtEpW7B3HHFaViI4s=
//{"status": 200, "error": false, "error_message": "", "channels": {"my-channel1":1,"my-channel":2}}
func newMessageCountsResponse(jsonBytes []byte, o *messageCountsOpts,
	status StatusResponse) (*MessageCountsResponse, StatusResponse, error) {

	resp := &MessageCountsResponse{}

	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyMessageCountsResp, status, e
	}

	if result, ok := value.(map[string]interface{}); ok {
		o.webpubsub.Config.Log.Println(result["channels"])
		if channels, ok1 := result["channels"].(map[string]interface{}); ok1 {
			if channels != nil {
				resp.Channels = make(map[string]int)
				for ch, v := range channels {
					resp.Channels[ch] = int(v.(float64))
				}
			} else {
				o.webpubsub.Config.Log.Printf("type assertion to map failed %v\n", result)
			}
		} else {
			o.webpubsub.Config.Log.Println("Assertion failed", reflect.TypeOf(result["channels"]))
		}
	} else {
		o.webpubsub.Config.Log.Printf("type assertion to map failed %v\n", value)
	}

	return resp, status, nil
}
