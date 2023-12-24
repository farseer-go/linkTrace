package linkTrace_elasticSearch

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/linkTrace"
	"github.com/farseer-go/mapper"
)

// 写入到ES
func saveTraceContextConsumer(subscribeName string, lstMessage collections.ListAny, remainingCount int) {
	trace.CurTraceContext.Get().Ignore()
	lstTraceContext := collections.NewList[TraceContextPO]()
	lstMessage.Foreach(func(item *any) {
		traceContext := (*item).(linkTrace.TraceContext)
		po := mapper.Single[TraceContextPO](traceContext)
		if !traceContext.Web.IsNil() {
			po.Web = &traceContext.Web
		}
		if !traceContext.Consumer.IsNil() {
			po.Consumer = &traceContext.Consumer
		}
		if !traceContext.Task.IsNil() {
			po.Task = &traceContext.Task
		}
		if !traceContext.WatchKey.IsNil() {
			po.WatchKey = &traceContext.WatchKey
		}
		lstTraceContext.Add(po)
	})
	err := ESContext.TraceContext.InsertList(lstTraceContext)
	flog.ErrorIfExists(err)
	return
}
