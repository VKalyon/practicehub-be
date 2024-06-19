package video

import (
	"fmt"
	"runtime"
	"time"

	"encore.dev/metrics"
	"encore.dev/rlog"
)

var MemoryUsage = metrics.NewGauge[float64]("memory_usage_in_bytes", metrics.GaugeConfig{})
var RequestResponseTime = metrics.NewGauge[float64]("request_response_time_in_ms", metrics.GaugeConfig{})

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

// call this function with a defer statement
func measureResponseTime() func() {
	start := time.Now()
	return func() {
		deltaTime := time.Since(start).Milliseconds()
		RequestResponseTime.Set(float64(deltaTime))
	}
}
