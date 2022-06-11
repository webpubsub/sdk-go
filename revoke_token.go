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

const revokeTokenPath = "/v3/pam/%s/grant/%s"

var emptyWPSRevokeTokenResponse *WPSRevokeTokenResponse

type revokeTokenBuilder struct {
	opts *revokeTokenOpts
}

func newRevokeTokenBuilder(webpubsub *WebPubSub) *revokeTokenBuilder {
	builder := revokeTokenBuilder{
		opts: &revokeTokenOpts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newRevokeTokenBuilderWithContext(webpubsub *WebPubSub, context Context) *revokeTokenBuilder {
	builder := revokeTokenBuilder{
		opts: &revokeTokenOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

func (b *revokeTokenBuilder) Token(token string) *revokeTokenBuilder {
	b.opts.Token = token

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *revokeTokenBuilder) QueryParam(queryParam map[string]string) *revokeTokenBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the Grant request.
func (b *revokeTokenBuilder) Execute() (*WPSRevokeTokenResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyWPSRevokeTokenResponse, status, err
	}

	return newWPSRevokeTokenResponse(rawJSON, b.opts, status)
}

type revokeTokenOpts struct {
	webpubsub *WebPubSub
	ctx       Context

	QueryParam map[string]string
	Token      string
}

func (o *revokeTokenOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *revokeTokenOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *revokeTokenOpts) context() Context {
	return o.ctx
}

func (o *revokeTokenOpts) validate() error {
	if o.config().PublishKey == "" {
		return newValidationError(o, StrMissingPubKey)
	}

	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.config().SecretKey == "" {
		return newValidationError(o, StrMissingSecretKey)
	}

	if o.Token == "" {
		return newValidationError(o, StrMissingToken)
	}
	return nil
}

func (o *revokeTokenOpts) buildPath() (string, error) {
	return fmt.Sprintf(revokeTokenPath, o.webpubsub.Config.SubscribeKey, utils.URLEncode(o.Token)), nil
}

func (o *revokeTokenOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *revokeTokenOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *revokeTokenOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *revokeTokenOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *revokeTokenOpts) httpMethod() string {
	return "DELETE"
}

func (o *revokeTokenOpts) isAuthRequired() bool {
	return true
}

func (o *revokeTokenOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *revokeTokenOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *revokeTokenOpts) operationType() OperationType {
	return WPSAccessManagerRevokeToken
}

func (o *revokeTokenOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *revokeTokenOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// WPSRevokeTokenResponse is the struct returned when the Execute function of Grant Token is called.
type WPSRevokeTokenResponse struct {
	status int `json:"status"`
}

func newWPSRevokeTokenResponse(jsonBytes []byte, o *revokeTokenOpts, status StatusResponse) (*WPSRevokeTokenResponse, StatusResponse, error) {
	resp := &WPSRevokeTokenResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyWPSRevokeTokenResponse, status, e
	}

	return resp, status, nil
}
