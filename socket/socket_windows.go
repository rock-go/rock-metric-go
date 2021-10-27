package socket

import (
	"fmt"
	"github.com/elastic/gosigar"
	"github.com/rock-go/rock/logger"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

const ERROR_INSUFFICIENT_BUFFER = 122

var (
	LazyDll = syscall.NewLazyDLL("Iphlpapi.dll")
	Proc    = LazyDll.NewProc("GetTcpTable2")
)

type Inet uint32

func (i Inet) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", i&255, i>>8&255, i>>16&255, i>>24&255)
}

type NtoHs uint32

func (n NtoHs) String() string {
	return fmt.Sprint(syscall.Ntohs(uint16(n)))
}

func (n NtoHs) Int() int {
	p, _ := strconv.Atoi(n.String())
	return p
}

type TCP_CONNECTION_OFFLOAD_STATE uint32

//状态枚举
var _MIB_TCP_STATE = map[uint32]string{
	1:  "CLOSED",
	2:  "LISTEN",
	3:  "SYN_SENT",
	4:  "SYN_RCVD",
	5:  "ESTABLISHED",
	6:  "FIN_WAIT1",
	7:  "FIN_WAIT2",
	8:  "CLOSE_WAIT",
	9:  "CLOSING",
	10: "LAST_ACK",
	11: "TIME_WAIT",
	12: "DELETE_TCB",
}

type MIB_TCP_STATE uint32

func (m MIB_TCP_STATE) String() string {
	return _MIB_TCP_STATE[uint32(m)]
}

type MIB_TCPROW2 struct {
	dwState        MIB_TCP_STATE
	dwLocalAddr    Inet
	dwLocalPort    NtoHs
	dwRemoteAddr   Inet
	dwRemotePort   NtoHs
	dwOwningPid    uint32
	dwOffloadState TCP_CONNECTION_OFFLOAD_STATE
}

type MIB_TCPTABLE2 struct {
	dwNumEntries uint32
	table        [1]MIB_TCPROW2
}

func GetSockets(summary *Summary, f string) {
	var mibtable2 MIB_TCPTABLE2
	size := unsafe.Sizeof(mibtable2)

	r, _, err := Proc.Call(uintptr(unsafe.Pointer(&mibtable2)), uintptr(unsafe.Pointer(&size)), 1)
	if err != nil && r != 0 {
		if r == ERROR_INSUFFICIENT_BUFFER {
			buf := make([]byte, size)
			r, _, err = Proc.Call(uintptr(unsafe.Pointer(&buf[0])), uintptr(unsafe.Pointer(&size)), 1)
			if r != 0 {
				logger.Errorf("get tcp table error: %v", err)
				return
			}

			var index = int(unsafe.Sizeof(mibtable2.dwNumEntries))
			var step = int(unsafe.Sizeof(mibtable2.table))
			dwNumEntries := *(*uint32)(unsafe.Pointer(&buf[0]))

			var sockets = make([]*Socket, 0)
			for i := 0; i < int(dwNumEntries); i++ {
				mib := *(*MIB_TCPROW2)(unsafe.Pointer(&buf[index]))
				index += step
				socket := format(mib)
				if filter(*socket, f) {
					sockets = append(sockets, socket)
					getStats(summary, *socket)
				}
			}
			summary.Sockets = sockets
		}
	}
}

func format(mib MIB_TCPROW2) *Socket {
	socket := Socket{
		State:      mib.dwState.String(),
		LocalIP:    mib.dwLocalAddr.String(),
		LocalPort:  mib.dwLocalPort.Int(),
		RemoteIP:   mib.dwRemoteAddr.String(),
		RemotePort: mib.dwRemotePort.Int(),
		Pid:        mib.dwOwningPid,
	}

	state := gosigar.ProcState{}
	err := state.Get(int(socket.Pid))
	if err != nil {
		logger.Debugf("get process name of pid [%d] error: %v", socket.Pid, err)
	}
	socket.Process = state.Name

	return &socket
}

func filter(socket Socket, s string) bool {
	if s == "" {
		return true
	}

	if strings.Contains(socket.LocalIP, s) ||
		strings.Contains(socket.RemoteIP, s) ||
		strings.Contains(socket.State, s) ||
		strings.Contains(strings.ToLower(socket.Process), strings.ToLower(s)) {
		return true
	}

	d, err := strconv.Atoi(s)
	if err != nil {
		return false
	}

	if socket.LocalPort == d || socket.RemotePort == d || socket.Pid == uint32(d) {
		return true
	}

	return false
}

func getStats(summary *Summary, socket Socket) {
	state := socket.State
	s := reflect.ValueOf(*summary)

	count := s.FieldByName(state).Int()
	count++

	f := reflect.ValueOf(summary).Elem()
	f.FieldByName(state).SetInt(count)
}
