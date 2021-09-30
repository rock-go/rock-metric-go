package command

import (
	"github.com/rock-go/rock/json"
	"github.com/rock-go/rock/lua"
)

type History struct {
	User    string `json:"user"`
	ID      string `json:"id"`
	Command string `json:"command"`
}

type HistoryMap map[string][]*History

func (hm HistoryMap) Byte() []byte {
	buf := json.NewBuffer()
	buf.Tab("")
	for key , hs := range hm {
		buf.Arr(key)
		for _ , item := range hs {
			buf.Tab("")
			buf.KV("user" , item.User)
			buf.KV("id" , item.ID)
			buf.KV("command" , item.Command)
			buf.End("},")
		}
		buf.End("],")
	}
	buf.End("}")

	return buf.Bytes()
}

func (hm HistoryMap) String() string {
	return lua.B2S(hm.Byte())
}