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

var emptyRemoveChannelMembersResponse *WPSRemoveChannelMembersResponse

const removeChannelMembersPath = "/v2/objects/%s/channels/%s/uuids"

const removeChannelMembersLimit = 100

type removeChannelMembersBuilder struct {
	opts *removeChannelMembersOpts
}

func newRemoveChannelMembersBuilder(webpubsub *WebPubSub) *removeChannelMembersBuilder {
	return newRemoveChannelMembersBuilderWithContext(webpubsub, webpubsub.ctx)
}

func newRemoveChannelMembersBuilderWithContext(webpubsub *WebPubSub,
	context Context) *removeChannelMembersBuilder {
	builder := removeChannelMembersBuilder{
		opts: &removeChannelMembersOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}
	builder.opts.Limit = removeChannelMembersLimit

	return &builder
}

func (b *removeChannelMembersBuilder) Include(include []WPSChannelMembersInclude) *removeChannelMembersBuilder {
	b.opts.Include = EnumArrayToStringArray(include)

	return b
}

func (b *removeChannelMembersBuilder) Channel(channel string) *removeChannelMembersBuilder {
	b.opts.Channel = channel

	return b
}

func (b *removeChannelMembersBuilder) Limit(limit int) *removeChannelMembersBuilder {
	b.opts.Limit = limit

	return b
}

func (b *removeChannelMembersBuilder) Start(start string) *removeChannelMembersBuilder {
	b.opts.Start = start

	return b
}

func (b *removeChannelMembersBuilder) End(end string) *removeChannelMembersBuilder {
	b.opts.End = end

	return b
}

func (b *removeChannelMembersBuilder) Count(count bool) *removeChannelMembersBuilder {
	b.opts.Count = count

	return b
}

func (b *removeChannelMembersBuilder) Filter(filter string) *removeChannelMembersBuilder {
	b.opts.Filter = filter

	return b
}

func (b *removeChannelMembersBuilder) Sort(sort []string) *removeChannelMembersBuilder {
	b.opts.Sort = sort

	return b
}

func (b *removeChannelMembersBuilder) Remove(channelMembersRemove []WPSChannelMembersRemove) *removeChannelMembersBuilder {
	b.opts.ChannelMembersRemove = channelMembersRemove

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *removeChannelMembersBuilder) QueryParam(queryParam map[string]string) *removeChannelMembersBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the removeChannelMembers request.
func (b *removeChannelMembersBuilder) Transport(tr http.RoundTripper) *removeChannelMembersBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the removeChannelMembers request.
func (b *removeChannelMembersBuilder) Execute() (*WPSRemoveChannelMembersResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyRemoveChannelMembersResponse, status, err
	}

	return newWPSRemoveChannelMembersResponse(rawJSON, b.opts, status)
}

type removeChannelMembersOpts struct {
	webpubsub            *WebPubSub
	Channel              string
	Limit                int
	Include              []string
	Start                string
	End                  string
	Filter               string
	Sort                 []string
	Count                bool
	QueryParam           map[string]string
	ChannelMembersRemove []WPSChannelMembersRemove
	Transport            http.RoundTripper

	ctx Context
}

func (o *removeChannelMembersOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *removeChannelMembersOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *removeChannelMembersOpts) context() Context {
	return o.ctx
}

func (o *removeChannelMembersOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.Channel == "" {
		return newValidationError(o, StrMissingChannel)
	}

	return nil
}

func (o *removeChannelMembersOpts) buildPath() (string, error) {
	return fmt.Sprintf(removeChannelMembersPath,
		o.webpubsub.Config.SubscribeKey, o.Channel), nil
}

func (o *removeChannelMembersOpts) buildQuery() (*url.Values, error) {

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

func (o *removeChannelMembersOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

// WPSChannelMembersRemoveChangeset is the Objects API input to add, remove or update membership
type WPSChannelMembersRemoveChangeset struct {
	Remove []WPSChannelMembersRemove `json:"delete"`
}

func (o *removeChannelMembersOpts) buildBody() ([]byte, error) {
	b := &WPSChannelMembersRemoveChangeset{
		Remove: o.ChannelMembersRemove,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.webpubsub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *removeChannelMembersOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *removeChannelMembersOpts) httpMethod() string {
	return "PATCH"
}

func (o *removeChannelMembersOpts) isAuthRequired() bool {
	return true
}

func (o *removeChannelMembersOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *removeChannelMembersOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *removeChannelMembersOpts) operationType() OperationType {
	return WPSRemoveChannelMembersOperation
}

func (o *removeChannelMembersOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *removeChannelMembersOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// WPSRemoveChannelMembersResponse is the Objects API Response for RemoveChannelMembers
type WPSRemoveChannelMembersResponse struct {
	status     int                 `json:"status"`
	Data       []WPSChannelMembers `json:"data"`
	TotalCount int                 `json:"totalCount"`
	Next       string              `json:"next"`
	Prev       string              `json:"prev"`
}

func newWPSRemoveChannelMembersResponse(jsonBytes []byte, o *removeChannelMembersOpts,
	status StatusResponse) (*WPSRemoveChannelMembersResponse, StatusResponse, error) {

	resp := &WPSRemoveChannelMembersResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyRemoveChannelMembersResponse, status, e
	}

	return resp, status, nil
}
