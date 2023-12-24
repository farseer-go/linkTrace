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
	lstTraceDetailDatabase := collections.NewList[TraceDetailDatabasePO]()
	lstTraceDetailEs := collections.NewList[TraceDetailEsPO]()
	lstTraceDetailEtcd := collections.NewList[TraceDetailEtcdPO]()
	lstTraceDetailHand := collections.NewList[TraceDetailHandPO]()
	lstTraceDetailHttp := collections.NewList[TraceDetailHttpPO]()
	lstTraceDetailMq := collections.NewList[TraceDetailMqPO]()
	lstTraceDetailRedis := collections.NewList[TraceDetailRedisPO]()

	lstMessage.Foreach(func(item *any) {
		// 上下文
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

		// 明细
		for _, detail := range traceContext.List {
			switch detailType := detail.(type) {
			case *linkTrace.TraceDetailDatabase:
				lstTraceDetailDatabase.Add(mapper.Single[TraceDetailDatabasePO](*detailType))
			case *linkTrace.TraceDetailEs:
				lstTraceDetailEs.Add(mapper.Single[TraceDetailEsPO](*detailType))
			case *linkTrace.TraceDetailEtcd:
				lstTraceDetailEtcd.Add(mapper.Single[TraceDetailEtcdPO](*detailType))
			case *linkTrace.TraceDetailHand:
				lstTraceDetailHand.Add(mapper.Single[TraceDetailHandPO](*detailType))
			case *linkTrace.TraceDetailHttp:
				lstTraceDetailHttp.Add(mapper.Single[TraceDetailHttpPO](*detailType))
			case *linkTrace.TraceDetailMq:
				lstTraceDetailMq.Add(mapper.Single[TraceDetailMqPO](*detailType))
			case *linkTrace.TraceDetailRedis:
				lstTraceDetailRedis.Add(mapper.Single[TraceDetailRedisPO](*detailType))
			}
		}
	})
	// 写入上下文
	_, err := CHContext.TraceContext.InsertList(lstTraceContext, 2000)
	flog.ErrorIfExists(err)

	// 写入明细
	if lstTraceDetailDatabase.Count() > 0 {
		_, err = CHContext.TraceDetailDatabase.InsertList(lstTraceDetailDatabase, 2000)
		flog.ErrorIfExists(err)
	}
	if lstTraceDetailEs.Count() > 0 {
		_, err = CHContext.TraceDetailEs.InsertList(lstTraceDetailEs, 2000)
		flog.ErrorIfExists(err)
	}
	if lstTraceDetailEtcd.Count() > 0 {
		_, err = CHContext.TraceDetailEtcd.InsertList(lstTraceDetailEtcd, 2000)
		flog.ErrorIfExists(err)
	}
	if lstTraceDetailHand.Count() > 0 {
		_, err = CHContext.TraceDetailHand.InsertList(lstTraceDetailHand, 2000)
		flog.ErrorIfExists(err)
	}
	if lstTraceDetailHttp.Count() > 0 {
		_, err = CHContext.TraceDetailHttp.InsertList(lstTraceDetailHttp, 2000)
		flog.ErrorIfExists(err)
	}
	if lstTraceDetailMq.Count() > 0 {
		_, err = CHContext.TraceDetailMq.InsertList(lstTraceDetailMq, 2000)
		flog.ErrorIfExists(err)
	}
	if lstTraceDetailRedis.Count() > 0 {
		_, err = CHContext.TraceDetailRedis.InsertList(lstTraceDetailRedis, 2000)
		flog.ErrorIfExists(err)
	}
	return
}
