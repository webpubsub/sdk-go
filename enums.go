package webpubsub

import (
	"fmt"
	"reflect"
)

// StatusCategory is used as an enum to catgorize the various status events
// in the APIs lifecycle
type StatusCategory int

// OperationType is used as an enum to catgorize the various operations
// in the APIs lifecycle
type OperationType int

// ReconnectionPolicy is used as an enum to catgorize the reconnection policies
type ReconnectionPolicy int

// WPSPushType is used as an enum to catgorize the available Push Types
type WPSPushType int

// WPSUUIDMetadataInclude is used as an enum to catgorize the available UUID include types
type WPSUUIDMetadataInclude int

// WPSChannelMetadataInclude is used as an enum to catgorize the available Channel include types
type WPSChannelMetadataInclude int

// WPSMembershipsInclude is used as an enum to catgorize the available Memberships include types
type WPSMembershipsInclude int

// WPSChannelMembersInclude is used as an enum to catgorize the available Members include types
type WPSChannelMembersInclude int

// WPSObjectsEvent is used as an enum to catgorize the available Object Events
type WPSObjectsEvent string

// WPSObjectsEventType is used as an enum to catgorize the available Object Event types
type WPSObjectsEventType string

// WPSMessageActionsEventType is used as an enum to catgorize the available Message Actions Event types
type WPSMessageActionsEventType string

// WPSPushEnvironment is used as an enum to catgorize the available Message Actions Event types
type WPSPushEnvironment string

const (
	//WPSPushEnvironmentDevelopment for development
	WPSPushEnvironmentDevelopment WPSPushEnvironment = "development"
	//WPSPushEnvironmentProduction for production
	WPSPushEnvironmentProduction = "production"
)

const (
	// WPSMessageActionsAdded is the enum when the event of type `added` occurs
	WPSMessageActionsAdded WPSMessageActionsEventType = "added"
	// WPSMessageActionsRemoved is the enum when the event of type `removed` occurs
	WPSMessageActionsRemoved = "removed"
)

const (
	// WPSObjectsMembershipEvent is the enum when the event of type `membership` occurs
	WPSObjectsMembershipEvent WPSObjectsEventType = "membership"
	// WPSObjectsChannelEvent is the enum when the event of type `channel` occurs
	WPSObjectsChannelEvent = "channel"
	// WPSObjectsUUIDEvent is the enum when the event of type `uuid` occurs
	WPSObjectsUUIDEvent = "uuid"
	// WPSObjectsNoneEvent is used for error handling
	WPSObjectsNoneEvent = "none"
)

const (
	// WPSObjectsEventRemove is the enum when the event `delete` occurs
	WPSObjectsEventRemove WPSObjectsEvent = "delete"
	// WPSObjectsEventSet is the enum when the event `set` occurs
	WPSObjectsEventSet = "set"
)

const (
	// WPSUUIDMetadataIncludeCustom is the enum equivalent to the value `custom` available UUID include types
	WPSUUIDMetadataIncludeCustom WPSUUIDMetadataInclude = 1 + iota
)

const (
	// WPSChannelMetadataIncludeCustom is the enum equivalent to the value `custom` available Channel include types
	WPSChannelMetadataIncludeCustom WPSChannelMetadataInclude = 1 + iota
)

func (s WPSUUIDMetadataInclude) String() string {
	return [...]string{"custom"}[s-1]
}

func (s WPSChannelMetadataInclude) String() string {
	return [...]string{"custom"}[s-1]
}

const (
	// WPSMembershipsIncludeCustom is the enum equivalent to the value `custom` available Memberships include types
	WPSMembershipsIncludeCustom WPSMembershipsInclude = 1 + iota
	// WPSMembershipsIncludeChannel is the enum equivalent to the value `channel` available Memberships include types
	WPSMembershipsIncludeChannel
	// WPSMembershipsIncludeChannelCustom is the enum equivalent to the value `channel.custom` available Memberships include types
	WPSMembershipsIncludeChannelCustom
)

func (s WPSMembershipsInclude) String() string {
	return [...]string{"custom", "channel", "channel.custom"}[s-1]
}

