package devicemonitor

import (
	"encoding/json"
	"fmt"
	"gpu-demonstration-api/models"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func GetCPUMemoryUtilizationJSON() ([]byte, error) {
	var c models.CPUStatus
	pct, err := GetCPUPercent()
	if err != nil {
		return nil, err
	}
	mem, err := GetVirtualMemory()
	if err != nil {
		return nil, err
	}
	c.CPUPercent = pct
	c.Memory = mem

	cpuBytes, err := json.Marshal(c)
	if err != nil {
		fmt.Println("devicemonitor.GetCPUMemoryUtilizationJSON: " + err.Error())
		return cpuBytes, err
	}
	return cpuBytes, nil
}

func GetCPUInfoJSON() ([]byte, error) {
	cpu, err := cpu.Info()
	if err != nil {
		fmt.Println("devicemonitor.GetCPUInfoJSON: " + err.Error())
		return nil, err
	}
	cpuBytes, err := json.Marshal(cpu)
	if err != nil {
		fmt.Println("devicemonitor.GetCPUInfoJSON: " + err.Error())
		return cpuBytes, err
	}
	return cpuBytes, nil
}

func GetCPUPercent() ([]float64, error) {
	var dur time.Duration
	percent, err := cpu.Percent(dur, false)
	if err != nil {
		fmt.Println("devicemonitor.GetCPUPercent: " + err.Error())
		return nil, err
	}
	return percent, nil
}

func GetVirtualMemory() (*mem.VirtualMemoryStat, error) {
	mem, err := mem.VirtualMemory()
	if err != nil {
		fmt.Println("devicemonitor.GetVirtualMemory: " + err.Error())
		return mem, err
	}

	return mem, nil
}
