package cpu

import (
	"github.com/elastic/gosigar"
	"runtime"
)

type Metric struct {
	lastSample *gosigar.Cpu
	nowSample  *gosigar.Cpu
}

func (m *Metric) CalcNormalPct() *CPU {
	return calCpuPct(m.lastSample, m.nowSample, 1)
}

// 计算百分比
func calCpuPct(s1, s2 *gosigar.Cpu, numCore int) *CPU {
	if s1 == nil || s2 == nil {
		return nil
	}

	timeDelta := s2.Total() - s1.Total()
	if timeDelta <= 0 {
		return nil
	}

	stat := &CPU{
		Architecture: runtime.GOARCH,
		CoreNum:      runtime.NumCPU(),
		User:         calcPct(s1.User, s2.User, timeDelta, numCore),
		System:       calcPct(s1.Sys, s2.Sys, timeDelta, numCore),
		Idle:         calcPct(s1.Idle, s2.Idle, timeDelta, numCore),
		IOWait:       calcPct(s1.Wait, s2.Wait, timeDelta, numCore),
		IRQ:          calcPct(s1.Irq, s2.Irq, timeDelta, numCore),
		Nice:         calcPct(s1.Nice, s2.Nice, timeDelta, numCore),
		SoftIRQ:      calcPct(s1.SoftIrq, s2.SoftIrq, timeDelta, numCore),
		Stolen:       calcPct(s1.Stolen, s2.Stolen, timeDelta, numCore),
		Total:        calcTotalPct(s1, s2, timeDelta, numCore),
	}

	return stat
}

func calcPct(v1, v2 uint64, timeDelta uint64, numCore int) float64 {
	cpuDelta := int64(v2 - v1)
	pct := float64(cpuDelta) / float64(timeDelta)
	return pct * float64(numCore)
}

func calcTotalPct(s1, s2 *gosigar.Cpu, timeDelta uint64, numCore int) float64 {
	idle := calcPct(s1.Idle, s2.Idle, timeDelta, numCore) + calcPct(s1.Wait, s2.Wait, timeDelta, numCore)
	return float64(numCore) - idle
}
