package webpubsub

// WPSUUID is the Objects API user struct
type WPSUUID struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	ExternalID string                 `json:"externalId"`
	ProfileURL string                 `json:"profileUrl"`
	Email      string                 `json:"email"`
	Updated    string                 `json:"updated"`
	ETag       string                 `json:"eTag"`
	Custom     map[string]interface{} `json:"custom"`
}

// WPSChannel is the Objects API space struct
type WPSChannel struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Updated     string                 `json:"updated"`
	ETag        string                 `json:"eTag"`
	Custom      map[string]interface{} `json:"custom"`
}

// WPSChannelMembers is the Objects API Members struct
type WPSChannelMembers struct {
	ID      string                 `json:"id"`
	UUID    WPSUUID                `json:"uuid"`
	Created string                 `json:"created"`
	Updated string                 `json:"updated"`
	ETag    string                 `json:"eTag"`
	Custom  map[string]interface{} `json:"custom"`
}

// WPSMemberships is the Objects API Memberships struct
type WPSMemberships struct {
	ID      string                 `json:"id"`
	Channel WPSChannel             `json:"channel"`
	Created string                 `json:"created"`
	Updated string                 `json:"updated"`
	ETag    string                 `json:"eTag"`
	Custom  map[string]interface{} `json:"custom"`
}

// WPSChannelMembersUUID is the Objects API Members input struct used to add members
type WPSChannelMembersUUID struct {
	ID string `json:"id"`
}

// WPSChannelMembersSet is the Objects API Members input struct used to add members
type WPSChannelMembersSet struct {
	UUID   WPSChannelMembersUUID  `json:"uuid"`
	Custom map[string]interface{} `json:"custom"`
}

// WPSChannelMembersRemove is the Objects API Members struct used to remove members
type WPSChannelMembersRemove struct {
	UUID WPSChannelMembersUUID `json:"uuid"`
}

// WPSMembershipsChannel is the Objects API Memberships input struct used to add members
type WPSMembershipsChannel struct {
	ID string `json:"id"`
}

// WPSMembershipsSet is the Objects API Memberships input struct used to add members
type WPSMembershipsSet struct {
	Channel WPSMembershipsChannel  `json:"channel"`
	Custom  map[string]interface{} `json:"custom"`
}

// WPSMembershipsRemove is the Objects API Memberships struct used to remove members
type WPSMembershipsRemove struct {
	Channel WPSMembershipsChannel `json:"channel"`
}

// WPSObjectsResponse is the Objects API collective Response struct of all methods.
type WPSObjectsResponse struct {
	Event       WPSObjectsEvent        `json:"event"` // enum value
	EventType   WPSObjectsEventType    `json:"type"`  // enum value
	Name        string                 `json:"name"`
	ID          string                 `json:"id"`          // the uuid if user related
	Channel     string                 `json:"channel"`     // the channel if space related
	Description string                 `json:"description"` // the description of what happened
	Timestamp   string                 `json:"timestamp"`   // the timetoken of the event
	ExternalID  string                 `json:"externalId"`
	ProfileURL  string                 `json:"profileUrl"`
	Email       string                 `json:"email"`
	Updated     string                 `json:"updated"`
	ETag        string                 `json:"eTag"`
	Custom      map[string]interface{} `json:"custom"`
	Data        map[string]interface{} `json:"data"`
}

// WPSManageMembershipsBody is the Objects API input to add, remove or update membership
type WPSManageMembershipsBody struct {
	Set    []WPSMembershipsSet    `json:"set"`
	Remove []WPSMembershipsRemove `json:"delete"`
}

// WPSManageChannelMembersBody is the Objects API input to add, remove or update members
type WPSManageChannelMembersBody struct {
	Set    []WPSChannelMembersSet    `json:"set"`
	Remove []WPSChannelMembersRemove `json:"delete"`
}
