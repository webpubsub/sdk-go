package webpubsub

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"sync"

	"github.com/webpubsub/sdk-go/v7/utils"
)

// Default constants
const (
	// Version :the version of the SDK
	Version = "7.0.3"
	// MaxSequence for publish messages
	MaxSequence = 65535
)

const (
	// StrMissingPubKey shows Missing Publish Key message
	StrMissingPubKey = "Missing Publish Key"
	// StrMissingSubKey shows Missing Subscribe Key message
	StrMissingSubKey = "Missing Subscribe Key"
	// StrMissingChannel shows Channel message
	StrMissingChannel = "Missing Channel"
	// StrMissingChannelGroup shows Channel Group message
	StrMissingChannelGroup = "Missing Channel Group"
	// StrMissingMessage shows Missing Message message
	StrMissingMessage = "Missing Message"
	// StrMissingSecretKey shows Missing Secret Key message
	StrMissingSecretKey = "Missing Secret Key"
	// StrMissingUUID shows Missing UUID message
	StrMissingUUID = "Missing UUID"
	// StrMissingDeviceID shows Missing Device ID message
	StrMissingDeviceID = "Missing Device ID"
	// StrMissingPushType shows Missing Push Type message
	StrMissingPushType = "Missing Push Type"
	// StrMissingPushTopic shows Missing Push Topic message
	StrMissingPushTopic = "Missing Push Topic"
	// StrChannelsTimetoken shows Missing Channels Timetoken message
	StrChannelsTimetoken = "Missing Channels Timetoken"
	// StrChannelsTimetokenLength shows Length of Channels Timetoken message
	StrChannelsTimetokenLength = "Length of Channels Timetoken and Channels do not match"
	// StrInvalidTTL shows Invalid TTL message
	StrInvalidTTL = "Invalid TTL"
	// StrMissingPushTitle shows `Push title missing` message
	StrMissingPushTitle = "Push title missing"
	// StrMissingFileID shows `Missing File ID` message
	StrMissingFileID = "Missing File ID"
	// StrMissingFileName shows `Missing File Name` message
	StrMissingFileName = "Missing File Name"
	// StrMissingToken shows `Missing PAMv3 token` message
	StrMissingToken = "Missing PAMv3 token"
)

// WebPubSub No server connection will be established when you create a new WebPubSub object.
// To establish a new connection use Subscribe() function of WebPubSub type.
type WebPubSub struct {
	sync.RWMutex

	Config               *Config
	nextPublishSequence  int
	publishSequenceMutex sync.RWMutex
	subscriptionManager  *SubscriptionManager
	telemetryManager     *TelemetryManager
	heartbeatManager     *HeartbeatManager
	client               *http.Client
	subscribeClient      *http.Client
	requestWorkers       *RequestWorkers
	jobQueue             chan *JobQItem
	ctx                  Context
	cancel               func()
	tokenManager         *TokenManager
}

// Publish is used to send a message to all subscribers of a channel.
func (pn *WebPubSub) Publish() *publishBuilder {
	return newPublishBuilder(pn)
}

// PublishWithContext function is used to send a message to all subscribers of a channel.
func (pn *WebPubSub) PublishWithContext(ctx Context) *publishBuilder {
	return newPublishBuilderWithContext(pn, ctx)
}

// Fire endpoint allows the client to send a message to WebPubSub Functions Event Handlers. These messages will go directly to any Event Handlers registered on the channel that you fire to and will trigger their execution.
func (pn *WebPubSub) Fire() *fireBuilder {
	return newFireBuilder(pn)
}

// FireWithContext endpoint allows the client to send a message to WebPubSub Functions Event Handlers. These messages will go directly to any Event Handlers registered on the channel that you fire to and will trigger their execution.
func (pn *WebPubSub) FireWithContext(ctx Context) *fireBuilder {
	return newFireBuilderWithContext(pn, ctx)
}

// Subscribe causes the client to create an open TCP socket to the WebPubSub Real-Time Network and begin listening for messages on a specified channel.
func (pn *WebPubSub) Subscribe() *subscribeBuilder {
	return newSubscribeBuilder(pn)
}

// History fetches historical messages of a channel.
func (pn *WebPubSub) History() *historyBuilder {
	return newHistoryBuilder(pn)
}

// HistoryWithContext fetches historical messages of a channel.
func (pn *WebPubSub) HistoryWithContext(ctx Context) *historyBuilder {
	return newHistoryBuilderWithContext(pn, ctx)
}

// Fetch fetches historical messages from multiple channels.
func (pn *WebPubSub) Fetch() *fetchBuilder {
	return newFetchBuilder(pn)
}

