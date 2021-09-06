package metric

import (
	"github.com/rock-go/rock-metric-go/account"
	"github.com/rock-go/rock-metric-go/base"
	"github.com/rock-go/rock-metric-go/command"
	"github.com/rock-go/rock-metric-go/cpu"
	"github.com/rock-go/rock-metric-go/diskIo"
	"github.com/rock-go/rock-metric-go/fileSystem"
	"github.com/rock-go/rock-metric-go/memory"
	"github.com/rock-go/rock-metric-go/network"
	"github.com/rock-go/rock-metric-go/process"
	"github.com/rock-go/rock-metric-go/service"
	"github.com/rock-go/rock-metric-go/socket"
	"github.com/rock-go/rock/lua"
	"github.com/rock-go/rock/xcall"
)

func LuaInjectApi(env xcall.Env) {
	kv := lua.NewUserKV()
	base.Inject(kv)
	cpu.Inject(kv)
	diskIo.Inject(kv)
	memory.Inject(kv)
	process.Inject(kv)
	service.Inject(kv)
	socket.Inject(kv)
	account.Inject(kv)
	command.Inject(kv)
	network.Inject(kv)
	fileSystem.Inject(kv)
	env.Set("metric", kv)
}
