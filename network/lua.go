package network

import (
	"github.com/rock-go/rock/lua"
)

func newLuaIFC(L *lua.LState) int {
	var addr string
	n := L.GetTop()
	if n > 0 {
		addr = L.CheckString(1)
	}

	d, err := GetDetail(addr)
	if err != nil {
		L.RaiseError("get network interfaces error: %v", err)
		return 0
	}

	L.Push(L.NewAnyData(&d))
	return 1
}

func Inject(kv lua.UserKV) {
	kv.Set("ifc", lua.NewFunction(newLuaIFC))
}
