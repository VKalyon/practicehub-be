package video

import (
	"fmt"
	"runtime"

	"encore.dev/metrics"
	"encore.dev/rlog"
)

var MemoryUsage = metrics.NewGauge[float64]("memory_usage_in_bytes", metrics.GaugeConfig{})

func measureMemory() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	MemoryUsage.Set(float64(m.Alloc))
	mem := fmt.Sprintf("Current memory usage is: %d bytes\n", m.Alloc)
	rlog.Info(mem)
}