// FetchWithContext fetches historical messages from multiple channels.
func (pn *WebPubSub) FetchWithContext(ctx Context) *fetchBuilder {
	return newFetchBuilderWithContext(pn, ctx)
}

// MessageCounts Returns the number of messages published on one or more channels since a given time.
func (pn *WebPubSub) MessageCounts() *messageCountsBuilder {
	return newMessageCountsBuilder(pn)
}

// MessageCountsWithContext Returns the number of messages published on one or more channels since a given time.
func (pn *WebPubSub) MessageCountsWithContext(ctx Context) *messageCountsBuilder {
	return newMessageCountsBuilderWithContext(pn, ctx)
}

// GetAllUUIDMetadata Returns a paginated list of UUID Metadata objects, optionally including the custom data object for each.
func (pn *WebPubSub) GetAllUUIDMetadata() *getAllUUIDMetadataBuilder {
	return newGetAllUUIDMetadataBuilder(pn)
}

// GetAllUUIDMetadataWithContext Returns a paginated list of UUID Metadata objects, optionally including the custom data object for each.
func (pn *WebPubSub) GetAllUUIDMetadataWithContext(ctx Context) *getAllUUIDMetadataBuilder {
	return newGetAllUUIDMetadataBuilderWithContext(pn, ctx)
}

// GetUUIDMetadata Returns metadata for the specified UUID, optionally including the custom data object for each.
func (pn *WebPubSub) GetUUIDMetadata() *getUUIDMetadataBuilder {
	return newGetUUIDMetadataBuilder(pn)
}

// GetUUIDMetadataWithContext Returns metadata for the specified UUID, optionally including the custom data object for each.
func (pn *WebPubSub) GetUUIDMetadataWithContext(ctx Context) *getUUIDMetadataBuilder {
	return newGetUUIDMetadataBuilderWithContext(pn, ctx)
}

// SetUUIDMetadata Set metadata for a UUID in the database, optionally including the custom data object for each.
func (pn *WebPubSub) SetUUIDMetadata() *setUUIDMetadataBuilder {
	return newSetUUIDMetadataBuilder(pn)
}

// SetUUIDMetadataWithContext Set metadata for a UUID in the database, optionally including the custom data object for each.
func (pn *WebPubSub) SetUUIDMetadataWithContext(ctx Context) *setUUIDMetadataBuilder {
	return newSetUUIDMetadataBuilderWithContext(pn, ctx)
}

// RemoveUUIDMetadata Removes the metadata from a specified UUID.
func (pn *WebPubSub) RemoveUUIDMetadata() *removeUUIDMetadataBuilder {
	return newRemoveUUIDMetadataBuilder(pn)
}

// RemoveUUIDMetadataWithContext Removes the metadata from a specified UUID.
func (pn *WebPubSub) RemoveUUIDMetadataWithContext(ctx Context) *removeUUIDMetadataBuilder {
	return newRemoveUUIDMetadataBuilderWithContext(pn, ctx)
}

// GetAllChannelMetadata Returns a paginated list of Channel Metadata objects, optionally including the custom data object for each.
func (pn *WebPubSub) GetAllChannelMetadata() *getAllChannelMetadataBuilder {
	return newGetAllChannelMetadataBuilder(pn)
}

// GetAllChannelMetadataWithContext Returns a paginated list of Channel Metadata objects, optionally including the custom data object for each.
func (pn *WebPubSub) GetAllChannelMetadataWithContext(ctx Context) *getAllChannelMetadataBuilder {
	return newGetAllChannelMetadataBuilderWithContext(pn, ctx)
}

// GetChannelMetadata Returns metadata for the specified Channel, optionally including the custom data object for each.
func (pn *WebPubSub) GetChannelMetadata() *getChannelMetadataBuilder {
	return newGetChannelMetadataBuilder(pn)
}

// GetChannelMetadataWithContext Returns metadata for the specified Channel, optionally including the custom data object for each.
func (pn *WebPubSub) GetChannelMetadataWithContext(ctx Context) *getChannelMetadataBuilder {
	return newGetChannelMetadataBuilderWithContext(pn, ctx)
}

// SetChannelMetadata Set metadata for a Channel in the database, optionally including the custom data object for each.
func (pn *WebPubSub) SetChannelMetadata() *setChannelMetadataBuilder {
	return newSetChannelMetadataBuilder(pn)
}

// SetChannelMetadataWithContext Set metadata for a Channel in the database, optionally including the custom data object for each.
func (pn *WebPubSub) SetChannelMetadataWithContext(ctx Context) *setChannelMetadataBuilder {
	return newSetChannelMetadataBuilderWithContext(pn, ctx)
}

