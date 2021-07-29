package network

import (
	"errors"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/utils"
	netStat "github.com/shirou/gopsutil/net"
	"net"
	"strings"
)

// Ifc 单个网卡interface
type Ifc struct {
	Name string

	Inet  []string
	Inet6 []string

	Mac string

	// 流量数据
	Flow
}

type Base struct {
	Name  string
	Inet  string
	Inet6 string
	Mac   string
}

var IfcMetrics map[string]Metric // 缓存所有网卡的统计信息
var Stats map[string]*netStat.IOCountersStat

func init() {
	IfcMetrics = make(map[string]Metric)
	Stats = make(map[string]*netStat.IOCountersStat)
}

// GetDetail 获取所有的统计数据
func GetDetail(filter string) (detail, error) {
	// 获取ip，Mac
	ifcs, err := net.Interfaces()
	if err != nil {
		logger.Errorf("get net interfaces error: %v", err)
		return nil, err
	}

	var interfaces detail

	// 网卡流量统计
	stats, err := netStat.IOCounters(true)
	if err != nil {
		logger.Errorf("get net io counters error: %v", err)
		return nil, err
	}

	for _, s := range stats {
		name := s.Name
		Stats[name] = &s
	}

	for _, ifc := range ifcs {
		if i, e := getIfcMetric(ifc, filter, true); e == nil {
			interfaces = append(interfaces, *i)
		}
	}

	return interfaces, nil
}

// GetBase 获取当前连接网络的网卡的基础信息
func GetBase(target string) (*Base, error) {
	var base Base

	addr, err := getSocketIP(target)
	if err != nil {
		return nil, err
	}

	// 获取ip，Mac
	ifcs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, ifc := range ifcs {
		name := ifc.Name
		base.Name = name
		base.Mac = ifc.HardwareAddr.String()

		ips, err := ifc.Addrs()
		if err != nil {
			continue
		}

		for _, ip := range ips {
			// inet6
			if !isInet4(ip) {
				base.Inet6 = ip.(*net.IPNet).IP.String()
			}
			// inet
			base.Inet = ip.(*net.IPNet).IP.String()
		}

		if strings.Contains(base.Inet, addr) || strings.Contains(base.Inet6, addr) {
			return &base, nil
		}
	}

	return nil, errors.New("not found")
}

// 获取单个网卡数据
func getIfcMetric(ifc net.Interface, f string, isDetail bool) (*Ifc, error) {
	var dev Ifc
	dev.Inet = make([]string, 0)
	dev.Inet6 = make([]string, 0)

	dev.Name = ifc.Name
	dev.Mac = ifc.HardwareAddr.String()

	ips, err := ifc.Addrs()
	if err != nil {
		return nil, err
	}

	// ip
	for _, ip := range ips {
		if isInet4(ip) {
			dev.Inet = append(dev.Inet, ip.String())
			continue
		}
		dev.Inet6 = append(dev.Inet6, ip.String())
	}

	// filter
	if !filter(dev, f) {
		return nil, errors.New("dropped")
	}

	if !isDetail {
		return &dev, nil
	}

	metric, ok := IfcMetrics[ifc.Name]
	if !ok {
		metric = Metric{
			lastSample: Sample{},
			nowSample:  Sample{},
		}
	}

	stat := Stats[ifc.Name]
	if stat != nil {
		metric.getMetric(*stat)
		IfcMetrics[ifc.Name] = metric
		flow := metric.calMetric()
		dev.Flow = flow
	}

	return &dev, nil
}

// 从设备名和ip地址过滤
func filter(dev Ifc, f string) bool {
	if strings.Contains(dev.Name, f) {
		return true
	}

	for _, ip := range dev.Inet {
		if strings.Contains(ip, f) {
			return true
		}
	}

	return false
}

func isInet4(ip net.Addr) bool {
	if strings.Contains(ip.String(), ":") {
		return false
	}
	return true
}

func getSocketIP(addr string) (string, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		logger.Errorf("get local socket addr by %s error: %v", addr, err)
		return "", err
	}

	defer func() {
		if err := conn.Close(); err != nil {
			logger.Errorf("conn close error: %v", err)
		}
	}()

	localAddr := conn.LocalAddr().(*net.UDPAddr).IP.String()
	return localAddr, nil
}

func Json(interfaces map[string]*Ifc) []byte {
	if interfaces == nil {
		return nil
	}

	var buf = lua.NewJsonBuffer("")
	buf.WriteVal("interface")
	buf.Write([]byte(":"))
	buf.Write([]byte("["))

	length := len(interfaces)
	k := 0
	for _, ifc := range interfaces {
		var ip4 string
		var ip6 string

		for _, ip := range ifc.Inet {
			ip4 = ip4 + ip + " "
		}

		for _, ip := range ifc.Inet6 {
			ip6 = ip6 + ip + " "
		}

		buf.EOF = false
		buf.Start("")
		buf.WriteKV("name", ifc.Name)
		buf.WriteKV("inet", ip4)
		buf.WriteKV("inet6", ip6)
		buf.WriteKV("mac", ifc.Mac)

		buf.WriteKV("in_bytes", utils.ToString(ifc.Flow.InBytes))
		buf.WriteKV("in_packets", utils.ToString(ifc.Flow.InPackets))
		buf.WriteKV("in_error", utils.ToString(ifc.Flow.InError))
		buf.WriteKV("in_dropped", utils.ToString(ifc.Flow.InDropped))
		buf.WriteKV("in_bytes_per_sec", utils.ToString(ifc.InBytesPerSec))
		buf.WriteKV("in_packet_per_sec", utils.ToString(ifc.Flow.InPacketsPerSec))
		buf.WriteKV("out_bytes", utils.ToString(ifc.Flow.OutBytes))
		buf.WriteKV("out_packets", utils.ToString(ifc.Flow.OutPackets))
		buf.WriteKV("out_error", utils.ToString(ifc.Flow.OutError))
		buf.WriteKV("out_dropped", utils.ToString(ifc.Flow.OutDropped))
		buf.WriteKV("out_bytes_per_sec", utils.ToString(ifc.Flow.OutBytesPerSec))
		buf.EOF = true
		buf.WriteKV("out_packets_per_sec", utils.ToString(ifc.Flow.OutPacketsPerSec))

		buf.End()

		if k < length-1 {
			buf.Write([]byte(","))
		}
		k++

	}

	buf.Write([]byte("]"))
	buf.End()

	return buf.Bytes()
}
