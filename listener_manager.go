package webpubsub

import (
	"sync"
)

// Listener type has all the `types` of response events
type Listener struct {
	Status              chan *WPSStatus
	Message             chan *WPSMessage
	Presence            chan *WPSPresence
	Signal              chan *WPSMessage
	UUIDEvent           chan *WPSUUIDEvent
	ChannelEvent        chan *WPSChannelEvent
	MembershipEvent     chan *WPSMembershipEvent
	MessageActionsEvent chan *WPSMessageActionsEvent
	File                chan *WPSFilesEvent
}

// NewListener initates the listener to facilitate the event handling
func NewListener() *Listener {
	return &Listener{
		Status:              make(chan *WPSStatus),
		Message:             make(chan *WPSMessage),
		Presence:            make(chan *WPSPresence),
		Signal:              make(chan *WPSMessage),
		UUIDEvent:           make(chan *WPSUUIDEvent),
		ChannelEvent:        make(chan *WPSChannelEvent),
		MembershipEvent:     make(chan *WPSMembershipEvent),
		MessageActionsEvent: make(chan *WPSMessageActionsEvent),
		File:                make(chan *WPSFilesEvent),
	}
}

// ListenerManager is used in the internal handling of listeners.
type ListenerManager struct {
	sync.RWMutex
	ctx                  Context
	listeners            map[*Listener]bool
	exitListener         chan bool
	exitListenerAnnounce chan bool
	webpubsub            *WebPubSub
}

func newListenerManager(ctx Context, pn *WebPubSub) *ListenerManager {
	return &ListenerManager{
		listeners:            make(map[*Listener]bool, 2),
		ctx:                  ctx,
		exitListener:         make(chan bool),
		exitListenerAnnounce: make(chan bool),
		webpubsub:            pn,
	}
}

func (m *ListenerManager) addListener(listener *Listener) {
	m.Lock()

	m.listeners[listener] = true
	m.Unlock()
}

func (m *ListenerManager) removeListener(listener *Listener) {
	m.webpubsub.Config.Log.Println("before removeListener")
	m.Lock()
	m.webpubsub.Config.Log.Println("in removeListener lock")
	delete(m.listeners, listener)
	m.Unlock()
	m.webpubsub.Config.Log.Println("after removeListener")
}

func (m *ListenerManager) removeAllListeners() {
	m.webpubsub.Config.Log.Println("in removeAllListeners")
	m.Lock()
	lis := m.listeners
	for l := range lis {
		delete(m.listeners, l)
	}
	m.Unlock()
}

func (m *ListenerManager) copyListeners() map[*Listener]bool {
	m.Lock()
	lis := make(map[*Listener]bool)
	for k, v := range m.listeners {
		lis[k] = v
	}
	m.Unlock()
	return lis
}

func (m *ListenerManager) announceStatus(status *WPSStatus) {
	go func() {
		lis := m.copyListeners()
	AnnounceStatusLabel:
		for l := range lis {
			select {
			case <-m.exitListener:
				m.webpubsub.Config.Log.Println("announceStatus exitListener")
				break AnnounceStatusLabel
			case l.Status <- status:
			}
		}
		m.webpubsub.Config.Log.Println("announceStatus exit")
	}()
}

func (m *ListenerManager) announceMessage(message *WPSMessage) {
	go func() {
		lis := m.copyListeners()
	AnnounceMessageLabel:
		for l := range lis {
			select {
			case <-m.exitListenerAnnounce:
				m.webpubsub.Config.Log.Println("announceMessage exitListenerAnnounce")
				break AnnounceMessageLabel
			case l.Message <- message:
			}
		}

	}()
}

func (m *ListenerManager) announceSignal(message *WPSMessage) {
	go func() {
		lis := m.copyListeners()

	AnnounceSignalLabel:
		for l := range lis {
			select {
			case <-m.exitListener:
				m.webpubsub.Config.Log.Println("announceSignal exitListener")
				break AnnounceSignalLabel

			case l.Signal <- message:
			}
		}
	}()
}

func (m *ListenerManager) announceUUIDEvent(message *WPSUUIDEvent) {
	go func() {
		lis := m.copyListeners()

	AnnounceUUIDEventLabel:
		for l := range lis {
			select {
			case <-m.exitListener:
				m.webpubsub.Config.Log.Println("announceUUIDEvent exitListener")
				break AnnounceUUIDEventLabel

			case l.UUIDEvent <- message:
				m.webpubsub.Config.Log.Println("l.UUIDEvent", message)
			}
		}
	}()
}

func (m *ListenerManager) announceChannelEvent(message *WPSChannelEvent) {
	go func() {
		lis := m.copyListeners()

	AnnounceChannelEventLabel:
		for l := range lis {
			m.webpubsub.Config.Log.Println("l.ChannelEvent", l)
			select {
			case <-m.exitListener:
				m.webpubsub.Config.Log.Println("announceChannelEvent exitListener")
				break AnnounceChannelEventLabel

			case l.ChannelEvent <- message:
				m.webpubsub.Config.Log.Println("l.ChannelEvent", message)
			}
		}
	}()
}

