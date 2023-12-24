package linkTrace_elasticSearch

import (
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/linkTrace"
	"github.com/farseer-go/linkTrace/eumTraceType"
	"time"
)

type TraceContextPO struct {
	TraceId       int64                      `es:"primaryKey"` // 上下文ID
	AppId         int64                      // 应用ID
	AppName       string                     // 应用名称
	AppIp         string                     // 应用IP
	ParentAppName string                     // 上游应用
	StartTs       int64                      // 调用开始时间戳
	EndTs         int64                      // 调用结束时间戳
	UseTs         time.Duration              // 总共使用时间毫秒
	TraceType     eumTraceType.Enum          // 状态码
	List          []trace.ITraceDetail       `es_type:"object"` // 调用的上下文
	Exception     *trace.ExceptionStack      `es_type:"object"` // 异常信息
	Web           *linkTrace.WebContext      `es_type:"object"` // Web请求上下文
	Consumer      *linkTrace.ConsumerContext `es_type:"object"` // 消费上下文
	Task          *linkTrace.TaskContext     `es_type:"object"` // 任务上下文
	WatchKey      *linkTrace.WatchKeyContext `es_type:"object"` // Etcd上下文
}
