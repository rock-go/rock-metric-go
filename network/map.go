package network

import (
	netStat "github.com/shirou/gopsutil/net"
	"sync"
)

type IfcMetricsMap struct {
	mu sync.Mutex
	dirty map[string]Metric
}

type StatMetricsMap struct {
	mu sync.Mutex
	dirty map[string]netStat.IOCountersStat
}


func (im *IfcMetricsMap) load(key string) Metric {
	im.mu.Lock()
	defer im.mu.Unlock()

	v , ok := im.dirty[key]
	if !ok {
		return Metric{}
	}

	return v
}

func (im *IfcMetricsMap) store(key string , val Metric){
	im.mu.Lock()
	defer im.mu.Unlock()
	im.dirty[key] = val
}

func (sm *StatMetricsMap) last(key string , ptr *netStat.IOCountersStat) bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	v , ok := sm.dirty[key]
	if !ok {
		return  false
	}

	*ptr = v
	return true
}

func (sm *StatMetricsMap) store(key string , val netStat.IOCountersStat) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.dirty[key] = val
}