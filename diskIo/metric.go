package diskIo

import (
	"github.com/elastic/gosigar"
	"github.com/rock-go/rock/logger"
	"runtime"
)

type Metric struct {
	lastSample Sample
	nowSample  Sample
}

// Sample 某次获取到的数据
type Sample struct {
	DiskStats map[string]*Disk
	cpu       *gosigar.Cpu
}

var _SC_CLK_TCK float64 = 100

func init() {
	if runtime.GOOS == "windows" {
		_SC_CLK_TCK = 1000
	}
}

func (m *Metric) getMetric(stats map[string]*Disk) error {
	var sample Sample
	var err error

	cpu := gosigar.Cpu{}
	err = cpu.Get()
	if err != nil {
		return err
	}

	sample.DiskStats = stats
	sample.cpu = &cpu

	m.lastSample = m.nowSample
	m.nowSample = sample

	return nil
}

func (m *Metric) calMetric() detail {
	if m.lastSample.DiskStats == nil || m.nowSample.DiskStats == nil {
		return nil
	}

	delta := 1000.0 * float64(m.nowSample.cpu.Total()-m.lastSample.cpu.Total()) / float64(runtime.NumCPU()) / _SC_CLK_TCK
	if delta <= 0 {
		logger.Errorf("calculate error")
		return nil
	}

	var d detail
	for name, nowStat := range m.nowSample.DiskStats {
		lastStat, ok := m.lastSample.DiskStats[name]
		if !ok {
			logger.Errorf("get %s disk last stat error", name)
			continue
		}
		d = append(d, *calStat(lastStat, nowStat, delta))
	}

	return d
}

func calStat(last *Disk, now *Disk, delta float64) *Disk {
	now.ReadPerSecBytes = float64(now.ReadBytes-last.ReadBytes) * 1000.0 / delta
	now.WritePerSecByte = float64(now.WriteBytes-last.WriteBytes) * 1000.0 / delta
	return now
}
