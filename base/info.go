package base

import (
	"github.com/elastic/gosigar"
	"github.com/rock-go/rock-metric-go/cpu"
	"github.com/rock-go/rock-metric-go/fileSystem"
	"github.com/rock-go/rock-metric-go/network"
	"github.com/rock-go/rock/json"
	"github.com/rock-go/rock/lua"
	"github.com/shirou/gopsutil/host"
	"runtime"
)

type BasicInfo struct {
	Inet      string  `json:"inet"`
	Inet6     string  `json:"inet6"`
	Mac       string  `json:"mac"`
	Arch      string  `json:"arch"`
	Platform  string  `json:"platform"`
	MemTotal  uint64  `json:"mem_total"`
	MemFree   uint64  `json:"mem_free"`
	SwapTotal uint64  `json:"swap_total"`
	SwapFree  uint64  `json:"swap_free"`
	CpuCore   int     `json:"cpu_core"`
	CpuUsage  float64 `json:"cpu_usage"`
	DiskTotal uint64  `json:"disk_total"`
	DiskPath  string  `json:"disk_path"`
	DiskFree  uint64  `json:"disk_free"`
}

var CpuSample cpu.Metric

func Get(target string) (*BasicInfo, error) {
	var info BasicInfo

	// 网卡信息
	ifc, err := network.GetBase(target)
	if err == nil {
		info.Inet = ifc.Inet
		info.Inet6 = ifc.Inet6
		info.Mac = ifc.Mac
	}

	// 操作系统
	platform, _, version, _ := host.PlatformInformation()
	info.Arch = runtime.GOARCH
	info.Platform = platform + " " + version

	// memory
	mem := gosigar.Mem{}
	err = mem.Get()
	if err == nil {
		info.MemTotal = mem.Total
		info.MemFree = mem.Free
	}

	// swap
	swap := gosigar.Swap{}
	err = swap.Get()
	if err == nil {
		info.SwapTotal = swap.Total
		info.SwapFree = swap.Free
	}

	// cpu
	info.CpuCore = runtime.NumCPU()
	cpuInfo := cpu.Get(&CpuSample)
	if cpuInfo != nil {
		info.CpuUsage = cpuInfo.Total
	}

	// disk usage
	path, total, free, _ := fileSystem.GetMax()
	info.DiskPath = path
	info.DiskTotal = total
	info.DiskFree = free

	return &info, nil
}

func (b *BasicInfo) Byte() []byte {
	buf := json.NewBuffer()
	buf.Tab("")
	buf.KV("inet" , b.Inet)
	buf.KV("inet6" , b.Inet6)
	buf.KV("mac" , b.Mac)
	buf.KV("inet"      ,b.Inet)
	buf.KV("inet6"     ,b.Inet6)
	buf.KV("mac"       ,b.Mac)
	buf.KV("arch"      ,b.Arch)
	buf.KV("platform"  ,b.Platform)
	buf.KL("mem_total" ,int64(b.MemTotal))
	buf.KL("mem_free"  ,int64(b.MemFree))
	buf.KL("swap_total",int64(b.SwapTotal))
	buf.KL("swap_free" ,int64(b.SwapFree))
	buf.KI("cpu_core"  ,b.CpuCore)
	buf.KL("cpu_usage" ,int64(b.CpuUsage))
	buf.KL("disk_total",int64(b.DiskTotal))
	buf.KV("disk_path" ,b.DiskPath)
	buf.KL("disk_free" ,int64(b.DiskFree))
	buf.End("}")

	return buf.Bytes()
}

func (b *BasicInfo) String() string {
	return lua.B2S(b.Byte())
}
