package linkTrace_elasticSearch

import (
	"github.com/farseer-go/elasticSearch"
	"github.com/farseer-go/fs/configure"
)

var ESContext *esContext

// EsContext 链路追踪上下文
type esContext struct {
	TraceContext        elasticSearch.IndexSet[TraceContextPO]        `es:"index=link_trace_yyyy_MM;alias=link_trace;shards=1;replicas=0;refresh=3"`
	TraceDetailDatabase elasticSearch.IndexSet[TraceDetailDatabasePO] `es:"index=trace_detail_database_yyyy_MM;alias=trace_detail_database;shards=1;replicas=0;refresh=3"`
	TraceDetailEs       elasticSearch.IndexSet[TraceDetailEsPO]       `es:"index=trace_detail_es_yyyy_MM;alias=trace_detail_es;shards=1;replicas=0;refresh=3"`
	TraceDetailEtcd     elasticSearch.IndexSet[TraceDetailEtcdPO]     `es:"index=trace_detail_etcd_yyyy_MM;alias=trace_detail_etcd;shards=1;replicas=0;refresh=3"`
	TraceDetailHand     elasticSearch.IndexSet[TraceDetailHandPO]     `es:"index=trace_detail_hand_yyyy_MM;alias=trace_detail_hand;shards=1;replicas=0;refresh=3"`
	TraceDetailHttp     elasticSearch.IndexSet[TraceDetailHttpPO]     `es:"index=trace_detail_http_yyyy_MM;alias=trace_detail_http;shards=1;replicas=0;refresh=3"`
	TraceDetailGrpc     elasticSearch.IndexSet[TraceDetailGrpcPO]     `es:"index=trace_detail_grpc_yyyy_MM;alias=trace_detail_http;shards=1;replicas=0;refresh=3"`
	TraceDetailMq       elasticSearch.IndexSet[TraceDetailMqPO]       `es:"index=trace_detail_mq_yyyy_MM;alias=trace_detail_mq;shards=1;replicas=0;refresh=3"`
	TraceDetailRedis    elasticSearch.IndexSet[TraceDetailRedisPO]    `es:"index=trace_detail_redis_yyyy_MM;alias=trace_detail_redis;shards=1;replicas=0;refresh=3"`
}

// initEsContext 初始化上下文
func initEsContext() {
	elasticSearch.RegisterInternalContext("LinkTrace", configure.GetString("LinkTrace.ES"))
	ESContext = elasticSearch.NewContext[esContext]("LinkTrace")
}
