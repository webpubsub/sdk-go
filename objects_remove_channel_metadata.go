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

var emptyWPSRemoveChannelMetadataResponse *WPSRemoveChannelMetadataResponse

const removeChannelMetadataPath = "/v2/objects/%s/channels/%s"

type removeChannelMetadataBuilder struct {
	opts *removeChannelMetadataOpts
}

func newRemoveChannelMetadataBuilder(webpubsub *WebPubSub) *removeChannelMetadataBuilder {
	builder := removeChannelMetadataBuilder{
		opts: &removeChannelMetadataOpts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newRemoveChannelMetadataBuilderWithContext(webpubsub *WebPubSub,
	context Context) *removeChannelMetadataBuilder {
	builder := removeChannelMetadataBuilder{
		opts: &removeChannelMetadataOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

func (b *removeChannelMetadataBuilder) Channel(channel string) *removeChannelMetadataBuilder {
	b.opts.Channel = channel

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *removeChannelMetadataBuilder) QueryParam(queryParam map[string]string) *removeChannelMetadataBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the removeChannelMetadata request.
func (b *removeChannelMetadataBuilder) Transport(tr http.RoundTripper) *removeChannelMetadataBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the removeChannelMetadata request.
func (b *removeChannelMetadataBuilder) Execute() (*WPSRemoveChannelMetadataResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyWPSRemoveChannelMetadataResponse, status, err
	}

	return newWPSRemoveChannelMetadataResponse(rawJSON, b.opts, status)
}

type removeChannelMetadataOpts struct {
	webpubsub  *WebPubSub
	Channel    string
	QueryParam map[string]string
	Transport  http.RoundTripper

	ctx Context
}

func (o *removeChannelMetadataOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *removeChannelMetadataOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *removeChannelMetadataOpts) context() Context {
	return o.ctx
}

func (o *removeChannelMetadataOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.Channel == "" {
		return newValidationError(o, StrMissingChannel)
	}

	return nil
}

func (o *removeChannelMetadataOpts) buildPath() (string, error) {
	return fmt.Sprintf(removeChannelMetadataPath,
		o.webpubsub.Config.SubscribeKey, o.Channel), nil
}

func (o *removeChannelMetadataOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *removeChannelMetadataOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *removeChannelMetadataOpts) buildBody() ([]byte, error) {
	return []byte{}, nil

}

func (o *removeChannelMetadataOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *removeChannelMetadataOpts) httpMethod() string {
	return "DELETE"
}

func (o *removeChannelMetadataOpts) isAuthRequired() bool {
	return true
}

func (o *removeChannelMetadataOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *removeChannelMetadataOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *removeChannelMetadataOpts) operationType() OperationType {
	return WPSRemoveChannelMetadataOperation
}

func (o *removeChannelMetadataOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *removeChannelMetadataOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// WPSRemoveChannelMetadataResponse is the Objects API Response for delete space
type WPSRemoveChannelMetadataResponse struct {
	status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func newWPSRemoveChannelMetadataResponse(jsonBytes []byte, o *removeChannelMetadataOpts,
	status StatusResponse) (*WPSRemoveChannelMetadataResponse, StatusResponse, error) {

	resp := &WPSRemoveChannelMetadataResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyWPSRemoveChannelMetadataResponse, status, e
	}

	return resp, status, nil
}