func (m *ListenerManager) announceMembershipEvent(message *WPSMembershipEvent) {
	go func() {
		lis := m.copyListeners()

	AnnounceMembershipEvent:
		for l := range lis {
			select {
			case <-m.exitListener:
				m.webpubsub.Config.Log.Println("announceMembershipEvent exitListener")
				break AnnounceMembershipEvent

			case l.MembershipEvent <- message:
				m.webpubsub.Config.Log.Println("l.MembershipEvent", message)
			}
		}
	}()
}

func (m *ListenerManager) announceMessageActionsEvent(message *WPSMessageActionsEvent) {
	go func() {
		lis := m.copyListeners()

	AnnounceMessageActionsEvent:
		for l := range lis {
			select {
			case <-m.exitListener:
				m.webpubsub.Config.Log.Println("announceMessageActionsEvent exitListener")
				break AnnounceMessageActionsEvent

			case l.MessageActionsEvent <- message:
				m.webpubsub.Config.Log.Println("l.MessageActionsEvent", message)
			}
		}
	}()
}

func (m *ListenerManager) announcePresence(presence *WPSPresence) {
	go func() {
		lis := m.copyListeners()

	AnnouncePresenceLabel:
		for l := range lis {
			select {
			case <-m.exitListener:
				m.webpubsub.Config.Log.Println("announcePresence exitListener")
				break AnnouncePresenceLabel

			case l.Presence <- presence:
			}
		}
	}()
}

func (m *ListenerManager) announceFile(file *WPSFilesEvent) {
	go func() {
		lis := m.copyListeners()

	AnnounceFileLabel:
		for l := range lis {
			select {
			case <-m.exitListener:
				m.webpubsub.Config.Log.Println("announceFile exitListener")
				break AnnounceFileLabel

			case l.File <- file:
			}
		}
	}()
}

// WPSStatus is the status struct
type WPSStatus struct {
	Category              StatusCategory
	Operation             OperationType
	ErrorData             error
	Error                 bool
	TLSEnabled            bool
	StatusCode            int
	UUID                  string
	AuthKey               string
	Origin                string
	ClientRequest         interface{} // Should be same for non-google environment
	AffectedChannels      []string
	AffectedChannelGroups []string
}

// WPSMessage is the Message Response for Subscribe
type WPSMessage struct {
	Message           interface{}
	UserMetadata      interface{}
	SubscribedChannel string
	ActualChannel     string
	Channel           string
	Subscription      string
	Publisher         string
	Timetoken         int64
}

// WPSPresence is the Message Response for Presence
type WPSPresence struct {
	Event             string
	UUID              string
	SubscribedChannel string
	ActualChannel     string
	Channel           string
	Subscription      string
	Occupancy         int
	Timetoken         int64
	Timestamp         int64
	UserMetadata      map[string]interface{}
	State             interface{}
	Join              []string
	Leave             []string
	Timeout           []string
	HereNowRefresh    bool
}

// WPSUUIDEvent is the Response for an User Event
type WPSUUIDEvent struct {
	Event             WPSObjectsEvent
	UUID              string
	Description       string
	Timestamp         string
	Name              string
	ExternalID        string
	ProfileURL        string
	Email             string
	Updated           string
	ETag              string
	Custom            map[string]interface{}
	SubscribedChannel string
	ActualChannel     string
	Channel           string
	Subscription      string
}

// WPSChannelEvent is the Response for a Space Event
type WPSChannelEvent struct {
	Event             WPSObjectsEvent
	ChannelID         string
	Description       string
	Timestamp         string
	Name              string
	Updated           string
	ETag              string
	Custom            map[string]interface{}
	SubscribedChannel string
	ActualChannel     string
	Channel           string
	Subscription      string
}

// WPSMembershipEvent is the Response for a Membership Event
type WPSMembershipEvent struct {
	Event             WPSObjectsEvent
	UUID              string
	ChannelID         string
	Description       string
	Timestamp         string
	Custom            map[string]interface{}
	SubscribedChannel string
	ActualChannel     string
	Channel           string
	Subscription      string
}

// WPSMessageActionsEvent is the Response for a Message Actions Event
type WPSMessageActionsEvent struct {
	Event             WPSMessageActionsEventType
	Data              WPSMessageActionsResponse
	SubscribedChannel string
	ActualChannel     string
	Channel           string
	Subscription      string
}

// WPSFilesEvent is the Response for a Files Event
type WPSFilesEvent struct {
	File              WPSFileMessageAndDetails
	UserMetadata      interface{}
	SubscribedChannel string
	ActualChannel     string
	Channel           string
	Subscription      string
	Publisher         string
	Timetoken         int64
}