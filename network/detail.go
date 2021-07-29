package network

import "github.com/rock-go/rock/lua"

type detail []Ifc

//关闭热反射
func (d *detail) DisableReflect() {}

func (d *detail) Get(L *lua.LState, key string) lua.LValue {
	if key == "json" {
		return lua.JsonMarshal(L, *d)
	}
	return lua.LNil
}
