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

var emptySetChannelMembersResponse *WPSSetChannelMembersResponse

const setChannelMembersPath = "/v2/objects/%s/channels/%s/uuids"

const setChannelMembersLimit = 100

type setChannelMembersBuilder struct {
	opts *setChannelMembersOpts
}

func newSetChannelMembersBuilder(webpubsub *WebPubSub) *setChannelMembersBuilder {
	return newSetChannelMembersBuilderWithContext(webpubsub, webpubsub.ctx)
}

func newSetChannelMembersBuilderWithContext(webpubsub *WebPubSub,
	context Context) *setChannelMembersBuilder {
	builder := setChannelMembersBuilder{
		opts: &setChannelMembersOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}
	builder.opts.Limit = setChannelMembersLimit

	return &builder
}

func (b *setChannelMembersBuilder) Include(include []WPSChannelMembersInclude) *setChannelMembersBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *setChannelMembersBuilder) Channel(channel string) *setChannelMembersBuilder {
	b.opts.Channel = channel

	return b
}

func (b *setChannelMembersBuilder) Limit(limit int) *setChannelMembersBuilder {
	b.opts.Limit = limit

	return b
}

func (b *setChannelMembersBuilder) Start(start string) *setChannelMembersBuilder {
	b.opts.Start = start

	return b
}

func (b *setChannelMembersBuilder) End(end string) *setChannelMembersBuilder {
	b.opts.End = end

	return b
}

func (b *setChannelMembersBuilder) Count(count bool) *setChannelMembersBuilder {
	b.opts.Count = count

	return b
}

func (b *setChannelMembersBuilder) Filter(filter string) *setChannelMembersBuilder {
	b.opts.Filter = filter

	return b
}

func (b *setChannelMembersBuilder) Sort(sort []string) *setChannelMembersBuilder {
	b.opts.Sort = sort

	return b
}

func (b *setChannelMembersBuilder) Set(channelMembersSet []WPSChannelMembersSet) *setChannelMembersBuilder {
	b.opts.ChannelMembersSet = channelMembersSet

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *setChannelMembersBuilder) QueryParam(queryParam map[string]string) *setChannelMembersBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the setChannelMembers request.
func (b *setChannelMembersBuilder) Transport(tr http.RoundTripper) *setChannelMembersBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the setChannelMembers request.
func (b *setChannelMembersBuilder) Execute() (*WPSSetChannelMembersResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptySetChannelMembersResponse, status, err
	}

	return newWPSSetChannelMembersResponse(rawJSON, b.opts, status)
}

type setChannelMembersOpts struct {
	webpubsub         *WebPubSub
	Channel           string
	Limit             int
	Include           []string
	Start             string
	End               string
	Filter            string
	Sort              []string
	Count             bool
	QueryParam        map[string]string
	ChannelMembersSet []WPSChannelMembersSet
	Transport         http.RoundTripper

	ctx Context
}

func (o *setChannelMembersOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *setChannelMembersOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *setChannelMembersOpts) context() Context {
	return o.ctx
}

func (o *setChannelMembersOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.Channel == "" {
		return newValidationError(o, StrMissingChannel)
	}

	return nil
}

func (o *setChannelMembersOpts) buildPath() (string, error) {
	return fmt.Sprintf(setChannelMembersPath,
		o.webpubsub.Config.SubscribeKey, o.Channel), nil
}

func (o *setChannelMembersOpts) buildQuery() (*url.Values, error) {

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

func (o *setChannelMembersOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

// WPSChannelMembersSetChangeset is the Objects API input to add, remove or update membership
type WPSChannelMembersSetChangeset struct {
	Set []WPSChannelMembersSet `json:"set"`
}

func (o *setChannelMembersOpts) buildBody() ([]byte, error) {
	b := &WPSChannelMembersSetChangeset{
		Set: o.ChannelMembersSet,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.webpubsub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *setChannelMembersOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *setChannelMembersOpts) httpMethod() string {
	return "PATCH"
}

func (o *setChannelMembersOpts) isAuthRequired() bool {
	return true
}

func (o *setChannelMembersOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *setChannelMembersOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *setChannelMembersOpts) operationType() OperationType {
	return WPSSetChannelMembersOperation
}

func (o *setChannelMembersOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *setChannelMembersOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// WPSSetChannelMembersResponse is the Objects API Response for SetChannelMembers
type WPSSetChannelMembersResponse struct {
	status     int                 `json:"status"`
	Data       []WPSChannelMembers `json:"data"`
	TotalCount int                 `json:"totalCount"`
	Next       string              `json:"next"`
	Prev       string              `json:"prev"`
}

func newWPSSetChannelMembersResponse(jsonBytes []byte, o *setChannelMembersOpts,
	status StatusResponse) (*WPSSetChannelMembersResponse, StatusResponse, error) {

	resp := &WPSSetChannelMembersResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptySetChannelMembersResponse, status, e
	}

	return resp, status, nil
}
