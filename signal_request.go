package webpubsub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"

	"github.com/webpubsub/sdk-go/v7/pnerr"
	"github.com/webpubsub/sdk-go/v7/utils"

	"net/http"
	"net/url"
	"strconv"
)

var emptySignalResponse *SignalResponse

const signalGetPath = "/signal/%s/%s/0/%s/%s/%s"
const signalPostPath = "/signal/%s/%s/0/%s/%s"

type signalBuilder struct {
	opts *signalOpts
}

func newSignalBuilder(webpubsub *WebPubSub) *signalBuilder {
	builder := signalBuilder{
		opts: &signalOpts{
			webpubsub: webpubsub,
		},
	}
	builder.opts.UsePost = false

	return &builder
}

func newSignalBuilderWithContext(webpubsub *WebPubSub,
	context Context) *signalBuilder {
	builder := signalBuilder{
		opts: &signalOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

// Channel sets the Channel for the Signal request.
func (b *signalBuilder) Channel(channel string) *signalBuilder {
	b.opts.Channel = channel
	return b
}

// Message sets the Payload for the Signal request.
func (b *signalBuilder) Message(msg interface{}) *signalBuilder {
	b.opts.Message = msg

	return b
}

// usePost sends the Signal request using HTTP POST. Not implemented
func (b *signalBuilder) usePost(post bool) *signalBuilder {
	b.opts.UsePost = post

	return b
}

// Transport sets the Transport for the objectAPICreateUsers request.
func (b *signalBuilder) Transport(tr http.RoundTripper) *signalBuilder {
	b.opts.Transport = tr
	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *signalBuilder) QueryParam(queryParam map[string]string) *signalBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the Signal request.
func (b *signalBuilder) Execute() (*SignalResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptySignalResponse, status, err
	}

	return newSignalResponse(rawJSON, b.opts, status)
}

type signalOpts struct {
	webpubsub  *WebPubSub
	Message    interface{}
	Channel    string
	UsePost    bool
	QueryParam map[string]string
	Transport  http.RoundTripper
	ctx        Context
}

func (o *signalOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *signalOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *signalOpts) context() Context {
	return o.ctx
}

func (o *signalOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.config().PublishKey == "" {
		return newValidationError(o, StrMissingPubKey)
	}

	return nil
}

func (o *signalOpts) buildPath() (string, error) {
	if o.UsePost == true {
		return fmt.Sprintf(signalPostPath,
			o.webpubsub.Config.PublishKey,
			o.webpubsub.Config.SubscribeKey,
			utils.URLEncode(o.Channel),
			"0"), nil
	}

	var msg string
	jsonEncBytes, errEnc := json.Marshal(o.Message)
	if errEnc != nil {
		o.webpubsub.Config.Log.Printf("ERROR: Publish error: %s\n", errEnc.Error())
		return "", errEnc
	}
	msg = string(jsonEncBytes)
	return fmt.Sprintf(signalGetPath,
		o.webpubsub.Config.PublishKey,
		o.webpubsub.Config.SubscribeKey,
		utils.URLEncode(o.Channel),
		"0",
		utils.URLEncode(msg),
	), nil
}

func (o *signalOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *signalOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *signalOpts) buildBody() ([]byte, error) {
	if o.UsePost {
		jsonEncBytes, errEnc := json.Marshal(o.Message)
		if errEnc != nil {
			o.webpubsub.Config.Log.Printf("ERROR: Signal error: %s\n", errEnc.Error())
			return []byte{}, errEnc
		}
		return jsonEncBytes, nil
	}
	return []byte{}, nil
}

func (o *signalOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *signalOpts) httpMethod() string {
	if o.UsePost {
		return "POST"
	}
	return "GET"
}

func (o *signalOpts) isAuthRequired() bool {
	return true
}

func (o *signalOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *signalOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *signalOpts) operationType() OperationType {
	return WPSSignalOperation
}

func (o *signalOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *signalOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// SignalResponse is the response to Signal request.
type SignalResponse struct {
	Timestamp int64
}

func newSignalResponse(jsonBytes []byte, o *signalOpts,
	status StatusResponse) (*SignalResponse, StatusResponse, error) {

	resp := &SignalResponse{}

	var value []interface{}

	err := json.Unmarshal(jsonBytes, &value)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptySignalResponse, status, e
	}

	if len(value) > 1 {
		timeString, ok := value[2].(string)
		if !ok {
			return emptySignalResponse, status, pnerr.NewResponseParsingError(fmt.Sprintf("Error unmarshalling response 2, %s %v", value[2], value), nil, nil)
		}
		timestamp, err := strconv.ParseInt(timeString, 10, 64)
		if err != nil {
			return emptySignalResponse, status, err
		}

		return &SignalResponse{
			Timestamp: timestamp,
		}, status, nil
	}

	return resp, status, nil
}
