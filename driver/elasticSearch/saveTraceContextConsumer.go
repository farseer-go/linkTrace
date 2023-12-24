package elasticSearch

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/linkTrace"
)

// 写入到ES
func saveTraceContextConsumer(subscribeName string, lstMessage collections.ListAny, remainingCount int) {
	trace.CurTraceContext.Get().Ignore()
	lstTraceContext := collections.NewList[linkTrace.TraceContext]()
	lstMessage.Foreach(func(item *any) {
		traceContext := *item
		lstTraceContext.Add(traceContext.(linkTrace.TraceContext))
	})
	err := ESContext.TraceContext.InsertList(lstTraceContext)
	flog.ErrorIfExists(err)
	return
}
