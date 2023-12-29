package linkTrace_clickhouse

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/linkTrace/eumTraceType"
	"time"
)

type TraceContextPO struct {
	TraceId           int64             `gorm:"not null;default:0;comment:上下文ID"`
	AppId             int64             `gorm:"not null;default:0;comment:应用ID"`
	AppName           string            `gorm:"not null;default:'';comment:应用名称"`
	AppIp             string            `gorm:"not null;default:'';comment:应用IP"`
	ParentAppName     string            `gorm:"not null;default:'';comment:上游应用"`
	StartTs           int64             `gorm:"not null;default:0;comment:调用开始时间戳（微秒）"`
	EndTs             int64             `gorm:"not null;default:0;comment:调用结束时间戳（微秒）"`
	UseTs             time.Duration     `gorm:"not null;default:0;comment:总共使用时间（微秒）"`
	UseDesc           string            `gorm:"not null;default:'';comment:总共使用时间（描述）"`
	TraceType         eumTraceType.Enum `gorm:"not null;comment:入口类型"`
	Exception         *ExceptionStackPO `gorm:"json;not null;comment:异常信息"`
	List              []any             `gorm:"json;not null;comment:调用的上下文"`
	WebContextPO      `gorm:"embedded;not null;comment:Web请求上下文"`
	ConsumerContextPO `gorm:"embedded;not null;comment:消费上下文"`
	TaskContextPO     `gorm:"embedded;not null;comment:任务上下文"`
	WatchKeyContextPO `gorm:"embedded;not null;comment:Etcd上下文"`
	CreateAt          dateTime.DateTime `gorm:"type:DateTime64(3);not null;comment:请求时间"`
}

type WebContextPO struct {
	WebDomain       string                                 `gorm:"not null;default:'';comment:请求域名"`
	WebPath         string                                 `gorm:"not null;default:'';comment:请求地址"`
	WebMethod       string                                 `gorm:"not null;default:'';comment:请求方式"`
	WebContentType  string                                 `gorm:"not null;default:'';comment:请求内容类型"`
	WebStatusCode   int                                    `gorm:"not null;default:0;comment:状态码"`
	WebHeaders      collections.Dictionary[string, string] `gorm:"type:String;json;not null;comment:请求头部"`
	WebRequestBody  string                                 `gorm:"not null;default:'';comment:请求参数"`
	WebResponseBody string                                 `gorm:"not null;default:'';comment:输出参数"`
	WebRequestIp    string                                 `gorm:"not null;default:'';comment:客户端IP"`
}

type ConsumerContextPO struct {
	ConsumerServer     string `gorm:"not null;default:'';comment:ConsumerServer"`
	ConsumerQueueName  string `gorm:"not null;default:'';comment:队列名称"`
	ConsumerRoutingKey string `gorm:"not null;default:'';comment:路由KEY"`
}

type TaskContextPO struct {
	TaskName    string `gorm:"not null;default:'';comment:任务名称"`
	TaskGroupId int64  `gorm:"not null;default:0;comment:任务组ID"`
	TaskId      int64  `gorm:"not null;default:0;comment:任务ID"`
}

type WatchKeyContextPO struct {
	WatchKey string `gorm:"not null;default:'';comment:KEY"`
}

type ExceptionStackPO struct {
	ExceptionCallFile     string `gorm:"not null;default:'';comment:调用者文件路径"`
	ExceptionCallLine     int    `gorm:"not null;default:0;comment:调用者行号"`
	ExceptionCallFuncName string `gorm:"not null;default:'';comment:调用者函数名称"`
	ExceptionIsException  bool   `gorm:"not null;default:false;comment:是否执行异常"`
	ExceptionMessage      string `gorm:"not null;default:'';comment:异常信息"`
}
