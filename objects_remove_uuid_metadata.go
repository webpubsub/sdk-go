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

var emptyWPSRemoveUUIDMetadataResponse *WPSRemoveUUIDMetadataResponse

const removeUUIDMetadataPath = "/v2/objects/%s/uuids/%s"

type removeUUIDMetadataBuilder struct {
	opts *removeUUIDMetadataOpts
}

func newRemoveUUIDMetadataBuilder(webpubsub *WebPubSub) *removeUUIDMetadataBuilder {
	builder := removeUUIDMetadataBuilder{
		opts: &removeUUIDMetadataOpts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newRemoveUUIDMetadataBuilderWithContext(webpubsub *WebPubSub,
	context Context) *removeUUIDMetadataBuilder {
	builder := removeUUIDMetadataBuilder{
		opts: &removeUUIDMetadataOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

func (b *removeUUIDMetadataBuilder) UUID(uuid string) *removeUUIDMetadataBuilder {
	b.opts.UUID = uuid

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *removeUUIDMetadataBuilder) QueryParam(queryParam map[string]string) *removeUUIDMetadataBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the removeUUIDMetadata request.
func (b *removeUUIDMetadataBuilder) Transport(tr http.RoundTripper) *removeUUIDMetadataBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the removeUUIDMetadata request.
func (b *removeUUIDMetadataBuilder) Execute() (*WPSRemoveUUIDMetadataResponse, StatusResponse, error) {
	if len(b.opts.UUID) <= 0 {
		b.opts.UUID = b.opts.webpubsub.Config.UUID
	}

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyWPSRemoveUUIDMetadataResponse, status, err
	}

	return newWPSRemoveUUIDMetadataResponse(rawJSON, b.opts, status)
}

type removeUUIDMetadataOpts struct {
	webpubsub  *WebPubSub
	UUID       string
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *removeUUIDMetadataOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *removeUUIDMetadataOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *removeUUIDMetadataOpts) context() Context {
	return o.ctx
}

func (o *removeUUIDMetadataOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *removeUUIDMetadataOpts) buildPath() (string, error) {
	return fmt.Sprintf(removeUUIDMetadataPath,
		o.webpubsub.Config.SubscribeKey, o.UUID), nil
}

func (o *removeUUIDMetadataOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *removeUUIDMetadataOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *removeUUIDMetadataOpts) buildBody() ([]byte, error) {
	return []byte{}, nil

}

func (o *removeUUIDMetadataOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *removeUUIDMetadataOpts) httpMethod() string {
	return "DELETE"
}

func (o *removeUUIDMetadataOpts) isAuthRequired() bool {
	return true
}

func (o *removeUUIDMetadataOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *removeUUIDMetadataOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *removeUUIDMetadataOpts) operationType() OperationType {
	return WPSRemoveUUIDMetadataOperation
}

func (o *removeUUIDMetadataOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *removeUUIDMetadataOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// WPSRemoveUUIDMetadataResponse is the Objects API Response for delete user
type WPSRemoveUUIDMetadataResponse struct {
	status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func newWPSRemoveUUIDMetadataResponse(jsonBytes []byte, o *removeUUIDMetadataOpts,
	status StatusResponse) (*WPSRemoveUUIDMetadataResponse, StatusResponse, error) {

	resp := &WPSRemoveUUIDMetadataResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyWPSRemoveUUIDMetadataResponse, status, e
	}

	return resp, status, nil
}
