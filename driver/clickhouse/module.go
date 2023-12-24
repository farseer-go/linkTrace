package linkTrace_clickhouse

import (
	"github.com/farseer-go/data/driver/clickhouse"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/linkTrace"
	"github.com/farseer-go/queue"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{data_clickhouse.Module{}, linkTrace.Module{}}
}

func (module Module) Initialize() {
}

func (module Module) PostInitialize() {
	// 启用了链路追踪后，才需要初始化ES和消费
	if configure.GetBool("LinkTrace.Enable") {
		initCHContext()
		queue.Subscribe("TraceContext", "SaveTraceContext", 1000, saveTraceContextConsumer)
	}
}
