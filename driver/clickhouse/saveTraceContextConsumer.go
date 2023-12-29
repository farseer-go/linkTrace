package linkTrace_clickhouse

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/linkTrace"
	"github.com/farseer-go/mapper"
	"time"
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
	lstTraceDetailGrpc := collections.NewList[TraceDetailGrpcPO]()
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
		po.UseDesc = po.UseTs.String()
		po.CreateAt = time.UnixMicro(po.StartTs)
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
				databasePO.UseDesc = databasePO.UseTs.String()
				databasePO.CreateAt = time.UnixMicro(databasePO.StartTs)
				lstTraceDetailDatabase.Add(databasePO)
			case *linkTrace.TraceDetailEs:
				esPO := mapper.Single[TraceDetailEsPO](*detailType)
				esPO.BaseTraceDetailPO = mapper.Single[BaseTraceDetailPO](detailType.BaseTraceDetail)
				_ = mapper.Auto(traceContext, &esPO.BaseTraceDetailPO)
				if !detailType.Exception.IsNil() {
					esPO.Exception = &detailType.Exception
				}
				esPO.UseDesc = esPO.UseTs.String()
				esPO.CreateAt = time.UnixMicro(esPO.StartTs)
				lstTraceDetailEs.Add(esPO)
			case *linkTrace.TraceDetailEtcd:
				etcdPO := mapper.Single[TraceDetailEtcdPO](*detailType)
				etcdPO.BaseTraceDetailPO = mapper.Single[BaseTraceDetailPO](detailType.BaseTraceDetail)
				_ = mapper.Auto(traceContext, &etcdPO.BaseTraceDetailPO)
				if !detailType.Exception.IsNil() {
					etcdPO.Exception = &detailType.Exception
				}
				etcdPO.UseDesc = etcdPO.UseTs.String()
				etcdPO.CreateAt = time.UnixMicro(etcdPO.StartTs)
				lstTraceDetailEtcd.Add(etcdPO)
			case *linkTrace.TraceDetailHand:
				handPO := mapper.Single[TraceDetailHandPO](*detailType)
				handPO.BaseTraceDetailPO = mapper.Single[BaseTraceDetailPO](detailType.BaseTraceDetail)
				_ = mapper.Auto(traceContext, &handPO.BaseTraceDetailPO)
				if !detailType.Exception.IsNil() {
					handPO.Exception = &detailType.Exception
				}
				handPO.UseDesc = handPO.UseTs.String()
				handPO.CreateAt = time.UnixMicro(handPO.StartTs)
				lstTraceDetailHand.Add(handPO)
			case *linkTrace.TraceDetailHttp:
				httpPO := mapper.Single[TraceDetailHttpPO](*detailType)
				httpPO.BaseTraceDetailPO = mapper.Single[BaseTraceDetailPO](detailType.BaseTraceDetail)
				_ = mapper.Auto(traceContext, &httpPO.BaseTraceDetailPO)
				if !detailType.Exception.IsNil() {
					httpPO.Exception = &detailType.Exception
				}
				httpPO.UseDesc = httpPO.UseTs.String()
				httpPO.CreateAt = time.UnixMicro(httpPO.StartTs)
				lstTraceDetailHttp.Add(httpPO)
			case *linkTrace.TraceDetailGrpc:
				grpcPO := mapper.Single[TraceDetailGrpcPO](*detailType)
				grpcPO.BaseTraceDetailPO = mapper.Single[BaseTraceDetailPO](detailType.BaseTraceDetail)
				_ = mapper.Auto(traceContext, &grpcPO.BaseTraceDetailPO)
				if !detailType.Exception.IsNil() {
					grpcPO.Exception = &detailType.Exception
				}
				grpcPO.UseDesc = grpcPO.UseTs.String()
				grpcPO.CreateAt = time.UnixMicro(grpcPO.StartTs)
				lstTraceDetailGrpc.Add(grpcPO)
			case *linkTrace.TraceDetailMq:
				mqPO := mapper.Single[TraceDetailMqPO](*detailType)
				mqPO.BaseTraceDetailPO = mapper.Single[BaseTraceDetailPO](detailType.BaseTraceDetail)
				_ = mapper.Auto(traceContext, &mqPO.BaseTraceDetailPO)
				if !detailType.Exception.IsNil() {
					mqPO.Exception = &detailType.Exception
				}
				mqPO.UseDesc = mqPO.UseTs.String()
				mqPO.CreateAt = time.UnixMicro(mqPO.StartTs)
				lstTraceDetailMq.Add(mqPO)
			case *linkTrace.TraceDetailRedis:
				redisPO := mapper.Single[TraceDetailRedisPO](*detailType)
				redisPO.BaseTraceDetailPO = mapper.Single[BaseTraceDetailPO](detailType.BaseTraceDetail)
				_ = mapper.Auto(traceContext, &redisPO.BaseTraceDetailPO)
				if !detailType.Exception.IsNil() {
					redisPO.Exception = &detailType.Exception
				}
				redisPO.UseDesc = redisPO.UseTs.String()
				redisPO.CreateAt = time.UnixMicro(redisPO.StartTs)
				lstTraceDetailRedis.Add(redisPO)
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
	if lstTraceDetailGrpc.Count() > 0 {
		_, err = CHContext.TraceDetailGrpc.InsertList(lstTraceDetailGrpc, 2000)
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
