package socket

import (
	"github.com/rock-go/rock/json"
	"github.com/rock-go/rock/lua"
)

type Socket struct {
	State      string `json:"state"`
	LocalIP    string `json:"local_ip"`
	LocalPort  int    `json:"local_port"`
	RemoteIP   string `json:"remote_ip"`
	RemotePort int    `json:"remote_port"`
	Pid        uint32 `json:"pid"`
	Process    string `json:"process"`
}

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

// GetSummary 获取全部的socket连接
func GetSummary(filter string) *Summary {
	var s Summary
	GetSockets(&s, filter)
	return &s
}

// GetSpecifiedSockets 获取指定的socket
func GetSpecifiedSockets() {

}

func (s *Summary) Byte() []byte {
	buf := json.NewBuffer()
	buf.Tab("")
	buf.KV("closed", s.CLOSED)
	buf.KV("listen", s.LISTEN)
	buf.KV("syn_sent", s.SYN_SENT)
	buf.KV("syn_rcvd", s.SYN_RCVD)
	buf.KV("established", s.ESTABLISHED)
	buf.KV("fin_wait1", s.FIN_WAIT1)
	buf.KV("fin_wait2", s.FIN_WAIT2)
	buf.KV("close_wait", s.CLOSE_WAIT)
	buf.KV("closing", s.CLOSING)
	buf.KV("last_ack", s.LAST_ACK)
	buf.KV("time_wait", s.TIME_WAIT)
	buf.KV("delete_tcb", s.DELETE_TCB)
	buf.Arr("sockets")

	for _, item := range s.Sockets {
		item.Marshal(buf)
	}
	buf.End("]}")

	return buf.Bytes()
}

func (s *Summary) String() string {
	return lua.B2S(s.Byte())
}

func (s *Socket) Marshal(buf *json.Buffer) {
	buf.Tab("")

	buf.KV("state", s.State)
	buf.KV("local_ip", s.LocalIP)
	buf.KV("local_port", s.LocalPort)
	buf.KV("remote_ip", s.RemoteIP)
	buf.KV("remote_port", s.RemotePort)
	buf.KV("pid", s.Pid)
	buf.KV("process_name", s.Process)

	buf.End("},")
}

func (s *Socket) Byte() []byte {
	buf := json.NewBuffer()
	s.Marshal(buf)
	buf.End("")
	return buf.Bytes()
}

func (s *Socket) String() string {
	return lua.B2S(s.Byte())
}
