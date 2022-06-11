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

var emptyRemoveMembershipsResponse *WPSRemoveMembershipsResponse

const removeMembershipsPath = "/v2/objects/%s/uuids/%s/channels"

const removeMembershipsLimit = 100

type removeMembershipsBuilder struct {
	opts *removeMembershipsOpts
}

func newRemoveMembershipsBuilder(webpubsub *WebPubSub) *removeMembershipsBuilder {
	return newRemoveMembershipsBuilderWithContext(webpubsub, webpubsub.ctx)
}

func newRemoveMembershipsBuilderWithContext(webpubsub *WebPubSub,
	context Context) *removeMembershipsBuilder {
	builder := removeMembershipsBuilder{
		opts: &removeMembershipsOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}
	builder.opts.Limit = removeMembershipsLimit

	return &builder
}

func (b *removeMembershipsBuilder) Include(include []WPSMembershipsInclude) *removeMembershipsBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *removeMembershipsBuilder) UUID(uuid string) *removeMembershipsBuilder {
	b.opts.UUID = uuid

	return b
}

func (b *removeMembershipsBuilder) Limit(limit int) *removeMembershipsBuilder {
	b.opts.Limit = limit

	return b
}

func (b *removeMembershipsBuilder) Start(start string) *removeMembershipsBuilder {
	b.opts.Start = start

	return b
}

func (b *removeMembershipsBuilder) End(end string) *removeMembershipsBuilder {
	b.opts.End = end

	return b
}

func (b *removeMembershipsBuilder) Count(count bool) *removeMembershipsBuilder {
	b.opts.Count = count

	return b
}

func (b *removeMembershipsBuilder) Filter(filter string) *removeMembershipsBuilder {
	b.opts.Filter = filter

	return b
}

func (b *removeMembershipsBuilder) Sort(sort []string) *removeMembershipsBuilder {
	b.opts.Sort = sort

	return b
}

func (b *removeMembershipsBuilder) Remove(membershipsRemove []WPSMembershipsRemove) *removeMembershipsBuilder {
	b.opts.MembershipsRemove = membershipsRemove

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *removeMembershipsBuilder) QueryParam(queryParam map[string]string) *removeMembershipsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the removeMemberships request.
func (b *removeMembershipsBuilder) Transport(tr http.RoundTripper) *removeMembershipsBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the removeMemberships request.
func (b *removeMembershipsBuilder) Execute() (*WPSRemoveMembershipsResponse, StatusResponse, error) {
	if len(b.opts.UUID) <= 0 {
		b.opts.UUID = b.opts.webpubsub.Config.UUID
	}

	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyRemoveMembershipsResponse, status, err
	}

	return newWPSRemoveMembershipsResponse(rawJSON, b.opts, status)
}

type removeMembershipsOpts struct {
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
	Transport         http.RoundTripper

	ctx Context
}

func (o *removeMembershipsOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *removeMembershipsOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *removeMembershipsOpts) context() Context {
	return o.ctx
}

func (o *removeMembershipsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *removeMembershipsOpts) buildPath() (string, error) {
	return fmt.Sprintf(removeMembershipsPath,
		o.webpubsub.Config.SubscribeKey, o.UUID), nil
}

func (o *removeMembershipsOpts) buildQuery() (*url.Values, error) {

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

func (o *removeMembershipsOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

// WPSMembershipsRemoveChangeSet is the Objects API input to add, remove or update members
type WPSMembershipsRemoveChangeSet struct {
	Remove []WPSMembershipsRemove `json:"delete"`
}

func (o *removeMembershipsOpts) buildBody() ([]byte, error) {
	b := &WPSMembershipsRemoveChangeSet{
		Remove: o.MembershipsRemove,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.webpubsub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil

}

func (o *removeMembershipsOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *removeMembershipsOpts) httpMethod() string {
	return "PATCH"
}

func (o *removeMembershipsOpts) isAuthRequired() bool {
	return true
}

func (o *removeMembershipsOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *removeMembershipsOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *removeMembershipsOpts) operationType() OperationType {
	return WPSRemoveMembershipsOperation
}

func (o *removeMembershipsOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *removeMembershipsOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// WPSRemoveMembershipsResponse is the Objects API Response for RemoveMemberships
type WPSRemoveMembershipsResponse struct {
	status     int              `json:"status"`
	Data       []WPSMemberships `json:"data"`
	TotalCount int              `json:"totalCount"`
	Next       string           `json:"next"`
	Prev       string           `json:"prev"`
}

func newWPSRemoveMembershipsResponse(jsonBytes []byte, o *removeMembershipsOpts,
	status StatusResponse) (*WPSRemoveMembershipsResponse, StatusResponse, error) {

	resp := &WPSRemoveMembershipsResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyRemoveMembershipsResponse, status, e
	}

	return resp, status, nil
}
