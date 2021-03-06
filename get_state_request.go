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
	"strings"

	"github.com/webpubsub/sdk-go/v7/pnerr"
	"github.com/webpubsub/sdk-go/v7/utils"
)

const getStatePath = "/v2/presence/sub-key/%s/channel/%s/uuid/%s"

var emptyGetStateResp *GetStateResponse

type getStateBuilder struct {
	opts *getStateOpts
}

func newGetStateBuilder(webpubsub *WebPubSub) *getStateBuilder {
	builder := getStateBuilder{
		opts: &getStateOpts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newGetStateBuilderWithContext(webpubsub *WebPubSub,
	context Context) *getStateBuilder {
	builder := getStateBuilder{
		opts: &getStateOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

// Channels sets the Channels for the Get State request.
func (b *getStateBuilder) Channels(ch []string) *getStateBuilder {
	b.opts.Channels = ch

	return b
}

// ChannelGroups sets the ChannelGroups for the Get State request.
func (b *getStateBuilder) ChannelGroups(cg []string) *getStateBuilder {
	b.opts.ChannelGroups = cg

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *getStateBuilder) QueryParam(queryParam map[string]string) *getStateBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// UUID sets the UUID for the Get State request.
func (b *getStateBuilder) UUID(uuid string) *getStateBuilder {
	b.opts.UUID = uuid

	return b
}

// Transport sets the Transport for the Get State request.
func (b *getStateBuilder) Transport(
	tr http.RoundTripper) *getStateBuilder {

	b.opts.Transport = tr

	return b
}

// Execute runs the the Get State request.
func (b *getStateBuilder) Execute() (
	*GetStateResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyGetStateResp, status, err
	}

	return newGetStateResponse(rawJSON, status)
}

type getStateOpts struct {
	webpubsub     *WebPubSub
	Channels      []string
	ChannelGroups []string
	UUID          string
	QueryParam    map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *getStateOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *getStateOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *getStateOpts) context() Context {
	return o.ctx
}

func (o *getStateOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if len(o.Channels) == 0 && len(o.ChannelGroups) == 0 {
		return newValidationError(o, "Missing Channel or Channel Group")
	}

	return nil
}

func (o *getStateOpts) buildPath() (string, error) {
	var channels []string

	for _, channel := range o.Channels {
		channels = append(channels, utils.PamEncode(channel))
	}

	uuid := o.UUID
	if uuid == "" {
		uuid = o.webpubsub.Config.UUID
	}

	return fmt.Sprintf(getStatePath,
		o.webpubsub.Config.SubscribeKey,
		strings.Join(channels, ","),
		utils.URLEncode(uuid)), nil
}

func (o *getStateOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)

	var groups []string

	for _, group := range o.ChannelGroups {
		groups = append(groups, utils.PamEncode(group))
	}

	q.Set("channel-group", strings.Join(groups, ","))
	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *getStateOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *getStateOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *getStateOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *getStateOpts) httpMethod() string {
	return "GET"
}

func (o *getStateOpts) isAuthRequired() bool {
	return true
}

func (o *getStateOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *getStateOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *getStateOpts) operationType() OperationType {
	return WPSGetStateOperation
}

func (o *getStateOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *getStateOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// GetStateResponse is the struct returned when the Execute function of GetState is called.
type GetStateResponse struct {
	State map[string]interface{}
	UUID  string
}

func newGetStateResponse(jsonBytes []byte, status StatusResponse) (
	*GetStateResponse, StatusResponse, error) {

	resp := &GetStateResponse{}

	var value interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyGetStateResp, status, e
	}

	v, ok := value.(map[string]interface{})
	if !ok {
		return emptyGetStateResp, status, errors.New("Response parsing error")
	}
	if v["error"] != nil {
		message := ""
		if v["message"] != nil {
			if msg, ok := v["message"].(string); ok {
				message = msg
			}
		}
		return emptyGetStateResp, status, errors.New(message)
	}

	if v["uuid"] != nil {
		resp.UUID = v["uuid"].(string)
	}
	m := make(map[string]interface{})
	if v["channel"] != nil {
		if channel, ok2 := v["channel"].(string); ok2 {
			if v["payload"] != nil {
				val, ok := v["payload"].(interface{})
				if !ok {
					return emptyGetStateResp, status, errors.New("Response parsing payload")
				}
				m[channel] = val
			} else {
				return emptyGetStateResp, status, errors.New("Response parsing channel")
			}
		} else {
			return emptyGetStateResp, status, errors.New("Response parsing channel 2")
		}
	} else {
		if v["payload"] != nil {
			val, ok := v["payload"].(map[string]interface{})
			if !ok {
				return emptyGetStateResp, status, errors.New("Response parsing payload 2")
			}
			channels, ok2 := val["channels"].(map[string]interface{})
			if !ok2 {
				return emptyGetStateResp, status, errors.New("Response parsing channels")
			}
			for ch, state := range channels {
				m[ch] = state
			}
		}

	}

	resp.State = m

	return resp, status, nil
}
