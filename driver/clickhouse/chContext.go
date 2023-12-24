package linkTrace_clickhouse

import (
	"github.com/farseer-go/data"
	"github.com/farseer-go/fs/configure"
)

var CHContext *chContext

// EsContext 链路追踪上下文
type chContext struct {
	TraceContext        data.TableSet[TraceContextPO]        `data:"name=link_trace;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,app_id,trace_id) PARTITION BY toYYYYMM(app_name)"`
	TraceDetailDatabase data.TableSet[TraceDetailDatabasePO] `data:"name=trace_detail_database;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,app_id,trace_id,db_name,table_name,connection_string) PARTITION BY toYYYYMM(app_name)"`
	TraceDetailEs       data.TableSet[TraceDetailEsPO]       `data:"name=trace_detail_es;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,app_id,trace_id,index_name,aliases_name) PARTITION BY toYYYYMM(app_name)"`
	TraceDetailEtcd     data.TableSet[TraceDetailEtcdPO]     `data:"name=trace_detail_etcd;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,app_id,trace_id,key) PARTITION BY toYYYYMM(app_name)"`
	TraceDetailHand     data.TableSet[TraceDetailHandPO]     `data:"name=trace_detail_hand;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,app_id,trace_id,name) PARTITION BY toYYYYMM(app_name)"`
	TraceDetailHttp     data.TableSet[TraceDetailHttpPO]     `data:"name=trace_detail_http;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,app_id,trace_id,method,url) PARTITION BY toYYYYMM(app_name)"`
	TraceDetailMq       data.TableSet[TraceDetailMqPO]       `data:"name=trace_detail_mq;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,app_id,trace_id,server,exchange,routing_key) PARTITION BY toYYYYMM(app_name)"`
	TraceDetailRedis    data.TableSet[TraceDetailRedisPO]    `data:"name=trace_detail_redis;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,app_id,trace_id,key,field) PARTITION BY toYYYYMM(app_name)"`
}

// initCHContext 初始化上下文
func initCHContext() {
	data.RegisterInternalContext("LinkTrace", configure.GetString("LinkTrace.CH"))
	CHContext = data.NewContext[chContext]("LinkTrace")
}
