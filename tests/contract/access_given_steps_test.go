package contract

import (
	"context"
	"fmt"
	"reflect"

	webpubsub "github.com/webpubsub/go/v7"
)

func theTTL(ctx context.Context, arg1 int) error {
	state := getAccessState(ctx)
	state.TTL = arg1
	return nil
}

func grantPermissionREAD(ctx context.Context) error {
	state := getAccessState(ctx)

	switch v := state.CurrentPermissions.(type) {
	case *webpubsub.ChannelPermissions:
		v.Read = true
	case *webpubsub.GroupPermissions:
		v.Read = true
	default:
		return fmt.Errorf("Not expected type %s", reflect.TypeOf(v).String())
	}
	return nil
}

func grantPermissionDELETE(ctx context.Context) error {
	state := getAccessState(ctx)

	switch v := state.CurrentPermissions.(type) {
	case *webpubsub.ChannelPermissions:
		v.Delete = true
	case *webpubsub.UUIDPermissions:
		v.Delete = true
	default:
		return fmt.Errorf("Not expected type %s", reflect.TypeOf(v).String())
	}
	return nil
}

func grantPermissionGET(ctx context.Context) error {
	state := getAccessState(ctx)

	switch v := state.CurrentPermissions.(type) {
	case *webpubsub.ChannelPermissions:
		v.Get = true
	case *webpubsub.UUIDPermissions:
		v.Get = true
	default:
		return fmt.Errorf("Not expected type %s", reflect.TypeOf(v).String())
	}
	return nil
}

func grantPermissionJOIN(ctx context.Context) error {
	state := getAccessState(ctx)

	switch v := state.CurrentPermissions.(type) {
	case *webpubsub.ChannelPermissions:
		v.Join = true
	default:
		return fmt.Errorf("Not expected type %s", reflect.TypeOf(v).String())
	}
	return nil
}

func grantPermissionMANAGE(ctx context.Context) error {
	state := getAccessState(ctx)

	switch v := state.CurrentPermissions.(type) {
	case *webpubsub.ChannelPermissions:
		v.Manage = true
	case *webpubsub.GroupPermissions:
		v.Manage = true
	default:
		return fmt.Errorf("Not expected type %s", reflect.TypeOf(v).String())
	}
	return nil
}

func grantPermissionUPDATE(ctx context.Context) error {
	state := getAccessState(ctx)

	switch v := state.CurrentPermissions.(type) {
	case *webpubsub.ChannelPermissions:
		v.Update = true
	case *webpubsub.UUIDPermissions:
		v.Update = true
	default:
		return fmt.Errorf("Not expected type %s", reflect.TypeOf(v).String())
	}
	return nil
}

func grantPermissionWRITE(ctx context.Context) error {
	state := getAccessState(ctx)

	switch v := state.CurrentPermissions.(type) {
	case *webpubsub.ChannelPermissions:
		v.Write = true
	default:
		return fmt.Errorf("Not expected type %s", reflect.TypeOf(v).String())
	}
	return nil
}

func theCHANNELResourceAccessPermissions(ctx context.Context, channel string) error {
	state := getAccessState(ctx)

	permissions := webpubsub.ChannelPermissions{}
	state.ChannelPermissions[channel] = &permissions
	state.CurrentPermissions = &permissions

	return nil
}

func theCHANNELPatternAccessPermissions(ctx context.Context, pattern string) error {
	state := getAccessState(ctx)

	permissions := webpubsub.ChannelPermissions{}
	state.ChannelPatternPermissions[pattern] = &permissions
	state.CurrentPermissions = &permissions

	return nil
}

func theCHANNEL_GROUPResourceAccessPermissions(ctx context.Context, id string) error {
	state := getAccessState(ctx)

	permissions := webpubsub.GroupPermissions{}
	state.ChannelGroupPermissions[id] = &permissions
	state.CurrentPermissions = &permissions

	return nil
}

func theCHANNEL_GROUPPatternAccessPermissions(ctx context.Context, pattern string) error {
	state := getAccessState(ctx)

	permissions := webpubsub.GroupPermissions{}
	state.ChannelGroupPatternPermissions[pattern] = &permissions
	state.CurrentPermissions = &permissions

	return nil
}

func theUUIDResourceAccessPermissions(ctx context.Context, id string) error {
	state := getAccessState(ctx)

	permissions := webpubsub.UUIDPermissions{}
	state.UUIDPermissions[id] = &permissions
	state.CurrentPermissions = &permissions

	return nil
}

func theUUIDPatternAccessPermissions(ctx context.Context, pattern string) error {
	state := getAccessState(ctx)

	permissions := webpubsub.UUIDPermissions{}
	state.UUIDPatternPermissions[pattern] = &permissions
	state.CurrentPermissions = &permissions

	return nil
}

func theAuthorizedUUID(ctx context.Context, uuid string) error {
	state := getAccessState(ctx)
	state.AuthorizedUUID = uuid

	return nil
}

const tokenWithEverything = "qEF2AkF0GmEI03xDdHRsGDxDcmVzpURjaGFuoWljaGFubmVsLTEY70NncnChb2NoYW5uZWxfZ3JvdXAtMQVDdXNyoENzcGOgRHV1aWShZnV1aWQtMRhoQ3BhdKVEY2hhbqFtXmNoYW5uZWwtXFMqJBjvQ2dycKF0XjpjaGFubmVsX2dyb3VwLVxTKiQFQ3VzcqBDc3BjoER1dWlkoWpedXVpZC1cUyokGGhEbWV0YaBEdXVpZHR0ZXN0LWF1dGhvcml6ZWQtdXVpZENzaWdYIPpU-vCe9rkpYs87YUrFNWkyNq8CVvmKwEjVinnDrJJc"

func iHaveAKnownTokenWithEverything(ctx context.Context) error {
	state := getAccessState(ctx)
	state.TokenString = tokenWithEverything
	return nil
}

func denyResourcePermissionGET(ctx context.Context) error {
	state := getAccessState(ctx)

	switch v := state.CurrentPermissions.(type) {
	case *webpubsub.ChannelPermissions:
		v.Get = false
	case *webpubsub.UUIDPermissions:
		v.Get = false
	default:
		return fmt.Errorf("Not expected type %s", reflect.TypeOf(v).String())
	}
	return nil
}

func aToken(ctx context.Context) error {
	state := getAccessState(ctx)
	state.TokenString = tokenWithEverything
	return nil
}

func aValidTokenWithPermissionsToPublishWithChannelChannel(ctx context.Context) error {
	state := getAccessState(ctx)
	state.TokenString = tokenWithEverything
	return nil
}

func anExpiredTokenWithPermissionsToPublishWithChannelChannel(ctx context.Context) error {
	state := getAccessState(ctx)
	state.TokenString = tokenWithEverything
	return nil
}

func theTokenString(ctx context.Context, token string) error {
	state := getAccessState(ctx)
	state.TokenString = token
	return nil
}
