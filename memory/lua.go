package memory

import (
	"github.com/rock-go/rock/lua"
)

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
	kv.Set("mem", lua.NewFunction(newLuaMem))
}
