package linkTrace

import (
	"github.com/farseer-go/elasticSearch"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/modules"
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
}

func (module Module) PostInitialize() {
	initEsContext()
	queue.Subscribe("TraceContext", "SaveTraceContext", 1000, saveTraceContextConsumer)
}
