package linkTrace

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/linkTrace/eumTraceType"
	"github.com/farseer-go/queue"
	"strings"
	"time"
)

type TraceContext struct {
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
	ignore        bool                 // 忽略这次的链路追踪
	Exception     trace.ExceptionStack // 异常信息
	Web           WebContext
	Consumer      ConsumerContext
	Task          TaskContext
	WatchKey      WatchKeyContext
}

type WebContext struct {
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

func (receiver WebContext) IsNil() bool {
	return receiver.Domain == "" && receiver.Path == "" && receiver.Method == "" && receiver.ContentType == "" && receiver.StatusCode == 0
}

type ConsumerContext struct {
	Server     string
	QueueName  string
	RoutingKey string
}

func (receiver ConsumerContext) IsNil() bool {
	return receiver.Server == "" && receiver.QueueName == "" && receiver.RoutingKey == ""
}

type TaskContext struct {
	TaskName    string
	TaskGroupId int64
	TaskId      int64
}

func (receiver TaskContext) IsNil() bool {
	return receiver.TaskName == "" && receiver.TaskGroupId == 0 && receiver.TaskId == 0
}

type WatchKeyContext struct {
	Key string
}

func (receiver WatchKeyContext) IsNil() bool {
	return receiver.Key == ""
}

func (receiver *TraceContext) SetBody(requestBody string, statusCode int, responseBody string) {
	receiver.Web.RequestBody = requestBody
	receiver.Web.StatusCode = statusCode
	receiver.Web.ResponseBody = responseBody
}

func (receiver *TraceContext) GetTraceId() int64 {
	return receiver.TraceId
}

func (receiver *TraceContext) GetStartTs() int64 {
	return receiver.StartTs
}

// End 结束当前链路
func (receiver *TraceContext) End() {
	if receiver.ignore {
		return
	}
	receiver.EndTs = time.Now().UnixMicro()
	receiver.UseTs = time.Duration(receiver.EndTs-receiver.StartTs) * time.Microsecond
	// 移除忽略的明细
	var newList []trace.ITraceDetail
	for _, detail := range receiver.List {
		if !detail.GetTraceDetail().IsIgnore() {
			newList = append(newList, detail)
		}
	}
	receiver.List = newList

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
func (receiver *TraceContext) GetList() []trace.ITraceDetail {
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
			tab := strings.Repeat("\t", receiver.List[i].GetLevel()-1)
			detail := receiver.List[i].GetTraceDetail()
			log := fmt.Sprintf("%s%s (%s)：%s", tab, flog.Blue(i+1), flog.Green(detail.UnTraceTs.String()), receiver.List[i].ToString())
			lst.Add(log)

			if detail.Exception.IsException {
				lst.Add(fmt.Sprintf("%s:%s %s 出错了：%s", detail.Exception.CallFile, flog.Blue(detail.Exception.CallLine), flog.Red(detail.Exception.CallFuncName), flog.Red(detail.Exception.ExceptionMessage)))
			}
		}

		if receiver.Exception.IsException {
			lst.Add(fmt.Sprintf("%s%s:%s %s %s", flog.Red("【异常】"), flog.Blue(receiver.Exception.CallFile), flog.Blue(receiver.Exception.CallLine), flog.Green(receiver.Exception.CallFuncName), flog.Red(receiver.Exception.ExceptionMessage)))
		}

		lst.Add("-----------------------------------------------------------------")
		logs := strings.Join(lst.ToArray(), "\n")
		switch receiver.TraceType {
		case eumTraceType.WebApi:
			flog.Printf("【%s链路追踪】TraceId:%s，耗时：%s，%s\n%s\n", receiver.TraceType.ToString(), flog.Green(parse.ToString(receiver.TraceId)), flog.Red(receiver.UseTs.String()), receiver.Web.Path, logs)
		case eumTraceType.MqConsumer, eumTraceType.QueueConsumer:
			flog.Printf("【%s链路追踪】TraceId:%s，耗时：%s，%s\n%s\n", receiver.TraceType.ToString(), flog.Green(parse.ToString(receiver.TraceId)), flog.Red(receiver.UseTs.String()), receiver.Consumer.QueueName, logs)
		case eumTraceType.Task, eumTraceType.FSchedule:
			flog.Printf("【%s链路追踪】TraceId:%s，耗时：%s，%s\n%s\n", receiver.TraceType.ToString(), flog.Green(parse.ToString(receiver.TraceId)), flog.Red(receiver.UseTs.String()), receiver.Task.TaskName, logs)
		default:
			flog.Printf("【%s链路追踪】TraceId:%s，耗时：%s\n%s\n", receiver.TraceType.ToString(), flog.Green(parse.ToString(receiver.TraceId)), flog.Red(receiver.UseTs.String()), logs)
		}
	}
}

func (receiver *TraceContext) Error(err error) {
	if err != nil {
		receiver.Exception.IsException = true
		receiver.Exception.ExceptionMessage = err.Error()
		receiver.Exception.CallFile, receiver.Exception.CallFuncName, receiver.Exception.CallLine = trace.GetCallerInfo()
	}
}
