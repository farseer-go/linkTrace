package linkTrace

import (
	"fmt"
	"strings"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/linkTrace/eumTraceType"
	"github.com/farseer-go/queue"
)

type TraceContext struct {
	TraceId       string                `json:"tid"` // 上下文ID
	AppId         string                `json:"aid"` // 应用ID
	AppName       string                `json:"an"`  // 应用名称
	AppIp         string                `json:"aip"` // 应用IP
	ParentAppName string                `json:"pn"`  // 上游应用
	TraceLevel    int                   `json:"tl"`  // 逐层递增（显示上下游顺序）
	StartTs       int64                 `json:"st"`  // 调用开始时间戳（微秒）
	EndTs         int64                 `json:"et"`  // 调用结束时间戳（微秒）
	UseTs         time.Duration         `json:"ut"`  // 总共使用时间（微秒）
	UseDesc       string                `json:"ud"`  // 总共使用时间（描述）
	TraceType     eumTraceType.Enum     `json:"tt"`  // 状态码
	List          []any                 `json:"l"`   // 调用的上下文trace.ITraceDetail
	TraceCount    int                   `json:"tc"`  // 追踪明细数量
	ignore        bool                  // 忽略这次的链路追踪
	Exception     *trace.ExceptionStack `json:"e"` // 异常信息
	WebContext
	ConsumerContext
	TaskContext
	WatchKeyContext
	CreateAt dateTime.DateTime // 请求时间
}

type WebContext struct {
	WebDomain       string                                 `json:"wd"`  // 请求域名
	WebPath         string                                 `json:"wp"`  // 请求地址
	WebMethod       string                                 `json:"wm"`  // 请求方式
	WebContentType  string                                 `json:"wct"` // 请求内容类型
	WebStatusCode   int                                    `json:"wsc"` // 状态码
	WebHeaders      collections.Dictionary[string, string] `json:"wh"`  // 请求头部
	WebRequestBody  string                                 `json:"wrb"` // 请求参数
	WebResponseBody string                                 `json:"wpb"` // 输出参数
	WebRequestIp    string                                 `json:"wip"` // 客户端IP
}

func (receiver WebContext) IsNil() bool {
	return receiver.WebDomain == "" && receiver.WebPath == "" && receiver.WebMethod == "" && receiver.WebContentType == "" && receiver.WebStatusCode == 0
}

type ConsumerContext struct {
	ConsumerServer     string `json:"cs"` // MQ服务器
	ConsumerQueueName  string `json:"cq"` // 队列名称
	ConsumerRoutingKey string `json:"cr"` // 路由KEY
}

func (receiver ConsumerContext) IsNil() bool {
	return receiver.ConsumerServer == "" && receiver.ConsumerQueueName == "" && receiver.ConsumerRoutingKey == ""
}

type TaskContext struct {
	TaskName      string                                 `json:"tn"`  // 任务名称
	TaskGroupName string                                 `json:"tgn"` // 任务组ID
	TaskId        int64                                  `json:"tid"` // 任务ID
	TaskData      collections.Dictionary[string, string] `json:"td"`  // 任务数据
}

func (receiver TaskContext) IsNil() bool {
	return receiver.TaskName == "" && receiver.TaskGroupName == "" && receiver.TaskId == 0
}

type WatchKeyContext struct {
	WatchKey string `json:"wk"` // KEY
}

func (receiver WatchKeyContext) IsNil() bool {
	return receiver.WatchKey == ""
}

func (receiver *TraceContext) SetBody(requestBody string, statusCode int, responseBody string) {
	receiver.WebContext.WebRequestBody = requestBody
	receiver.WebContext.WebStatusCode = statusCode
	receiver.WebContext.WebResponseBody = responseBody
}

func (receiver *TraceContext) SetResponseBody(responseBody string) {
	receiver.WebContext.WebResponseBody = responseBody
}

func (receiver *TraceContext) GetTraceId() string {
	return receiver.TraceId
}
func (receiver *TraceContext) GetTraceLevel() int { return receiver.TraceLevel }
func (receiver *TraceContext) GetStartTs() int64 {
	return receiver.StartTs
}

