package linkTrace_clickhouse

import (
	"github.com/farseer-go/data"
	"github.com/farseer-go/fs/configure"
)

var CHContext *chContext

// EsContext 链路追踪上下文
type chContext struct {
	TraceContext data.TableSet[TraceContextPO] `data:"name=link_trace;migrate"`
}

// initCHContext 初始化上下文
func initCHContext() {
	data.RegisterInternalContext("LinkTrace", configure.GetString("LinkTrace.CH"))
	CHContext = data.NewContext[chContext]("LinkTrace")
}
