package main

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

func getCpusUsage() []float64 {
	prec, err := cpu.Percent(time.Second, false)

	if err != nil {
		log.Println("Error getting CPU usage:", err)
		return []float64{}
	}

	if len(prec) == 0 {
		return []float64{}
	}

	return prec
}

func getRamUsage() (int, int) {
	memory, err := mem.VirtualMemory()

	if err != nil {
		log.Println("Error getting RAM usage:", err)
		return -1, -1
	}

	return int(memory.Used), int(memory.Total)
}

func getDiskUsage() (int, int) {
	usageStat, err := disk.Usage("/")

	if err != nil {
		log.Println("Error getting disk usage:", err)
		return -1, -1
	}

	return int(usageStat.Used), int(usageStat.Total)
}

func getNetworkUsage() (int, int) {
	counters, err := net.IOCounters(false)

	if err != nil {
		log.Println("Error getting network stats:", err)
		return -1, -1
	}

	if len(counters) == 0 || len(counters) == 0 {
		return -1, -1
	}

	return int(counters[0].BytesSent), int(counters[0].BytesRecv)
}

func getMetric() *Metric {
	useMemory, totalMemory := getRamUsage()
	useDisk, totalDisk := getDiskUsage()
	networkSent, networkReceived := getNetworkUsage()

	return &Metric{
		Cpus:     getCpusUsage(),
		UseRam:    useMemory,
		TotalRam: totalMemory,
		UseDisks:  []int{useDisk},
		TotalDisks: []int{totalDisk},
		NetworkSend:    networkSent,
		NetworkReceive: networkReceived,
		Time:    time.Now(),
	}
}