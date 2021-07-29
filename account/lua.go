package account

import (
	"github.com/rock-go/rock/lua"
)

func (a *Accounts) Get(L *lua.LState , key string) lua.LValue {
	if key == "json" {
		return lua.JsonMarshal(L , a)
	}

	return lua.LNil
}

func (g *Groups) Get(L *lua.LState , key string) lua.LValue {
	if key == "json" {
		return lua.JsonMarshal(L , g)
	}

	return lua.LNil
}

func getAccountByLua(L *lua.LState) int {
	accounts := GetAccounts()
	L.Push(L.NewAnyData(&accounts))
	return 1
}

func getGroupByLua(L *lua.LState) int {
	groups := GetGroups()
	L.Push(L.NewAnyData(&groups))
	return 1
}

func Inject(kv lua.UserKV) {
	kv.Set("account" , lua.NewFunction(getAccountByLua))
	kv.Set("group" , lua.NewFunction(getGroupByLua))
}
