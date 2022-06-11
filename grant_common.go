package webpubsub

import (
	"bytes"
	"encoding/base64"
	"strings"

	cbor "github.com/brianolson/cbor_go"
)

// WPSGrantBitMask is the type for perms BitMask
type WPSGrantBitMask int64

const (
	// WPSRead Read Perms
	WPSRead WPSGrantBitMask = 1
	// WPSWrite Write Perms
	WPSWrite = 2
	// WPSManage Manage Perms
	WPSManage = 4
	// WPSDelete Delete Perms
	WPSDelete = 8
	// WPSCreate Create Perms
	WPSCreate = 16
	// WPSGet Get Perms
	WPSGet = 32
	// WPSUpdate Update Perms
	WPSUpdate = 64
	// WPSJoin Join Perms
	WPSJoin = 128
)

// WPSGrantType grant types
type WPSGrantType int

const (
	// WPSReadEnabled Read Enabled. Applies to Subscribe, History, Presence, Objects
	WPSReadEnabled WPSGrantType = 1 + iota
	// WPSWriteEnabled Write Enabled. Applies to Publish, Objects
	WPSWriteEnabled
	// WPSManageEnabled Manage Enabled. Applies to Channel-Groups, Objects
	WPSManageEnabled
	// WPSDeleteEnabled Delete Enabled. Applies to History, Objects
	WPSDeleteEnabled
	// WPSCreateEnabled Create Enabled. Applies to Objects
	WPSCreateEnabled
	// WPSGetEnabled Get Enabled. Applies to Objects
	WPSGetEnabled
	// WPSUpdateEnabled Update Enabled. Applies to Objects
	WPSUpdateEnabled
	// WPSJoinEnabled Join Enabled. Applies to Objects
	WPSJoinEnabled
)

// WPSResourceType grant types
type WPSResourceType int

const (
	// WPSChannels for channels
	WPSChannels WPSResourceType = 1 + iota
	// WPSGroups for groups
	WPSGroups
	// WPSUsers for users
	WPSUUIDs
)

// ChannelPermissions contains all the acceptable perms for channels
type ChannelPermissions struct {
	Read   bool
	Write  bool
	Delete bool
	Get    bool
	Manage bool
	Update bool
	Join   bool
}

// GroupPermissions contains all the acceptable perms for groups
type GroupPermissions struct {
	Read   bool
	Manage bool
}

type UUIDPermissions struct {
	Get    bool
	Update bool
	Delete bool
}

// WPSPAMEntityData is the struct containing the access details of the channels.
type WPSPAMEntityData struct {
	Name          string
	AuthKeys      map[string]*WPSAccessManagerKeyData
	ReadEnabled   bool
	WriteEnabled  bool
	ManageEnabled bool
	DeleteEnabled bool
	GetEnabled    bool
	UpdateEnabled bool
	JoinEnabled   bool
	TTL           int
}

// WPSAccessManagerKeyData is the struct containing the access details of the channel groups.
type WPSAccessManagerKeyData struct {
	ReadEnabled   bool
	WriteEnabled  bool
	ManageEnabled bool
	DeleteEnabled bool
	GetEnabled    bool
	UpdateEnabled bool
	JoinEnabled   bool
	TTL           int
}

// GetPermissions decodes the CBORToken
func GetPermissions(token string) (WPSGrantTokenDecoded, error) {
	token = strings.Replace(token, "-", "+", -1)
	token = strings.Replace(token, "_", "/", -1)
	if i := len(token) % 4; i != 0 {
		token += strings.Repeat("=", 4-i)
	}

	var cborObject WPSGrantTokenDecoded
	value, decodeErr := base64.StdEncoding.DecodeString(token)
	if decodeErr != nil {
		return cborObject, decodeErr
	}

	c := cbor.NewDecoder(bytes.NewReader(value))
	err1 := c.Decode(&cborObject)
	if err1 != nil {
		return cborObject, err1
	}

	return cborObject, nil
}

type WPSToken struct {
	Version        int
	Timestamp      int64
	TTL            int
	AuthorizedUUID string
	Resources      WPSTokenResources
	Patterns       WPSTokenResources
	Meta           map[string]interface{}
}

type WPSTokenResources struct {
	Channels      map[string]ChannelPermissions
	ChannelGroups map[string]GroupPermissions
	UUIDs         map[string]UUIDPermissions
}

func ParseToken(token string) (*WPSToken, error) {
	permissions, err := GetPermissions(token)

	if err != nil {
		return nil, err
	}

	resources := grantResourcesToWPSTokenResources(permissions.Resources)
	patterns := grantResourcesToWPSTokenResources(permissions.Patterns)

	return &WPSToken{
		Version:        permissions.Version,
		Meta:           permissions.Meta,
		TTL:            permissions.TTL,
		Timestamp:      permissions.Timestamp,
		AuthorizedUUID: permissions.AuthorizedUUID,
		Resources:      resources,
		Patterns:       patterns,
	}, nil
}

