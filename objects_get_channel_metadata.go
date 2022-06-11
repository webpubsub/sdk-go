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

	"github.com/webpubsub/sdk-go/v7/pnerr"
)

var emptyWPSGetChannelMetadataResponse *WPSGetChannelMetadataResponse

const getChannelMetadataPath = "/v2/objects/%s/channels/%s"

type getChannelMetadataBuilder struct {
	opts *getChannelMetadataOpts
}

func newGetChannelMetadataBuilder(webpubsub *WebPubSub) *getChannelMetadataBuilder {
	builder := getChannelMetadataBuilder{
		opts: &getChannelMetadataOpts{
			webpubsub: webpubsub,
		},
	}
	return &builder
}

func newGetChannelMetadataBuilderWithContext(webpubsub *WebPubSub,
	context Context) *getChannelMetadataBuilder {
	builder := getChannelMetadataBuilder{
		opts: &getChannelMetadataOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

func (b *getChannelMetadataBuilder) Include(include []WPSChannelMetadataInclude) *getChannelMetadataBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *getChannelMetadataBuilder) Channel(channel string) *getChannelMetadataBuilder {
	b.opts.Channel = channel

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getChannelMetadataBuilder) QueryParam(queryParam map[string]string) *getChannelMetadataBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getChannelMetadata request.
func (b *getChannelMetadataBuilder) Transport(tr http.RoundTripper) *getChannelMetadataBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getChannelMetadata request.
func (b *getChannelMetadataBuilder) Execute() (*WPSGetChannelMetadataResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyWPSGetChannelMetadataResponse, status, err
	}

	return newWPSGetChannelMetadataResponse(rawJSON, b.opts, status)
}

type getChannelMetadataOpts struct {
	webpubsub  *WebPubSub
	Channel    string
	Include    []string
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *getChannelMetadataOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *getChannelMetadataOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *getChannelMetadataOpts) context() Context {
	return o.ctx
}

func (o *getChannelMetadataOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.Channel == "" {
		return newValidationError(o, StrMissingChannel)
	}

	return nil
}

func (o *getChannelMetadataOpts) buildPath() (string, error) {
	return fmt.Sprintf(getChannelMetadataPath,
		o.webpubsub.Config.SubscribeKey, o.Channel), nil
}

func (o *getChannelMetadataOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)

	if o.Include != nil {
		SetQueryParamAsCommaSepString(q, o.Include, "include")
	}
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getChannelMetadataOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *getChannelMetadataOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getChannelMetadataOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *getChannelMetadataOpts) httpMethod() string {
	return "GET"
}

func (o *getChannelMetadataOpts) isAuthRequired() bool {
	return true
}

func (o *getChannelMetadataOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *getChannelMetadataOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *getChannelMetadataOpts) operationType() OperationType {
	return WPSGetChannelMetadataOperation
}

func (o *getChannelMetadataOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *getChannelMetadataOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// WPSGetChannelMetadataResponse is the Objects API Response for Get Space
type WPSGetChannelMetadataResponse struct {
	status int        `json:"status"`
	Data   WPSChannel `json:"data"`
}

func newWPSGetChannelMetadataResponse(jsonBytes []byte, o *getChannelMetadataOpts,
	status StatusResponse) (*WPSGetChannelMetadataResponse, StatusResponse, error) {

	resp := &WPSGetChannelMetadataResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyWPSGetChannelMetadataResponse, status, e
	}

	return resp, status, nil
}
