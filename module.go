package linkTrace

import (
	"github.com/farseer-go/elasticSearch"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/queue"
)

// Enable 是否启用
var defConfig config

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{queue.Module{}, elasticSearch.Module{}}
}

func (module Module) PreInitialize() {
	defConfig = configure.ParseConfig[config]("LinkTrace")
	// 使用了链路追踪组件，则要把空组件移除后，重新注册
	container.Remove[trace.IManager]()
	container.Register(func() trace.IManager {
		return &traceManager{}
	})
}

func (module Module) PostInitialize() {
	// 启用了链路追踪后，才需要初始化ES和消费
	if defConfig.Enable {
		initEsContext()
		queue.Subscribe("TraceContext", "SaveTraceContext", 1000, saveTraceContextConsumer)
	}
}
