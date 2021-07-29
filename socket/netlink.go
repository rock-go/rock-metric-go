//+build linux

package socket

import (
	"github.com/elastic/gosigar/sys/linux"
	"github.com/rock-go/rock/logger"
	"os"
	"sync/atomic"
)

type NetlinkSession struct {
	readBuffer []byte
	sq         uint32
}

func NewNetlinkSession() *NetlinkSession {
	return &NetlinkSession{
		readBuffer: make([]byte, os.Getpagesize()),
		sq:         0,
	}
}

// GetSocketList 从内核获取socket连接
func (s *NetlinkSession) GetSocketList() ([]*linux.InetDiagMsg, error) {
	req := linux.NewInetDiagReq()
	req.Header.Seq = atomic.AddUint32(&s.sq, 1)
	sockets, err := linux.NetlinkInetDiagWithBuf(req, s.readBuffer, nil)
	if err != nil {
		logger.Errorf("get sockets error: %v", err)
		return nil, err
	}

	return sockets, err
}
