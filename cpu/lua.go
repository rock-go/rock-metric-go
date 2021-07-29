package cpu

import (
	"github.com/rock-go/rock/lua"
)

var sample = Metric{}

func (cpu *CPU) Get(L *lua.LState, key string) lua.LValue {
	if key == "json" {
		return lua.JsonMarshal(L, cpu)
	}
	return lua.LNil
}

func newLuaCpu(L *lua.LState) int {
	info := Get(&sample)
	L.Push(L.NewAnyData(info))
	return 1
}

func Inject(kv lua.UserKV) {
	kv.Set("cpu", lua.NewFunction(newLuaCpu))
}
