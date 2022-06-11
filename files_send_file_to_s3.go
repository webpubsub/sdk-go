package webpubsub

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/webpubsub/sdk-go/v7/pnerr"
	"github.com/webpubsub/sdk-go/v7/utils"
)

var emptySendFileToS3Response *WPSSendFileToS3Response

type sendFileToS3Builder struct {
	opts *sendFileToS3Opts
}

func newSendFileToS3Builder(webpubsub *WebPubSub) *sendFileToS3Builder {
	builder := sendFileToS3Builder{
		opts: &sendFileToS3Opts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newSendFileToS3BuilderWithContext(webpubsub *WebPubSub,
	context Context) *sendFileToS3Builder {
	builder := sendFileToS3Builder{
		opts: &sendFileToS3Opts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

func (b *sendFileToS3Builder) CipherKey(cipherKey string) *sendFileToS3Builder {
	b.opts.CipherKey = cipherKey

	return b
}

func (b *sendFileToS3Builder) FileUploadRequestData(fileUploadRequestData WPSFileUploadRequest) *sendFileToS3Builder {
	b.opts.FileUploadRequestData = fileUploadRequestData

	return b
}

func (b *sendFileToS3Builder) File(f *os.File) *sendFileToS3Builder {
	b.opts.File = f

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *sendFileToS3Builder) QueryParam(queryParam map[string]string) *sendFileToS3Builder {
	b.opts.QueryParam = queryParam

	return b
}

// Transport sets the Transport for the sendFileToS3 request.
func (b *sendFileToS3Builder) Transport(tr http.RoundTripper) *sendFileToS3Builder {
	b.opts.Transport = tr
	return b
}

// Execute runs the sendFileToS3 request.
func (b *sendFileToS3Builder) Execute() (*WPSSendFileToS3Response, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptySendFileToS3Response, status, err
	}

	return newWPSSendFileToS3Response(rawJSON, b.opts, status)
}

type sendFileToS3Opts struct {
	webpubsub *WebPubSub

	File                  *os.File
	FileUploadRequestData WPSFileUploadRequest
	QueryParam            map[string]string
	CipherKey             string
	Transport             http.RoundTripper

	ctx Context
}

func (o *sendFileToS3Opts) config() Config {
	return *o.webpubsub.Config
}

func (o *sendFileToS3Opts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *sendFileToS3Opts) context() Context {
	return o.ctx
}

func (o *sendFileToS3Opts) validate() error {
	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	return nil
}

func (o *sendFileToS3Opts) buildPath() (string, error) {
	return o.FileUploadRequestData.URL, nil
}

func (o *sendFileToS3Opts) buildQuery() (*url.Values, error) {
	return &url.Values{}, nil
}

func (o *sendFileToS3Opts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *sendFileToS3Opts) buildBody() ([]byte, error) {
	return []byte{}, nil
}

func (o *sendFileToS3Opts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {

	fileInfo, _ := o.File.Stat()
	s := fileInfo.Size()
	buffer := make([]byte, 512)
	_, err := o.File.Read(buffer)
	if err != nil {
		return bytes.Buffer{}, nil, s, err
	}
	o.File.Seek(0, 0)
	contentType := http.DetectContentType(buffer)

	var fileBody bytes.Buffer
	writer := multipart.NewWriter(&fileBody)

	for _, v := range o.FileUploadRequestData.FormFields {
		o.webpubsub.Config.Log.Printf("FormFields: Key: %s Value: %s\n", v.Key, v.Value)
		if v.Key == "Content-Type" {
			v.Value = contentType
		}
		_ = writer.WriteField(v.Key, v.Value)
	}

	filePart, errFilePart := writer.CreateFormFile("file", fileInfo.Name())

	if errFilePart != nil {
		o.webpubsub.Config.Log.Printf("ERROR: writer CreateFormFile: %s\n", errFilePart.Error())
		return bytes.Buffer{}, writer, s, errFilePart
	}

	if o.CipherKey != "" {
		utils.EncryptFile(o.CipherKey, []byte{}, filePart, o.File)
	} else if o.webpubsub.Config.CipherKey != "" {
		utils.EncryptFile(o.webpubsub.Config.CipherKey, []byte{}, filePart, o.File)
	} else {
		_, errIOCopy := io.Copy(filePart, o.File)

		if errIOCopy != nil {
			o.webpubsub.Config.Log.Printf("ERROR: io Copy error: %s\n", errIOCopy.Error())
			return bytes.Buffer{}, writer, s, errIOCopy
		}
	}

	errWriterClose := writer.Close()
	if errWriterClose != nil {
		o.webpubsub.Config.Log.Printf("ERROR: Writer close: %s\n", errWriterClose.Error())
		return bytes.Buffer{}, writer, s, errWriterClose
	}

	return fileBody, writer, s, nil

}

func (o *sendFileToS3Opts) httpMethod() string {
	return "POSTFORM"
}

func (o *sendFileToS3Opts) isAuthRequired() bool {
	return true
}

func (o *sendFileToS3Opts) requestTimeout() int {
	return o.webpubsub.Config.FileUploadRequestTimeout
}

func (o *sendFileToS3Opts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *sendFileToS3Opts) operationType() OperationType {
	return WPSSendFileToS3Operation
}

func (o *sendFileToS3Opts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *sendFileToS3Opts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// WPSSendFileToS3Response is the File Upload API Response for Get Spaces
type WPSSendFileToS3Response struct {
}

func newWPSSendFileToS3Response(jsonBytes []byte, o *sendFileToS3Opts,
	status StatusResponse) (*WPSSendFileToS3Response, StatusResponse, error) {

	resp := &WPSSendFileToS3Response{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptySendFileToS3Response, status, e
	}
	o.webpubsub.Config.Log.Printf("newWPSSendFileToS3Response status.StatusCode==> %d", status.StatusCode)

	return resp, status, nil
}
