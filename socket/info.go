package socket

type Socket struct {
	State      string `json:"state"`
	LocalIP    string `json:"local_ip"`
	LocalPort  int    `json:"local_port"`
	RemoteIP   string `json:"remote_ip"`
	RemotePort int    `json:"remote_port"`
	Pid        uint32 `json:"pid"`
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
