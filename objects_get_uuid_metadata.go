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

var emptyWPSGetUUIDMetadataResponse *WPSGetUUIDMetadataResponse

const getUUIDMetadataPath = "/v2/objects/%s/uuids/%s"

type getUUIDMetadataBuilder struct {
	opts *getUUIDMetadataOpts
}

func newGetUUIDMetadataBuilder(webpubsub *WebPubSub) *getUUIDMetadataBuilder {
	builder := getUUIDMetadataBuilder{
		opts: &getUUIDMetadataOpts{
			webpubsub: webpubsub,
		},
	}
	return &builder
}

func newGetUUIDMetadataBuilderWithContext(webpubsub *WebPubSub,
	context Context) *getUUIDMetadataBuilder {
	builder := getUUIDMetadataBuilder{
		opts: &getUUIDMetadataOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

func (b *getUUIDMetadataBuilder) Include(include []WPSUUIDMetadataInclude) *getUUIDMetadataBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *getUUIDMetadataBuilder) UUID(uuid string) *getUUIDMetadataBuilder {
	b.opts.UUID = uuid

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getUUIDMetadataBuilder) QueryParam(queryParam map[string]string) *getUUIDMetadataBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getUUIDMetadata request.
func (b *getUUIDMetadataBuilder) Transport(tr http.RoundTripper) *getUUIDMetadataBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getUUIDMetadata request.
func (b *getUUIDMetadataBuilder) Execute() (*WPSGetUUIDMetadataResponse, StatusResponse, error) {
	if len(b.opts.UUID) <= 0 {
		b.opts.UUID = b.opts.webpubsub.Config.UUID
	}

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyWPSGetUUIDMetadataResponse, status, err
	}

	return newWPSGetUUIDMetadataResponse(rawJSON, b.opts, status)
}

type getUUIDMetadataOpts struct {
	webpubsub  *WebPubSub
	UUID       string
	Include    []string
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *getUUIDMetadataOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *getUUIDMetadataOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *getUUIDMetadataOpts) context() Context {
	return o.ctx
}

func (o *getUUIDMetadataOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *getUUIDMetadataOpts) buildPath() (string, error) {
	return fmt.Sprintf(getUUIDMetadataPath,
		o.webpubsub.Config.SubscribeKey, o.UUID), nil
}

func (o *getUUIDMetadataOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)

	if o.Include != nil {
		SetQueryParamAsCommaSepString(q, o.Include, "include")
	}

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getUUIDMetadataOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *getUUIDMetadataOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getUUIDMetadataOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *getUUIDMetadataOpts) httpMethod() string {
	return "GET"
}

func (o *getUUIDMetadataOpts) isAuthRequired() bool {
	return true
}

func (o *getUUIDMetadataOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *getUUIDMetadataOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *getUUIDMetadataOpts) operationType() OperationType {
	return WPSGetUUIDMetadataOperation
}

func (o *getUUIDMetadataOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *getUUIDMetadataOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// WPSGetUUIDMetadataResponse is the Objects API Response for Get User
type WPSGetUUIDMetadataResponse struct {
	status int     `json:"status"`
	Data   WPSUUID `json:"data"`
}

func newWPSGetUUIDMetadataResponse(jsonBytes []byte, o *getUUIDMetadataOpts,
	status StatusResponse) (*WPSGetUUIDMetadataResponse, StatusResponse, error) {

	resp := &WPSGetUUIDMetadataResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyWPSGetUUIDMetadataResponse, status, e
	}

	return resp, status, nil
}
