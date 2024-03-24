package linkTrace

import (
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/queue"
	"strings"
	"time"
)

// Enable 是否启用
var defConfig config

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{queue.Module{}}
}

func (module Module) PreInitialize() {
	defConfig = configure.ParseConfig[config]("LinkTrace")

	// 使用了链路追踪组件，则要把空组件移除后，重新注册
	container.Remove[trace.IManager]()
	container.Register(func() trace.IManager {
		return &traceManager{}
	})

	// 启用了链路追踪后，才需要初始化ES和消费
	if defConfig.Enable {
		FopsServer = strings.ToLower(configure.GetString("Fops.Server"))
		if !strings.HasPrefix(FopsServer, "http") {
			panic("[farseer.yaml]Fops.Server配置不正确，示例：https://fops.fsgit.com")
		}
		if !strings.HasSuffix(FopsServer, "/") {
			FopsServer += "/"
		}
		queue.Subscribe("TraceContext", "SaveTraceContext", 1000,5*time.Second, SaveTraceContextConsumer)
	}
}
