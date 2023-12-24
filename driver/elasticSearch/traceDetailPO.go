package linkTrace_elasticSearch

import (
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/fs/trace/eumCallType"
	"time"
)

type BaseTraceDetailPO struct {
	TraceId        int64                // 上下文ID
	AppId          int64                // 应用ID
	AppName        string               // 应用名称
	AppIp          string               // 应用IP
	ParentAppName  string               // 上游应用
	DetailId       int64                `es:"primaryKey"`
	ParentDetailId int64                // 父级明细ID
	Level          int                  // 当前层级（入口为0层）
	MethodName     string               // 调用方法
	CallType       eumCallType.Enum     // 调用类型
	Timeline       time.Duration        // 从入口开始统计
	UnTraceTs      time.Duration        // 上一次结束到现在开始之间未Trace的时间
	StartTs        int64                // 调用开始时间戳
	EndTs          int64                // 调用停止时间戳
	UseTs          time.Duration        // 总共使用时间毫秒
	Exception      trace.ExceptionStack `es_type:"object"` // 异常信息
}

type TraceDetailDatabasePO struct {
	BaseTraceDetailPO
	DbName           string // 数据库名
	TableName        string // 表名
	Sql              string // SQL
	ConnectionString string // 连接字符串
	RowsAffected     int64  // 影响行数
}

type TraceDetailEsPO struct {
	trace.BaseTraceDetail
	IndexName   string // 索引名称
	AliasesName string // 别名
}
type TraceDetailEtcdPO struct {
	trace.BaseTraceDetail
	Key     string // etcd key
	LeaseID int64  // LeaseID
}

// TraceDetailHandPO 手动埋点
type TraceDetailHandPO struct {
	trace.BaseTraceDetail
	Name string // 名称
}
type TraceDetailHttpPO struct {
	trace.BaseTraceDetail
	Method string // post/get/put/delete
	Url    string // 请求url
}
type TraceDetailMqPO struct {
	trace.BaseTraceDetail
	Server     string // MQ服务器地址
	Exchange   string // 交换器名称
	RoutingKey string // 路由key
}
type TraceDetailRedisPO struct {
	trace.BaseTraceDetail
	Key   string // redis key
	Field string // hash field
}
