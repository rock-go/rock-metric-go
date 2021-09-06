package diskIo

import (
	"errors"
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

func (d *detail) DisableReflect() {}
