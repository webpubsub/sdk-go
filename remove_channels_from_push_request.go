package webpubsub

import (
	"bytes"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/webpubsub/sdk-go/v7/utils"
)

const removeChannelsFromPushPath = "/v1/push/sub-key/%s/devices/%s"
const removeChannelsFromPushPathAPNS2 = "/v2/push/sub-key/%s/devices-apns2/%s"

var emptyRemoveChannelsFromPushResponse *RemoveChannelsFromPushResponse

type removeChannelsFromPushBuilder struct {
	opts *removeChannelsFromPushOpts
}

func newRemoveChannelsFromPushBuilder(webpubsub *WebPubSub) *removeChannelsFromPushBuilder {
	builder := removeChannelsFromPushBuilder{
		opts: &removeChannelsFromPushOpts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newRemoveChannelsFromPushBuilderWithContext(webpubsub *WebPubSub, context Context) *removeChannelsFromPushBuilder {
	builder := removeChannelsFromPushBuilder{
		opts: &removeChannelsFromPushOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

// Channels sets the channels to remove from Push Notifications
func (b *removeChannelsFromPushBuilder) Channels(channels []string) *removeChannelsFromPushBuilder {
	b.opts.Channels = channels
	return b
}

// PushType sets the PushType for the RemovePushNotificationsFromChannels request.
func (b *removeChannelsFromPushBuilder) PushType(pushType WPSPushType) *removeChannelsFromPushBuilder {
	b.opts.PushType = pushType
	return b
}

// DeviceIDForPush sets the DeviceIDForPush for the RemovePushNotificationsFromChannels request.
func (b *removeChannelsFromPushBuilder) DeviceIDForPush(deviceID string) *removeChannelsFromPushBuilder {
	b.opts.DeviceIDForPush = deviceID
	return b
}

// Topic sets the topic of for APNS2 Push Notifcataions
func (b *removeChannelsFromPushBuilder) Topic(topic string) *removeChannelsFromPushBuilder {
	b.opts.Topic = topic
	return b
}

// Environment sets the environment of for APNS2 Push Notifcataions
func (b *removeChannelsFromPushBuilder) Environment(env WPSPushEnvironment) *removeChannelsFromPushBuilder {
	b.opts.Environment = env
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *removeChannelsFromPushBuilder) QueryParam(queryParam map[string]string) *removeChannelsFromPushBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the RemovePushNotificationsFromChannels request.
func (b *removeChannelsFromPushBuilder) Execute() (*RemoveChannelsFromPushResponse, StatusResponse, error) {
	_, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyRemoveChannelsFromPushResponse, status, err
	}

	return emptyRemoveChannelsFromPushResponse, status, err
}

type removeChannelsFromPushOpts struct {
	webpubsub *WebPubSub

	Channels        []string
	QueryParam      map[string]string
	PushType        WPSPushType
	DeviceIDForPush string
	Topic           string
	Environment     WPSPushEnvironment

	Transport http.RoundTripper

	ctx Context
}

func (o *removeChannelsFromPushOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *removeChannelsFromPushOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *removeChannelsFromPushOpts) context() Context {
	return o.ctx
}

func (o *removeChannelsFromPushOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if len(o.Channels) == 0 {
		return newValidationError(o, StrMissingChannel)
	}

	if o.DeviceIDForPush == "" {
		return newValidationError(o, StrMissingDeviceID)
	}

	if o.PushType == WPSPushTypeNone {
		return newValidationError(o, StrMissingPushType)
	}

	if o.PushType == WPSPushTypeAPNS2 && (o.Topic == "") {
		return newValidationError(o, StrMissingPushTopic)
	}

	return nil
}

// RemoveChannelsFromPushResponse is the struct returned when the Execute function of RemovePushNotificationsFromChannels is called.
type RemoveChannelsFromPushResponse struct{}

func (o *removeChannelsFromPushOpts) buildPath() (string, error) {
	if o.PushType == WPSPushTypeAPNS2 {
		return fmt.Sprintf(removeChannelsFromPushPathAPNS2,
			o.webpubsub.Config.SubscribeKey,
			utils.URLEncode(o.DeviceIDForPush)), nil

	}
	return fmt.Sprintf(removeChannelsFromPushPath,
		o.webpubsub.Config.SubscribeKey,
		utils.URLEncode(o.DeviceIDForPush)), nil
}

func (o *removeChannelsFromPushOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)
	q.Set("type", o.PushType.String())

	var channels []string

	for _, v := range o.Channels {
		channels = append(channels, v)
	}

	q.Set("remove", strings.Join(channels, ","))
	SetPushEnvironment(q, o.Environment)
	SetPushTopic(q, o.Topic)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *removeChannelsFromPushOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *removeChannelsFromPushOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *removeChannelsFromPushOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *removeChannelsFromPushOpts) httpMethod() string {
	return "GET"
}

func (o *removeChannelsFromPushOpts) isAuthRequired() bool {
	return true
}

func (o *removeChannelsFromPushOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *removeChannelsFromPushOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *removeChannelsFromPushOpts) operationType() OperationType {
	return WPSRemoveGroupOperation
}

func (o *removeChannelsFromPushOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *removeChannelsFromPushOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}
