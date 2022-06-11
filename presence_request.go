package webpubsub

import (
	"strings"
)

type presenceBuilder struct {
	opts *presenceOpts
}

type presenceOpts struct {
	webpubsub *WebPubSub

	channels      []string
	channelGroups []string
	connected     bool
	ctx           Context
	queryParam    map[string]string
	state         map[string]interface{}
}

func newPresenceBuilder(webpubsub *WebPubSub) *presenceBuilder {
	builder := presenceBuilder{
		opts: &presenceOpts{
			webpubsub: webpubsub,
		},
	}

	return &builder
}

func newPresenceBuilderWithContext(webpubsub *WebPubSub, context Context) *presenceBuilder {
	builder := presenceBuilder{
		opts: &presenceOpts{
			webpubsub: webpubsub,
			ctx:       context,
		},
	}

	return &builder
}

// Channels sets the Channels for the Presence request.
func (b *presenceBuilder) Channels(ch []string) *presenceBuilder {
	b.opts.channels = ch

	return b
}

// ChannelGroups sets the ChannelGroups for the Presence request.
func (b *presenceBuilder) ChannelGroups(cg []string) *presenceBuilder {
	b.opts.channelGroups = cg

	return b
}

// Channels sets the Channels for the Presence request.
func (b *presenceBuilder) Connected(connected bool) *presenceBuilder {
	b.opts.connected = connected

	return b
}

// QueryParam accepts a map, the keys and values of the map are passed as the query string parameters of the URL called by the API.
func (b *presenceBuilder) QueryParam(queryParam map[string]string) *presenceBuilder {
	b.opts.queryParam = queryParam

	return b
}

// State sets the State for the Set State request.
func (b *presenceBuilder) State(state map[string]interface{}) *presenceBuilder {
	b.opts.state = state
	return b
}

func (b *presenceBuilder) Execute() {
	if b.opts.connected {
		for _, ch := range b.opts.channels {
			if strings.Contains(ch, "-pnpres") {
				ch = strings.Replace(ch, "-pnpres", "", -1)
			}
			b.opts.webpubsub.heartbeatManager.Lock()
			b.opts.webpubsub.heartbeatManager.heartbeatChannels[ch] = newSubscriptionItem(ch)
			b.opts.webpubsub.heartbeatManager.Unlock()
		}
		for _, cg := range b.opts.channelGroups {
			if strings.Contains(cg, "-pnpres") {
				cg = strings.Replace(cg, "-pnpres", "", -1)
			}
			b.opts.webpubsub.heartbeatManager.Lock()
			b.opts.webpubsub.heartbeatManager.heartbeatGroups[cg] = newSubscriptionItem(cg)
			b.opts.webpubsub.heartbeatManager.Unlock()
		}
		b.opts.webpubsub.heartbeatManager.state = b.opts.state
		b.opts.webpubsub.heartbeatManager.queryParam = b.opts.queryParam
		b.opts.webpubsub.heartbeatManager.startHeartbeatTimer(true)
	} else {
		b.opts.webpubsub.heartbeatManager.Lock()
		b.opts.webpubsub.heartbeatManager.heartbeatChannels = make(map[string]*SubscriptionItem)
		b.opts.webpubsub.heartbeatManager.heartbeatGroups = make(map[string]*SubscriptionItem)
		b.opts.webpubsub.heartbeatManager.state = nil
		b.opts.webpubsub.heartbeatManager.queryParam = nil

		b.opts.webpubsub.heartbeatManager.Unlock()
	}
}
