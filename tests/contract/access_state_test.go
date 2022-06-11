package contract

import webpubsub "github.com/webpubsub/go/v7"

type accessStateKey struct{}

type accessState struct {
	CurrentPermissions             interface{}
	ChannelPermissions             map[string]*webpubsub.ChannelPermissions
	ChannelPatternPermissions      map[string]*webpubsub.ChannelPermissions
	ChannelGroupPermissions        map[string]*webpubsub.GroupPermissions
	ChannelGroupPatternPermissions map[string]*webpubsub.GroupPermissions
	UUIDPermissions                map[string]*webpubsub.UUIDPermissions
	UUIDPatternPermissions         map[string]*webpubsub.UUIDPermissions
	TTL                            int
	TokenString                    string
	AuthorizedUUID                 string
	GrantTokenResult               webpubsub.WPSGrantTokenResponse
	ParsedToken                    *webpubsub.WPSToken
	ResourcePermissions            interface{}
	RevokeTokenResult              webpubsub.WPSRevokeTokenResponse
}

func newAccessState(pn *webpubsub.WebPubSub) *accessState {
	return &accessState{
		TTL:                            0,
		ChannelPermissions:             make(map[string]*webpubsub.ChannelPermissions),
		ChannelPatternPermissions:      make(map[string]*webpubsub.ChannelPermissions),
		ChannelGroupPermissions:        make(map[string]*webpubsub.GroupPermissions),
		ChannelGroupPatternPermissions: make(map[string]*webpubsub.GroupPermissions),
		UUIDPermissions:                make(map[string]*webpubsub.UUIDPermissions),
		UUIDPatternPermissions:         make(map[string]*webpubsub.UUIDPermissions),
		CurrentPermissions:             nil}
}
