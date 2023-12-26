package linkTrace_elasticSearch

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/linkTrace/eumTraceType"
	"time"
)

type TraceContextPO struct {
	TraceId       int64                `es:"primaryKey"` // 上下文ID
	AppId         int64                // 应用ID
	AppName       string               // 应用名称
	AppIp         string               // 应用IP
	ParentAppName string               // 上游应用
	StartTs       int64                // 调用开始时间戳（微秒）
	EndTs         int64                // 调用结束时间戳
	UseTs         time.Duration        // 总共使用时间微秒
	TraceType     eumTraceType.Enum    // 状态码
	List          []trace.ITraceDetail `es_type:"object"` // 调用的上下文
	Exception     *ExceptionStackPO    `es_type:"object"` // 异常信息
	Web           WebContextPO         `es_type:"object"` // Web请求上下文
	Consumer      ConsumerContextPO    `es_type:"object"` // 消费上下文
	Task          TaskContextPO        `es_type:"object"` // 任务上下文
	WatchKey      WatchKeyContextPO    `es_type:"object"` // Etcd上下文
}

type WebContextPO struct {
	Domain       string                                 // 请求域名
	Path         string                                 `es_type:"text"` // 请求地址
	Method       string                                 // 请求方式
	ContentType  string                                 // 请求内容类型
	StatusCode   int                                    // 状态码
	Headers      collections.Dictionary[string, string] `es_type:"flattened"` // 请求头部
	RequestBody  string                                 `es_type:"text"`      // 请求参数
	ResponseBody string                                 `es_type:"text"`      // 输出参数
	RequestIp    string                                 // 客户端IP
}

type ConsumerContextPO struct {
	Server     string
	QueueName  string
	RoutingKey string
}

type TaskContextPO struct {
	TaskName    string
	TaskGroupId int64
	TaskId      int64
}

type WatchKeyContextPO struct {
	Key string
}

type ExceptionStackPO struct {
	CallFile         string // 调用者文件路径
	CallLine         int    // 调用者行号
	CallFuncName     string // 调用者函数名称
	IsException      bool   // 是否执行异常
	ExceptionMessage string // 异常信息
}
