package diskIo

import (
	"errors"
	"github.com/rock-go/rock/json"
	"github.com/rock-go/rock/lua"
	"github.com/shirou/gopsutil/disk"
)

type detail []Disk

func newDiskIoDetail() (detail, error) {
	var d detail
	counters, err := disk.IOCounters()
	if err != nil {
		//logger.Errorf("get disk io stats error: %v", err)
		return d, errors.New("get disk io state error: " + err.Error())
	}

	var stats = make(map[string]*Disk)
	for name, counter := range counters {
		stat := Disk{
			Name:            counter.Name,
			SerialNumber:    counter.SerialNumber,
			IoTime:          counter.IoTime,
			ReadBytes:       counter.ReadBytes,
			ReadPerSecBytes: 0,
			WriteBytes:      counter.WriteBytes,
			WritePerSecByte: 0,
		}
		stats[name] = &stat
	}

	err = CacheMetric.getMetric(stats)
	if err != nil {
		//logger.Errorf("get disk io stats error: %v", err)
		return d, errors.New("get disk io state error: " + err.Error())
	}

	d = CacheMetric.calMetric()
	return d, nil

}

func (d *detail) Byte() []byte {
	buf := json.NewBuffer()
	buf.Arr("")
	for _ , item := range *d {
		buf.Tab("")
		buf.KV("name"    , item.Name)
		buf.KV("serial"  , item.SerialNumber)
		buf.KL("io_time" , int64(item.IoTime))
		buf.KL("read"    , int64(item.ReadBytes))
		buf.KL("write"   , int64(item.WriteBytes))
		buf.KF64("read_sec_bytes" , item.ReadPerSecBytes)
		buf.KF64("write_sec_byte" , item.WritePerSecByte)
		buf.End("},")
	}
	buf.End("]")
	return buf.Bytes()
}

func (d *detail) String() string {
	return lua.B2S(d.Byte())
}