// RemoveChannelMetadata Removes the metadata from a specified channel.
func (pn *WebPubSub) RemoveChannelMetadata() *removeChannelMetadataBuilder {
	return newRemoveChannelMetadataBuilder(pn)
}

// RemoveChannelMetadataWithContext Removes the metadata from a specified channel.
func (pn *WebPubSub) RemoveChannelMetadataWithContext(ctx Context) *removeChannelMetadataBuilder {
	return newRemoveChannelMetadataBuilderWithContext(pn, ctx)
}

// GetMemberships The method returns a list of channel memberships for a user. This method doesn't return a user's subscriptions.
func (pn *WebPubSub) GetMemberships() *getMembershipsBuilderV2 {
	return newGetMembershipsBuilderV2(pn)
}

// GetMembershipsWithContext The method returns a list of channel memberships for a user. This method doesn't return a user's subscriptions.
func (pn *WebPubSub) GetMembershipsWithContext(ctx Context) *getMembershipsBuilderV2 {
	return newGetMembershipsBuilderV2WithContext(pn, ctx)
}

// GetChannelMembers The method returns a list of members in a channel. The list will include user metadata for members that have additional metadata stored in the database.
func (pn *WebPubSub) GetChannelMembers() *getChannelMembersBuilderV2 {
	return newGetChannelMembersBuilderV2(pn)
}

// GetChannelMembersWithContext The method returns a list of members in a channel. The list will include user metadata for members that have additional metadata stored in the database.
func (pn *WebPubSub) GetChannelMembersWithContext(ctx Context) *getChannelMembersBuilderV2 {
	return newGetChannelMembersBuilderV2WithContext(pn, ctx)
}

// SetChannelMembers This method sets members in a channel.
func (pn *WebPubSub) SetChannelMembers() *setChannelMembersBuilder {
	return newSetChannelMembersBuilder(pn)
}

// SetChannelMembersWithContext This method sets members in a channel.
func (pn *WebPubSub) SetChannelMembersWithContext(ctx Context) *setChannelMembersBuilder {
	return newSetChannelMembersBuilderWithContext(pn, ctx)
}

// RemoveChannelMembers Remove members from a Channel.
func (pn *WebPubSub) RemoveChannelMembers() *removeChannelMembersBuilder {
	return newRemoveChannelMembersBuilder(pn)
}

// RemoveChannelMembersWithContext Remove members from a Channel.
func (pn *WebPubSub) RemoveChannelMembersWithContext(ctx Context) *removeChannelMembersBuilder {
	return newRemoveChannelMembersBuilderWithContext(pn, ctx)
}

// SetMemberships Set channel memberships for a UUID.
func (pn *WebPubSub) SetMemberships() *setMembershipsBuilder {
	return newSetMembershipsBuilder(pn)
}

// SetMembershipsWithContext Set channel memberships for a UUID.
func (pn *WebPubSub) SetMembershipsWithContext(ctx Context) *setMembershipsBuilder {
	return newSetMembershipsBuilderWithContext(pn, ctx)
}

// RemoveMemberships Remove channel memberships for a UUID.
func (pn *WebPubSub) RemoveMemberships() *removeMembershipsBuilder {
	return newRemoveMembershipsBuilder(pn)
}

// RemoveMembershipsWithContext Remove channel memberships for a UUID.
func (pn *WebPubSub) RemoveMembershipsWithContext(ctx Context) *removeMembershipsBuilder {
	return newRemoveMembershipsBuilderWithContext(pn, ctx)
}

// ManageChannelMembers The method Set and Remove channel memberships for a user.
func (pn *WebPubSub) ManageChannelMembers() *manageChannelMembersBuilderV2 {
	return newManageChannelMembersBuilderV2(pn)
}

// ManageChannelMembersWithContext The method Set and Remove channel memberships for a user.
func (pn *WebPubSub) ManageChannelMembersWithContext(ctx Context) *manageChannelMembersBuilderV2 {
	return newManageChannelMembersBuilderV2WithContext(pn, ctx)
}

// ManageMemberships Manage the specified UUID's memberships. You can Add, Remove, and Update a UUID's memberships.
func (pn *WebPubSub) ManageMemberships() *manageMembershipsBuilderV2 {
	return newManageMembershipsBuilderV2(pn)
}

// ManageMembershipsWithContext Manage the specified UUID's memberships. You can Add, Remove, and Update a UUID's memberships.
func (pn *WebPubSub) ManageMembershipsWithContext(ctx Context) *manageMembershipsBuilderV2 {
	return newManageMembershipsBuilderV2WithContext(pn, ctx)
}

