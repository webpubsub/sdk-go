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

	"github.com/webpubsub/sdk-go/v7/pnerr"
	"github.com/webpubsub/sdk-go/v7/utils"
)

const allChannelGroupPath = "/v1/channel-registration/sub-key/%s/channel-group/%s"

var emptyAllChannelGroupResponse *AllChannelGroupResponse

type allChannelGroupBuilder struct {
	opts *allChannelGroupOpts
}

func newAllChannelGroupBuilder(webpubsub *WebPubSub) *allChannelGroupBuilder {
	builder := allChannelGroupBuilder{
		opts: &allChannelGroupOpts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newAllChannelGroupBuilderWithContext(webpubsub *WebPubSub,
	context Context) *allChannelGroupBuilder {
	builder := allChannelGroupBuilder{
		opts: &allChannelGroupOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

// ChannelGroup sets the channel group to list channels.
func (b *allChannelGroupBuilder) ChannelGroup(
	cg string) *allChannelGroupBuilder {
	b.opts.ChannelGroup = cg
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *allChannelGroupBuilder) QueryParam(queryParam map[string]string) *allChannelGroupBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the ListChannelsInChannelGroup request.
func (b *allChannelGroupBuilder) Execute() (
	*AllChannelGroupResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyAllChannelGroupResponse, status, err
	}

	return newAllChannelGroupResponse(rawJSON, status)
}

type allChannelGroupOpts struct {
	webpubsub *WebPubSub

	ChannelGroup string
	QueryParam   map[string]string
	Transport    http.RoundTripper

	ctx Context
}

func (o *allChannelGroupOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *allChannelGroupOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *allChannelGroupOpts) context() Context {
	return o.ctx
}

func (o *allChannelGroupOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.ChannelGroup == "" {
		return newValidationError(o, StrMissingChannelGroup)
	}

	return nil
}

func (o *allChannelGroupOpts) buildPath() (string, error) {
	return fmt.Sprintf(allChannelGroupPath,
		o.webpubsub.Config.SubscribeKey,
		utils.URLEncode(o.ChannelGroup)), nil
}

func (o *allChannelGroupOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)
	SetQueryParam(q, o.QueryParam)
	return q, nil
}

func (o *allChannelGroupOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *allChannelGroupOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *allChannelGroupOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *allChannelGroupOpts) httpMethod() string {
	return "GET"
}

func (o *allChannelGroupOpts) isAuthRequired() bool {
	return true
}

func (o *allChannelGroupOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *allChannelGroupOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *allChannelGroupOpts) operationType() OperationType {
	return WPSChannelsForGroupOperation
}

func (o *allChannelGroupOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *allChannelGroupOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// AllChannelGroupResponse is the struct returned when the Execute function of List All Channel Groups is called.
type AllChannelGroupResponse struct {
	Channels     []string
	ChannelGroup string
}

func newAllChannelGroupResponse(jsonBytes []byte, status StatusResponse) (
	*AllChannelGroupResponse, StatusResponse, error) {
	resp := &AllChannelGroupResponse{}

	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyAllChannelGroupResponse, status, e
	}

	if parsedValue, ok := value.(map[string]interface{}); ok {
		if payload, ok := parsedValue["payload"].(map[string]interface{}); ok {
			if group, ok := payload["group"].(string); ok {
				resp.ChannelGroup = group
			}

			if channels, ok := payload["channels"].([]interface{}); ok {
				parsedChannels := []string{}

				for _, channel := range channels {
					if ch, ok := channel.(string); ok {
						parsedChannels = append(parsedChannels, ch)
					}
				}

				resp.Channels = parsedChannels
			}
		}
	}

	return resp, status, nil
}
