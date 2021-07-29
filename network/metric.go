package network

import (
	"github.com/shirou/gopsutil/net"
	"time"
)

type Metric struct {
	lastSample Sample
	nowSample  Sample
}

// Sample 单次获取的样本数据
type Sample struct {
	timeGet int64 // 获取本次数据的时间
	net.IOCountersStat
}

type Flow struct {
	InBytes         uint64
	InPackets       uint64
	InError         uint64
	InDropped       uint64
	InBytesPerSec   float64
	InPacketsPerSec float64

	OutBytes         uint64
	OutPackets       uint64
	OutError         uint64
	OutDropped       uint64
	OutBytesPerSec   float64
	OutPacketsPerSec float64
}

func (m *Metric) getMetric(stat net.IOCountersStat) {
	var sample Sample
	sample.timeGet = time.Now().Unix()
	sample.IOCountersStat = stat

	m.lastSample = m.nowSample
	m.nowSample = sample
}

func (m *Metric) calMetric() Flow {
	flow := Flow{
		InBytes:          m.nowSample.BytesRecv,
		InPackets:        m.nowSample.PacketsRecv,
		InError:          m.nowSample.Errin,
		InDropped:        m.nowSample.Dropin,
		InBytesPerSec:    0,
		InPacketsPerSec:  0,
		OutBytes:         m.nowSample.BytesSent,
		OutPackets:       m.nowSample.PacketsSent,
		OutError:         m.nowSample.Errout,
		OutDropped:       m.nowSample.Dropout,
		OutBytesPerSec:   0,
		OutPacketsPerSec: 0,
	}

	if m.lastSample.timeGet == 0 {
		return flow
	}

	delta := float64(m.nowSample.timeGet - m.lastSample.timeGet)
	flow.InBytesPerSec = float64(m.nowSample.BytesRecv-m.lastSample.BytesRecv) / delta
	flow.InPacketsPerSec = float64(m.nowSample.PacketsRecv-m.lastSample.PacketsRecv) / delta
	flow.OutBytesPerSec = float64(m.nowSample.BytesSent-m.lastSample.BytesSent) / delta
	flow.OutPacketsPerSec = float64(m.nowSample.PacketsSent-m.lastSample.PacketsSent) / delta

	return flow
}
