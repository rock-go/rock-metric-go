package socket

import (
	"github.com/elastic/gosigar"
	"github.com/rock-go/rock/logger"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var (
	InodeProc = make(map[uint32]int)
	lock      sync.Mutex
)

func GetSockets(s *Summary, f string) {
	GetInodes()

	var netlink NetlinkSession
	socketList, err := netlink.GetSocketList()
	if err != nil {
		logger.Errorf("get socket list by netlink error: %v", err)
		return
	}
	lock.Lock()
	defer lock.Unlock()

	var sockets = make([]*Socket, 0)
	for _, socketNetlink := range socketList {
		socket := &Socket{
			State:      TCPState(socketNetlink.State).String(),
			LocalIP:    socketNetlink.SrcIP().String(),
			LocalPort:  socketNetlink.SrcPort(),
			RemoteIP:   socketNetlink.DstIP().String(),
			RemotePort: socketNetlink.DstPort(),
			Pid:        uint32(InodeProc[socketNetlink.Inode]),
		}

		if filter(*socket, f) {
			sockets = append(sockets, socket)
			getStats(s, *socket)
		}
	}
	s.Sockets = sockets
}

func filter(socket Socket, s string) bool {
	if s == "" {
		return true
	}

	if strings.Contains(socket.LocalIP, s) ||
		strings.Contains(socket.RemoteIP, s) ||
		strings.Contains(socket.State, s) {
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

type TCPState uint8

// https://github.com/torvalds/linux/blob/5924bbecd0267d87c24110cbe2041b5075173a25/include/net/tcp_states.h#L16
const (
	TCP_ESTABLISHED TCPState = iota + 1
	TCP_SYN_SENT
	TCP_SYN_RECV
	TCP_FIN_WAIT1
	TCP_FIN_WAIT2
	TCP_TIME_WAIT
	TCP_CLOSE
	TCP_CLOSE_WAIT
	TCP_LAST_ACK
	TCP_LISTEN
	TCP_CLOSING /* Now a valid state */
)

var tcpStateNames = map[TCPState]string{
	TCP_ESTABLISHED: "ESTABLISHED",
	TCP_SYN_SENT:    "SYN_SENT",
	TCP_SYN_RECV:    "SYN_RCVD",
	TCP_FIN_WAIT1:   "FIN_WAIT1",
	TCP_FIN_WAIT2:   "FIN_WAIT2",
	TCP_TIME_WAIT:   "TIME_WAIT",
	TCP_CLOSE:       "CLOSED",
	TCP_CLOSE_WAIT:  "CLOSE_WAIT",
	TCP_LAST_ACK:    "LAST_ACK",
	TCP_LISTEN:      "LISTEN",
	TCP_CLOSING:     "CLOSING",
}

func (s TCPState) String() string {
	if state, found := tcpStateNames[s]; found {
		return state
	}
	return "UNKNOWN"
}

// GetInodes 获取所有的inode和pid对应关系
func GetInodes() {
	pids := gosigar.ProcList{}
	err := pids.Get()
	if err != nil {
		logger.Errorf("get pid list error: %v", err)
		return
	}

	for _, pid := range pids.List {
		GetInodesByPid(pid)
	}
}

// GetInodesByPid 通过pid获取inode，返回该inode对应的pid
func GetInodesByPid(pid int) {
	path := "/proc" + "/" + strconv.Itoa(pid) + "/fd/"
	d, err := os.Open(path)
	if err != nil {
		return
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return
	}

	lock.Lock()
	defer lock.Unlock()

	for _, name := range names {
		pathLink := path + name
		target, err := os.Readlink(pathLink)
		if err != nil {
			continue
		}

		setValue(InodeProc, pid, target)
	}
}

func setValue(inodeProc map[uint32]int, pid int, targetLink string) {
	if !strings.HasPrefix(targetLink, "socket:[") {
		return
	}

	inode, err := strconv.ParseInt(targetLink[8:len(targetLink)-1], 10, 64)
	if err != nil {
		return
	}

	// todo, lock
	inodeProc[uint32(inode)] = pid
}
