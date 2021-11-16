package fileSystem

import (
	"github.com/rock-go/rock/audit"
	"github.com/rock-go/rock/lua"
)


func newLuaFileSystem(L *lua.LState) int {
	audit.RecoverByCodeVM(L, audit.Subject("file system info error"))

	fs, err := newFileSystemDetail()
	if err != nil {
		panic(err)
		return 0
	}

	L.Push(L.NewAnyData(fs))
	return 1
}

func Inject(kv lua.UserKV) {
	kv.Set("file_system", lua.NewFunction(newLuaFileSystem))
}
