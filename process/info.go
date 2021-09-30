package process

import (
	"fmt"
	"github.com/elastic/gosigar"
	"github.com/rock-go/rock/json"
	"github.com/rock-go/rock/lua"
	"regexp"
	"strings"
	"sync"
)

type Summary struct {
	Idle     uint32 `json:"idle"`
	Running  uint32 `json:"running"`
	Sleeping uint32 `json:"sleeping"`
	Stopped  uint32 `json:"stopped"`
	Total    uint32 `json:"total"`
	Unknown  uint32 `json:"unknown"`
	Zombie   uint32 `json:"zombie"`

	Process []*Process `json:"process"`
}

var (
	ProcMetric = Metric{
		cfg:         &Config{},
		ProcReg:     nil,
		ProcMapLast: nil,
		mutex:       sync.RWMutex{},
	}
)

func init() {
	ProcMetric.cfg.pattern = []string{".*"}
	ProcMetric.ProcMapLast = make(map[int]*Process)
	ProcMetric.InodeProc = make(map[int64]uint32)
	ProcMetric.limiter = newLimiter(1)

	ProcMetric.ProcReg = []regexp.Regexp{}
	for _, pattern := range ProcMetric.cfg.pattern {
		reg, err := regexp.Compile(pattern)
		if err != nil {
			fmt.Printf("regx compile pattern [%s] error: %v", pattern, err)
			return
		}
		ProcMetric.ProcReg = append(ProcMetric.ProcReg, *reg)
	}
}

func GetSummary(filter string) *Summary {
	var summary = Summary{}

	ProcMetric.limiter.Handler()

	pidList := gosigar.ProcList{}
	err := pidList.Get()
	if err != nil {
		return &summary
	}

	for _, pid := range pidList.List {
		proc := ProcMetric.GetProc(pid, filter)
		if proc == nil {
			continue
		}

		summary.Process = append(summary.Process, proc)
		summary.getStatsCount(proc)
	}

	return &summary
}

// GetByPid 通过pid获取单个进程信息
func GetByPid(pid int) *Process {
	return ProcMetric.GetProc(pid, "")
}

func getSummaryFromMap(metric *Metric, filter string) *Summary {
	var summary = Summary{}

	for _, proc := range metric.ProcMapLast {
		if strings.Contains(strings.ToLower(proc.Name), filter) {
			summary.Process = append(summary.Process, proc)
			summary.getStatsCount(proc)
		}
	}

	return &summary
}

func (s *Summary) getStatsCount(proc *Process) {
	if proc == nil {
		return
	}

	switch proc.State {
	case "sleeping":
		s.Sleeping++
	case "running":
		s.Running++
	case "idle":
		s.Idle++
	case "stopped":
		s.Stopped++
	case "zombie":
		s.Zombie++
	default:
		s.Unknown++
	}

	s.Total++
}

func (s *Summary) Byte() []byte {
	buf := json.NewBuffer()
	buf.Tab("")
	buf.KV("idle"    ,  s.Idle     )
	buf.KV("running" ,  s.Running  )
	buf.KV("sleeping",  s.Sleeping )
	buf.KV("stopped" ,  s.Stopped  )
	buf.KV("total"   ,  s.Total    )
	buf.KV("unknown" ,  s.Unknown  )
	buf.KV("zombie"  ,  s.Zombie   )
	buf.Arr("process")

	for _ , item := range s.Process {
		item.Marshal(buf)
	}

	buf.End("]}")

	return buf.Bytes()
}

func (s *Summary) String() string {
	return lua.B2S(s.Byte())
}