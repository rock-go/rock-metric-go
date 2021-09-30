package process

import (
	"github.com/elastic/gosigar"
	"github.com/rock-go/rock/json"
	"github.com/rock-go/rock/lua"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Process struct {
	Name       string   `json:"name"`
	State      string   `json:"state"`
	Pid        uint32   `json:"pid"`
	Ppid       uint32   `json:"ppid"`
	Pgid       uint32   `json:"pgid"`
	Cmdline    string   `json:"cmdline"`
	Username   string   `json:"username"`
	Cwd        string   `json:"cwd"`
	Executable string   `json:"executable"` // linux
	Args       []string `json:"args"`

	//CPU，单位 毫秒
	UserTicks    uint64  `json:"user_ticks"`
	TotalPct     float64 `json:"total_pct"`
	TotalNormPct float64 `json:"total_norm_pct"`
	SystemTicks  uint64  `json:"system_ticks"`
	TotalTicks   uint64  `json:"total_ticks"`
	StartTime    string  `json:"start_time"`

	//Memory
	MemSize  uint64  `json:"mem_size"`
	RssBytes uint64  `json:"rss_bytes"`
	RssPct   float64 `json:"rss_pct"`
	Share    uint64  `json:"share"`

	inodes []uint64 // 进程相关的inode

	SampleTime uint64 // 数据采集时间
}

func newProcess(pid int) (*Process, error) {
	var err error

	state := gosigar.ProcState{}
	err = state.Get(pid)
	if err != nil {
		return nil, err
	}

	proc := &Process{
		Name:     state.Name,
		State:    getProcState(byte(state.State)),
		Pid:      uint32(pid),
		Ppid:     uint32(state.Ppid),
		Pgid:     uint32(state.Pgid),
		Username: state.Username,
	}

	exe := gosigar.ProcExe{}
	_ = exe.Get(pid)

	proc.Cwd = exe.Cwd
	proc.Executable = exe.Name
	return proc, nil
}

func (p *Process) GetDetail() {
	p.SampleTime = uint64(time.Now().UnixNano())

	p.getMemUsage()
	p.getCpuUsage()
	p.getArgs()
}

func (p *Process) getMemUsage() {
	mem := gosigar.ProcMem{}
	_ = mem.Get(int(p.Pid))
	p.MemSize = mem.Size
	p.RssBytes = mem.Resident
	p.Share = mem.Share

	stats := gosigar.Mem{}
	err := stats.Get()
	if err != nil {
		return
	}

	p.RssPct = float64(mem.Resident) / float64(stats.Total)
}

func (p *Process) getCpuUsage() {
	cpu := gosigar.ProcTime{}
	err := cpu.Get(int(p.Pid))
	if err != nil {
		return
	}

	p.UserTicks = cpu.User
	p.SystemTicks = cpu.Sys
	p.TotalTicks = cpu.Total
	p.StartTime = time.Unix(0, int64(cpu.StartTime)).String()
}

func (p *Process) getArgs() {
	args := gosigar.ProcArgs{}
	err := args.Get(int(p.Pid))
	if err != nil {
		return
	}

	p.Args = args.List
	p.Cmdline = strings.Join(p.Args, " ")
}

func getProcState(b byte) string {
	switch b {
	case 'S':
		return "sleeping"
	case 'R':
		return "running"
	case 'D':
		return "idle"
	case 'T':
		return "stopped"
	case 'Z':
		return "zombie"
	}

	return "unknown"
}

// GetInodes 获取进程inodes,如果该进程有socket连接
func (p *Process) GetInodes() {
	path := filepath.Join("/proc", strconv.Itoa(int(p.Pid)), "fd")
	d, err := os.Open(path)
	if err != nil {
		return
	}

	names, err := d.Readdirnames(-1)
	if err != nil {
		return
	}

	targets := make([]string, len(names))
	for i, name := range names {
		target, err := os.Readlink(filepath.Join(path, name))
		if err == nil {
			targets[i] = target
		}
	}

	var inodes []uint64
	for _, fd := range targets {
		if strings.HasPrefix(fd, "socket:[") {
			inode, err := strconv.ParseInt(fd[8:len(fd)-1], 10, 64)
			if err != nil {
				continue
			}

			inodes = append(inodes, uint64(inode))
		}
	}
	p.inodes = inodes

	return
}

func (p *Process) Marshal(buf *json.Buffer) {
	buf.Tab("")
	buf.KV("name"           , p.Name)
	buf.KV("state"          , p.State)
	buf.KV("pid"            , p.Pid)
	buf.KV("ppid"           , p.Ppid)
	buf.KV("pgid"           , p.Pgid)
	buf.KV("cmdline"        , p.Cmdline)
	buf.KV("username"       , p.Username)
	buf.KV("cwd"            , p.Cwd)
	buf.KV("executable"     , p.Executable)
	buf.KV("args"           , p.Args)
	buf.KV("user_ticks"     , p.UserTicks)
	buf.KV("total_pct"      , p.TotalPct)
	buf.KV("total_norm_pct" , p.TotalNormPct)
	buf.KV("system_ticks"   , p.SystemTicks)
	buf.KV("total_ticks"    , p.TotalTicks)
	buf.KV("start_time"     , p.StartTime)
	buf.KV("user_ticks"     , p.UserTicks)
	buf.KV("total_pct"      , p.TotalPct)
	buf.KV("total_norm_pct" , p.TotalNormPct)
	buf.KV("system_ticks"   , p.SystemTicks)
	buf.KV("total_ticks"    , p.TotalTicks)
	buf.KV("start_time"     , p.StartTime)

	buf.KV("mem_size"       , p.MemSize  )
	buf.KV("rss_bytes"      , p.RssBytes )
	buf.KV("rss_pct"        , p.RssPct   )
	buf.KV("share"          , p.Share    )
	buf.KV("inodes"         , p.inodes)
	buf.KV("sample_time"    , p.SampleTime)
	buf.End("},")
}

func (p *Process) Byte() []byte {
	buf := json.NewBuffer()
	p.Marshal(buf)
	buf.End("")
	return buf.Bytes()
}

func (p *Process) String() string {
	return lua.B2S(p.Byte())
}
