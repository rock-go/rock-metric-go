package cpu

import (
	"github.com/elastic/gosigar"
)

type CPU struct {
	Architecture string
	CoreNum      int

	User   float64
	System float64
	Idle   float64
	// linux
	IOWait  float64
	IRQ     float64
	Nice    float64
	SoftIRQ float64
	Stolen  float64

	Total float64
}

//var Sample Metric // 缓存每次获取到的数据

func Get(m *Metric) *CPU {
	cpu := gosigar.Cpu{}
	err := cpu.Get()
	if err != nil {
		return nil
	}

	m.lastSample = m.nowSample
	m.nowSample = &cpu
	return m.CalcNormalPct()
}

func (cpu *CPU) DisableReflect() {}
