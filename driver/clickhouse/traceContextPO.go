package linkTrace_clickhouse

import (
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/linkTrace"
	"github.com/farseer-go/linkTrace/eumTraceType"
	"time"
)

type TraceContextPO struct {
	TraceId       int64                      `gorm:"not null;default:0;comment:上下文ID"`
	AppId         int64                      `gorm:"not null;default:0;comment:应用ID"`
	AppName       string                     `gorm:"not null;default:'';comment:应用名称"`
	AppIp         string                     `gorm:"not null;default:'';comment:应用IP"`
	ParentAppName string                     `gorm:"not null;default:'';comment:上游应用"`
	StartTs       int64                      `gorm:"not null;default:0;comment:调用开始时间戳"`
	EndTs         int64                      `gorm:"not null;default:0;comment:调用结束时间戳"`
	UseTs         time.Duration              `gorm:"not null;default:0;comment:总共使用时间毫秒"`
	TraceType     eumTraceType.Enum          `gorm:"not null;comment:状态码"`
	ignore        bool                       `gorm:"not null;default:false;comment:忽略这次的链路追踪"`
	Exception     *trace.ExceptionStack      `gorm:"json;not null;comment:异常信息"`
	List          []trace.ITraceDetail       `gorm:"json;not null;comment:调用的上下文"`
	Web           *linkTrace.WebContext      `gorm:"json;not null;comment:Web请求上下文"`
	Consumer      *linkTrace.ConsumerContext `gorm:"json;not null;comment:消费上下文"`
	Task          *linkTrace.TaskContext     `gorm:"json;not null;comment:任务上下文"`
	WatchKey      *linkTrace.WatchKeyContext `gorm:"json;not null;comment:Etcd上下文"`
}
