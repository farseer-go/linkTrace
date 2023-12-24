package linkTrace_elasticSearch

import (
	"github.com/farseer-go/elasticSearch"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/linkTrace"
)

var ESContext *esContext

// EsContext 链路追踪上下文
type esContext struct {
	TraceContext elasticSearch.IndexSet[linkTrace.TraceContext] `es:"index=link_trace_yyyy_MM;alias=link_trace;shards=1;replicas=0;refresh=3"`
}

// initEsContext 初始化上下文
func initEsContext() {
	elasticSearch.RegisterInternalContext("LinkTrace", configure.GetString("LinkTrace.ES"))
	ESContext = elasticSearch.NewContext[esContext]("LinkTrace")
}