// Signal The signal() function is used to send a signal to all subscribers of a channel.
func (pn *WebPubSub) Signal() *signalBuilder {
	return newSignalBuilder(pn)
}

// SignalWithContext The signal() function is used to send a signal to all subscribers of a channel.
func (pn *WebPubSub) SignalWithContext(ctx Context) *signalBuilder {
	return newSignalBuilderWithContext(pn, ctx)
}

// SetState The state API is used to set/get key/value pairs specific to a subscriber UUID. State information is supplied as a JSON object of key/value pairs.
func (pn *WebPubSub) SetState() *setStateBuilder {
	return newSetStateBuilder(pn)
}

// SetStateWithContext The state API is used to set/get key/value pairs specific to a subscriber UUID. State information is supplied as a JSON object of key/value pairs.
func (pn *WebPubSub) SetStateWithContext(ctx Context) *setStateBuilder {
	return newSetStateBuilderWithContext(pn, ctx)
}

// Grant This function establishes access permissions for WebPubSub Access Manager (PAM) by setting the read or write attribute to true. A grant with read or write set to false (or not included) will revoke any previous grants with read or write set to true.
func (pn *WebPubSub) Grant() *grantBuilder {
	return newGrantBuilder(pn)
}

// GrantWithContext This function establishes access permissions for WebPubSub Access Manager (PAM) by setting the read or write attribute to true. A grant with read or write set to false (or not included) will revoke any previous grants with read or write set to true.
func (pn *WebPubSub) GrantWithContext(ctx Context) *grantBuilder {
	return newGrantBuilderWithContext(pn, ctx)
}

// GrantToken Use the Grant Token method to generate an auth token with embedded access control lists. The client sends the auth token to WebPubSub along with each request.
func (pn *WebPubSub) GrantToken() *grantTokenBuilder {
	return newGrantTokenBuilder(pn)
}

// GrantTokenWithContext Use the Grant Token method to generate an auth token with embedded access control lists. The client sends the auth token to WebPubSub along with each request.
func (pn *WebPubSub) GrantTokenWithContext(ctx Context) *grantTokenBuilder {
	return newGrantTokenBuilderWithContext(pn, ctx)
}

// RevokeToken Use the Grant Token method to generate an auth token with embedded access control lists. The client sends the auth token to WebPubSub along with each request.
func (pn *WebPubSub) RevokeToken() *revokeTokenBuilder {
	return newRevokeTokenBuilder(pn)
}

// RevokeTokenWithContext Use the Grant Token method to generate an auth token with embedded access control lists. The client sends the auth token to WebPubSub along with each request.
func (pn *WebPubSub) RevokeTokenWithContext(ctx Context) *revokeTokenBuilder {
	return newRevokeTokenBuilderWithContext(pn, ctx)
}

// AddMessageAction Add an action on a published message. Returns the added action in the response.
func (pn *WebPubSub) AddMessageAction() *addMessageActionsBuilder {
	return newAddMessageActionsBuilder(pn)
}

// AddMessageActionWithContext Add an action on a published message. Returns the added action in the response.
func (pn *WebPubSub) AddMessageActionWithContext(ctx Context) *addMessageActionsBuilder {
	return newAddMessageActionsBuilderWithContext(pn, ctx)
}

// GetMessageActions Get a list of message actions in a channel. Returns a list of actions in the response.
func (pn *WebPubSub) GetMessageActions() *getMessageActionsBuilder {
	return newGetMessageActionsBuilder(pn)
}

// GetMessageActionsWithContext Get a list of message actions in a channel. Returns a list of actions in the response.
func (pn *WebPubSub) GetMessageActionsWithContext(ctx Context) *getMessageActionsBuilder {
	return newGetMessageActionsBuilderWithContext(pn, ctx)
}

// RemoveMessageAction Remove a peviously added action on a published message. Returns an empty response.
func (pn *WebPubSub) RemoveMessageAction() *removeMessageActionsBuilder {
	return newRemoveMessageActionsBuilder(pn)
}

// RemoveMessageActionWithContext Remove a peviously added action on a published message. Returns an empty response.
func (pn *WebPubSub) RemoveMessageActionWithContext(ctx Context) *removeMessageActionsBuilder {
	return newRemoveMessageActionsBuilderWithContext(pn, ctx)
}

// SetToken Stores a single token in the Token Management System for use in API calls.
func (pn *WebPubSub) SetToken(token string) {
	pn.tokenManager.StoreToken(token)
}

