package linkTrace_clickhouse

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/linkTrace/eumTraceType"
	"time"
)

type TraceContextPO struct {
	TraceId           int64                `gorm:"not null;default:0;comment:上下文ID"`
	AppId             int64                `gorm:"not null;default:0;comment:应用ID"`
	AppName           string               `gorm:"not null;default:'';comment:应用名称"`
	AppIp             string               `gorm:"not null;default:'';comment:应用IP"`
	ParentAppName     string               `gorm:"not null;default:'';comment:上游应用"`
	StartTs           int64                `gorm:"not null;default:0;comment:调用开始时间戳（微秒）"`
	EndTs             int64                `gorm:"not null;default:0;comment:调用结束时间戳"`
	UseTs             time.Duration        `gorm:"not null;default:0;comment:总共使用时间微秒"`
	TraceType         eumTraceType.Enum    `gorm:"not null;comment:状态码"`
	Exception         *ExceptionStackPO    `gorm:"json;not null;comment:异常信息"`
	List              []trace.ITraceDetail `gorm:"json;not null;comment:调用的上下文"`
	WebContextPO      `gorm:"embedded;embeddedPrefix:web_;not null;comment:Web请求上下文"`
	ConsumerContextPO `gorm:"embedded;embeddedPrefix:consumer_;not null;comment:消费上下文"`
	TaskContextPO     `gorm:"embedded;embeddedPrefix:task_;not null;comment:任务上下文"`
	WatchKeyContextPO `gorm:"embedded;embeddedPrefix:watchkey_;not null;comment:Etcd上下文"`
}

type WebContextPO struct {
	Domain       string                                 `gorm:"not null;default:'';comment:请求域名"`
	Path         string                                 `gorm:"not null;default:'';comment:请求地址"`
	Method       string                                 `gorm:"not null;default:'';comment:请求方式"`
	ContentType  string                                 `gorm:"not null;default:'';comment:请求内容类型"`
	StatusCode   int                                    `gorm:"not null;default:0;comment:状态码"`
	Headers      collections.Dictionary[string, string] `gorm:"type:String;json;not null;comment:请求头部"`
	RequestBody  string                                 `gorm:"not null;default:'';comment:请求参数"`
	ResponseBody string                                 `gorm:"not null;default:'';comment:输出参数"`
	RequestIp    string                                 `gorm:"not null;default:'';comment:客户端IP"`
}

type ConsumerContextPO struct {
	Server     string `gorm:"not null;default:'';comment:Server"`
	QueueName  string `gorm:"not null;default:'';comment:队列名称"`
	RoutingKey string `gorm:"not null;default:'';comment:路由KEY"`
}

type TaskContextPO struct {
	TaskName    string `gorm:"not null;default:'';comment:任务名称"`
	TaskGroupId int64  `gorm:"not null;default:0;comment:任务组ID"`
	TaskId      int64  `gorm:"not null;default:0;comment:任务ID"`
}

type WatchKeyContextPO struct {
	Key string `gorm:"not null;default:'';comment:KEY"`
}

type ExceptionStackPO struct {
	CallFile         string `gorm:"not null;default:'';comment:调用者文件路径"`
	CallLine         int    `gorm:"not null;default:0;comment:调用者行号"`
	CallFuncName     string `gorm:"not null;default:'';comment:调用者函数名称"`
	IsException      bool   `gorm:"not null;default:false;comment:是否执行异常"`
	ExceptionMessage string `gorm:"not null;default:'';comment:异常信息"`
}
