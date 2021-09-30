package process

import (
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
)

var (
	NumCpu = runtime.NumCPU()
)

// Config 备用字段
type Config struct {
	pattern   []string
	TopNByCPU int
	TopNByMem int
}

type Metric struct {
	cfg *Config

	ProcReg []regexp.Regexp

	ProcMapLast map[int]*Process // 上次获取到的进程map
	InodeProc   map[int64]uint32 // inode与进程号对应关系

	limiter *Limiter

	mutex sync.RWMutex
}

func (m *Metric) GetProc(pid int, filter string) *Process {
	proc, err := newProcess(pid)
	if err != nil {
		return nil
	}

	if filter == "all" {
		goto DO
	}

	if !strings.Contains(strings.ToLower(proc.Name), filter) {
		return nil
	}

	//if !m.matchProc(proc.Name) {
	//	return nil
	//}
DO:
	proc.GetDetail()

	// 计算CPU使用
	m.mutex.Lock()
	last := m.ProcMapLast[pid]
	proc.TotalNormPct, proc.TotalPct = getProcCpuPct(last, proc)
	m.ProcMapLast[pid] = proc
	m.mutex.Unlock()

	return proc
}

// 校验进程名是否与配置的正则匹配
func (m *Metric) matchProc(name string) bool {
	for _, reg := range m.ProcReg {
		if reg.MatchString(name) {
			return true
		}
	}

	return false
}

func getProcCpuPct(p0, p1 *Process) (normalPct, pct float64) {

	// 首次取值last为nil,跳过计算
	if p0 != nil && p1 != nil {
		timeDelta := float64(p1.SampleTime-p0.SampleTime) / 1000000
		totalCpuDeltaMillis := int64(p1.TotalTicks - p0.TotalTicks)

		pct := float64(totalCpuDeltaMillis) / timeDelta
		normalPct := pct / float64(NumCpu)

		return normalPct, pct
	}

	return 0, 0
}

// 筛选top N,取cpu和memory top N的并集
func (m *Metric) includeTopN(s *Summary) *Summary {

	if m.cfg.TopNByCPU <= 0 && m.cfg.TopNByMem <= 0 {
		return s
	}

	var res Summary
	var numProc int

	// cpu top N
	if len(s.Process) < m.cfg.TopNByCPU {
		numProc = len(s.Process)
	}

	sort.SliceStable(s.Process, func(i, j int) bool {
		return s.Process[i].TotalPct > s.Process[j].TotalPct
	})
	numProc = m.cfg.TopNByCPU
	res.Process = append(res.Process, s.Process[:numProc]...)

	// memory top N
	if len(s.Process) < m.cfg.TopNByMem {
		numProc = len(s.Process)
	}

	sort.SliceStable(s.Process, func(i, j int) bool {
		return s.Process[i].RssBytes > s.Process[j].RssBytes
	})

	// 取并集, 将不在res的proc存入res
	for _, proc := range s.Process[:numProc] {
		if !hasProc(&res, proc) {
			res.Process = append(res.Process, proc)
		}
	}

	return &res
}

// 判断proc是否在processes内
func hasProc(s *Summary, proc *Process) bool {
	for _, p := range s.Process {
		if p.Pid == proc.Pid {
			return true
		}
	}
	return false
}
