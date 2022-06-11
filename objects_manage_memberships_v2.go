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

var emptyManageMembershipsResponse *WPSManageMembershipsResponse

const manageMembershipsPathV2 = "/v2/objects/%s/uuids/%s/channels"

const manageMembershipsLimitV2 = 100

type manageMembershipsBuilderV2 struct {
	opts *manageMembershipsOptsV2
}

func newManageMembershipsBuilderV2(webpubsub *WebPubSub) *manageMembershipsBuilderV2 {
	return newManageMembershipsBuilderV2WithContext(webpubsub, webpubsub.ctx)
}

func newManageMembershipsBuilderV2WithContext(webpubsub *WebPubSub,
	context Context) *manageMembershipsBuilderV2 {
	builder := manageMembershipsBuilderV2{
		opts: &manageMembershipsOptsV2{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}
	builder.opts.Limit = manageMembershipsLimitV2

	return &builder
}

func (b *manageMembershipsBuilderV2) Include(include []WPSMembershipsInclude) *manageMembershipsBuilderV2 {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *manageMembershipsBuilderV2) UUID(uuid string) *manageMembershipsBuilderV2 {
	b.opts.UUID = uuid

	return b
}

func (b *manageMembershipsBuilderV2) Limit(limit int) *manageMembershipsBuilderV2 {
	b.opts.Limit = limit

	return b
}

func (b *manageMembershipsBuilderV2) Start(start string) *manageMembershipsBuilderV2 {
	b.opts.Start = start

	return b
}

func (b *manageMembershipsBuilderV2) End(end string) *manageMembershipsBuilderV2 {
	b.opts.End = end

	return b
}

func (b *manageMembershipsBuilderV2) Count(count bool) *manageMembershipsBuilderV2 {
	b.opts.Count = count

	return b
}

func (b *manageMembershipsBuilderV2) Filter(filter string) *manageMembershipsBuilderV2 {
	b.opts.Filter = filter

	return b
}

func (b *manageMembershipsBuilderV2) Sort(sort []string) *manageMembershipsBuilderV2 {
	b.opts.Sort = sort

	return b
}

func (b *manageMembershipsBuilderV2) Set(membershipsSet []WPSMembershipsSet) *manageMembershipsBuilderV2 {
	b.opts.MembershipsSet = membershipsSet

	return b
}

func (b *manageMembershipsBuilderV2) Remove(membershipsRemove []WPSMembershipsRemove) *manageMembershipsBuilderV2 {
	b.opts.MembershipsRemove = membershipsRemove

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *manageMembershipsBuilderV2) QueryParam(queryParam map[string]string) *manageMembershipsBuilderV2 {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the manageMemberships request.
func (b *manageMembershipsBuilderV2) Transport(tr http.RoundTripper) *manageMembershipsBuilderV2 {
	b.opts.Transport = tr
	return b
}

// Execute runs the manageMemberships request.
func (b *manageMembershipsBuilderV2) Execute() (*WPSManageMembershipsResponse, StatusResponse, error) {
	if len(b.opts.UUID) <= 0 {
		b.opts.UUID = b.opts.webpubsub.Config.UUID
	}

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyManageMembershipsResponse, status, err
	}

	return newWPSManageMembershipsResponse(rawJSON, b.opts, status)
}

type manageMembershipsOptsV2 struct {
	webpubsub         *WebPubSub
	UUID              string
	Limit             int
	Include           []string
	Start             string
	End               string
	Filter            string
	Sort              []string
	Count             bool
	QueryParam        map[string]string
	MembershipsRemove []WPSMembershipsRemove
	MembershipsSet    []WPSMembershipsSet
	Transport         http.RoundTripper

	ctx Context
}

func (o *manageMembershipsOptsV2) config() Config {
	return *o.webpubsub.Config
}

func (o *manageMembershipsOptsV2) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *manageMembershipsOptsV2) context() Context {
	return o.ctx
}

func (o *manageMembershipsOptsV2) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *manageMembershipsOptsV2) buildPath() (string, error) {
	return fmt.Sprintf(manageMembershipsPathV2,
		o.webpubsub.Config.SubscribeKey, o.UUID), nil
}

func (o *manageMembershipsOptsV2) buildQuery() (*url.Values, error) {

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

func (o *manageMembershipsOptsV2) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *manageMembershipsOptsV2) buildBody() ([]byte, error) {
	b := &WPSManageMembershipsBody{
		Set:    o.MembershipsSet,
		Remove: o.MembershipsRemove,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.webpubsub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *manageMembershipsOptsV2) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *manageMembershipsOptsV2) httpMethod() string {
	return "PATCH"
}

func (o *manageMembershipsOptsV2) isAuthRequired() bool {
	return true
}

func (o *manageMembershipsOptsV2) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *manageMembershipsOptsV2) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *manageMembershipsOptsV2) operationType() OperationType {
	return WPSManageMembershipsOperation
}

func (o *manageMembershipsOptsV2) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *manageMembershipsOptsV2) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// WPSManageMembershipsResponse is the Objects API Response for ManageMemberships
type WPSManageMembershipsResponse struct {
	status     int              `json:"status"`
	Data       []WPSMemberships `json:"data"`
	TotalCount int              `json:"totalCount"`
	Next       string           `json:"next"`
	Prev       string           `json:"prev"`
}

func newWPSManageMembershipsResponse(jsonBytes []byte, o *manageMembershipsOptsV2,
	status StatusResponse) (*WPSManageMembershipsResponse, StatusResponse, error) {

	resp := &WPSManageMembershipsResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyManageMembershipsResponse, status, e
	}

	return resp, status, nil
}
