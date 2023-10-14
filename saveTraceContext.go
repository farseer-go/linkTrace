package linkTrace

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/mapper"
)

func saveTraceContextConsumer(subscribeName string, lstMessage collections.ListAny, remainingCount int) {
	lstTraceContext := mapper.ToList[TraceContext](lstMessage)
	err := ESContext.TraceContext.InsertList(lstTraceContext)
	_ = flog.Error(err)
	return
}
