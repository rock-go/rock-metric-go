# 说明

rock-metric-go 模块基于rock-go框架开发，用于获取系统相关信息

# 使用

## 导入

```go
import metric "github.com/rock-go/rock-metric-go"
```

## 注册

```go
rock.Inject(env, metric.LuaInjectApi)
```

## lua 脚本调用
获取的是系统各个信息的对象， 类型为AnyData ,但是封装了对应json对象 

每个函数被调用时仅执行一次，只获取一次数据，获取到的数据为userdata，其值为golang的[]byte数据类型，json格式

```lua
-- 声明metric对象
local metric= rock.metric

local cpu = metric.cpu()
rock.ERR(cpu.json)

```

### 基础信息

```lua
local base = metric.base("114.114.114.114:53")
-- 参数说明：该参数为ip:port格式，用于获取当前系统网络连接的网卡信息，默认为8.8.8.8:53
```

#### 结果字段

cpu使用率需要至少两组数据进行计算，因此第一次获取可能为0

```go
package basicInfo

type BasicInfo struct {
	Inet      string  `json:"inet"`       // ipv4地址
	Inet6     string  `json:"inet6"`      // ipv6地址
	Mac       string  `json:"mac"`        // mac地址
	Arch      string  `json:"arch"`       // 处理器架构
	Platform  string  `json:"platform"`   // 操作系统平台
	MemTotal  uint64  `json:"mem_total"`  // 内存总量，单位byte
	MemFree   uint64  `json:"mem_free"`   // 可用内存，单位byte
	SwapTotal uint64  `json:"swap_total"` // Swap总量，单位byte
	SwapFree  uint64  `json:"swap_free"`  // 可用Swap，单位byte
	CpuCore   int     `json:"cpu_core"`   // CPU核心数
	CpuUsage  float64 `json:"cpu_usage"`  // CPU用量
	DiskTotal uint64  `json:"disk_total"` // 最大容量磁盘的总量
	DiskPath  string  `json:"disk_path"`  // 最大容量磁盘的路径
	DiskFree  uint64  `json:"disk_free"`  // 最大容量磁盘可用容量
}
```

### 系统账户信息

获取当前系统所有用户和用户组信息，返回userdata数据，其值为[]byte类型，格式为json

```lua
local accounts = sysinfo.account()
local groups = sysinfo.groups()
```

#### 结果字段

##### linux

linux的账户名通过解析/etc/passwd和/etc/groups文件来获取

```go
package account

type Account struct {
	LoginName string `json:"login_name"` // 登录字段
	UID       string `json:"uid"`
	GID       string `json:"gid"`
	UserName  string `json:"user_name"`
	HomeDir   string `json:"home_dir"`
	Shell     string `json:"shell"`
	Raw       string `json:"raw"` // 原始数据
}

type Group struct {
	GroupName string `json:"group_name"`
	GID       string `json:"gid"`
	Raw       string `json:"raw"`
}
```

##### windows

```go
package account

type Account struct {
	AccountType        uint32    `json:"account_type"`
	Caption            string    `json:"caption"` // 简单描述
	Description        string    `json:"description"`
	Disabled           bool      `json:"disabled"`
	Domain             string    `json:"domain"`
	FullName           string    `json:"full_name"`
	InstallDate        time.Time `json:"install_date"`
	LocalAccount       bool      `json:"local_account"`
	Lockout            bool      `json:"lockout"` // 是否锁定
	Name               string    `json:"name"`
	PasswordChangeable bool      `json:"password_changeable"`
	PasswordExpires    bool      `json:"password_expires"`
	PasswordRequired   bool      `json:"password_required"`
	SID                string    `json:"sid"`
	SIDType            uint8     `json:"sid_type"`
	Status             string    `json:"status"`
}

type Group struct {
	Caption      string    `json:"caption"`
	Description  string    `json:"description"`
	Domain       string    `json:"domain"`
	InstallDate  time.Time `json:"install_date"`
	LocalAccount bool      `json:"local_account"`
	Name         string    `json:"name"`
	Sid          string    `json:"sid"`
	SidType      uint8     `json:"sid_type"`
	Status       string    `json:"status"`
}
```

### CPU信息

CPU使用百分比需要至少两次数据进行计算，因此第一次获取到的可能为0，一般周期性获取。返回的值为userdata，其值为[]byte类型，json格式，字段为各状态的百分比

```lua
local cpu = sysinfo.cpu()
```

