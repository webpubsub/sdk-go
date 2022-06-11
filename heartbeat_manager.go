package webpubsub

import (
	"fmt"
	"sync"
	"time"
)

// HeartbeatManager is a struct that assists in running of the heartbeat.
type HeartbeatManager struct {
	sync.RWMutex

	heartbeatChannels map[string]*SubscriptionItem
	heartbeatGroups   map[string]*SubscriptionItem
	webpubsub         *WebPubSub

	hbLoopMutex               sync.RWMutex
	hbTimer                   *time.Ticker
	hbDone                    chan bool
	ctx                       Context
	runIndependentOfSubscribe bool
	hbRunning                 bool
	queryParam                map[string]string
	state                     map[string]interface{}
}

func newHeartbeatManager(pn *WebPubSub, context Context) *HeartbeatManager {
	return &HeartbeatManager{
		heartbeatChannels: make(map[string]*SubscriptionItem),
		heartbeatGroups:   make(map[string]*SubscriptionItem),
		ctx:               context,
		webpubsub:         pn,
	}
}

// Destroy stops the running heartbeat.
func (m *HeartbeatManager) Destroy() {
	m.stopHeartbeat(true, true)
}

func (m *HeartbeatManager) nonIndependentHeartbeatLoop() {
	timeNow := time.Now().Unix()

	m.webpubsub.subscriptionManager.hbDataMutex.RLock()
	reqSentAt := m.webpubsub.subscriptionManager.requestSentAt
	m.webpubsub.subscriptionManager.hbDataMutex.RUnlock()

	if reqSentAt > 0 {
		timediff := int64(m.webpubsub.Config.HeartbeatInterval) - (timeNow - reqSentAt)
		m.webpubsub.Config.Log.Println(fmt.Sprintf("heartbeat timediff: %d", timediff))
		m.webpubsub.subscriptionManager.hbDataMutex.Lock()
		m.webpubsub.subscriptionManager.requestSentAt = 0
		m.webpubsub.subscriptionManager.hbDataMutex.Unlock()
		if timediff > 10 {
			m.Lock()
			m.hbTimer.Stop()
			m.Unlock()

			m.webpubsub.Config.Log.Println(fmt.Sprintf("heartbeat sleeping timediff: %d", timediff))
			waitTimer := time.NewTicker(time.Duration(timediff) * time.Second)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				waitTimerCh := waitTimer.C
				for {
					select {
					case <-m.hbDone:
						wg.Done()
						m.webpubsub.Config.Log.Println("nonIndependentHeartbeatLoop: breaking out to HeartbeatLabel")
						return
					case <-waitTimerCh:
						m.webpubsub.Config.Log.Println("nonIndependentHeartbeatLoop: waitTimerCh done")
						wg.Done()
						return
					}
				}
			}()
			wg.Wait()
			m.webpubsub.Config.Log.Println("heartbeat sleep end")

			m.Lock()
			m.hbTimer = time.NewTicker(time.Duration(m.webpubsub.Config.HeartbeatInterval) * time.Second)
			m.Unlock()
		}
	}
	m.performHeartbeatLoop()
}

func (m *HeartbeatManager) readHeartBeatTimer(runIndependentOfSubscribe bool) {
	go func() {

		defer m.hbLoopMutex.Unlock()
		defer func() {
			m.Lock()
			m.hbDone = nil
			m.Unlock()
		}()
	HeartbeatLabel:
		for {
			m.RLock()
			timerCh := m.hbTimer.C
			m.RUnlock()

			select {
			case <-timerCh:
				if runIndependentOfSubscribe {
					m.performHeartbeatLoop()
				} else {
					m.nonIndependentHeartbeatLoop()
				}
			case <-m.hbDone:
				m.webpubsub.Config.Log.Println("heartbeat: loop after stop")
				break HeartbeatLabel
			}
		}
	}()
}

