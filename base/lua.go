package base

import (
	"github.com/rock-go/rock/lua"
)

func (bi *BasicInfo) Get(L *lua.LState, key string) lua.LValue {
	if key == "json" {
		return lua.JsonMarshal(L, bi)
	}
	return lua.LNil
}

func luaBaseGet(L *lua.LState) int {
	target := "8.8.8.8:53"
	n := L.GetTop()
	if n > 0 {
		target = L.CheckString(1)
	}

	info, err := Get(target)
	if err != nil {
		L.RaiseError("get basic info error: %v", err)
		return 0
	}

	data := L.NewAnyData(info)
	L.Push(data)
	return 1
}

func Inject(kv lua.UserKV) {
	kv.Set("base", lua.NewFunction(luaBaseGet))
}
