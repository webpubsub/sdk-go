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

const removeAllPushChannelsForDevicePath = "/v1/push/sub-key/%s/devices/%s/remove"
const removeAllPushChannelsForDevicePathAPNS2 = "/v2/push/sub-key/%s/devices-apns2/%s/remove"

var emptyRemoveAllPushChannelsForDeviceResponse *RemoveAllPushChannelsForDeviceResponse

type removeAllPushChannelsForDeviceBuilder struct {
	opts *removeAllPushChannelsForDeviceOpts
}

func newRemoveAllPushChannelsForDeviceBuilder(webpubsub *WebPubSub) *removeAllPushChannelsForDeviceBuilder {
	builder := removeAllPushChannelsForDeviceBuilder{
		opts: &removeAllPushChannelsForDeviceOpts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newRemoveAllPushChannelsForDeviceBuilderWithContext(
	webpubsub *WebPubSub, context Context) *removeAllPushChannelsForDeviceBuilder {
	builder := removeAllPushChannelsForDeviceBuilder{
		opts: &removeAllPushChannelsForDeviceOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

// PushType sets the PushType for the RemoveAllPushNotifications request.
func (b *removeAllPushChannelsForDeviceBuilder) PushType(pushType WPSPushType) *removeAllPushChannelsForDeviceBuilder {
	b.opts.PushType = pushType
	return b
}

// DeviceIDForPush sets the device id for RemoveAllPushNotifications request.
func (b *removeAllPushChannelsForDeviceBuilder) DeviceIDForPush(deviceID string) *removeAllPushChannelsForDeviceBuilder {
	b.opts.DeviceIDForPush = deviceID
	return b
}

// Topic sets the topic of for APNS2 Push Notifcataions
func (b *removeAllPushChannelsForDeviceBuilder) Topic(topic string) *removeAllPushChannelsForDeviceBuilder {
	b.opts.Topic = topic
	return b
}

// Environment sets the environment of for APNS2 Push Notifcataions
func (b *removeAllPushChannelsForDeviceBuilder) Environment(env WPSPushEnvironment) *removeAllPushChannelsForDeviceBuilder {
	b.opts.Environment = env
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *removeAllPushChannelsForDeviceBuilder) QueryParam(queryParam map[string]string) *removeAllPushChannelsForDeviceBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the RemoveAllPushNotifications request.
func (b *removeAllPushChannelsForDeviceBuilder) Execute() (
	*RemoveAllPushChannelsForDeviceResponse, StatusResponse, error) {
	_, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyRemoveAllPushChannelsForDeviceResponse, status, err
	}

	return emptyRemoveAllPushChannelsForDeviceResponse, status, err
}

type removeAllPushChannelsForDeviceOpts struct {
	webpubsub *WebPubSub

	PushType        WPSPushType
	QueryParam      map[string]string
	DeviceIDForPush string
	Topic           string
	Environment     WPSPushEnvironment

	Transport http.RoundTripper

	ctx Context
}

func (o *removeAllPushChannelsForDeviceOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *removeAllPushChannelsForDeviceOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *removeAllPushChannelsForDeviceOpts) context() Context {
	return o.ctx
}

func (o *removeAllPushChannelsForDeviceOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
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

// RemoveAllPushChannelsForDeviceResponse is the struct returned when the Execute function of RemoveAllPushNotifications is called.
type RemoveAllPushChannelsForDeviceResponse struct{}

func (o *removeAllPushChannelsForDeviceOpts) buildPath() (string, error) {
	if o.PushType == WPSPushTypeAPNS2 {
		return fmt.Sprintf(removeAllPushChannelsForDevicePathAPNS2,
			o.webpubsub.Config.SubscribeKey,
			utils.URLEncode(o.DeviceIDForPush)), nil
	}

	return fmt.Sprintf(removeAllPushChannelsForDevicePath,
		o.webpubsub.Config.SubscribeKey,
		utils.URLEncode(o.DeviceIDForPush)), nil
}

func (o *removeAllPushChannelsForDeviceOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)
	q.Set("type", o.PushType.String())
	SetPushEnvironment(q, o.Environment)
	SetPushTopic(q, o.Topic)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *removeAllPushChannelsForDeviceOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *removeAllPushChannelsForDeviceOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *removeAllPushChannelsForDeviceOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *removeAllPushChannelsForDeviceOpts) httpMethod() string {
	return "GET"
}

func (o *removeAllPushChannelsForDeviceOpts) isAuthRequired() bool {
	return true
}

func (o *removeAllPushChannelsForDeviceOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *removeAllPushChannelsForDeviceOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *removeAllPushChannelsForDeviceOpts) operationType() OperationType {
	return WPSRemoveGroupOperation
}

func (o *removeAllPushChannelsForDeviceOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *removeAllPushChannelsForDeviceOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}
