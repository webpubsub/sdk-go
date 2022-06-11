package webpubsub

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/webpubsub/go/v7/tests/stubs"
)

func TestExponentialExhaustion(t *testing.T) {
	assert := assert.New(t)
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/time/0",
		Query:              "",
		ResponseBody:       `[15078947309567840]`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})

	config := NewConfig(GenerateUUID())

	pn := NewWebPubSub(config)
	pn.Config.MaximumReconnectionRetries = 2
	pn.Config.NonSubscribeRequestTimeout = 2
	pn.Config.ConnectTimeout = 2
	pn.Config.WPSReconnectionPolicy = WPSExponentialPolicycy

	pn.SetClient(interceptor.GetClient())
	t1 := time.Now()
	r := newReconnectionManager(pn)
	reconnectionExhausted := false
	r.HandleOnMaxReconnectionExhaustion(func() {
		reconnectionExhausted = true
	})

	r.startHeartbeatTimer()
	t2 := time.Now()
	diff := t2.Unix() - t1.Unix()
	assert.True((diff >= 11) && (diff <= 12))
	assert.True(reconnectionExhausted)
	r.stopHeartbeatTimer()
}

func TestLinearExhaustion(t *testing.T) {
	assert := assert.New(t)
	interceptor := stubs.NewInterceptor()
	interceptor.AddStub(&stubs.Stub{
		Method:             "GET",
		Path:               "/time/0",
		Query:              "",
		ResponseBody:       `[15078947309567840]`,
		IgnoreQueryKeys:    []string{"uuid", "pnsdk"},
		ResponseStatusCode: 200,
	})

	config := NewConfig(GenerateUUID())

	pn := NewWebPubSub(config)
	pn.Config.MaximumReconnectionRetries = 1
	pn.Config.WPSReconnectionPolicy = WPSLinearPolicy
	pn.SetClient(interceptor.GetClient())
	t1 := time.Now()
	r := newReconnectionManager(pn)
	reconnectionExhausted := false
	r.HandleOnMaxReconnectionExhaustion(func() {
		reconnectionExhausted = true
	})

	r.startHeartbeatTimer()
	t2 := time.Now()
	diff := t2.Unix() - t1.Unix()
	assert.True((diff >= 10) && (diff <= 11))
	assert.True(reconnectionExhausted)
	r.stopHeartbeatTimer()
}

func TestReconnect(t *testing.T) {
	assert := assert.New(t)

	config := NewConfig(GenerateUUID())
	pn := NewWebPubSub(config)
	pn.Config.MaximumReconnectionRetries = 1
	pn.Config.WPSReconnectionPolicy = WPSLinearPolicy
	r := newReconnectionManager(pn)
	r.FailedCalls = 1
	reconnected := false
	doneReconnected := make(chan bool)
	r.HandleReconnection(func() {
		reconnected = true
		doneReconnected <- true
	})
	go r.startHeartbeatTimer()
	<-doneReconnected
	assert.True(reconnected)
	r.stopHeartbeatTimer()
}
