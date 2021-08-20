package base

import (
	"context"
	"github.com/rock-go/rock/logger"
	"github.com/rock-go/rock/lua"
	"github.com/smallnest/rpcx/client"
)

type Receipt struct {
	Successful bool `json:"successful"`
}

func TestClient(data interface{}) {
	cli := client.NewClient(client.DefaultOption)
	err := cli.Connect("tcp", "172.31.61.168:8082")
	if err != nil {
		logger.Errorf("%v", err)
		return
	}

	in := data
	resp := &Receipt{}
	err = cli.Call(context.TODO(), "Report", "Base", in, resp)
	if err != nil {
		logger.Errorf("%v", err)
	}
	logger.Errorf("%v", resp)
}

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
	TestClient(data.Value)
	L.Push(data)
	return 1
}

func Inject(kv lua.UserKV) {
	kv.Set("base", lua.NewFunction(luaBaseGet))
}
