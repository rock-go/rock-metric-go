package service

import (
	"github.com/rock-go/rock/lua"
)

func GetServiceByLua(L *lua.LState) int {
	var filter string
	n := L.GetTop()
	if n > 0 {
		filter = L.CheckString(1)
	}

	ss := GetDetail(filter)
	L.Push(L.NewAnyData(ss))
	return 1
}

func Inject(kv lua.UserKV) {
	kv.Set("service" , lua.NewFunction(GetServiceByLua))
}