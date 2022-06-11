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

	"github.com/webpubsub/go/v7/pnerr"
)

var emptyWPSRemoveMessageActionsResponse *WPSRemoveMessageActionsResponse

const removeMessageActionsPath = "/v1/message-actions/%s/channel/%s/message/%s/action/%s"

type removeMessageActionsBuilder struct {
	opts *removeMessageActionsOpts
}

func newRemoveMessageActionsBuilder(webpubsub *WebPubSub) *removeMessageActionsBuilder {
	builder := removeMessageActionsBuilder{
		opts: &removeMessageActionsOpts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newRemoveMessageActionsBuilderWithContext(webpubsub *WebPubSub,
	context Context) *removeMessageActionsBuilder {
	builder := removeMessageActionsBuilder{
		opts: &removeMessageActionsOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

func (b *removeMessageActionsBuilder) Channel(channel string) *removeMessageActionsBuilder {
	b.opts.Channel = channel

	return b
}

func (b *removeMessageActionsBuilder) MessageTimetoken(timetoken string) *removeMessageActionsBuilder {
	b.opts.MessageTimetoken = timetoken

	return b
}

func (b *removeMessageActionsBuilder) ActionTimetoken(timetoken string) *removeMessageActionsBuilder {
	b.opts.ActionTimetoken = timetoken

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *removeMessageActionsBuilder) QueryParam(queryParam map[string]string) *removeMessageActionsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the removeMessageActions request.
func (b *removeMessageActionsBuilder) Transport(tr http.RoundTripper) *removeMessageActionsBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the removeMessageActions request.
func (b *removeMessageActionsBuilder) Execute() (*WPSRemoveMessageActionsResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyWPSRemoveMessageActionsResponse, status, err
	}

	return newWPSRemoveMessageActionsResponse(rawJSON, b.opts, status)
}

type removeMessageActionsOpts struct {
	webpubsub *WebPubSub

	Channel          string
	MessageTimetoken string
	ActionTimetoken  string
	Custom           map[string]interface{}
	QueryParam       map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *removeMessageActionsOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *removeMessageActionsOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *removeMessageActionsOpts) context() Context {
	return o.ctx
}

func (o *removeMessageActionsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *removeMessageActionsOpts) buildPath() (string, error) {
	return fmt.Sprintf(removeMessageActionsPath,
		o.webpubsub.Config.SubscribeKey, o.Channel, o.MessageTimetoken, o.ActionTimetoken), nil
}

func (o *removeMessageActionsOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *removeMessageActionsOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *removeMessageActionsOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *removeMessageActionsOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *removeMessageActionsOpts) httpMethod() string {
	return "DELETE"
}

func (o *removeMessageActionsOpts) isAuthRequired() bool {
	return true
}

func (o *removeMessageActionsOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *removeMessageActionsOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *removeMessageActionsOpts) operationType() OperationType {
	return WPSRemoveMessageActionsOperation
}

func (o *removeMessageActionsOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *removeMessageActionsOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// WPSRemoveMessageActionsResponse is the Objects API Response for create space
type WPSRemoveMessageActionsResponse struct {
	status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func newWPSRemoveMessageActionsResponse(jsonBytes []byte, o *removeMessageActionsOpts,
	status StatusResponse) (*WPSRemoveMessageActionsResponse, StatusResponse, error) {

	resp := &WPSRemoveMessageActionsResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyWPSRemoveMessageActionsResponse, status, e
	}

	return resp, status, nil
}
