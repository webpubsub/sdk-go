package webpubsub

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/webpubsub/sdk-go/v7/utils"
)

const deleteChannelGroupPath = "/v1/channel-registration/sub-key/%s/channel-group/%s/remove"

var emptyDeleteChannelGroupResponse *DeleteChannelGroupResponse

type deleteChannelGroupBuilder struct {
	opts *deleteChannelGroupOpts
}

func newDeleteChannelGroupBuilder(webpubsub *WebPubSub) *deleteChannelGroupBuilder {
	builder := deleteChannelGroupBuilder{
		opts: &deleteChannelGroupOpts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newDeleteChannelGroupBuilderWithContext(
	webpubsub *WebPubSub, context Context) *deleteChannelGroupBuilder {
	builder := deleteChannelGroupBuilder{
		opts: &deleteChannelGroupOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

// ChannelGroup sets the channel group to delete.
func (b *deleteChannelGroupBuilder) ChannelGroup(
	cg string) *deleteChannelGroupBuilder {
	b.opts.ChannelGroup = cg
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *deleteChannelGroupBuilder) QueryParam(queryParam map[string]string) *deleteChannelGroupBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the DeleteChannelGroup request.
func (b *deleteChannelGroupBuilder) Execute() (
	*DeleteChannelGroupResponse, StatusResponse, error) {
	_, status, err := executeRequest(b.opts)

	if err != nil {
		return emptyDeleteChannelGroupResponse, status, err
	}

	return emptyDeleteChannelGroupResponse, status, nil
}

type deleteChannelGroupOpts struct {
	webpubsub    *WebPubSub
	ChannelGroup string
	Transport    http.RoundTripper
	QueryParam   map[string]string
	ctx          Context
}

func (o *deleteChannelGroupOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *deleteChannelGroupOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *deleteChannelGroupOpts) context() Context {
	return o.ctx
}

func (o *deleteChannelGroupOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.ChannelGroup == "" {
		return newValidationError(o, StrMissingChannelGroup)
	}

	return nil
}

// DeleteChannelGroupResponse is response structure for Delete Channel Group function
type DeleteChannelGroupResponse struct{}

func (o *deleteChannelGroupOpts) buildPath() (string, error) {
	return fmt.Sprintf(deleteChannelGroupPath,
		o.webpubsub.Config.SubscribeKey,
		utils.URLEncode(o.ChannelGroup)), nil
}

func (o *deleteChannelGroupOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *deleteChannelGroupOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *deleteChannelGroupOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *deleteChannelGroupOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *deleteChannelGroupOpts) httpMethod() string {
	return "GET"
}

func (o *deleteChannelGroupOpts) isAuthRequired() bool {
	return true
}

func (o *deleteChannelGroupOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *deleteChannelGroupOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *deleteChannelGroupOpts) operationType() OperationType {
	return WPSRemoveGroupOperation
}

func (o *deleteChannelGroupOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *deleteChannelGroupOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}
