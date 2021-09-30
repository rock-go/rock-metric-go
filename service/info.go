package service

import (
	"github.com/rock-go/rock/json"
	"github.com/rock-go/rock/lua"
)

type Service struct {
	Name        string `json:"name"`
	StartType   string `json:"start_type"`
	ExecPath    string `json:"exec_path"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	State       string `json:"state"`
	Pid         uint32 `json:"pid"`
	ExitCode    uint32 `json:"exit_code"`
}

type Services []*Service

func GetDetail(pattern string) Services {
	return Services(GetService(pattern))
}

func (ss *Services) Byte() []byte {
	buf := json.NewBuffer()
	buf.Arr("")

	for _ , item := range *ss {
		buf.Tab("")
		buf.KV("name", item.Name)
		buf.KV("start_type", item.StartType)
		buf.KV("exec_path", item.ExecPath)
		buf.KV("display_name", item.DisplayName)
		buf.KV("description", item.Description)
		buf.KV("state", item.State)
		buf.KV("pid", item.Pid)
		buf.KV("exit_code", item.ExitCode)
		buf.End("},")
	}

	buf.End("]")
	return buf.Bytes()
}

func (ss Services) String() string {
	return lua.B2S(ss.Byte())
}