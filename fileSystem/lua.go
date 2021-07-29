package fileSystem

import (
	"github.com/rock-go/rock/lua"
)

func (d *detail) Get(L *lua.LState, key string) lua.LValue {
	if key == "json" {
		return lua.JsonMarshal(L, d)
	}

	return lua.LNil

}

func newLuaFileSystem(L *lua.LState) int {
	fs, err := newFileSystemDetail()
	if err != nil {
		L.RaiseError("get file system info error: %v", err)
		return 0
	}

	L.Push(L.NewAnyData(fs))
	return 1
}

func Inject(kv lua.UserKV) {
	kv.Set("file_system", lua.NewFunction(newLuaFileSystem))
}
