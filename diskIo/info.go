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

var CacheMetric Metric // 缓存每次获取的数据
