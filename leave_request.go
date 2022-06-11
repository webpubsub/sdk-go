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

const leavePath = "/v2/presence/sub-key/%s/channel/%s/leave"

type leaveBuilder struct {
	opts *leaveOpts
}

func newLeaveBuilder(webpubsub *WebPubSub) *leaveBuilder {
	builder := leaveBuilder{
		opts: &leaveOpts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newLeaveBuilderWithContext(webpubsub *WebPubSub, context Context) *leaveBuilder {
	builder := leaveBuilder{
		opts: &leaveOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

// Channels sets the channel names in the Unsubscribe request.
func (b *leaveBuilder) Channels(channels []string) *leaveBuilder {
	b.opts.Channels = channels
	return b
}

// ChannelGroups sets the channel group names in the Unsubscribe request.
func (b *leaveBuilder) ChannelGroups(groups []string) *leaveBuilder {
	b.opts.ChannelGroups = groups
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *leaveBuilder) QueryParam(queryParam map[string]string) *leaveBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the Leave request.
func (b *leaveBuilder) Execute() (StatusResponse, error) {
	_, status, err := executeRequest(b.opts)
	if err != nil {
		return status, err
	}

	return status, nil
}

type leaveOpts struct {
	Channels      []string
	ChannelGroups []string
	QueryParam    map[string]string

	webpubsub *WebPubSub
	ctx       Context
}

func (o *leaveOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *leaveOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *leaveOpts) httpMethod() string {
	return "GET"
}

func (o *leaveOpts) buildPath() (string, error) {
	channels := utils.JoinChannels(o.Channels)

	if string(channels) == "" {
		channels = []byte(",")
	}

	return fmt.Sprintf(leavePath,
		o.webpubsub.Config.SubscribeKey,
		channels), nil
}

func (o *leaveOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *leaveOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)

	if len(o.ChannelGroups) > 0 {
		channelGroup := utils.JoinChannels(o.ChannelGroups)
		q.Set("channel-group", string(channelGroup))
	}
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *leaveOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *leaveOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *leaveOpts) context() Context {
	return o.ctx
}

func (o *leaveOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if len(o.Channels) == 0 && len(o.ChannelGroups) == 0 {
		return newValidationError(o, "Missing Channel or Channel Group")
	}

	return nil
}

func (o *leaveOpts) operationType() OperationType {
	return WPSUnsubscribeOperation
}

func (o *leaveOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *leaveOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}
