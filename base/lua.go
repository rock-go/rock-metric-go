package base

import (
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/node"
)

type Receipt struct {
	Successful bool `json:"successful"`
}

func getByLua(L *lua.LState) int {
	info , err :=  Get(node.LoadAddr())
	if err != nil {
		L.RaiseError("%s get basic info err: %v" , err)
		return 0
	}

	L.Push(L.NewAnyData(info))
	return 1
}

func Inject(kv lua.UserKV) {
	kv.Set("base", lua.NewFunction(getByLua))
}