#### 结果字段

```go
package cpu

type CPU struct {
	Architecture string
	CoreNum      int

	User    float64
	System  float64
	Idle    float64
	IOWait  float64
	IRQ     float64
	Nice    float64
	SoftIRQ float64
	Stolen  float64
	Total   float64
}
```

### 磁盘IO

磁盘IO至少需要两次数据进行计算，因此第一次可能为0，一般周期性调用获取统计值。

```lua
local io = sysinfo.diskio()
```

#### 结果字段

```go
package diskIo

type Disk struct {
	Name            string
	SerialNumber    string
	IoTime          uint64
	ReadBytes       uint64
	ReadPerSecBytes float64
	WriteBytes      uint64
	WritePerSecByte float64
}
```

### 文件系统

获取挂载的文件系统使用情况。

```lua
local fs = sysinfo.fs()
```

#### 结果字段

```go
package fileSystem

type FileSystem struct {
	Name       string
	Type       string
	MountPoint string // linux
	Available  uint64
	Free       uint64
	Used       uint64
	UsedPct    float64
	Total      uint64
}
```

### 内存

获取内存的使用情况

```lua
local mem = sysinfo.mem()
```

#### 结果字段

```go
package memory

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
```

### 网卡流量

获取网卡使用情况，返回多个网卡信息

```lua
local ifc = sysinfo.ifc()
```

#### 结果字段

```go
package network

// Ifc 单个网卡interface
type Ifc struct {
	Name string

	Inet  []string
	Inet6 []string

	Mac string

	// 流量数据
	Flow
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
```

### 进程信息

获取进程的相关信息

```lua
local process = sysinfo.process()
local chrome = sysinfo.process("chrome")
-- 参数：默认获取所有进程信息；参数为string类型时，返回包含该字段的进程信息
local procPid = sysinfo.process_by_pid(8312)
-- 参数：int类型，通过pid获取进程信息
```

#### 结果字段

```go
package process

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
```

### 服务

获取系统的服务信息

```lua
local service = sysinfo.service("abs")
-- 参数：默认为获取所有服务信息，带有string类型参数时，返回名称中包含该字段的服务
```

#### 返回字段

```go

package service

type Service struct {
	Name        string `json:"name"`
	StartType   string `json:"start_type"`
	ExecPath    string `json:"exec_path"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	State       string `json:"state"`
	Pid         uint32 `json:"pid"`
	ExitCode    uint32 `json:"exit_code"`
}
```

### Socket连接

获取当前系统的网络连接

```lua
--local sockets = sysinfo.socket()
local sockets = sysinfo.socket("172.31.61.67")
-- 参数：不带参数时，默认获取所有的网络连接；当包含参数（int或string）时，获取包含该参数的连接
```

#### 结果字段

```go

package socket

type Socket struct {
	State      string `json:"state"`
	LocalIP    string `json:"local_ip"`
	LocalPort  int    `json:"local_port"`
	RemoteIP   string `json:"remote_ip"`
	RemotePort int    `json:"remote_port"`
	Pid        uint32 `json:"pid"`
}

// Summary 返回如下内容
type Summary struct {
	CLOSED      int `json:"closed"`
	LISTEN      int `json:"listen"`
	SYN_SENT    int `json:"syn_sent"`
	SYN_RCVD    int `json:"syn_rcvd"`
	ESTABLISHED int `json:"established"`
	FIN_WAIT1   int `json:"fin_wait1"`
	FIN_WAIT2   int `json:"fin_wait2"`
	CLOSE_WAIT  int `json:"close_wait"`
	CLOSING     int `json:"closing"`
	LAST_ACK    int `json:"last_ack"`
	TIME_WAIT   int `json:"time_wait"`
	DELETE_TCB  int `json:"delete_tcb, omitempty"`

	Sockets []*Socket `json:"sockets"`
}
```

### linux 历史命令

获取Linux的历史命令，目前从用户主目录的.bash_history文件中获取。

```lua
local h = sysinfo.history("root")
-- 参数：最多一个参数，string类型，指定获取一个用户的历史命令；默认获取全部
```

#### 返回结果

返回的结果为json格式，包含指定的用户或所有用户的历史命令记录

```go
package command

type History struct {
	User    string `json:"user"`
	ID      string `json:"id"`
	Command string `json:"command"`
}

// return map[string][]*History
```

### lua脚本调试

```lua
local groups = sysinfo.group()
sysinfo.debug(groups)
```
