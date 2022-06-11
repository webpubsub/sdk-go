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

	"github.com/webpubsub/sdk-go/v7/pnerr"
)

var emptyWPSGetAllUUIDMetadataResponse *WPSGetAllUUIDMetadataResponse

const getAllUUIDMetadataPath = "/v2/objects/%s/uuids"

const getAllUUIDMetadataLimitV2 = 100

type getAllUUIDMetadataBuilder struct {
	opts *getAllUUIDMetadataOpts
}

func newGetAllUUIDMetadataBuilder(webpubsub *WebPubSub) *getAllUUIDMetadataBuilder {
	return newGetAllUUIDMetadataBuilderWithContext(webpubsub, webpubsub.ctx)
}

func newGetAllUUIDMetadataBuilderWithContext(webpubsub *WebPubSub,
	context Context) *getAllUUIDMetadataBuilder {
	builder := getAllUUIDMetadataBuilder{
		opts: &getAllUUIDMetadataOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}
	builder.opts.Limit = getAllUUIDMetadataLimitV2

	return &builder
}

func (b *getAllUUIDMetadataBuilder) Include(include []WPSUUIDMetadataInclude) *getAllUUIDMetadataBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *getAllUUIDMetadataBuilder) Limit(limit int) *getAllUUIDMetadataBuilder {
	b.opts.Limit = limit

	return b
}

func (b *getAllUUIDMetadataBuilder) Start(start string) *getAllUUIDMetadataBuilder {
	b.opts.Start = start

	return b
}

func (b *getAllUUIDMetadataBuilder) End(end string) *getAllUUIDMetadataBuilder {
	b.opts.End = end

	return b
}

func (b *getAllUUIDMetadataBuilder) Filter(filter string) *getAllUUIDMetadataBuilder {
	b.opts.Filter = filter

	return b
}

func (b *getAllUUIDMetadataBuilder) Sort(sort []string) *getAllUUIDMetadataBuilder {
	b.opts.Sort = sort

	return b
}

func (b *getAllUUIDMetadataBuilder) Count(count bool) *getAllUUIDMetadataBuilder {
	b.opts.Count = count

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getAllUUIDMetadataBuilder) QueryParam(queryParam map[string]string) *getAllUUIDMetadataBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the getAllUUIDMetadata request.
func (b *getAllUUIDMetadataBuilder) Transport(tr http.RoundTripper) *getAllUUIDMetadataBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the getAllUUIDMetadata request.
func (b *getAllUUIDMetadataBuilder) Execute() (*WPSGetAllUUIDMetadataResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyWPSGetAllUUIDMetadataResponse, status, err
	}

	return newWPSGetAllUUIDMetadataResponse(rawJSON, b.opts, status)
}

type getAllUUIDMetadataOpts struct {
	webpubsub *WebPubSub

	Limit      int
	Include    []string
	Start      string
	End        string
	Filter     string
	Sort       []string
	Count      bool
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *getAllUUIDMetadataOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *getAllUUIDMetadataOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *getAllUUIDMetadataOpts) context() Context {
	return o.ctx
}

func (o *getAllUUIDMetadataOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *getAllUUIDMetadataOpts) buildPath() (string, error) {
	return fmt.Sprintf(getAllUUIDMetadataPath,
		o.webpubsub.Config.SubscribeKey), nil
}

func (o *getAllUUIDMetadataOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)

	if o.Include != nil {
		SetQueryParamAsCommaSepString(q, o.Include, "include")
	}

	q.Set("limit", strconv.Itoa(o.Limit))

	if o.Start != "" {
		q.Set("start", o.Start)
	}

	if o.Count {
		q.Set("count", "1")
	} else {
		q.Set("count", "0")
	}

	if o.End != "" {
		q.Set("end", o.End)
	}
	if o.Filter != "" {
		q.Set("filter", o.Filter)
	}
	if o.Sort != nil {
		SetQueryParamAsCommaSepString(q, o.Sort, "sort")
	}

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getAllUUIDMetadataOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *getAllUUIDMetadataOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getAllUUIDMetadataOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *getAllUUIDMetadataOpts) httpMethod() string {
	return "GET"
}

func (o *getAllUUIDMetadataOpts) isAuthRequired() bool {
	return true
}

func (o *getAllUUIDMetadataOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *getAllUUIDMetadataOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *getAllUUIDMetadataOpts) operationType() OperationType {
	return WPSGetAllUUIDMetadataOperation
}

func (o *getAllUUIDMetadataOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *getAllUUIDMetadataOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// WPSGetAllUUIDMetadataResponse is the Objects API Response for Get Users
type WPSGetAllUUIDMetadataResponse struct {
	status     int       `json:"status"`
	Data       []WPSUUID `json:"data"`
	TotalCount int       `json:"totalCount"`
	Next       string    `json:"next"`
	Prev       string    `json:"prev"`
}

func newWPSGetAllUUIDMetadataResponse(jsonBytes []byte, o *getAllUUIDMetadataOpts,
	status StatusResponse) (*WPSGetAllUUIDMetadataResponse, StatusResponse, error) {

	resp := &WPSGetAllUUIDMetadataResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyWPSGetAllUUIDMetadataResponse, status, e
	}

	return resp, status, nil
}
