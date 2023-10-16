package linkTrace

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
)

// 写入到ES
func saveTraceContextConsumer(subscribeName string, lstMessage collections.ListAny, remainingCount int) {
	lstTraceContext := collections.NewList[TraceContext]()
	lstMessage.Foreach(func(item *any) {
		traceContext := *item
		lstTraceContext.Add(traceContext.(TraceContext))
	})
	err := ESContext.TraceContext.InsertList(lstTraceContext)
	flog.ErrorIfExists(err)
	return
}
