package memory

import (
	"github.com/rock-go/rock/lua"
)

func (m *Memory) DisableReflect() {}

func (m *Memory) Get(L *lua.LState , key string) lua.LValue {
	if key == "json" {
		return lua.JsonMarshal(L , m)
	}
	return lua.LNil
}

func newLuaMem(L *lua.LState) int {
	mem, err := GetMemDetail()
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	L.Push(L.NewAnyData(mem))
	return 1
}

func Inject(kv lua.UserKV) {
	kv.Set("mem" , lua.NewFunction(newLuaMem))
}
