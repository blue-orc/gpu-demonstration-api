package devicemonitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NVIDIA/gpu-monitoring-tools/bindings/go/nvml"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	DMONHEADER = `
				# gpu   pwr  temp    sm   mem   enc   dec  mclk  pclk
				# Idx     W     C     %     %     %     %   MHz   MHz
			`
)

const (
	DEVICEINFO = `
				UUID           : {{.UUID}}
				Model          : {{or .Model "N/A"}}
				Path           : {{.Path}}
				Power          : {{if .Power}}{{.Power}} W{{else}}N/A{{end}}
				Memory         : {{if .Memory}}{{.Memory}} MiB{{else}}N/A{{end}}
				CudaComputeCap : {{if .CudaComputeCapability.Major}}{{.CudaComputeCapability.Major}}.{{.CudaComputeCapability.Minor}}{{else}}N/A{{end}}
				CPU Affinity   : {{if .CPUAffinity}}NUMA node{{.CPUAffinity}}{{else}}N/A{{end}}
				Bus ID         : {{.PCI.BusID}}
				BAR1           : {{if .PCI.BAR1}}{{.PCI.BAR1}} MiB{{else}}N/A{{end}}
				Bandwidth      : {{if .PCI.Bandwidth}}{{.PCI.Bandwidth}} MB/s{{else}}N/A{{end}}
				Cores          : {{if .Clocks.Cores}}{{.Clocks.Cores}} MHz{{else}}N/A{{end}}
				Memory         : {{if .Clocks.Memory}}{{.Clocks.Memory}} MHz{{else}}N/A{{end}}
				P2P Available  : {{if not .Topology}}None{{else}}{{range .Topology}}
									{{.BusID}} - {{(.Link.String)}}{{end}}{{end}}
				---------------------------------------------------------------------
				`
)

var CurrentStatus *nvml.DeviceStatus

func GetCurrentStatusJSON() ([]byte, error) {
	csBytes, err := json.Marshal(CurrentStatus)
	if err != nil {
		return csBytes, err
	}
	return csBytes, nil
}

func GetGPUInfo() ([]*nvml.Device, error) {
	var devices []*nvml.Device
	count, err := nvml.GetDeviceCount()
	if err != nil {
		fmt.Println("Error getting device count:", err)
		return devices, errors.New(fmt.Sprintf("Error getting device count:", err))
	}

	for i := uint(0); i < count; i++ {
		device, err := nvml.NewDevice(i)
		if err != nil {
			fmt.Println("Error getting device %d: %v\n", i, err)
			return devices, errors.New(fmt.Sprintf("Error getting device %d: %v\n", i, err))
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func Init() {
	nvml.Init()
	defer nvml.Shutdown()

	count, err := nvml.GetDeviceCount()
	if err != nil {
		log.Panicln("Error getting device count:", err)
	}

	var devices []*nvml.Device
	for i := uint(0); i < count; i++ {
		device, err := nvml.NewDevice(i)
		if err != nil {
			log.Panicf("Error getting device %d: %v\n", i, err)
		}
		devices = append(devices, device)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for i, device := range devices {
				st, err := device.Status()
				if err != nil {
					log.Panicf("Error getting device %d status: %v\n", i, err)
				}
				CurrentStatus = st
				//fmt.Printf("%5d %5d %5d %5d %5d %5d %5d %5d %5d\n",
				//	i, *st.Power, *st.Temperature, *st.Utilization.GPU, *st.Utilization.Memory,
				//	*st.Utilization.Encoder, *st.Utilization.Decoder, *st.Clocks.Memory, *st.Clocks.Cores)
			}
		case <-sigs:
			return
		}
	}
}
