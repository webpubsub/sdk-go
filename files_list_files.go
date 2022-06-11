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

var emptyListFilesResponse *WPSListFilesResponse

const listFilesPath = "/v1/files/%s/channels/%s/files"

const listFilesLimit = 100

type listFilesBuilder struct {
	opts *listFilesOpts
}

func newListFilesBuilder(webpubsub *WebPubSub) *listFilesBuilder {
	builder := listFilesBuilder{
		opts: &listFilesOpts{
			webpubsub: webpubsub,
		},
	}
	builder.opts.Limit = listFilesLimit

	return &builder
}

func newListFilesBuilderWithContext(webpubsub *WebPubSub,
	context Context) *listFilesBuilder {
	builder := listFilesBuilder{
		opts: &listFilesOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

func (b *listFilesBuilder) Limit(limit int) *listFilesBuilder {
	b.opts.Limit = limit

	return b
}

func (b *listFilesBuilder) Next(next string) *listFilesBuilder {
	b.opts.Next = next

	return b
}

func (b *listFilesBuilder) Channel(channel string) *listFilesBuilder {
	b.opts.Channel = channel

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *listFilesBuilder) QueryParam(queryParam map[string]string) *listFilesBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the listFiles request.
func (b *listFilesBuilder) Transport(tr http.RoundTripper) *listFilesBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the listFiles request.
func (b *listFilesBuilder) Execute() (*WPSListFilesResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyListFilesResponse, status, err
	}

	return newWPSListFilesResponse(rawJSON, b.opts, status)
}

type listFilesOpts struct {
	webpubsub *WebPubSub

	Limit      int
	Next       string
	Channel    string
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *listFilesOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *listFilesOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *listFilesOpts) context() Context {
	return o.ctx
}

func (o *listFilesOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *listFilesOpts) buildPath() (string, error) {
	return fmt.Sprintf(listFilesPath,
		o.webpubsub.Config.SubscribeKey, o.Channel), nil
}

func (o *listFilesOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)

	q.Set("limit", strconv.Itoa(o.Limit))

	if o.Next != "" {
		q.Set("next", o.Next)
	}

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *listFilesOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *listFilesOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *listFilesOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *listFilesOpts) httpMethod() string {
	return "GET"
}

func (o *listFilesOpts) isAuthRequired() bool {
	return true
}

func (o *listFilesOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *listFilesOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *listFilesOpts) operationType() OperationType {
	return WPSListFilesOperation
}

func (o *listFilesOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *listFilesOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// WPSListFilesResponse is the File Upload API Response for Get Spaces
type WPSListFilesResponse struct {
	status int           `json:"status"`
	Data   []WPSFileInfo `json:"data"`
	Count  int           `json:"count"`
	Next   string        `json:"next"`
}

func newWPSListFilesResponse(jsonBytes []byte, o *listFilesOpts,
	status StatusResponse) (*WPSListFilesResponse, StatusResponse, error) {

	resp := &WPSListFilesResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyListFilesResponse, status, e
	}

	return resp, status, nil
}
