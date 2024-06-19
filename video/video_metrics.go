package video

import (
	"fmt"
	"runtime"
	"time"

	"encore.dev/metrics"
	"encore.dev/rlog"
)

type RequestLabel struct {
	RequestType     string
	RequestResource string
}

var MemoryUsage = metrics.NewGauge[float64]("server_memory_usage_bytes", metrics.GaugeConfig{})
var RequestResponseTime = metrics.NewCounterGroup[RequestLabel, float64]("http_request_duration_ms_total", metrics.CounterConfig{})
var Requests = metrics.NewCounterGroup[RequestLabel, uint64]("http_requests_total", metrics.CounterConfig{})

func measureMemory() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		MemoryUsage.Set(float64(m.Alloc))
		mem := fmt.Sprintf("Current memory usage is: %d bytes\n", m.Alloc)
		rlog.Info(mem)
	}
}

func startResponseTime() time.Time {
	return time.Now()
}

// call this function with a defer statement and pass startResponseTime to it
func measureResponseTime(requestLabel RequestLabel, startTime time.Time) {
	deltaTime := time.Since(startTime).Milliseconds()
	RequestResponseTime.With(requestLabel).Add(float64(deltaTime))
	Requests.With(requestLabel).Add(1)
}
