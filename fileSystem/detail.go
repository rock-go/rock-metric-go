package fileSystem

import (
	"github.com/elastic/gosigar"
	"github.com/rock-go/rock/json"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
)

type detail []FileSystem

// GetDetail 获取所有文件系统的详细统计
func newFileSystemDetail() (detail, error) {
	fsList := gosigar.FileSystemList{}
	err := fsList.Get()
	if err != nil {
		logger.Errorf("get file system list error: %v", err)
		return nil, err
	}

	fsStats := make([]FileSystem, len(fsList.List))
	for i, fs := range fsList.List {
		fsStats[i] = *getFSStat(fs)
	}

	return fsStats, nil
}

func (d detail) Byte() []byte {
	buf := json.NewBuffer()
	buf.Arr("")

	for _, item := range d {
		buf.Tab("")
		buf.KV("name", item.Name)
		buf.KV("type", item.Type)
		buf.KV("mount", item.MountPoint)
		buf.KL("available", int64(item.Available))
		buf.KL("free", int64(item.Free))
		buf.KL("used", int64(item.Used))
		buf.KF64("used_pct", item.UsedPct)
		buf.KL("total", int64(item.Total))
		buf.End("},")
	}
	buf.End("]")
	return buf.Bytes()
}

func (d detail) String() string {
	return lua.B2S(d.Byte())
}