func (m *HeartbeatManager) startHeartbeatTimer(runIndependentOfSubscribe bool) {
	m.RLock()
	hbRunning := m.hbRunning
	m.RUnlock()
	if hbRunning && !runIndependentOfSubscribe {
		return
	}
	m.stopHeartbeat(runIndependentOfSubscribe, true)

	m.Lock()
	m.hbRunning = true
	m.Unlock()

	m.webpubsub.Config.Log.Println("heartbeat: new timer", m.webpubsub.Config.HeartbeatInterval)
	m.webpubsub.Config.Lock()
	presenceTimeout := m.webpubsub.Config.PresenceTimeout
	heartbeatInterval := m.webpubsub.Config.HeartbeatInterval
	m.webpubsub.Config.Unlock()
	if presenceTimeout <= 0 && heartbeatInterval <= 0 {
		return
	}

	m.hbLoopMutex.Lock()
	m.Lock()
	m.hbDone = make(chan bool)
	m.hbTimer = time.NewTicker(time.Duration(m.webpubsub.Config.HeartbeatInterval) * time.Second)
	m.Unlock()

	if runIndependentOfSubscribe {
		m.performHeartbeatLoop()
	}

	m.readHeartBeatTimer(runIndependentOfSubscribe)

}

func (m *HeartbeatManager) stopHeartbeat(runIndependentOfSubscribe bool, skipRuncheck bool) {
	if !skipRuncheck {
		m.RLock()
		hbRunning := m.hbRunning
		m.RUnlock()

		if hbRunning && !runIndependentOfSubscribe {
			return
		}
	}
	m.webpubsub.Config.Log.Println("heartbeat: loop: stopping...")

	m.Lock()
	if m.hbTimer != nil {
		m.hbTimer.Stop()
		m.webpubsub.Config.Log.Println("heartbeat: loop: timer stopped")
	}

	if m.hbDone != nil {
		m.hbDone <- true
		m.webpubsub.Config.Log.Println("heartbeat: loop: done channel stopped")
	}
	m.hbRunning = false
	m.Unlock()
	m.webpubsub.subscriptionManager.hbDataMutex.Lock()
	m.webpubsub.subscriptionManager.requestSentAt = 0
	m.webpubsub.subscriptionManager.hbDataMutex.Unlock()
}

func (m *HeartbeatManager) prepareList(subItem map[string]*SubscriptionItem) []string {
	response := []string{}

	for _, v := range subItem {
		response = append(response, v.name)
	}
	return response
}

func (m *HeartbeatManager) performHeartbeatLoop() error {
	var stateStorage map[string]interface{}

	m.RLock()
	presenceChannels := m.prepareList(m.heartbeatChannels)
	presenceGroups := m.prepareList(m.heartbeatGroups)
	stateStorage = m.state
	queryParam := m.queryParam
	m.webpubsub.Config.Log.Println("performHeartbeatLoop: count presenceChannels, presenceGroups", len(presenceChannels), len(presenceGroups))
	m.RUnlock()

	if (len(presenceChannels) == 0) && (len(presenceGroups) == 0) {
		m.webpubsub.Config.Log.Println("performHeartbeatLoop: count presenceChannels, presenceGroups nil")
		presenceChannels = m.webpubsub.subscriptionManager.stateManager.prepareChannelList(false)
		presenceGroups = m.webpubsub.subscriptionManager.stateManager.prepareGroupList(false)
		stateStorage = m.webpubsub.subscriptionManager.stateManager.createStatePayload()
		queryParam = nil

		m.webpubsub.Config.Log.Println("performHeartbeatLoop: count sub presenceChannels, presenceGroups", len(presenceChannels), len(presenceGroups))
	}

	if len(presenceChannels) <= 0 && len(presenceGroups) <= 0 {
		m.webpubsub.Config.Log.Println("heartbeat: no channels left")
		go m.stopHeartbeat(true, true)
		return nil
	}

	_, status, err := newHeartbeatBuilder(m.webpubsub).
		Channels(presenceChannels).
		ChannelGroups(presenceGroups).
		State(stateStorage).
		QueryParam(queryParam).
		Execute()

	if err != nil {

		pnStatus := &WPSStatus{
			Operation: WPSHeartBeatOperation,
			Category:  WPSBadRequestCategory,
			Error:     true,
			ErrorData: err,
		}
		m.webpubsub.Config.Log.Println("performHeartbeatLoop: err", err, pnStatus)

		m.webpubsub.subscriptionManager.listenerManager.announceStatus(pnStatus)

		return err
	}

	pnStatus := &WPSStatus{
		Category:   WPSUnknownCategory,
		Error:      false,
		Operation:  WPSHeartBeatOperation,
		StatusCode: status.StatusCode,
	}
	m.webpubsub.Config.Log.Println("performHeartbeatLoop: err", err, pnStatus)

	m.webpubsub.subscriptionManager.listenerManager.announceStatus(pnStatus)

	return nil
}
