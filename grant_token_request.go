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

const grantTokenPath = "/v3/pam/%s/grant"

var emptyWPSGrantTokenResponse *WPSGrantTokenResponse

type grantTokenBuilder struct {
	opts *grantTokenOpts
}

func newGrantTokenBuilder(webpubsub *WebPubSub) *grantTokenBuilder {
	builder := grantTokenBuilder{
		opts: &grantTokenOpts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newGrantTokenBuilderWithContext(webpubsub *WebPubSub, context Context) *grantTokenBuilder {
	builder := grantTokenBuilder{
		opts: &grantTokenOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

// TTL in minutes for which granted permissions are valid.
//
// Min: 1
// Max: 525600
// Default: 1440
//
// Setting value to 0 will apply the grant indefinitely (forever grant).
func (b *grantTokenBuilder) TTL(ttl int) *grantTokenBuilder {
	b.opts.TTL = ttl
	b.opts.setTTL = true

	return b
}

func (b *grantTokenBuilder) AuthorizedUUID(uuid string) *grantTokenBuilder {
	b.opts.AuthorizedUUID = uuid

	return b
}

//Channels sets the Channels for the Grant request.
func (b *grantTokenBuilder) Channels(channels map[string]ChannelPermissions) *grantTokenBuilder {
	b.opts.Channels = channels

	return b
}

// ChannelGroups sets the ChannelGroups for the Grant request.
func (b *grantTokenBuilder) ChannelGroups(groups map[string]GroupPermissions) *grantTokenBuilder {
	b.opts.ChannelGroups = groups

	return b
}

func (b *grantTokenBuilder) UUIDs(uuids map[string]UUIDPermissions) *grantTokenBuilder {
	b.opts.UUIDs = uuids

	return b
}

// Channels sets the Channels for the Grant request.
func (b *grantTokenBuilder) ChannelsPattern(channels map[string]ChannelPermissions) *grantTokenBuilder {
	b.opts.ChannelsPattern = channels

	return b
}

// ChannelGroups sets the ChannelGroups for the Grant request.
func (b *grantTokenBuilder) ChannelGroupsPattern(groups map[string]GroupPermissions) *grantTokenBuilder {
	b.opts.ChannelGroupsPattern = groups

	return b
}

func (b *grantTokenBuilder) UUIDsPattern(uuids map[string]UUIDPermissions) *grantTokenBuilder {
	b.opts.UUIDsPattern = uuids

	return b
}

// Meta sets the Meta for the Grant request.
func (b *grantTokenBuilder) Meta(meta map[string]interface{}) *grantTokenBuilder {
	b.opts.Meta = meta

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *grantTokenBuilder) QueryParam(queryParam map[string]string) *grantTokenBuilder {
	b.opts.QueryParam = queryParam

	return b
}

// Execute runs the Grant request.
func (b *grantTokenBuilder) Execute() (*WPSGrantTokenResponse, StatusResponse, error) {
	rawJSON, status, err := executeRequest(b.opts)
	if err != nil {
		return emptyWPSGrantTokenResponse, status, err
	}

	return newGrantTokenResponse(b, rawJSON, status)
}

type grantTokenOpts struct {
	webpubsub *WebPubSub
	ctx       Context

	AuthKeys             []string
	Channels             map[string]ChannelPermissions
	ChannelGroups        map[string]GroupPermissions
	UUIDs                map[string]UUIDPermissions
	ChannelsPattern      map[string]ChannelPermissions
	ChannelGroupsPattern map[string]GroupPermissions
	UUIDsPattern         map[string]UUIDPermissions
	QueryParam           map[string]string
	Meta                 map[string]interface{}
	AuthorizedUUID       string

	// Max: 525600
	// Min: 1
	// Default: 1440
	// Setting 0 will apply the grant indefinitely
	TTL int

	// nil hacks
	setTTL bool
}

func (o *grantTokenOpts) config() Config {
	return *o.webpubsub.Config
}

func (o *grantTokenOpts) client() *http.Client {
	return o.webpubsub.GetClient()
}

func (o *grantTokenOpts) context() Context {
	return o.ctx
}

func (o *grantTokenOpts) validate() error {
	if o.config().PublishKey == "" {
		return newValidationError(o, StrMissingPubKey)
	}

	if o.config().SubscribeKey == "" {
		return newValidationError(o, StrMissingSubKey)
	}

	if o.config().SecretKey == "" {
		return newValidationError(o, StrMissingSecretKey)
	}

	if o.TTL <= 0 {
		return newValidationError(o, StrInvalidTTL)
	}

	return nil
}

func (o *grantTokenOpts) buildPath() (string, error) {
	return fmt.Sprintf(grantTokenPath, o.webpubsub.Config.SubscribeKey), nil
}

type grantBody struct {
	TTL         int             `json:"ttl"`
	Permissions PermissionsBody `json:"permissions"`
}

func (o *grantTokenOpts) setBitmask(value bool, bitmask WPSGrantBitMask, bm int64) int64 {
	if value {
		bm |= int64(bitmask)
	}
	o.webpubsub.Config.Log.Println(fmt.Sprintf("bmVal: %t %d %d", value, bitmask, bm))
	return bm
}

func (o *grantTokenOpts) parseResourcePermissions(resource interface{}, resourceType WPSResourceType) map[string]int64 {
	bmVal := int64(0)
	switch resourceType {
	case WPSChannels:
		resourceWithPerms := resource.(map[string]ChannelPermissions)
		resourceWithPermsLen := len(resourceWithPerms)
		if resourceWithPermsLen > 0 {
			r := make(map[string]int64, resourceWithPermsLen)
			for k, v := range resourceWithPerms {
				bmVal = int64(0)
				bmVal = o.setBitmask(v.Read, WPSRead, bmVal)
				bmVal = o.setBitmask(v.Write, WPSWrite, bmVal)
				bmVal = o.setBitmask(v.Delete, WPSDelete, bmVal)
				bmVal = o.setBitmask(v.Join, WPSJoin, bmVal)
				bmVal = o.setBitmask(v.Update, WPSUpdate, bmVal)
				bmVal = o.setBitmask(v.Manage, WPSManage, bmVal)
				bmVal = o.setBitmask(v.Get, WPSGet, bmVal)
				o.webpubsub.Config.Log.Println("bmVal ChannelPermissions:", bmVal)
				r[k] = bmVal
			}
			return r
		}
		return make(map[string]int64)

	case WPSGroups:
		resourceWithPerms := resource.(map[string]GroupPermissions)
		resourceWithPermsLen := len(resourceWithPerms)
		if resourceWithPermsLen > 0 {
			r := make(map[string]int64, resourceWithPermsLen)
			for k, v := range resourceWithPerms {
				bmVal = int64(0)
				bmVal = o.setBitmask(v.Read, WPSRead, bmVal)
				bmVal = o.setBitmask(v.Manage, WPSManage, bmVal)
				o.webpubsub.Config.Log.Println("bmVal GroupPermissions:", bmVal)
				r[k] = bmVal
			}
			return r
		}
		return make(map[string]int64)

	case WPSUUIDs:
		resourceWithPerms := resource.(map[string]UUIDPermissions)
		resourceWithPermsLen := len(resourceWithPerms)
		if resourceWithPermsLen > 0 {
			r := make(map[string]int64, resourceWithPermsLen)
			for k, v := range resourceWithPerms {
				bmVal = int64(0)
				bmVal = o.setBitmask(v.Get, WPSGet, bmVal)
				bmVal = o.setBitmask(v.Update, WPSUpdate, bmVal)
				bmVal = o.setBitmask(v.Delete, WPSDelete, bmVal)
				o.webpubsub.Config.Log.Println("bmVal UUIDPermissions:", bmVal)
				r[k] = bmVal
			}
			return r
		}
		return make(map[string]int64)
	default:
		return make(map[string]int64)
	}

}

func (o *grantTokenOpts) buildQuery() (*url.Values, error) {
	q := defaultQuery(o.webpubsub.Config.UUID, o.webpubsub.telemetryManager)

	SetQueryParam(q, o.QueryParam)

	return q, nil
}

func (o *grantTokenOpts) jobQueue() chan *JobQItem {
	return o.webpubsub.jobQueue
}

func (o *grantTokenOpts) buildBody() ([]byte, error) {

	meta := o.Meta

	if meta == nil {
		meta = make(map[string]interface{})
	}

	permissions := PermissionsBody{
		Resources: GrantResources{
			Channels: o.parseResourcePermissions(o.Channels, WPSChannels),
			Groups:   o.parseResourcePermissions(o.ChannelGroups, WPSGroups),
			UUIDs:    o.parseResourcePermissions(o.UUIDs, WPSUUIDs),
			Users:    make(map[string]int64),
			Spaces:   make(map[string]int64),
		},
		Patterns: GrantResources{
			Channels: o.parseResourcePermissions(o.ChannelsPattern, WPSChannels),
			Groups:   o.parseResourcePermissions(o.ChannelGroupsPattern, WPSGroups),
			UUIDs:    o.parseResourcePermissions(o.UUIDsPattern, WPSUUIDs),
			Users:    make(map[string]int64),
			Spaces:   make(map[string]int64),
		},
		Meta:           meta,
		AuthorizedUUID: o.AuthorizedUUID,
	}

	o.webpubsub.Config.Log.Println("permissions: ", permissions)

	ttl := -1
	if o.setTTL {
		if o.TTL >= -1 {
			ttl = o.TTL
		}
	}

	b := &grantBody{
		TTL:         ttl,
		Permissions: permissions,
	}

	jsonEncBytes, errEnc := json.Marshal(b)

	if errEnc != nil {
		o.webpubsub.Config.Log.Printf("ERROR: Serialization error: %s\n", errEnc.Error())
		return []byte{}, errEnc
	}
	return jsonEncBytes, nil
}

func (o *grantTokenOpts) buildBodyMultipartFileUpload() (bytes.Buffer, *multipart.Writer, int64, error) {
	return bytes.Buffer{}, nil, 0, errors.New("Not required")
}

func (o *grantTokenOpts) httpMethod() string {
	return "POST"
}

func (o *grantTokenOpts) isAuthRequired() bool {
	return true
}

func (o *grantTokenOpts) requestTimeout() int {
	return o.webpubsub.Config.NonSubscribeRequestTimeout
}

func (o *grantTokenOpts) connectTimeout() int {
	return o.webpubsub.Config.ConnectTimeout
}

func (o *grantTokenOpts) operationType() OperationType {
	return WPSAccessManagerGrantToken
}

func (o *grantTokenOpts) telemetryManager() *TelemetryManager {
	return o.webpubsub.telemetryManager
}

func (o *grantTokenOpts) tokenManager() *TokenManager {
	return o.webpubsub.tokenManager
}

// WPSGrantTokenData is the struct used to decode the server response
type WPSGrantTokenData struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

// WPSGrantTokenResponse is the struct returned when the Execute function of Grant Token is called.
type WPSGrantTokenResponse struct {
	status  int               `json:"status"`
	Data    WPSGrantTokenData `json:"data"`
	service string            `json:"service"`
}

func newGrantTokenResponse(b *grantTokenBuilder, jsonBytes []byte, status StatusResponse) (*WPSGrantTokenResponse, StatusResponse, error) {
	resp := &WPSGrantTokenResponse{}

	err := json.Unmarshal(jsonBytes, &resp)
	if err != nil {
		e := pnerr.NewResponseParsingError("Error unmarshalling response",
			ioutil.NopCloser(bytes.NewBufferString(string(jsonBytes))), err)

		return emptyWPSGrantTokenResponse, status, e
	}

	b.opts.webpubsub.tokenManager.StoreToken(resp.Data.Token)

	return resp, status, nil
}
