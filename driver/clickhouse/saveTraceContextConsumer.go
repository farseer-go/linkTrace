package linkTrace_clickhouse

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/linkTrace"
	"github.com/farseer-go/mapper"
)

// 写入到clickhouse
func saveTraceContextConsumer(subscribeName string, lstMessage collections.ListAny, remainingCount int) {
	trace.CurTraceContext.Get().Ignore()
	lstTraceContext := collections.NewList[TraceContextPO]()
	lstMessage.Foreach(func(item *any) {
		traceContext := (*item).(linkTrace.TraceContext)
		po := mapper.Single[TraceContextPO](traceContext)
		lstTraceContext.Add(po)
	})
	_, err := CHContext.TraceContext.InsertList(lstTraceContext, 2000)
	flog.ErrorIfExists(err)
	return
}
