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
)

var emptyWPSAddMessageActionsResponse *WPSAddMessageActionsResponse

const addMessageActionsPath = "/v1/message-actions/%s/channel/%s/message/%s"

type addMessageActionsBuilder struct {
	opts *addMessageActionsOpts
}

func newAddMessageActionsBuilder(webpubsub *WebPubSub) *addMessageActionsBuilder {
	builder := addMessageActionsBuilder{
		opts: &addMessageActionsOpts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newAddMessageActionsBuilderWithContext(webpubsub *WebPubSub,
	context Context) *addMessageActionsBuilder {
	builder := addMessageActionsBuilder{
		opts: &addMessageActionsOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

// MessageAction struct is used to create a Message Action
type MessageAction struct {
	ActionType  string `json:"type"`
	ActionValue string `json:"value"`
}

func (b *addMessageActionsBuilder) Channel(channel string) *addMessageActionsBuilder {
	b.opts.Channel = channel

	return b
}

func (b *addMessageActionsBuilder) MessageTimetoken(timetoken string) *addMessageActionsBuilder {
	b.opts.MessageTimetoken = timetoken

	return b
}

func (b *addMessageActionsBuilder) Action(action MessageAction) *addMessageActionsBuilder {
	b.opts.Action = action

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *addMessageActionsBuilder) QueryParam(queryParam map[string]string) *addMessageActionsBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the addMessageActions request.
func (b *addMessageActionsBuilder) Transport(tr http.RoundTripper) *addMessageActionsBuilder {
	b.opts.Transport = tr
	return b
}

// Execute runs the addMessageActions request.
func (b *addMessageActionsBuilder) Execute() (*WPSAddMessageActionsResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyWPSAddMessageActionsResponse, status, err
	}

	return newWPSAddMessageActionsResponse(rawJSON, b.opts, status)
}

type addMessageActionsOpts struct {
	webpubsub *WebPubSub

	Channel          string
	MessageTimetoken string
	Action           MessageAction
	QueryParam       map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *addMessageActionsOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *addMessageActionsOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *addMessageActionsOpts) context() Context {
	return o.ctx
}

func (o *addMessageActionsOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *addMessageActionsOpts) buildPath() (string, error) {
	return fmt.Sprintf(addMessageActionsPath,
		o.webpubsub.Config.SubscribeKey, o.Channel, o.MessageTimetoken), nil
}

func (o *addMessageActionsOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *addMessageActionsOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *addMessageActionsOpts) buildBody() ([]byte, error) {
	jsonEncBytes, errEnc := json.Marshal(o.Action)

	if errEnc != nil {
		o.webpubsub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil

}

func (o *addMessageActionsOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *addMessageActionsOpts) httpMethod() string {
	return "POST"
}

func (o *addMessageActionsOpts) isAuthRequired() bool {
	return true
}

func (o *addMessageActionsOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *addMessageActionsOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *addMessageActionsOpts) operationType() OperationType {
	return WPSAddMessageActionsOperation
}

func (o *addMessageActionsOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *addMessageActionsOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// WPSMessageActionsResponse Message Actions response.
type WPSMessageActionsResponse struct {
	ActionType       string `json:"type"`
	ActionValue      string `json:"value"`
	ActionTimetoken  string `json:"actionTimetoken"`
	MessageTimetoken string `json:"messageTimetoken"`
	UUID             string `json:"uuid"`
}

// WPSAddMessageActionsResponse is the Add Message Actions API Response
type WPSAddMessageActionsResponse struct {
	status int                       `json:"status"`
	Data   WPSMessageActionsResponse `json:"data"`
}

func newWPSAddMessageActionsResponse(jsonBytes []byte, o *addMessageActionsOpts,
	status StatusResponse) (*WPSAddMessageActionsResponse, StatusResponse, error) {

	resp := &WPSAddMessageActionsResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyWPSAddMessageActionsResponse, status, e
	}

	return resp, status, nil
}