func grantResourcesToWPSTokenResources(grantResources GrantResources) WPSTokenResources {
	tokenResources := WPSTokenResources{
		Channels:      make(map[string]ChannelPermissions),
		ChannelGroups: make(map[string]GroupPermissions),
		UUIDs:         make(map[string]UUIDPermissions),
	}
	for k, v := range grantResources.Channels {
		tokenResources.Channels[k] = parseGrantPerms(v, WPSChannels).(ChannelPermissions)
	}
	for k, v := range grantResources.Groups {
		tokenResources.ChannelGroups[k] = parseGrantPerms(v, WPSGroups).(GroupPermissions)
	}
	for k, v := range grantResources.UUIDs {
		tokenResources.UUIDs[k] = parseGrantPerms(v, WPSUUIDs).(UUIDPermissions)
	}
	return tokenResources
}

// ParseGrantResources parses the token for permissions and adds them along the other values to the GrantResourcesWithPermissions struct
func ParseGrantResources(res GrantResources, token string, timetoken int64, ttl int) *GrantResourcesWithPermissions {
	channels := make(map[string]ChannelPermissionsWithToken, len(res.Channels))

	for k, v := range res.Channels {
		channels[k] = ChannelPermissionsWithToken{
			Permissions:  parseGrantPerms(v, WPSChannels).(ChannelPermissions),
			BitMaskPerms: v,
			Token:        token,
			Timestamp:    timetoken,
			TTL:          ttl,
		}
	}

	groups := make(map[string]GroupPermissionsWithToken, len(res.Groups))
	for k, v := range res.Groups {
		groups[k] = GroupPermissionsWithToken{
			Permissions:  parseGrantPerms(v, WPSGroups).(GroupPermissions),
			BitMaskPerms: v,
			Token:        token,
			Timestamp:    timetoken,
			TTL:          ttl,
		}
	}

	g := GrantResourcesWithPermissions{
		Channels: channels,
		Groups:   groups,
	}
	return &g
}

func parseGrantPerms(i int64, resourceType WPSResourceType) interface{} {
	read := i&int64(WPSRead) != 0
	write := i&int64(WPSWrite) != 0
	manage := i&int64(WPSManage) != 0
	delete := i&int64(WPSDelete) != 0
	get := i&int64(WPSGet) != 0
	update := i&int64(WPSUpdate) != 0
	join := i&int64(WPSJoin) != 0

	switch resourceType {
	case WPSChannels:
		return ChannelPermissions{
			Read:   read,
			Write:  write,
			Delete: delete,
			Update: update,
			Get:    get,
			Join:   join,
			Manage: manage,
		}
	case WPSGroups:
		return GroupPermissions{
			Read:   read,
			Manage: manage,
		}
	default:
		return UUIDPermissions{
			Get:    get,
			Update: update,
			Delete: delete,
		}
	}
}

// ChannelPermissionsWithToken is used for channels resource type permissions
type ChannelPermissionsWithToken struct {
	Permissions  ChannelPermissions
	BitMaskPerms int64
	Token        string
	Timestamp    int64
	TTL          int
}

// GroupPermissionsWithToken is used for groups resource type permissions
type GroupPermissionsWithToken struct {
	Permissions  GroupPermissions
	BitMaskPerms int64
	Token        string
	Timestamp    int64
	TTL          int
}

// GrantResourcesWithPermissions is used as a common struct to store all resource type permissions
type GrantResourcesWithPermissions struct {
	Channels        map[string]ChannelPermissionsWithToken
	Groups          map[string]GroupPermissionsWithToken
	ChannelsPattern map[string]ChannelPermissionsWithToken
	GroupsPattern   map[string]GroupPermissionsWithToken
}

// PermissionsBody is the struct used to decode the server response
type PermissionsBody struct {
	Resources      GrantResources         `json:"resources"`
	Patterns       GrantResources         `json:"patterns"`
	Meta           map[string]interface{} `json:"meta"`
	AuthorizedUUID string                 `json:"uuid,omitempty"`
}

// GrantResources is the struct used to decode the server response
type GrantResources struct {
	Channels map[string]int64 `json:"channels" cbor:"chan"`
	Groups   map[string]int64 `json:"groups" cbor:"grp"`
	UUIDs    map[string]int64 `json:"uuids" cbor:"uuid"`
	Users    map[string]int64 `json:"users" cbor:"usr"`
	Spaces   map[string]int64 `json:"spaces" cbor:"spc"`
}

// WPSGrantTokenDecoded is the struct used to decode the server response
type WPSGrantTokenDecoded struct {
	Resources      GrantResources         `cbor:"res"`
	Patterns       GrantResources         `cbor:"pat"`
	Meta           map[string]interface{} `cbor:"meta"`
	Signature      []byte                 `cbor:"sig"`
	Version        int                    `cbor:"v"`
	Timestamp      int64                  `cbor:"t"`
	TTL            int                    `cbor:"ttl"`
	AuthorizedUUID string                 `cbor:"uuid"`
}