const (
	// WPSChannelMembersIncludeCustom is the enum equivalent to the value `custom` available Members include types
	WPSChannelMembersIncludeCustom WPSChannelMembersInclude = 1 + iota
	// WPSChannelMembersIncludeUUID is the enum equivalent to the value `uuid` available Members include types
	WPSChannelMembersIncludeUUID
	// WPSChannelMembersIncludeUUIDCustom is the enum equivalent to the value `uuid.custom` available Members include types
	WPSChannelMembersIncludeUUIDCustom
)

func (s WPSChannelMembersInclude) String() string {
	//return [...]string{"custom", "user", "user.custom", "uuid", "uuid.custom"}[s-1]
	return [...]string{"custom", "uuid", "uuid.custom"}[s-1]
}

// WPSMessageType is used as an enum to catgorize the Subscribe response.
type WPSMessageType int

const (
	// WPSNonePolicy is to be used when selecting the no Reconnection Policy
	// ReconnectionPolicy is set in the config.
	WPSNonePolicy ReconnectionPolicy = 1 + iota
	// WPSLinearPolicy is to be used when selecting the Linear Reconnection Policy
	// ReconnectionPolicy is set in the config.
	WPSLinearPolicy
	// WPSExponentialPolicycy is to be used when selecting the Exponential Reconnection Policy
	// ReconnectionPolicy is set in the config.
	WPSExponentialPolicycy
)

const (
	// WPSMessageTypeSignal is to identify Signal the Subscribe response
	WPSMessageTypeSignal WPSMessageType = 1 + iota
	// WPSMessageTypeObjects is to identify Objects the Subscribe response
	WPSMessageTypeObjects
	// WPSMessageTypeMessageActions is to identify Actions the Subscribe response
	WPSMessageTypeMessageActions
	// WPSMessageTypeFile is to identify Files the Subscribe response
	WPSMessageTypeFile
)

const (
	// WPSUnknownCategory as the StatusCategory means an unknown status category event occurred.
	WPSUnknownCategory StatusCategory = 1 + iota
	// WPSTimeoutCategory as the StatusCategory means the request timeout has reached.
	WPSTimeoutCategory
	// WPSConnectedCategory as the StatusCategory means the channel is subscribed to receive messages.
	WPSConnectedCategory
	// WPSDisconnectedCategory as the StatusCategory means a disconnection occurred due to network issues.
	WPSDisconnectedCategory
	// WPSCancelledCategory as the StatusCategory means the context was cancelled.
	WPSCancelledCategory
	// WPSLoopStopCategory as the StatusCategory means the subscribe loop was stopped.
	WPSLoopStopCategory
	// WPSAcknowledgmentCategory as the StatusCategory is the Acknowledgement of an operation (like Unsubscribe).
	WPSAcknowledgmentCategory
	// WPSBadRequestCategory as the StatusCategory means the request was malformed.
	WPSBadRequestCategory
	// WPSAccessDeniedCategory as the StatusCategory means that PAM is enabled and the channel is not granted R/W access.
	WPSAccessDeniedCategory
	// WPSNoStubMatchedCategory as the StatusCategory means an unknown status category event occurred.
	WPSNoStubMatchedCategory
	// WPSReconnectedCategory as the StatusCategory means that the network was reconnected (after a disconnection).
	// Applicable on for WPSLinearPolicy and WPSExponentialPolicycy.
	WPSReconnectedCategory
	// WPSReconnectionAttemptsExhausted as the StatusCategory means that the reconnection attempts
	// to reconnect to the network were exhausted. All channels would be unsubscribed at this point.
	// Applicable on for WPSLinearPolicy and WPSExponentialPolicycy.
	// Reconnection attempts are set in the config: MaximumReconnectionRetries.
	WPSReconnectionAttemptsExhausted
	// WPSRequestMessageCountExceededCategory is fired when the MessageQueueOverflowCount limit is exceeded by the number of messages received in a single subscribe request
	WPSRequestMessageCountExceededCategory
)