// End 结束当前链路
func (receiver *TraceContext) End(err error) {
	// 清空当前上下文
	trace.CurTraceContext.Remove()

	if receiver.ignore {
		return
	}
	receiver.EndTs = time.Now().UnixMicro()
	receiver.UseTs = time.Duration(receiver.EndTs-receiver.StartTs) * time.Microsecond
	receiver.UseDesc = receiver.UseTs.String()
	// 移除忽略的明细
	var newList []any
	for _, detail := range receiver.List {
		if !detail.(trace.ITraceDetail).GetTraceDetail().IsIgnore() {
			newList = append(newList, detail)
		}
	}
	receiver.List = newList
	receiver.TraceCount = len(newList)

	// 启用了链路追踪后，把数据写入到本地队列中
	if defConfig.Enable {
		queue.Push("TraceContext", *receiver)
	}

	// 打印日志
	receiver.printLog()
}

func (receiver *TraceContext) Ignore() {
	receiver.ignore = true
}

// GetList 获取链路明细
func (receiver *TraceContext) GetList() []any {
	return receiver.List
}

// AddDetail 添加链路明细
func (receiver *TraceContext) AddDetail(detail trace.ITraceDetail) {
	receiver.List = append(receiver.List, detail)
}

// printLog 打印日志
func (receiver *TraceContext) printLog() {
	// 打印日志
	if defConfig.PrintLog {
		lst := collections.NewList[string]()
		for i := 0; i < len(receiver.List); i++ {
			tab := strings.Repeat("\t", receiver.List[i].(trace.ITraceDetail).GetLevel()-1)
			detail := receiver.List[i].(trace.ITraceDetail).GetTraceDetail()
			log := fmt.Sprintf("%s%s (%s)：%s", tab, flog.Blue(i+1), flog.Green(detail.UnTraceTs.String()), receiver.List[i].(trace.ITraceDetail).ToString())
			lst.Add(log)

			if detail.Exception != nil && detail.Exception.ExceptionIsException {
				lst.Add(fmt.Sprintf("%s:%s %s 出错了：%s", detail.Exception.ExceptionCallFile, flog.Blue(detail.Exception.ExceptionCallLine), flog.Red(detail.Exception.ExceptionCallFuncName), flog.Red(detail.Exception.ExceptionMessage)))
			}
		}

		if receiver.Exception != nil && receiver.Exception.ExceptionIsException {
			lst.Add(fmt.Sprintf("%s%s:%s %s %s", flog.Red("【异常】"), flog.Blue(receiver.Exception.ExceptionCallFile), flog.Blue(receiver.Exception.ExceptionCallLine), flog.Green(receiver.Exception.ExceptionCallFuncName), flog.Red(receiver.Exception.ExceptionMessage)))
		}

		lst.Add("-----------------------------------------------------------------")
		logs := strings.Join(lst.ToArray(), "\n")
		switch receiver.TraceType {
		case eumTraceType.WebApi:
			flog.Printf("【%s链路追踪】TraceId:%s，耗时：%s，%s\n%s\n", receiver.TraceType.ToString(), flog.Green(receiver.TraceId), flog.Red(receiver.UseTs.String()), receiver.WebContext.WebPath, logs)
		case eumTraceType.MqConsumer, eumTraceType.QueueConsumer, eumTraceType.EventConsumer:
			flog.Printf("【%s链路追踪】TraceId:%s，耗时：%s，%s\n%s\n", receiver.TraceType.ToString(), flog.Green(receiver.TraceId), flog.Red(receiver.UseTs.String()), receiver.ConsumerContext.ConsumerQueueName, logs)
		case eumTraceType.Task, eumTraceType.FSchedule:
			flog.Printf("【%s链路追踪】TraceId:%s，耗时：%s，%s %s\n%s\n", receiver.TraceType.ToString(), flog.Green(receiver.TraceId), flog.Red(receiver.UseTs.String()), receiver.TaskContext.TaskName, receiver.TaskContext.TaskGroupName, logs)
		default:
			flog.Printf("【%s链路追踪】TraceId:%s，耗时：%s\n%s\n", receiver.TraceType.ToString(), flog.Green(receiver.TraceId), flog.Red(receiver.UseTs.String()), logs)
		}
	}
}

func (receiver *TraceContext) Error(err error) {
	if err != nil {
		receiver.Exception = &trace.ExceptionStack{
			ExceptionIsException: true,
			ExceptionMessage:     err.Error(),
		}
		receiver.Exception.ExceptionCallFile, receiver.Exception.ExceptionCallFuncName, receiver.Exception.ExceptionCallLine = trace.GetCallerInfo()
	}
}

func (receiver *TraceContext) GetAppInfo() (string, string, string, string, string) {
	return receiver.TraceId, receiver.AppName, receiver.AppId, receiver.AppIp, receiver.ParentAppName
}
