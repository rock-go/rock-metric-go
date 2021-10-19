package fileSystem

import (
	"github.com/elastic/gosigar"
	"github.com/rock-go/rock/logger"
	"math"
)

type FileSystem struct {
	Name       string
	Type       string
	MountPoint string // linux
	Available  uint64
	Free       uint64
	Used       uint64
	UsedPct    float64
	Total      uint64
}

// GetMax 获取最大的分区系统的简要信息
func GetMax() (path string, total uint64, free uint64, err error) {
	fsList := gosigar.FileSystemList{}
	err = fsList.Get()
	if err != nil {
		return "", 0, 0, err
	}

	for _, fs := range fsList.List {
		fsStat := getFSStat(fs)
		if fsStat == nil {
			continue
		}

		if total < fsStat.Total {
			path = fs.DevName
			total = fsStat.Total
			free = fsStat.Free
		}
	}

	return path, total, free, nil
}

// 通过基本信息获取使用统计
func getFSStat(fs gosigar.FileSystem) *FileSystem {
	stat := gosigar.FileSystemUsage{}
	if err := stat.Get(fs.DirName); err != nil {
		logger.Errorf("get file system stat error: %v", err)
		return nil
	}

	fssStat := &FileSystem{
		Name:       fs.DevName,
		Type:       fs.SysTypeName,
		MountPoint: fs.DirName,
		Available:  stat.Avail,
		Free:       stat.Free,
		Used:       stat.Used,
		UsedPct:    float64(stat.Used) / float64(stat.Used+stat.Avail),
		Total:      stat.Total,
	}

	if math.IsNaN(fssStat.UsedPct) || math.IsInf(fssStat.UsedPct, 1) {
		fssStat.UsedPct = 0
	}

	return fssStat
}
