package webpubsub

import (
	"fmt"
	"sync"
	"time"
)

const timestampDivider = 1000

const cleanUpInterval = 1
const cleanUpIntervalMultiplier = 1000

// LatencyEntry is the struct to store the timestamp and latency values.
type LatencyEntry struct {
	D int64
	L float64
}

// Operations is the struct to store the latency values of different operations.
type Operations struct {
	latencies []LatencyEntry
}

// TelemetryManager is the struct to store the Telemetry details.
type TelemetryManager struct {
	sync.RWMutex

	operations map[string][]LatencyEntry

	ctx Context

	cleanUpTimer *time.Ticker

	maxLatencyDataAge int
	IsRunning         bool
}

func newTelemetryManager(maxLatencyDataAge int, ctx Context) *TelemetryManager {
	manager := &TelemetryManager{
		maxLatencyDataAge: maxLatencyDataAge,
		operations:        make(map[string][]LatencyEntry),
		ctx:               ctx,
	}

	go manager.startCleanUpTimer()

	return manager
}

// OperationLatency returns a map of the stored latencies by operation.
func (m *TelemetryManager) OperationLatency() map[string]string {
	operationLatencies := make(map[string]string)

	//var ops map[string][]LatencyEntry
	m.RLock()

	for endpointName := range m.operations {
		queryKey := fmt.Sprintf("l_%s", endpointName)

		endpointAverageLatency := averageLatencyFromData(
			m.operations[endpointName])

		if endpointAverageLatency > 0 {
			operationLatencies[queryKey] = fmt.Sprint(endpointAverageLatency)
		}
	}
	m.RUnlock()

	return operationLatencies
}

// StoreLatency stores the latency values of the different operations.
func (m *TelemetryManager) StoreLatency(latency float64, t OperationType) {
	if latency > float64(0) && t != WPSSubscribeOperation {
		endpointName := telemetryEndpointNameForOperation(t)

		storeTimestamp := time.Now().Unix()

		m.Lock()
		m.operations[endpointName] = append(m.operations[endpointName], LatencyEntry{
			D: storeTimestamp,
			L: latency,
		})
		m.Unlock()
	}
}

// CleanUpTelemetryData cleans up telemetry data of all operations.
func (m *TelemetryManager) CleanUpTelemetryData() {
	currentTimestamp := time.Now().Unix()

	m.Lock()
	for endpoint, latencies := range m.operations {
		index := 0

		for _, latency := range latencies {
			if currentTimestamp-latency.D > int64(m.maxLatencyDataAge) {
				m.operations[endpoint] = append(m.operations[endpoint][:index],
					m.operations[endpoint][index+1:]...)
				continue
			}
			index++
		}

		if len(m.operations[endpoint]) == 0 {
			delete(m.operations, endpoint)
		}
	}
	m.ctx.Done()
	m.Unlock()
}

func (m *TelemetryManager) startCleanUpTimer() {
	m.cleanUpTimer = time.NewTicker(
		time.Duration(
			cleanUpInterval*cleanUpIntervalMultiplier) * time.Millisecond)

	go func() {
	CleanUpTimerLabel:
		for {
			timerCh := m.cleanUpTimer.C

			select {
			case <-timerCh:
				m.CleanUpTelemetryData()
			case <-m.ctx.Done():
				m.cleanUpTimer.Stop()
				break CleanUpTimerLabel
			}
		}
	}()
}

func telemetryEndpointNameForOperation(t OperationType) string {
	var endpoint string

	switch t {
	case WPSPublishOperation:
		endpoint = "pub"
		break
	case WPSMessageCountsOperation:
		endpoint = "mc"
		break
	case WPSHistoryOperation:
		fallthrough
	case WPSFetchMessagesOperation:
		fallthrough
	case WPSDeleteMessagesOperation:
		endpoint = "hist"
		break
	case WPSUnsubscribeOperation:
		fallthrough
	case WPSWhereNowOperation:
		fallthrough
	case WPSHereNowOperation:
		fallthrough
	case WPSHeartBeatOperation:
		fallthrough
	case WPSSetStateOperation:
		fallthrough
	case WPSGetStateOperation:
		endpoint = "pres"
		break
	case WPSAddChannelsToChannelGroupOperation:
		fallthrough
	case WPSRemoveChannelFromChannelGroupOperation:
		fallthrough
	case WPSChannelsForGroupOperation:
		fallthrough
	case WPSRemoveGroupOperation:
		endpoint = "cg"
		break
	case WPSAccessManagerRevoke:
		fallthrough
	case WPSAccessManagerGrant:
		endpoint = "pam"
		break
	case WPSAccessManagerGrantToken:
		fallthrough
	case WPSAccessManagerRevokeToken:
		endpoint = "pamv3"
		break
	case WPSSignalOperation:
		endpoint = "sig"
		break
	case WPSGetMessageActionsOperation:
		fallthrough
	case WPSAddMessageActionsOperation:
		fallthrough
	case WPSRemoveMessageActionsOperation:
		endpoint = "msga"
		break
	case WPSHistoryWithActionsOperation:
		endpoint = "hist"
		break
	case WPSCreateUserOperation:
		fallthrough
	case WPSGetUsersOperation:
		fallthrough
	case WPSGetUserOperation:
		fallthrough
	case WPSUpdateUserOperation:
		fallthrough
	case WPSDeleteUserOperation:
		fallthrough
	case WPSGetSpaceOperation:
		fallthrough
	case WPSGetSpacesOperation:
		fallthrough
	case WPSCreateSpaceOperation:
		fallthrough
	case WPSDeleteSpaceOperation:
		fallthrough
	case WPSUpdateSpaceOperation:
		fallthrough
	case WPSGetMembershipsOperation:
		fallthrough
	case WPSGetChannelMembersOperation:
		fallthrough
	case WPSManageMembershipsOperation:
		fallthrough
	case WPSManageMembersOperation:
		fallthrough
	case WPSSetChannelMembersOperation:
		fallthrough
	case WPSSetMembershipsOperation:
		fallthrough
	case WPSRemoveChannelMetadataOperation:
		fallthrough
	case WPSRemoveUUIDMetadataOperation:
		fallthrough
	case WPSGetAllChannelMetadataOperation:
		fallthrough
	case WPSGetAllUUIDMetadataOperation:
		fallthrough
	case WPSGetUUIDMetadataOperation:
		fallthrough
	case WPSRemoveMembershipsOperation:
		fallthrough
	case WPSRemoveChannelMembersOperation:
		fallthrough
	case WPSSetUUIDMetadataOperation:
		fallthrough
	case WPSGetChannelMetadataOperation:
		fallthrough
	case WPSSetChannelMetadataOperation:
		endpoint = "obj"
		break
	case WPSDeleteFileOperation:
		fallthrough
	case WPSDownloadFileOperation:
		fallthrough
	case WPSGetFileURLOperation:
		fallthrough
	case WPSListFilesOperation:
		fallthrough
	case WPSSendFileOperation:
		fallthrough
	case WPSSendFileToS3Operation:
		fallthrough
	case WPSPublishFileMessageOperation:
		endpoint = "file"
		break
	default:
		endpoint = "time"
		break
	}

	return endpoint
}

func averageLatencyFromData(endpointLatencies []LatencyEntry) float64 {
	var totalLatency float64

	for _, latency := range endpointLatencies {
		totalLatency += latency.L
	}

	return totalLatency / float64(len(endpointLatencies))
}