const (
	// WPSSubscribeOperation is the enum used for the Subcribe operation.
	WPSSubscribeOperation OperationType = 1 + iota
	// WPSUnsubscribeOperation is the enum used for the Unsubcribe operation.
	WPSUnsubscribeOperation
	// WPSPublishOperation is the enum used for the Publish operation.
	WPSPublishOperation
	// WPSFireOperation is the enum used for the Fire operation.
	WPSFireOperation
	// WPSHistoryOperation is the enum used for the History operation.
	WPSHistoryOperation
	// WPSFetchMessagesOperation is the enum used for the Fetch operation.
	WPSFetchMessagesOperation
	// WPSWhereNowOperation is the enum used for the Where Now operation.
	WPSWhereNowOperation
	// WPSHereNowOperation is the enum used for the Here Now operation.
	WPSHereNowOperation
	// WPSHeartBeatOperation is the enum used for the Heartbeat operation.
	WPSHeartBeatOperation
	// WPSSetStateOperation is the enum used for the Set State operation.
	WPSSetStateOperation
	// WPSGetStateOperation is the enum used for the Get State operation.
	WPSGetStateOperation
	// WPSAddChannelsToChannelGroupOperation is the enum used for the Add Channels to Channel Group operation.
	WPSAddChannelsToChannelGroupOperation
	// WPSRemoveChannelFromChannelGroupOperation is the enum used for the Remove Channels from Channel Group operation.
	WPSRemoveChannelFromChannelGroupOperation
	// WPSRemoveGroupOperation is the enum used for the Remove Channel Group operation.
	WPSRemoveGroupOperation
	// WPSChannelsForGroupOperation is the enum used for the List Channels of Channel Group operation.
	WPSChannelsForGroupOperation
	// WPSPushNotificationsEnabledChannelsOperation is the enum used for the List Channels with Push Notifications enabled operation.
	WPSPushNotificationsEnabledChannelsOperation
	// WPSAddPushNotificationsOnChannelsOperation is the enum used for the Add Channels to Push Notifications operation.
	WPSAddPushNotificationsOnChannelsOperation
	// WPSRemovePushNotificationsFromChannelsOperation is the enum used for the Remove Channels from Push Notifications operation.
	WPSRemovePushNotificationsFromChannelsOperation
	// WPSRemoveAllPushNotificationsOperation is the enum used for the Remove All Channels from Push Notifications operation.
	WPSRemoveAllPushNotificationsOperation
	// WPSTimeOperation is the enum used for the Time operation.
	WPSTimeOperation
	// WPSAccessManagerGrant is the enum used for the Access Manager Grant operation.
	WPSAccessManagerGrant
	// WPSAccessManagerRevoke is the enum used for the Access Manager Revoke operation.
	WPSAccessManagerRevoke
	// WPSDeleteMessagesOperation is the enum used for the Delete Messages from History operation.
	WPSDeleteMessagesOperation
	// WPSMessageCountsOperation is the enum used for History with messages operation.
	WPSMessageCountsOperation
	// WPSSignalOperation is the enum used for Signal opertaion.
	WPSSignalOperation
	// WPSCreateUserOperation is the enum used to create users in the Object API.
	// ENUM ORDER needs to be maintained for Objects AIP
	WPSCreateUserOperation
	// WPSGetUsersOperation is the enum used to get users in the Object API.
	WPSGetUsersOperation
	// WPSGetUserOperation is the enum used to get user in the Object API.
	WPSGetUserOperation
	// WPSUpdateUserOperation is the enum used to update users in the Object API.
	WPSUpdateUserOperation
	// WPSDeleteUserOperation is the enum used to delete users in the Object API.
	WPSDeleteUserOperation
	// WPSGetSpaceOperation is the enum used to get space in the Object API.
	WPSGetSpaceOperation
	// WPSGetSpacesOperation is the enum used to get spaces in the Object API.
	WPSGetSpacesOperation
	// WPSCreateSpaceOperation is the enum used to create space in the Object API.
	WPSCreateSpaceOperation
	// WPSDeleteSpaceOperation is the enum used to delete space in the Object API.
	WPSDeleteSpaceOperation
	// WPSUpdateSpaceOperation is the enum used to update space in the Object API.
	WPSUpdateSpaceOperation
	// WPSGetMembershipsOperation is the enum used to get memberships in the Object API.
	WPSGetMembershipsOperation
	// WPSGetChannelMembersOperation is the enum used to get members in the Object API.
	WPSGetChannelMembersOperation
	// WPSManageMembershipsOperation is the enum used to manage memberships in the Object API.
	WPSManageMembershipsOperation
	// WPSManageMembersOperation is the enum used to manage members in the Object API.
	// ENUM ORDER needs to be maintained for Objects API.
	WPSManageMembersOperation
	// WPSSetChannelMembersOperation is the enum used to Set Members in the Object API.
	WPSSetChannelMembersOperation
	// WPSSetMembershipsOperation is the enum used to Set Memberships in the Object API.
	WPSSetMembershipsOperation
	// WPSRemoveChannelMetadataOperation is the enum used to Remove Channel Metadata in the Object API.
	WPSRemoveChannelMetadataOperation
	// WPSRemoveUUIDMetadataOperation is the enum used to Remove UUID Metadata in the Object API.
	WPSRemoveUUIDMetadataOperation
	// WPSGetAllChannelMetadataOperation is the enum used to Get All Channel Metadata in the Object API.
	WPSGetAllChannelMetadataOperation
	// WPSGetAllUUIDMetadataOperation is the enum used to Get All UUID Metadata in the Object API.
	WPSGetAllUUIDMetadataOperation
	// WPSGetUUIDMetadataOperation is the enum used to Get UUID Metadata in the Object API.
	WPSGetUUIDMetadataOperation
	// WPSRemoveMembershipsOperation is the enum used to Remove Memberships in the Object API.
	WPSRemoveMembershipsOperation
	// WPSRemoveChannelMembersOperation is the enum used to Remove Members in the Object API.
	WPSRemoveChannelMembersOperation
	// WPSSetUUIDMetadataOperation is the enum used to Set UUID Metadata in the Object API.
	WPSSetUUIDMetadataOperation
	// WPSSetChannelMetadataOperation is the enum used to Set Channel Metadata in the Object API.
	WPSSetChannelMetadataOperation
	// WPSGetChannelMetadataOperation is the enum used to Get Channel Metadata in the Object API.
	WPSGetChannelMetadataOperation
	// WPSAccessManagerGrantToken is the enum used for Grant v3 requests.
	WPSAccessManagerGrantToken
	// WPSGetMessageActionsOperation is the enum used for Message Actions Get requests.
	WPSGetMessageActionsOperation
	// WPSHistoryWithActionsOperation is the enum used for History with Actions requests.
	WPSHistoryWithActionsOperation
	// WPSAddMessageActionsOperation is the enum used for Message Actions Add requests.
	WPSAddMessageActionsOperation
	// WPSRemoveMessageActionsOperation is the enum used for Message Actions Remove requests.
	WPSRemoveMessageActionsOperation
	// WPSDeleteFileOperation is the enum used for DeleteFile requests.
	WPSDeleteFileOperation
	// WPSDownloadFileOperation is the enum used for DownloadFile requests.
	WPSDownloadFileOperation
	// WPSGetFileURLOperation is the enum used for GetFileURL requests.
	WPSGetFileURLOperation
	// WPSListFilesOperation is the enum used for ListFiles requests.
	WPSListFilesOperation
	// WPSSendFileOperation is the enum used for SendFile requests.
	WPSSendFileOperation
	// WPSSendFileToS3Operation is the enum used for v requests.
	WPSSendFileToS3Operation
	// WPSPublishFileMessageOperation is the enum used for PublishFileMessage requests.
	WPSPublishFileMessageOperation
	// WPSAccessManagerRevokeToken is the enum used for Grant Token remove requests.
	WPSAccessManagerRevokeToken
)

