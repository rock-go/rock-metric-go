package diskIo

import (
	"github.com/rock-go/rock/lua"
)

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
