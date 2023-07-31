package fSchedule

import (
	"fmt"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/fs/timingWheel"
	"github.com/farseer-go/webapi"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{webapi.Module{}}
}

func (module Module) PreInitialize() {
	// 服务端配置
	defaultServer = serverVO{
		Address: configure.GetSlice("FSchedule.Server.Address"),
		Token:   configure.GetString("FSchedule.Server.Token"),
	}

	// 客户端配置
	NewClient()
	timingWheel.Start()
}

func (module Module) PostInitialize() {
	webapi.Area("/api/", func() {
		webapi.RegisterPOST("/check", Check)
		webapi.RegisterPOST("/invoke", Invoke)
		webapi.RegisterPOST("/status", Status)
		webapi.RegisterPOST("/kill", Kill)
	})
	webapi.UseApiResponse()
	webapi.UsePprof()
	go webapi.Run(fmt.Sprintf("%s:%d", defaultClient.ClientIp, defaultClient.ClientPort))

	fs.AddInitCallback("注册客户端", func() {
		defaultClient.RegistryClient()
	})
}

func (module Module) Shutdown() {
	defaultClient.LogoutClient()
}
