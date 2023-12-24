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
	err := ESContext.TraceContext.InsertList(lstTraceContext)
	flog.ErrorIfExists(err)

	// 写入明细
	if lstTraceDetailDatabase.Count() > 0 {
		err = ESContext.TraceDetailDatabase.InsertList(lstTraceDetailDatabase)
		flog.ErrorIfExists(err)
	}
	if lstTraceDetailEs.Count() > 0 {
		err = ESContext.TraceDetailEs.InsertList(lstTraceDetailEs)
		flog.ErrorIfExists(err)
	}
	if lstTraceDetailEtcd.Count() > 0 {
		err = ESContext.TraceDetailEtcd.InsertList(lstTraceDetailEtcd)
		flog.ErrorIfExists(err)
	}
	if lstTraceDetailHand.Count() > 0 {
		err = ESContext.TraceDetailHand.InsertList(lstTraceDetailHand)
		flog.ErrorIfExists(err)
	}
	if lstTraceDetailHttp.Count() > 0 {
		err = ESContext.TraceDetailHttp.InsertList(lstTraceDetailHttp)
		flog.ErrorIfExists(err)
	}
	if lstTraceDetailMq.Count() > 0 {
		err = ESContext.TraceDetailMq.InsertList(lstTraceDetailMq)
		flog.ErrorIfExists(err)
	}
	if lstTraceDetailRedis.Count() > 0 {
		err = ESContext.TraceDetailRedis.InsertList(lstTraceDetailRedis)
		flog.ErrorIfExists(err)
	}
	return
}
