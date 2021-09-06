package socket

import (
	"github.com/rock-go/rock/lua"
)

func (s *Summary) Get(L *lua.LState, key string) lua.LValue {
	if key == "json" {
		return lua.JsonMarshal(L, s)
	}
	return lua.LNil
}

func getSocketByLua(L *lua.LState) int {
	var f string
	n := L.GetTop()
	if n > 0 {
		f = L.CheckAny(1).String()
	}

	sockets := GetSummary(f)
	L.Push(L.NewAnyData(sockets))
	return 1
}

func Inject(kv lua.UserKV) {
	kv.Set("socket", lua.NewFunction(getSocketByLua))
}
