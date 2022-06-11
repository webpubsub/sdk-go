package webpubsub

import (
	"fmt"
	"math"
	"sync"
	"time"
)

const (
	reconnectionInterval              = 10
	reconnectionMinExponentialBackoff = 1
	reconnectionMaxExponentialBackoff = 32
)

// ReconnectionManager is used to store the properties required in running the Reconnection Manager.
type ReconnectionManager struct {
	sync.RWMutex

	timerMutex sync.RWMutex

	ExponentialMultiplier       int
	FailedCalls                 int
	Milliseconds                int
	OnReconnection              func()
	OnMaxReconnectionExhaustion func()
	DoneTimer                   chan bool
	hbRunning                   bool
	webpubsub                   *WebPubSub
	exitReconnectionManager     chan bool
}

func newReconnectionManager(webpubsub *WebPubSub) *ReconnectionManager {
	manager := &ReconnectionManager{}

	manager.webpubsub = webpubsub

	manager.ExponentialMultiplier = 1
	manager.FailedCalls = 0
	manager.Milliseconds = 1000
	manager.exitReconnectionManager = make(chan bool)
	manager.hbRunning = false

	return manager
}

// HandleReconnection sets the handler that will be called when the network reconnects after a disconnect.
func (m *ReconnectionManager) HandleReconnection(handler func()) {
	m.Lock()
	m.OnReconnection = handler
	m.Unlock()
}

// HandleOnMaxReconnectionExhaustion sets the handler that will be called when the max reconnection attempts are exhausted.
func (m *ReconnectionManager) HandleOnMaxReconnectionExhaustion(handler func()) {
	m.Lock()
	m.OnMaxReconnectionExhaustion = handler
	m.Unlock()
}

func (m *ReconnectionManager) startPolling() {

	if m.webpubsub.Config.WPSReconnectionPolicy == WPSNonePolicy {
		m.webpubsub.Config.Log.Println("Reconnection policy is disabled, please handle reconnection manually.")
		return
	}

	m.Lock()
	m.ExponentialMultiplier = 1
	m.FailedCalls = 0
	hbRunning := m.hbRunning
	m.Unlock()

	if !hbRunning {
		m.webpubsub.Config.Log.Println(fmt.Sprintf("Reconnection policy: %d, retries: %d", m.webpubsub.Config.WPSReconnectionPolicy, m.webpubsub.Config.MaximumReconnectionRetries))

		m.startHeartbeatTimer()
	} else {
		m.webpubsub.Config.Log.Println("hb already running")
	}

}

func (m *ReconnectionManager) startHeartbeatTimer() {

	timerInterval := reconnectionInterval

	for {

		m.Lock()
		m.hbRunning = true
		failedCalls := m.FailedCalls
		m.Unlock()
		_, status, err := m.webpubsub.Time().Execute()
		if status.Error == nil {
			if failedCalls > 0 {
				timerInterval = reconnectionInterval
				m.Lock()
				m.FailedCalls = 0
				m.Unlock()
				m.webpubsub.Config.Log.Println(fmt.Sprintf("Network reconnected"))
				m.OnReconnection()
			}
		} else {
			if m.webpubsub.Config.WPSReconnectionPolicy == WPSExponentialPolicycy {
				timerInterval = m.getExponentialInterval()
			}
			m.Lock()
			m.FailedCalls++
			m.webpubsub.Config.Log.Println(fmt.Sprintf("Network disconnected, reconnection try %d of %d\n %v %v", m.FailedCalls, m.webpubsub.Config.MaximumReconnectionRetries, status, err))
			m.ExponentialMultiplier++

			failedCalls := m.FailedCalls
			retries := m.webpubsub.Config.MaximumReconnectionRetries
			m.Unlock()
			if retries != -1 && failedCalls >= retries {
				m.webpubsub.Config.Log.Printf(fmt.Sprintf("Network connection retry limit (%d) exceeded", retries))
				m.Lock()
				m.hbRunning = false
				m.Unlock()
				m.OnMaxReconnectionExhaustion()
				return
			}
		}

		select {
		case <-time.After(time.Duration(timerInterval) * time.Second):
		case <-m.webpubsub.ctx.Done():
			m.webpubsub.Config.Log.Printf(fmt.Sprintf("webpubsub.ctx.Done\n"))
			m.Lock()
			m.hbRunning = false
			m.Unlock()
			return
		case <-m.exitReconnectionManager:
			m.webpubsub.Config.Log.Printf(fmt.Sprintf("exitReconnectionManager\n"))
			return
		}
	}
}

func (m *ReconnectionManager) getExponentialInterval() int {
	timerInterval := int(math.Pow(2, float64(m.ExponentialMultiplier)) - 1)
	if timerInterval > reconnectionMaxExponentialBackoff {
		timerInterval = reconnectionMinExponentialBackoff

		m.Lock()
		m.ExponentialMultiplier = 1
		m.webpubsub.Config.Log.Printf(fmt.Sprintf("timerInterval > MaxExponentialBackoff at: %d\n", m.ExponentialMultiplier))
		m.Unlock()

	} else if timerInterval < 1 {
		timerInterval = reconnectionMinExponentialBackoff
		m.Lock()
		m.ExponentialMultiplier = 1
		m.webpubsub.Config.Log.Printf(fmt.Sprintf("timerInterval < 1 at: %d\n", m.ExponentialMultiplier))
		m.Unlock()
	}
	return timerInterval
}

func (m *ReconnectionManager) stopHeartbeatTimer() {
	m.webpubsub.Config.Log.Printf("stopHeartbeatTimer")
	m.Lock()
	if m.hbRunning {
		m.hbRunning = false
		m.exitReconnectionManager <- true
	}
	m.Unlock()
	m.webpubsub.Config.Log.Printf("stopHeartbeatTimer true")
}
