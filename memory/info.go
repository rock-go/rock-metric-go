package memory

import (
	sysInfo "github.com/elastic/go-sysinfo"
	sysInfoTypes "github.com/elastic/go-sysinfo/types"
	"github.com/elastic/gosigar"
	"github.com/rock-go/rock/json"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"runtime"
)

type Memory struct {
	Total        uint64  `json:"total"`
	Free         uint64  `json:"free"`
	UsedPct      float64 `json:"used_pct"`
	SwapTotal    uint64  `json:"swap_total"`
	SwapFree     uint64  `json:"swap_free"`
	SwapUsedPct  float64 `json:"swap_used_pct"`
	SwapInPages  uint64  `json:"swap_in_pages"`
	SwapOutPages uint64  `json:"swap_out_pages"`
}

func GetMemDetail() (*Memory, error) {
	var mem Memory

	memStat := gosigar.Mem{}
	err := memStat.Get()
	if err != nil {
		logger.Errorf("get memory stats error: %v", err)
		return &mem, err
	}

	mem.Total = memStat.Total
	mem.Free = memStat.Free
	mem.UsedPct = float64(memStat.Total-memStat.Free) / float64(memStat.Total)

	swapStat := gosigar.Swap{}
	err = swapStat.Get()
	if err != nil {
		logger.Errorf("get swap stats error: %v", err)
		return &mem, err
	}

	mem.SwapTotal = memStat.Total
	mem.SwapFree = memStat.Free
	mem.SwapUsedPct = float64(mem.SwapTotal-mem.SwapFree) / float64(memStat.Total)

	if runtime.GOOS == "linux" {
		vmStat, err := GetVMStat()
		if err != nil {
			logger.Errorf("get linux vm stats error: %v", err)
			return &mem, err
		}

		mem.SwapInPages = vmStat.Pswpin
		mem.SwapOutPages = vmStat.Pswpout
	}

	return &mem, nil
}

// GetVMStat linux vmstat 统计
func GetVMStat() (*sysInfoTypes.VMStatInfo, error) {

	h, err := sysInfo.Host()
	if err != nil {
		logger.Errorf("get process info error: %v", err)
		return nil, err
	}

	if vmStatHandle, ok := h.(sysInfoTypes.VMStat); ok {
		info, err := vmStatHandle.VMStat()
		if err != nil {
			logger.Errorf("get vmStat info error: %v", err)
			return nil, err
		}
		return info, nil
	}

	return nil, err
}

func (mem *Memory) Byte() []byte {
	buf := json.NewBuffer()
	buf.Tab("")
	buf.KL("total"          , int64(mem.Total))
	buf.KL("free"           , int64(mem.Free))
	buf.KL("used_pct"       , int64(mem.UsedPct))
	buf.KL("swap_total"     , int64(mem.SwapTotal))
	buf.KL("swap_free"      , int64(mem.SwapFree))
	buf.KL("swap_used_pct"  , int64(mem.SwapUsedPct))
	buf.KL("swap_in_pages"  , int64(mem.SwapInPages))
	buf.KL("swap_out_pages" , int64(mem.SwapOutPages))
	buf.End("}")
	return buf.Bytes()
}

func (mem *Memory) String() string {
	return lua.B2S(mem.Byte())
}