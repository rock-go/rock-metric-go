package cpu

import (
	"github.com/elastic/gosigar"
	"github.com/rock-go/rock/json"
	"github.com/rock-go/rock/lua"
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

func (v *CPU) Byte() []byte {
	buf := json.NewBuffer()
	buf.Tab("")

	buf.KV("arch"     , v.Architecture)
	buf.KV("core_num" , v.CoreNum)
	buf.KV("user"     , v.User)
	buf.KV("system"   , v.System)
	buf.KV("Idle"     , v.Idle)
	buf.KV("io_wait"  , v.IOWait)
	buf.KV("io_wait"  , v.IOWait)
	buf.KV("io_wait"  , v.IOWait)
	buf.KV("irq"      , v.IRQ)
	buf.KV("nice"     , v.Nice)
	buf.KV("softirq"  , v.SoftIRQ)
	buf.KV("stolen"   , v.Stolen)
	buf.KV("total"    , v.Total)
	buf.End("}")

	return  buf.Bytes()
}

func (v *CPU) String() string {
	return lua.B2S(v.Byte())
}