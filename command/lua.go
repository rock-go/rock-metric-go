package command

import (
	"github.com/rock-go/rock/lua"
)


func getHistoryByLua(L *lua.LState) int {
	var username string
	n := L.GetTop()
	if n > 0 {
		username = L.CheckSocket(1)
	}

	hm := GetHistory(username)

	L.Push(L.NewAnyData(&hm))
	return 1
}

func Inject(kv lua.UserKV) {
	kv.Set("history", lua.NewFunction(getHistoryByLua))
}
