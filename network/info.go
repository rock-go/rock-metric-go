package network

import (
	"errors"
	"github.com/rock-go/rock/json"
	"github.com/rock-go/rock/logger"
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
func GetBase(addr string) (*Base, error) {
	var base Base

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
	if f == "all" {
		return true
	}

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

func Json(interfaces []Ifc) []byte {
	if interfaces == nil {
		return nil
	}

	var buf = json.NewBuffer()
	buf.Arr("")

	length := len(interfaces)
	for i := 0; i < length ; i++ {
		ifc := interfaces[i]
		ipv4 := strings.Join(ifc.Inet , " ")
		ipv6 := strings.Join(ifc.Inet6 , " ")
		buf.Tab("")
		buf.KV("name"                , ifc.Name)
		buf.KV("inet"                , ipv4)
		buf.KV("inet6"               , ipv6)
		buf.KV("mac"                 , ifc.Mac)
		buf.KV("in_bytes"            , ifc.Flow.InBytes)
		buf.KV("in_packets"          , ifc.Flow.InPackets)
		buf.KV("in_error"            , ifc.Flow.InError)
		buf.KV("in_dropped"          , ifc.Flow.InDropped)
		buf.KV("in_bytes_per_sec"    , ifc.InBytesPerSec)
		buf.KV("in_packet_per_sec"   , ifc.Flow.InPacketsPerSec)
		buf.KV("out_bytes"           , ifc.Flow.OutBytes)
		buf.KV("out_packets"         , ifc.Flow.OutPackets)
		buf.KV("out_error"           , ifc.Flow.OutError)
		buf.KV("out_dropped"         , ifc.Flow.OutDropped)
		buf.KV("out_bytes_per_sec"   , ifc.Flow.OutBytesPerSec)
		buf.KV("out_packets_per_sec" , ifc.Flow.OutPacketsPerSec)
		buf.End("},")
	}

	buf.End("]")
	return buf.Bytes()
}