// ResetTokenManager resets the token manager.
func (pn *WebPubSub) ResetTokenManager() {
	pn.tokenManager.CleanUp()
}

// Unsubscribe When subscribed to a single channel, this function causes the client to issue a leave from the channel and close any open socket to the WebPubSub Network. For multiplexed channels, the specified channel(s) will be removed and the socket remains open until there are no more channels remaining in the list.
func (pn *WebPubSub) Unsubscribe() *unsubscribeBuilder {
	return newUnsubscribeBuilder(pn)
}

// AddListener lets you add a new listener.
func (pn *WebPubSub) AddListener(listener *Listener) {
	pn.subscriptionManager.AddListener(listener)
}

// RemoveListener lets you remove new listener.
func (pn *WebPubSub) RemoveListener(listener *Listener) {
	pn.subscriptionManager.RemoveListener(listener)
}

// GetListeners gets all the existing isteners.
func (pn *WebPubSub) GetListeners() map[*Listener]bool {
	return pn.subscriptionManager.GetListeners()
}

// Leave unsubscribes from a channel.
func (pn *WebPubSub) Leave() *leaveBuilder {
	return newLeaveBuilder(pn)
}

// LeaveWithContext unsubscribes from a channel.
func (pn *WebPubSub) LeaveWithContext(ctx Context) *leaveBuilder {
	return newLeaveBuilderWithContext(pn, ctx)
}

// Presence lets you subscribe to a presence channel.
func (pn *WebPubSub) Presence() *presenceBuilder {
	return newPresenceBuilder(pn)
}

// PresenceWithContext lets you subscribe to a presence channel.
func (pn *WebPubSub) PresenceWithContext(ctx Context) *presenceBuilder {
	return newPresenceBuilderWithContext(pn, ctx)
}

// Heartbeat You can send presence heartbeat notifications without subscribing to a channel. These notifications are sent periodically and indicate whether a client is connected or not.
func (pn *WebPubSub) Heartbeat() *heartbeatBuilder {
	return newHeartbeatBuilder(pn)
}

// HeartbeatWithContext You can send presence heartbeat notifications without subscribing to a channel. These notifications are sent periodically and indicate whether a client is connected or not.
func (pn *WebPubSub) HeartbeatWithContext(ctx Context) *heartbeatBuilder {
	return newHeartbeatBuilderWithContext(pn, ctx)
}

// SetClient Set a client for transactional requests (Non Subscribe).
func (pn *WebPubSub) SetClient(c *http.Client) {
	pn.Lock()
	pn.client = c
	pn.Unlock()
}

// GetClient Get a client for transactional requests (Non Subscribe).
func (pn *WebPubSub) GetClient() *http.Client {
	pn.Lock()
	defer pn.Unlock()

	if pn.client == nil {
		if pn.Config.UseHTTP2 {
			pn.client = NewHTTP2Client(pn.Config.ConnectTimeout,
				pn.Config.SubscribeRequestTimeout)
		} else {
			pn.client = NewHTTP1Client(pn.Config.ConnectTimeout,
				pn.Config.NonSubscribeRequestTimeout,
				pn.Config.MaxIdleConnsPerHost)
		}
	}

	return pn.client
}

// SetSubscribeClient Set a client for transactional requests.
func (pn *WebPubSub) SetSubscribeClient(client *http.Client) {
	pn.Lock()
	pn.subscribeClient = client
	pn.Unlock()
}

// GetSubscribeClient Get a client for transactional requests.
func (pn *WebPubSub) GetSubscribeClient() *http.Client {
	pn.Lock()
	defer pn.Unlock()
	if pn.subscribeClient == nil {

		if pn.Config.UseHTTP2 {
			pn.subscribeClient = NewHTTP2Client(pn.Config.ConnectTimeout,
				pn.Config.SubscribeRequestTimeout)
		} else {
			pn.subscribeClient = NewHTTP1Client(pn.Config.ConnectTimeout,
				pn.Config.SubscribeRequestTimeout, pn.Config.MaxIdleConnsPerHost)
		}

	}

	return pn.subscribeClient
}

// GetSubscribedChannels gets a list of all subscribed channels.
func (pn *WebPubSub) GetSubscribedChannels() []string {
	return pn.subscriptionManager.getSubscribedChannels()
}

// GetSubscribedGroups gets a list of all subscribed channel groups.
func (pn *WebPubSub) GetSubscribedGroups() []string {
	return pn.subscriptionManager.getSubscribedGroups()
}

// UnsubscribeAll Unsubscribe from all channels and all channel groups.
func (pn *WebPubSub) UnsubscribeAll() {
	pn.subscriptionManager.unsubscribeAll()
}

