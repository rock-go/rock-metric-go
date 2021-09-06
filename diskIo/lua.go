package diskIo

import (
	"github.com/rock-go/rock/lua"
)

func (d *detail) Get(L *lua.LState, key string) lua.LValue {
	if key == "json" {
		return lua.JsonMarshal(L, d)
	}

	return lua.LNil
}

func newLuaDiskIO(L *lua.LState) int {
	d, e := newDiskIoDetail()
	if e != nil {
		L.RaiseError("marshal disk io stats to json error: %v", e)
		return 0
	}
	L.Push(L.NewAnyData(&d))
	return 1
}

func Inject(kv lua.UserKV) {
	kv.Set("disk_io", lua.NewFunction(newLuaDiskIO))
}
