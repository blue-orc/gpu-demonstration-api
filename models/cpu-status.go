package models

import (
	"github.com/shirou/gopsutil/mem"
)

type CPUStatus struct {
	Memory     *mem.VirtualMemoryStat
	CPUPercent []float64
}
