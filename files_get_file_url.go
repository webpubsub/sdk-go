package webpubsub

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
)

var emptyGetFileURLResponse *WPSGetFileURLResponse

const getFileURLPath = "/v1/files/%s/channels/%s/files/%s/%s"

type getFileURLBuilder struct {
	opts *getFileURLOpts
}

func newGetFileURLBuilder(webpubsub *WebPubSub) *getFileURLBuilder {
	builder := getFileURLBuilder{
		opts: &getFileURLOpts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newGetFileURLBuilderWithContext(webpubsub *WebPubSub,
	context Context) *getFileURLBuilder {
	builder := getFileURLBuilder{
		opts: &getFileURLOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

func (b *getFileURLBuilder) Channel(channel string) *getFileURLBuilder {
	b.opts.Channel = channel

	return b
}

func (b *getFileURLBuilder) ID(id string) *getFileURLBuilder {
	b.opts.ID = id

	return b
}

func (b *getFileURLBuilder) Name(name string) *getFileURLBuilder {
	b.opts.Name = name

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getFileURLBuilder) QueryParam(queryParam map[string]string) *getFileURLBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getFileURL request.
func (b *getFileURLBuilder) Transport(tr http.RoundTripper) *getFileURLBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getFileURL request.
func (b *getFileURLBuilder) Execute() (*WPSGetFileURLResponse, StatusResponse, error) {
	u, _ := buildURL(b.opts)

	resp := &WPSGetFileURLResponse{
		URL: u.RequestURI(),
	}
	stat := StatusResponse{
		AffectedChannels: []string{b.opts.Channel},
		AuthKey:          b.opts.config().AuthKey,
		Category:         WPSUnknownCategory,
		Operation:        WPSGetFileURLOperation,
		StatusCode:       200,
		TLSEnabled:       b.opts.config().Secure,
		Origin:           b.opts.config().Origin,
		UUID:             b.opts.config().UUID,
	}
	return resp, stat, nil
}

type getFileURLOpts struct {
	webpubsub *WebPubSub

	Channel    string
	ID         string
	Name       string
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *getFileURLOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *getFileURLOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *getFileURLOpts) context() Context {
	return o.ctx
}

func (o *getFileURLOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.Channel == "" {
		return newValidationError(o, StrMissingChannel)
	}

	if o.Name == "" {
		return newValidationError(o, StrMissingFileName)
	}

	if o.ID == "" {
		return newValidationError(o, StrMissingFileID)
	}

	return nil
}

func (o *getFileURLOpts) buildPath() (string, error) {
	return fmt.Sprintf(getFileURLPath,
		o.webpubsub.Config.SubscribeKey, o.Channel, o.ID, o.Name), nil
}

func (o *getFileURLOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getFileURLOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *getFileURLOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getFileURLOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *getFileURLOpts) httpMethod() string {
	return "GET"
}

func (o *getFileURLOpts) isAuthRequired() bool {
	return true
}

func (o *getFileURLOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *getFileURLOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *getFileURLOpts) operationType() OperationType {
	return WPSGetFileURLOperation
}

func (o *getFileURLOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *getFileURLOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// WPSGetFileURLResponse is the File Upload API Response for Get Spaces
type WPSGetFileURLResponse struct {
	URL string `json:"location"`
}
