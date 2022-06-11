package webpubsub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

	"github.com/webpubsub/go/v7/pnerr"
)

var emptyWPSGetMessageActionsResponse *WPSGetMessageActionsResponse

const getMessageActionsPath = "/v1/message-actions/%s/channel/%s"

type getMessageActionsBuilder struct {
	opts *getMessageActionsOpts
}

func newGetMessageActionsBuilder(webpubsub *WebPubSub) *getMessageActionsBuilder {
	builder := getMessageActionsBuilder{
		opts: &getMessageActionsOpts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newGetMessageActionsBuilderWithContext(webpubsub *WebPubSub,
	context Context) *getMessageActionsBuilder {
	builder := getMessageActionsBuilder{
		opts: &getMessageActionsOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

func (b *getMessageActionsBuilder) Channel(channel string) *getMessageActionsBuilder {
	b.opts.Channel = channel

	return b
}

func (b *getMessageActionsBuilder) Start(timetoken string) *getMessageActionsBuilder {
	b.opts.Start = timetoken

	return b
}

func (b *getMessageActionsBuilder) End(timetoken string) *getMessageActionsBuilder {
	b.opts.End = timetoken

	return b
}

func (b *getMessageActionsBuilder) Limit(limit int) *getMessageActionsBuilder {
	b.opts.Limit = limit

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getMessageActionsBuilder) QueryParam(queryParam map[string]string) *getMessageActionsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getMessageActions request.
func (b *getMessageActionsBuilder) Transport(tr http.RoundTripper) *getMessageActionsBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getMessageActions request.
func (b *getMessageActionsBuilder) Execute() (*WPSGetMessageActionsResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyWPSGetMessageActionsResponse, status, err
	}

	return newWPSGetMessageActionsResponse(rawJSON, b.opts, status)
}

type getMessageActionsOpts struct {
	webpubsub *WebPubSub

	Channel    string
	Start      string
	End        string
	Limit      int
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *getMessageActionsOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *getMessageActionsOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *getMessageActionsOpts) context() Context {
	return o.ctx
}

func (o *getMessageActionsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *getMessageActionsOpts) buildPath() (string, error) {
	return fmt.Sprintf(getMessageActionsPath,
		o.webpubsub.Config.SubscribeKey, o.Channel), nil
}

func (o *getMessageActionsOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)

	if o.Start != "" {
		q.Set("start", o.Start)
	}

	if o.End != "" {
		q.Set("end", o.End)
	}

	if o.Limit > 0 {
		q.Set("limit", strconv.Itoa(o.Limit))
	}

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getMessageActionsOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *getMessageActionsOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getMessageActionsOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *getMessageActionsOpts) httpMethod() string {
	return "GET"
}

func (o *getMessageActionsOpts) isAuthRequired() bool {
	return true
}

func (o *getMessageActionsOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *getMessageActionsOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *getMessageActionsOpts) operationType() OperationType {
	return WPSGetMessageActionsOperation
}

func (o *getMessageActionsOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *getMessageActionsOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// WPSGetMessageActionsMore is the struct used when the WPSGetMessageActionsResponse has more link
type WPSGetMessageActionsMore struct {
	URL   string `json:"url"`
	Start string `json:"start"`
	End   string `json:"end"`
	Limit int    `json:"limit"`
}

// WPSGetMessageActionsResponse is the GetMessageActions API Response
type WPSGetMessageActionsResponse struct {
	status int                         `json:"status"`
	Data   []WPSMessageActionsResponse `json:"data"`
	More   WPSGetMessageActionsMore    `json:"more"`
}

func newWPSGetMessageActionsResponse(jsonBytes []byte, o *getMessageActionsOpts,
	status StatusResponse) (*WPSGetMessageActionsResponse, StatusResponse, error) {

	resp := &WPSGetMessageActionsResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyWPSGetMessageActionsResponse, status, e
	}

	return resp, status, nil
}
