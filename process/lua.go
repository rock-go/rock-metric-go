package process

import (
	"github.com/rock-go/rock/lua"
)

func (s *Summary) Get(L *lua.LState , key string) lua.LValue {
	if key == "json" {
		return lua.JsonMarshal(L , s)
	}
	return lua.LNil
}

func (p *Process) Get(L *lua.LState , key string) lua.LValue {
	if key == "json" {
		return lua.JsonMarshal(L , p)
	}

	return lua.LNil
}

func newLuaProc(L *lua.LState) int {
	var pattern string
	n := L.GetTop()
	if n > 0 {
		pattern = L.CheckString(1)
	}

	proc := GetSummary(pattern)
	L.Push(L.NewAnyData(proc))
	return 1
}

func getProcByPid(L *lua.LState) int {
	n := L.GetTop()
	if n < 1 {
		L.RaiseError("need 1 arg, got 0")
		return 0
	}

	pid := L.CheckInt(1)
	val := GetByPid(pid)
	L.Push(L.NewAnyData(val))
	return 1
}

func Inject(kv lua.UserKV) {
	kv.Set("process" , lua.NewFunction(newLuaProc))
	kv.Set("process_by_pid" , lua.NewFunction(getProcByPid))
}