const (
	// WPSPushTypeNone is used as an enum to for selecting `none` as the WPSPushType
	WPSPushTypeNone WPSPushType = 1 + iota
	// WPSPushTypeGCM is used as an enum to for selecting `GCM` as the WPSPushType
	WPSPushTypeGCM
	// WPSPushTypeAPNS is used as an enum to for selecting `APNS` as the WPSPushType
	WPSPushTypeAPNS
	// WPSPushTypeMPNS is used as an enum to for selecting `MPNS` as the WPSPushType
	WPSPushTypeMPNS
	// WPSPushTypeAPNS2 is used as an enum to for selecting `APNS2` as the WPSPushType
	WPSPushTypeAPNS2
)

func (p WPSPushType) String() string {
	switch p {
	case WPSPushTypeAPNS:
		return "apns"

	case WPSPushTypeGCM:
		return "gcm"

	case WPSPushTypeMPNS:
		return "mpns"

	case WPSPushTypeAPNS2:
		return "apns2"

	default:
		return "none"

	}
}

var operations = [...]string{
	"Subscribe",
	"Unsubscribe",
	"Publish",
	"History",
	"Fetch Messages",
	"Where Now",
	"Here Now",
	"Heartbeat",
	"Set State",
	"Get State",
	"Add Channel To Channel Group",
	"Remove Channel From Channel Group",
	"Remove Channel Group",
	"List Channels In Channel Group",
	"List Push Enabled Channels",
	"Add Push From Channel",
	"Remove Push From Channel",
	"Remove All Push Notifications",
	"Time",
	"Grant",
	"Revoke",
	"Delete messages",
	"Signal",
	"Create User",
	"Get Users",
	"Fetch User",
	"Update User",
	"Delete User",
	"Get Space",
	"Get Spaces",
	"Create Space",
	"Delete Space",
	"Update Space",
	"GetMemberships",
	"GetChannelMembers",
	"ManageMemberships",
	"ManageMembers",
	"SetChannelMembers",
	"SetMemberships",
	"RemoveChannelMetadata",
	"RemoveUUIDMetadata",
	"GetAllChannelMetadata",
	"GetAllUUIDMetadata",
	"GetUUIDMetadata",
	"RemoveMemberships",
	"RemoveChannelMembers",
	"SetUUIDMetadata",
	"SetChannelMetadata",
	"GetChannelMetadata",
	"Grant Token",
	"GetMessageActions",
	"HistoryWithActions",
	"AddMessageActions",
	"RemoveMessageActions",
	"DeleteFile",
	"DownloadFile",
	"GetFileURL",
	"ListFiles",
	"SendFile",
	"SendFileToS3",
	"PublishFile",
}

