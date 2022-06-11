package webpubsub

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"

	"github.com/webpubsub/sdk-go/v7/utils"
)

var emptyDownloadFileResponse *WPSDownloadFileResponse

const downloadFilePath = "/v1/files/%s/channels/%s/files/%s/%s"

const downloadFileLimit = 100

type downloadFileBuilder struct {
	opts *downloadFileOpts
}

func newDownloadFileBuilder(webpubsub *WebPubSub) *downloadFileBuilder {
	builder := downloadFileBuilder{
		opts: &downloadFileOpts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newDownloadFileBuilderWithContext(webpubsub *WebPubSub,
	context Context) *downloadFileBuilder {
	builder := downloadFileBuilder{
		opts: &downloadFileOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

func (b *downloadFileBuilder) Channel(channel string) *downloadFileBuilder {
	b.opts.Channel = channel

	return b
}

func (b *downloadFileBuilder) CipherKey(cipherKey string) *downloadFileBuilder {
	b.opts.CipherKey = cipherKey

	return b
}

func (b *downloadFileBuilder) ID(id string) *downloadFileBuilder {
	b.opts.ID = id

	return b
}

func (b *downloadFileBuilder) Name(name string) *downloadFileBuilder {
	b.opts.Name = name

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *downloadFileBuilder) QueryParam(queryParam map[string]string) *downloadFileBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the downloadFile request.
func (b *downloadFileBuilder) Transport(tr http.RoundTripper) *downloadFileBuilder {
	b.opts.Transport = tr
	return b
}

func (b *downloadFileBuilder) Execute() (*WPSDownloadFileResponse, StatusResponse, error) {
	u, _ := buildURL(b.opts)
	stat := StatusResponse{
		AffectedChannels: []string{b.opts.Channel},
		AuthKey:          b.opts.config().AuthKey,
		Category:         WPSUnknownCategory,
		Operation:        WPSGetFileURLOperation,
		StatusCode:       200,
		TLSEnabled:       b.opts.config().Secure,
		Origin:           b.opts.config().Origin,
		UUID:             b.opts.config().UUID,
	}
	b.opts.webpubsub.Config.Log.Printf("u.RequestURI(): %s", u.RequestURI())
	resp, err := b.opts.client().Get(u.RequestURI())
	if err != nil {
		b.opts.webpubsub.Config.Log.Printf("err %s", err)
		return nil, stat, err
	}
	if resp.StatusCode != 200 {
		stat.StatusCode = resp.StatusCode
		return nil, stat, err
	}
	contentLenEnc, err := strconv.ParseInt(string(resp.Header.Get("Content-Length")), 10, 64)
	if err != nil {
		b.opts.webpubsub.Config.Log.Printf("err in parsing content length %s", err)
		return nil, stat, err
	}

	var respDL *WPSDownloadFileResponse
	if b.opts.CipherKey != "" {
		r, w := io.Pipe()
		utils.DecryptFile(b.opts.CipherKey, contentLenEnc, resp.Body, w)
		respDL = &WPSDownloadFileResponse{
			File: r,
		}

	} else if b.opts.webpubsub.Config.CipherKey != "" {
		r, w := io.Pipe()
		utils.DecryptFile(b.opts.webpubsub.Config.CipherKey, contentLenEnc, resp.Body, w)
		respDL = &WPSDownloadFileResponse{
			File: r,
		}

	} else {
		respDL = &WPSDownloadFileResponse{
			File: resp.Body,
		}
	}
	return respDL, stat, nil
}

type downloadFileOpts struct {
	webpubsub *WebPubSub

	Channel    string
	CipherKey  string
	ID         string
	Name       string
	QueryParam map[string]string

	Transport http.RoundTripper

	ctx Context
}

func (o *downloadFileOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *downloadFileOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *downloadFileOpts) context() Context {
	return o.ctx
}

func (o *downloadFileOpts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}
	if o.Channel == "" {
		return newValidationError(o, StrMissingChannel)
	}

	if o.Name == "" {
		return newValidationError(o, StrMissingFileName)
	}

	if o.ID == "" {
		return newValidationError(o, StrMissingFileID)
	}

	return nil
}

func (o *downloadFileOpts) buildPath() (string, error) {
	return fmt.Sprintf(downloadFilePath,
		o.webpubsub.Config.SubscribeKey, o.Channel, o.ID, o.Name), nil
}

func (o *downloadFileOpts) buildQuery() (*url.Values, error) {

	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *downloadFileOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *downloadFileOpts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *downloadFileOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *downloadFileOpts) httpMethod() string {
	return "GET"
}

func (o *downloadFileOpts) isAuthRequired() bool {
	return true
}

func (o *downloadFileOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *downloadFileOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *downloadFileOpts) operationType() OperationType {
	return WPSDownloadFileOperation
}

func (o *downloadFileOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *downloadFileOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// WPSDownloadFileResponse is the File Upload API Response for Get Spaces
type WPSDownloadFileResponse struct {
	status int       `json:"status"`
	File   io.Reader `json:"data"`
}

func newWPSDownloadFileResponse(jsonBytes []byte, o *downloadFileOpts,
	status StatusResponse) (*WPSDownloadFileResponse, StatusResponse, error) {

	resp := &WPSDownloadFileResponse{}

	return resp, status, nil
}
