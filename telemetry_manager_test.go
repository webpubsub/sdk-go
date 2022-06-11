package webpubsub

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAverageLatency(t *testing.T) {
	assert := assert.New(t)

	endpointLatencies := []LatencyEntry{
		LatencyEntry{D: int64(100), L: float64(10)},
		LatencyEntry{D: int64(100), L: float64(20)},
		LatencyEntry{D: int64(100), L: float64(30)},
		LatencyEntry{D: int64(100), L: float64(40)},
		LatencyEntry{D: int64(100), L: float64(50)}}

	averageLatency := averageLatencyFromData(endpointLatencies)
	assert.Equal(float64(30), averageLatency)
}

func TestCleanUp(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := contextWithCancel(backgroundContext)
	manager := newTelemetryManager(1, ctx)

	for i := 0; i < 10; i++ {
		manager.StoreLatency(float64(i), WPSPublishOperation)
	}

	// await for store timestamp expired
	time.Sleep(2 * time.Second)

	manager.CleanUpTelemetryData()

	assert.Equal(0, len(manager.OperationLatency()))
}

func TestValidQueries(t *testing.T) {
	assert := assert.New(t)
	ctx, _ := contextWithCancel(backgroundContext)
	manager := newTelemetryManager(60, ctx)

	manager.StoreLatency(float64(1), WPSPublishOperation)
	manager.StoreLatency(float64(2), WPSPublishOperation)
	manager.StoreLatency(float64(3), WPSPublishOperation)

	manager.StoreLatency(float64(4), WPSHistoryOperation)
	manager.StoreLatency(float64(5), WPSHistoryOperation)
	manager.StoreLatency(float64(6), WPSHistoryOperation)

	manager.StoreLatency(float64(7), WPSRemoveGroupOperation)
	manager.StoreLatency(float64(8), WPSRemoveGroupOperation)
	manager.StoreLatency(float64(9), WPSRemoveGroupOperation)

	queries := manager.OperationLatency()

	assert.Equal("2", queries["l_pub"])
	assert.Equal("5", queries["l_hist"])
	assert.Equal("8", queries["l_cg"])
}