func (c StatusCategory) String() string {
	switch c {
	case WPSUnknownCategory:
		return "Unknown"

	case WPSTimeoutCategory:
		return "Timeout"

	case WPSConnectedCategory:
		return "Connected"

	case WPSDisconnectedCategory:
		return "Disconnected"

	case WPSCancelledCategory:
		return "Cancelled"

	case WPSLoopStopCategory:
		return "Loop Stop"

	case WPSAcknowledgmentCategory:
		return "Acknowledgment"

	case WPSBadRequestCategory:
		return "Bad Request"

	case WPSAccessDeniedCategory:
		return "Access Denied"

	case WPSReconnectedCategory:
		return "Reconnected"

	case WPSReconnectionAttemptsExhausted:
		return "Reconnection Attempts Exhausted"

	case WPSNoStubMatchedCategory:
		return "No Stub Matched"

	default:
		return "No Stub Matched"

	}
}

func (t OperationType) String() string {
	switch t {
	case WPSSubscribeOperation:
		return "Subscribe"

	case WPSUnsubscribeOperation:
		return "Unsubscribe"

	case WPSPublishOperation:
		return "Publish"

	case WPSFireOperation:
		return "Fire"

	case WPSHistoryOperation:
		return "History"

	case WPSFetchMessagesOperation:
		return "Fetch Messages"

	case WPSWhereNowOperation:
		return "Where Now"

	case WPSHereNowOperation:
		return "Here Now"

	case WPSHeartBeatOperation:
		return "Heartbeat"

	case WPSSetStateOperation:
		return "Set State"

	case WPSGetStateOperation:
		return "Get State"

	case WPSAddChannelsToChannelGroupOperation:
		return "Add Channel To Channel Group"

	case WPSRemoveChannelFromChannelGroupOperation:
		return "Remove Channel From Channel Group"

	case WPSRemoveGroupOperation:
		return "Remove Channel Group"

	case WPSChannelsForGroupOperation:
		return "List Channels In Channel Group"

	case WPSPushNotificationsEnabledChannelsOperation:
		return "List Push Enabled Channels"

	case WPSAddPushNotificationsOnChannelsOperation:
		return "Add Push From Channel"

	case WPSRemovePushNotificationsFromChannelsOperation:
		return "Remove Push From Channel"

	case WPSRemoveAllPushNotificationsOperation:
		return "Remove All Push Notifications"

	case WPSTimeOperation:
		return "Time"

	case WPSAccessManagerGrant:
		return "Grant"

	case WPSAccessManagerRevoke:
		return "Revoke"

	case WPSDeleteMessagesOperation:
		return "Delete messages"

	case WPSSignalOperation:
		return "Signal"

	case WPSCreateUserOperation:
		return "Create User"
	case WPSGetUsersOperation:
		return "Get Users"
	case WPSGetUserOperation:
		return "Fetch Users"
	case WPSUpdateUserOperation:
		return "Update User"
	case WPSDeleteUserOperation:
		return "Delete User"
	case WPSGetSpaceOperation:
		return "Get Space"
	case WPSGetSpacesOperation:
		return "Get Spaces"
	case WPSCreateSpaceOperation:
		return "Create Space"
	case WPSDeleteSpaceOperation:
		return "Delete Space"
	case WPSUpdateSpaceOperation:
		return "Update Space"
	case WPSGetMembershipsOperation:
		return "Get Memberships V2"
	case WPSGetChannelMembersOperation:
		return "Get Members V2"
	case WPSManageMembershipsOperation:
		return "Manage Memberships V2"
	case WPSManageMembersOperation:
		return "Manage Members V2"
	case WPSSetChannelMembersOperation:
		return "Set Members V2"
	case WPSSetMembershipsOperation:
		return "Set Memberships V2"
	case WPSRemoveChannelMetadataOperation:
		return "Remove Channel Metadata V2"
	case WPSRemoveUUIDMetadataOperation:
		return "Remove Metadata V2"
	case WPSGetAllChannelMetadataOperation:
		return "Get All Channel Metadata V2"
	case WPSGetAllUUIDMetadataOperation:
		return "Get All UUID Metadata V2"
	case WPSGetUUIDMetadataOperation:
		return "Get UUID Metadata V2"
	case WPSRemoveMembershipsOperation:
		return "Remove Memberships V2"
	case WPSRemoveChannelMembersOperation:
		return "Remove Members V2"
	case WPSSetUUIDMetadataOperation:
		return "Set UUID Metadata V2"
	case WPSSetChannelMetadataOperation:
		return "Set Channel Metadata V2"
	case WPSGetChannelMetadataOperation:
		return "Get Channel Metadata V2"
	case WPSAccessManagerGrantToken:
		return "Grant Token"
	case WPSGetMessageActionsOperation:
		return "Get Message Actions"
	case WPSHistoryWithActionsOperation:
		return "History With Actions"
	case WPSAddMessageActionsOperation:
		return "Add Message Actions"
	case WPSRemoveMessageActionsOperation:
		return "Remove Message Actions"
	case WPSDeleteFileOperation:
		return "Delete File"
	case WPSDownloadFileOperation:
		return "Download File"
	case WPSGetFileURLOperation:
		return "Get File URL"
	case WPSListFilesOperation:
		return "List Files"
	case WPSSendFileOperation:
		return "Send File"
	case WPSSendFileToS3Operation:
		return "Send File To S3"
	case WPSPublishFileMessageOperation:
		return "Publish File"
	default:
		return "No Category Matched"
	}
}

// EnumArrayToStringArray converts a string enum to an array
func EnumArrayToStringArray(include interface{}) []string {
	s := []string{}
	switch fmt.Sprintf("%s", reflect.TypeOf(include)) {
	case "[]webpubsub.WPSChannelMembersInclude":
		for _, v := range include.([]WPSChannelMembersInclude) {
			s = append(s, fmt.Sprintf("%s", v))
		}
	case "[]webpubsub.WPSMembershipsInclude":
		for _, v := range include.([]WPSMembershipsInclude) {
			s = append(s, fmt.Sprintf("%s", v))
		}
	case "[]webpubsub.WPSUUIDMetadataInclude":
		for _, v := range include.([]WPSUUIDMetadataInclude) {
			s = append(s, fmt.Sprintf("%s", v))
		}
	case "[]webpubsub.WPSChannelMetadataInclude":
		for _, v := range include.([]WPSChannelMetadataInclude) {
			s = append(s, fmt.Sprintf("%s", v))
		}
	}
	return s
}
