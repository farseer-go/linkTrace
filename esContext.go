package linkTrace

import (
	"github.com/farseer-go/elasticSearch"
	"github.com/farseer-go/fs/configure"
)

var ESContext *esContext

// EsContext 链路追踪上下文
type esContext struct {
	TraceContext elasticSearch.IndexSet[TraceContext] `es:"index=linktrace_yyyy_MM;alias=linktrace;shards=1;replicas=0;refresh=3"`
}

// initEsContext 初始化上下文
func initEsContext() {
	elasticSearch.RegisterInternalContext("LinkTrace", configure.GetString("LinkTrace.ES"))
	ESContext = elasticSearch.NewContext[esContext]("LinkTrace")
}