// ListPushProvisions Request for all channels on which push notification has been enabled using specified pushToken.
func (pn *WebPubSub) ListPushProvisions() *listPushProvisionsRequestBuilder {
	return newListPushProvisionsRequestBuilder(pn)
}

// ListPushProvisionsWithContext Request for all channels on which push notification has been enabled using specified pushToken.
func (pn *WebPubSub) ListPushProvisionsWithContext(ctx Context) *listPushProvisionsRequestBuilder {
	return newListPushProvisionsRequestBuilderWithContext(pn, ctx)
}

// AddPushNotificationsOnChannels Enable push notifications on provided set of channels.
func (pn *WebPubSub) AddPushNotificationsOnChannels() *addPushNotificationsOnChannelsBuilder {
	return newAddPushNotificationsOnChannelsBuilder(pn)
}

// AddPushNotificationsOnChannelsWithContext Enable push notifications on provided set of channels.
func (pn *WebPubSub) AddPushNotificationsOnChannelsWithContext(ctx Context) *addPushNotificationsOnChannelsBuilder {
	return newAddPushNotificationsOnChannelsBuilderWithContext(pn, ctx)
}

// RemovePushNotificationsFromChannels Disable push notifications on provided set of channels.
func (pn *WebPubSub) RemovePushNotificationsFromChannels() *removeChannelsFromPushBuilder {
	return newRemoveChannelsFromPushBuilder(pn)
}

// RemovePushNotificationsFromChannelsWithContext Disable push notifications on provided set of channels.
func (pn *WebPubSub) RemovePushNotificationsFromChannelsWithContext(ctx Context) *removeChannelsFromPushBuilder {
	return newRemoveChannelsFromPushBuilderWithContext(pn, ctx)
}

// RemoveAllPushNotifications Disable push notifications from all channels registered with the specified pushToken.
func (pn *WebPubSub) RemoveAllPushNotifications() *removeAllPushChannelsForDeviceBuilder {
	return newRemoveAllPushChannelsForDeviceBuilder(pn)
}

// RemoveAllPushNotificationsWithContext Disable push notifications from all channels registered with the specified pushToken.
func (pn *WebPubSub) RemoveAllPushNotificationsWithContext(ctx Context) *removeAllPushChannelsForDeviceBuilder {
	return newRemoveAllPushChannelsForDeviceBuilderWithContext(pn, ctx)
}

// AddChannelToChannelGroup This function adds a channel to a channel group.
func (pn *WebPubSub) AddChannelToChannelGroup() *addChannelToChannelGroupBuilder {
	return newAddChannelToChannelGroupBuilder(pn)
}

// AddChannelToChannelGroupWithContext This function adds a channel to a channel group.
func (pn *WebPubSub) AddChannelToChannelGroupWithContext(ctx Context) *addChannelToChannelGroupBuilder {
	return newAddChannelToChannelGroupBuilderWithContext(pn, ctx)
}

// RemoveChannelFromChannelGroup This function removes the channels from the channel group.
func (pn *WebPubSub) RemoveChannelFromChannelGroup() *removeChannelFromChannelGroupBuilder {
	return newRemoveChannelFromChannelGroupBuilder(pn)
}

// RemoveChannelFromChannelGroupWithContext This function removes the channels from the channel group.
func (pn *WebPubSub) RemoveChannelFromChannelGroupWithContext(ctx Context) *removeChannelFromChannelGroupBuilder {
	return newRemoveChannelFromChannelGroupBuilderWithContext(pn, ctx)
}

// DeleteChannelGroup This function removes the channel group.
func (pn *WebPubSub) DeleteChannelGroup() *deleteChannelGroupBuilder {
	return newDeleteChannelGroupBuilder(pn)
}

// DeleteChannelGroupWithContext This function removes the channel group.
func (pn *WebPubSub) DeleteChannelGroupWithContext(ctx Context) *deleteChannelGroupBuilder {
	return newDeleteChannelGroupBuilderWithContext(pn, ctx)
}

// ListChannelsInChannelGroup This function lists all the channels of the channel group.
func (pn *WebPubSub) ListChannelsInChannelGroup() *allChannelGroupBuilder {
	return newAllChannelGroupBuilder(pn)
}

// ListChannelsInChannelGroupWithContext This function lists all the channels of the channel group.
func (pn *WebPubSub) ListChannelsInChannelGroupWithContext(ctx Context) *allChannelGroupBuilder {
	return newAllChannelGroupBuilderWithContext(pn, ctx)
}

