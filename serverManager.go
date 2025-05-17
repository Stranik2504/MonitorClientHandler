package ClientHandler

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// getCpusUsage возвращает загрузку процессоров в процентах.
//
// @return срез float64 с процентом загрузки каждого CPU
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

// getRamUsage возвращает используемую и общую оперативную память.
//
// @return используемая память (int), общая память (int)
func getRamUsage() (int, int) {
	memory, err := mem.VirtualMemory()

	if err != nil {
		log.Println("Error getting RAM usage:", err)
		return -1, -1
	}

	return int(memory.Used), int(memory.Total)
}

// getDiskUsage возвращает используемое и общее дисковое пространство.
//
// @return используемое место (int), общий объём (int)
func getDiskUsage() (int, int) {
	usageStat, err := disk.Usage("/")

	if err != nil {
		log.Println("Error getting disk usage:", err)
		return -1, -1
	}

	return int(usageStat.Used), int(usageStat.Total)
}

// getNetworkUsage возвращает количество отправленных и полученных байт по сети.
//
// @return отправлено байт (int), получено байт (int)
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

// getMetric собирает и возвращает метрики системы.
//
// @return указатель на структуру Metric с актуальными данными
func getMetric() *Metric {
	useMemory, totalMemory := getRamUsage()
	useDisk, totalDisk := getDiskUsage()
	networkSent, networkReceived := getNetworkUsage()

	return &Metric{
		Cpus:           getCpusUsage(),
		UseRam:         useMemory,
		TotalRam:       totalMemory,
		UseDisks:       []int{useDisk},
		TotalDisks:     []int{totalDisk},
		NetworkSend:    networkSent,
		NetworkReceive: networkReceived,
		Time:           time.Now(),
	}
}
