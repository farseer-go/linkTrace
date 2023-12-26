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
		if !traceContext.Exception.IsNil() {
			exceptionStackPO := mapper.Single[ExceptionStackPO](traceContext.Exception)
			po.Exception = &exceptionStackPO
		}
		lstTraceContext.Add(po)

		// 明细
		for _, detail := range traceContext.List {
			switch detailType := detail.(type) {
			case *linkTrace.TraceDetailDatabase:
				databasePO := mapper.Single[TraceDetailDatabasePO](*detailType)
				databasePO.BaseTraceDetailPO = mapper.Single[BaseTraceDetailPO](detailType.BaseTraceDetail)
				_ = mapper.Auto(traceContext, &databasePO.BaseTraceDetailPO)
				if !detailType.Exception.IsNil() {
					databasePO.Exception = &detailType.Exception
				}
				lstTraceDetailDatabase.Add(databasePO)
			case *linkTrace.TraceDetailEs:
				esPO := mapper.Single[TraceDetailEsPO](*detailType)
				esPO.BaseTraceDetailPO = mapper.Single[BaseTraceDetailPO](detailType.BaseTraceDetail)
				_ = mapper.Auto(traceContext, &esPO.BaseTraceDetailPO)
				if !detailType.Exception.IsNil() {
					esPO.Exception = &detailType.Exception
				}
				lstTraceDetailEs.Add(esPO)
			case *linkTrace.TraceDetailEtcd:
				etcdPO := mapper.Single[TraceDetailEtcdPO](*detailType)
				etcdPO.BaseTraceDetailPO = mapper.Single[BaseTraceDetailPO](detailType.BaseTraceDetail)
				_ = mapper.Auto(traceContext, &etcdPO.BaseTraceDetailPO)
				if !detailType.Exception.IsNil() {
					etcdPO.Exception = &detailType.Exception
				}
				lstTraceDetailEtcd.Add(etcdPO)
			case *linkTrace.TraceDetailHand:
				handPO := mapper.Single[TraceDetailHandPO](*detailType)
				handPO.BaseTraceDetailPO = mapper.Single[BaseTraceDetailPO](detailType.BaseTraceDetail)
				_ = mapper.Auto(traceContext, &handPO.BaseTraceDetailPO)
				if !detailType.Exception.IsNil() {
					handPO.Exception = &detailType.Exception
				}
				lstTraceDetailHand.Add(handPO)
			case *linkTrace.TraceDetailHttp:
				httpPO := mapper.Single[TraceDetailHttpPO](*detailType)
				httpPO.BaseTraceDetailPO = mapper.Single[BaseTraceDetailPO](detailType.BaseTraceDetail)
				_ = mapper.Auto(traceContext, &httpPO.BaseTraceDetailPO)
				if !detailType.Exception.IsNil() {
					httpPO.Exception = &detailType.Exception
				}
				lstTraceDetailHttp.Add(httpPO)
			case *linkTrace.TraceDetailMq:
				mqPO := mapper.Single[TraceDetailMqPO](*detailType)
				mqPO.BaseTraceDetailPO = mapper.Single[BaseTraceDetailPO](detailType.BaseTraceDetail)
				_ = mapper.Auto(traceContext, &mqPO.BaseTraceDetailPO)
				if !detailType.Exception.IsNil() {
					mqPO.Exception = &detailType.Exception
				}
				lstTraceDetailMq.Add(mqPO)
			case *linkTrace.TraceDetailRedis:
				redisPO := mapper.Single[TraceDetailRedisPO](*detailType)
				redisPO.BaseTraceDetailPO = mapper.Single[BaseTraceDetailPO](detailType.BaseTraceDetail)
				_ = mapper.Auto(traceContext, &redisPO.BaseTraceDetailPO)
				if !detailType.Exception.IsNil() {
					redisPO.Exception = &detailType.Exception
				}
				lstTraceDetailRedis.Add(redisPO)
			}
		}
	})
	// 写入上下文
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