// GetState The state API is used to set/get key/value pairs specific to a subscriber UUID. State information is supplied as a JSON object of key/value pairs.
func (pn *WebPubSub) GetState() *getStateBuilder {
	return newGetStateBuilder(pn)
}

// GetStateWithContext The state API is used to set/get key/value pairs specific to a subscriber UUID. State information is supplied as a JSON object of key/value pairs.
func (pn *WebPubSub) GetStateWithContext(ctx Context) *getStateBuilder {
	return newGetStateBuilderWithContext(pn, ctx)
}

// HereNow You can obtain information about the current state of a channel including a list of unique user-ids currently subscribed to the channel and the total occupancy count of the channel by calling the HereNow() function in your application.
func (pn *WebPubSub) HereNow() *hereNowBuilder {
	return newHereNowBuilder(pn)
}

// HereNowWithContext You can obtain information about the current state of a channel including a list of unique user-ids currently subscribed to the channel and the total occupancy count of the channel by calling the HereNow() function in your application.
func (pn *WebPubSub) HereNowWithContext(ctx Context) *hereNowBuilder {
	return newHereNowBuilderWithContext(pn, ctx)
}

// WhereNow You can obtain information about the current list of channels to which a UUID is subscribed to by calling the WhereNow() function in your application.
func (pn *WebPubSub) WhereNow() *whereNowBuilder {
	return newWhereNowBuilder(pn)
}

// WhereNowWithContext You can obtain information about the current list of channels to which a UUID is subscribed to by calling the WhereNow() function in your application.
func (pn *WebPubSub) WhereNowWithContext(ctx Context) *whereNowBuilder {
	return newWhereNowBuilderWithContext(pn, ctx)
}

// Time This function will return a 17 digit precision Unix epoch.
func (pn *WebPubSub) Time() *timeBuilder {
	return newTimeBuilder(pn)
}

// TimeWithContext This function will return a 17 digit precision Unix epoch.
func (pn *WebPubSub) TimeWithContext(ctx Context) *timeBuilder {
	return newTimeBuilderWithContext(pn, ctx)
}

// CreatePushPayload This method creates the push payload for use in the appropriate endpoint calls.
func (pn *WebPubSub) CreatePushPayload() *publishPushHelperBuilder {
	return newPublishPushHelperBuilder(pn)
}

// CreatePushPayloadWithContext This method creates the push payload for use in the appropriate endpoint calls.
func (pn *WebPubSub) CreatePushPayloadWithContext(ctx Context) *publishPushHelperBuilder {
	return newPublishPushHelperBuilderWithContext(pn, ctx)
}

// DeleteMessages Removes the messages from the history of a specific channel.
func (pn *WebPubSub) DeleteMessages() *historyDeleteBuilder {
	return newHistoryDeleteBuilder(pn)
}

// DeleteMessagesWithContext Removes the messages from the history of a specific channel.
func (pn *WebPubSub) DeleteMessagesWithContext(ctx Context) *historyDeleteBuilder {
	return newHistoryDeleteBuilderWithContext(pn, ctx)
}

// SendFile Clients can use this SDK method to upload a file and publish it to other users in a channel.
func (pn *WebPubSub) SendFile() *sendFileBuilder {
	return newSendFileBuilder(pn)
}

// SendFileWithContext Clients can use this SDK method to upload a file and publish it to other users in a channel.
func (pn *WebPubSub) SendFileWithContext(ctx Context) *sendFileBuilder {
	return newSendFileBuilderWithContext(pn, ctx)
}

// ListFiles Provides the ability to fetch all files in a channel.
func (pn *WebPubSub) ListFiles() *listFilesBuilder {
	return newListFilesBuilder(pn)
}

// ListFilesWithContext Provides the ability to fetch all files in a channel.
func (pn *WebPubSub) ListFilesWithContext(ctx Context) *listFilesBuilder {
	return newListFilesBuilderWithContext(pn, ctx)
}

// GetFileURL gets the URL of the file.
func (pn *WebPubSub) GetFileURL() *getFileURLBuilder {
	return newGetFileURLBuilder(pn)
}

// GetFileURLWithContext gets the URL of the file.
func (pn *WebPubSub) GetFileURLWithContext(ctx Context) *getFileURLBuilder {
	return newGetFileURLBuilderWithContext(pn, ctx)
}

// DownloadFile Provides the ability to fetch an individual file.
func (pn *WebPubSub) DownloadFile() *downloadFileBuilder {
	return newDownloadFileBuilder(pn)
}

// DownloadFileWithContext Provides the ability to fetch an individual file.
func (pn *WebPubSub) DownloadFileWithContext(ctx Context) *downloadFileBuilder {
	return newDownloadFileBuilderWithContext(pn, ctx)
}

