package fileSystem

import "github.com/elastic/gosigar"

type detail []FileSystem

// GetDetail 获取所有文件系统的详细统计
func newFileSystemDetail() (detail, error) {
	fsList := gosigar.FileSystemList{}
	err := fsList.Get()
	if err != nil {
		return nil, err
	}

	fsStats := make([]FileSystem, len(fsList.List))
	for i, fs := range fsList.List {
		fsStats[i] = *getFSStat(fs)
	}

	return fsStats, nil
}

func (d *detail) DisableReflect() {}