// DeleteFile Provides the ability to delete an individual file.
func (pn *WebPubSub) DeleteFile() *deleteFileBuilder {
	return newDeleteFileBuilder(pn)
}

// DeleteFileWithContext Provides the ability to delete an individual file
func (pn *WebPubSub) DeleteFileWithContext(ctx Context) *deleteFileBuilder {
	return newDeleteFileBuilderWithContext(pn, ctx)
}

// PublishFileMessage Provides the ability to publish the asccociated messages with the uploaded file in case of failure to auto publish. To be used only in the case of failure.
func (pn *WebPubSub) PublishFileMessage() *publishFileMessageBuilder {
	return newPublishFileMessageBuilder(pn)
}

// PublishFileMessageWithContext Provides the ability to publish the asccociated messages with the uploaded file in case of failure to auto publish. To be used only in the case of failure.
func (pn *WebPubSub) PublishFileMessageWithContext(ctx Context) *publishFileMessageBuilder {
	return newPublishFileMessageBuilderWithContext(pn, ctx)
}

// Destroy stops all open requests, removes listeners, closes heartbeats, and cleans up.
func (pn *WebPubSub) Destroy() {
	pn.Config.Log.Println("Calling Destroy")
	pn.UnsubscribeAll()
	pn.cancel()

	if pn.subscriptionManager != nil {
		pn.subscriptionManager.Destroy()
		pn.Config.Log.Println("after subscription manager Destroy")
	}

	pn.Config.Log.Println("calling subscriptionManager Destroy")
	if pn.heartbeatManager != nil {
		pn.heartbeatManager.Destroy()
		pn.Config.Log.Println("after heartbeat manager Destroy")
	}

	pn.Config.Log.Println("After Destroy")
	pn.Config.Log.Println("calling RemoveAllListeners")
	pn.subscriptionManager.RemoveAllListeners()
	pn.Config.Log.Println("after RemoveAllListeners")
	close(pn.jobQueue)
	pn.Config.Log.Println("after close jobQueue")
	pn.requestWorkers.Close()
	pn.Config.Log.Println("after close requestWorkers")
	pn.tokenManager.CleanUp()
	pn.client.CloseIdleConnections()

}

func (pn *WebPubSub) getPublishSequence() int {
	pn.publishSequenceMutex.Lock()
	defer pn.publishSequenceMutex.Unlock()

	if pn.nextPublishSequence == MaxSequence {
		pn.nextPublishSequence = 1
	} else {
		pn.nextPublishSequence++
	}

	return pn.nextPublishSequence
}

func GenerateUUID() string {
	return utils.UUID()
}

// NewWebPubSub instantiates a WebPubSub instance with default values.
func NewWebPubSub(pnconf *Config) *WebPubSub {
	ctx, cancel := contextWithCancel(backgroundContext)

	if pnconf.Log == nil {
		pnconf.Log = log.New(ioutil.Discard, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
	pnconf.Log.Println(fmt.Sprintf("WebPubSub Go v4 SDK: %s\npnconf: %v\n%s\n%s\n%s", Version, pnconf, runtime.Version(), runtime.GOARCH, runtime.GOOS))

	utils.CheckUUID(pnconf.UUID)
	pn := &WebPubSub{
		Config:              pnconf,
		nextPublishSequence: 0,
		ctx:                 ctx,
		cancel:              cancel,
	}

	pn.subscriptionManager = newSubscriptionManager(pn, ctx)
	pn.heartbeatManager = newHeartbeatManager(pn, ctx)
	pn.telemetryManager = newTelemetryManager(pnconf.MaximumLatencyDataAge, ctx)
	pn.jobQueue = make(chan *JobQItem)
	pn.requestWorkers = pn.newNonSubQueueProcessor(pnconf.MaxWorkers, ctx)
	pn.tokenManager = newTokenManager(pn, ctx)

	return pn
}

func (pn *WebPubSub) newNonSubQueueProcessor(maxWorkers int, ctx Context) *RequestWorkers {
	workers := make(chan chan *JobQItem, maxWorkers)

	pn.Config.Log.Printf("Init RequestWorkers: workers %d", maxWorkers)

	p := &RequestWorkers{
		WorkersChannel: workers,
		MaxWorkers:     maxWorkers,
	}
	p.Start(pn, ctx)
	return p
}

// NewWebPubSubDemo returns an instance with demo keys
func NewWebPubSubDemo() *WebPubSub {
	return NewWebPubSub(NewDemoConfig())
}
