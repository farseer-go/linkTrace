package linkTrace

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/sonyflake"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/fs/trace/eumCallType"
	"github.com/farseer-go/fs/trace/eumTraceType"
	"github.com/farseer-go/queue"
)

type traceManager struct {
}

// EntryWebApi Webapi入口
func (*traceManager) EntryWebApi(domain string, path string, method string, contentType string, header map[string]string, requestIp string) *trace.TraceContext {
	headerDictionary := collections.NewDictionaryFromMap(header)
	traceId := parse.ToString(headerDictionary.GetValue("Trace-Id"))
	traceLevel := parse.ToInt(headerDictionary.GetValue("Trace-Level"))
	if traceId == "" {
		traceId = strconv.FormatInt(sonyflake.GenerateId(), 10)
	} else {
		traceLevel++ // 来自上游的请求，自动+1层
	}
	context := &trace.TraceContext{
		AppId:         strconv.FormatInt(core.AppId, 10),
		AppName:       core.AppName,
		AppIp:         core.AppIp,
		ParentAppName: headerDictionary.GetValue("Trace-App-Name"),
		TraceId:       traceId,
		TraceLevel:    traceLevel,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.WebApi,
		CreateAt:      dateTime.Now(),
		WebContext: trace.WebContext{
			WebDomain:      domain,
			WebPath:        path,
			WebMethod:      method,
			WebContentType: contentType,
			WebHeaders:     headerDictionary.ToDictionary(),
			WebRequestBody: "",
			WebRequestIp:   requestIp,
		},
	}
	trace.CurTraceContext.Set(context)
	trace.ScopeLevel.Set([]trace.BaseTraceDetail{})
	return context
}

// EntryWebSocket WebSocket入口
func (*traceManager) EntryWebSocket(domain string, path string, header map[string]string, requestIp string) *trace.TraceContext {
	headerDictionary := collections.NewDictionaryFromMap(header)
	parentTraceId := parse.ToString(headerDictionary.GetValue("Trace-Id"))
	traceLevel := parse.ToInt(headerDictionary.GetValue("Trace-Level"))
	if parentTraceId == "" {
		parentTraceId = strconv.FormatInt(sonyflake.GenerateId(), 10)
	} else {
		traceLevel++ // 来自上游的请求，自动+1层
	}
	context := &trace.TraceContext{
		AppId:         strconv.FormatInt(core.AppId, 10),
		AppName:       core.AppName,
		AppIp:         core.AppIp,
		ParentAppName: headerDictionary.GetValue("Trace-App-Name"),
		TraceId:       parentTraceId,
		TraceLevel:    traceLevel,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.WebSocket,
		CreateAt:      dateTime.Now(),
		WebContext: trace.WebContext{
			WebDomain:      domain,
			WebPath:        path,
			WebMethod:      "WEBSOCKET",
			WebContentType: "",
			WebHeaders:     headerDictionary.ToDictionary(),
			WebRequestBody: "Conn",
			WebRequestIp:   requestIp,
		},
	}
	trace.CurTraceContext.Set(context)
	trace.ScopeLevel.Set([]trace.BaseTraceDetail{})
	return context
}

// EntryMqConsumer mq 消费埋点
func (*traceManager) EntryMqConsumer(parentTraceId, parentAppName, server string, queueName string, routingKey string) *trace.TraceContext {
	// 如果来自上游，则要自动+1层，默认为0
	var traceLevel int
	if parentTraceId == "" {
		parentTraceId = strconv.FormatInt(sonyflake.GenerateId(), 10)
	} else {
		traceLevel++ // 来自上游的请求，自动+1层
	}
	context := &trace.TraceContext{
		AppId:         strconv.FormatInt(core.AppId, 10),
		AppName:       core.AppName,
		AppIp:         core.AppIp,
		ParentAppName: parentAppName,
		TraceId:       parentTraceId,
		TraceLevel:    traceLevel,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.MqConsumer,
		CreateAt:      dateTime.Now(),
		ConsumerContext: trace.ConsumerContext{
			ConsumerServer:     server,
			ConsumerQueueName:  queueName,
			ConsumerRoutingKey: routingKey,
		},
	}
	trace.CurTraceContext.Set(context)
	trace.ScopeLevel.Set([]trace.BaseTraceDetail{})
	return context
}

// EntryQueueConsumer queue 消费埋点
func (*traceManager) EntryQueueConsumer(queueName, subscribeName string) *trace.TraceContext {
	context := &trace.TraceContext{
		AppId:         strconv.FormatInt(core.AppId, 10),
		AppName:       core.AppName,
		AppIp:         core.AppIp,
		ParentAppName: "",
		TraceId:       strconv.FormatInt(sonyflake.GenerateId(), 10),
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.QueueConsumer,
		CreateAt:      dateTime.Now(),
		ConsumerContext: trace.ConsumerContext{
			ConsumerServer:    fmt.Sprintf("本地Queue/%s/%s/%v", core.AppName, core.AppIp, core.AppId),
			ConsumerQueueName: queueName + "/" + subscribeName,
		},
	}
	trace.CurTraceContext.Set(context)
	trace.ScopeLevel.Set([]trace.BaseTraceDetail{})
	return context
}

// EntryEventConsumer event 事件消费埋点
func (receiver *traceManager) EntryEventConsumer(server, eventName, subscribeName string) *trace.TraceContext {
	// 事件消费，一般是由其它入口的程序触发的，所以这里先看能不能取到之前的上下文
	var traceId string
	var traceLevel int
	var parentAppName string
	if cur := trace.CurTraceContext.Get(); cur != nil {
		traceId, parentAppName, _, _, _ = cur.GetAppInfo()
		traceLevel = cur.TraceLevel + 1
	} else {
		traceId = strconv.FormatInt(sonyflake.GenerateId(), 10)
		parentAppName = core.AppName
	}
	context := &trace.TraceContext{
		AppId:         strconv.FormatInt(core.AppId, 10),
		AppName:       core.AppName,
		AppIp:         core.AppIp,
		ParentAppName: parentAppName,
		TraceId:       traceId,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.EventConsumer,
		TraceLevel:    traceLevel,
		CreateAt:      dateTime.Now(),
		ConsumerContext: trace.ConsumerContext{
			ConsumerServer:    server,
			ConsumerQueueName: eventName + "/" + subscribeName,
		},
	}
	trace.CurTraceContext.Set(context)
	trace.ScopeLevel.Set([]trace.BaseTraceDetail{})
	return context
}

// EntryTask 创建本地任务入口
func (*traceManager) EntryTask(taskName string) *trace.TraceContext {
	traceId := strconv.FormatInt(sonyflake.GenerateId(), 10)
	context := &trace.TraceContext{
		AppId:         strconv.FormatInt(core.AppId, 10),
		AppName:       core.AppName,
		AppIp:         core.AppIp,
		ParentAppName: "",
		TraceId:       traceId,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.Task,
		CreateAt:      dateTime.Now(),
		TaskContext: trace.TaskContext{
			TaskName: taskName,
		},
	}
	trace.CurTraceContext.Set(context)
	trace.ScopeLevel.Set([]trace.BaseTraceDetail{})
	return context
}

// EntryTaskGroup 创建本地任务入口（调度中心专用）
func (*traceManager) EntryTaskGroup(taskName string, taskGroupName string, taskId int64) *trace.TraceContext {
	traceId := strconv.FormatInt(sonyflake.GenerateId(), 10)
	context := &trace.TraceContext{
		AppId:         strconv.FormatInt(core.AppId, 10),
		AppName:       core.AppName,
		AppIp:         core.AppIp,
		ParentAppName: "",
		TraceId:       traceId,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.Task,
		CreateAt:      dateTime.Now(),
		TaskContext: trace.TaskContext{
			TaskName:      fmt.Sprintf("%s，任务组=%s，任务ID=%v", taskName, taskGroupName, taskId),
			TaskGroupName: taskGroupName,
			TaskId:        taskId,
		},
	}
	trace.CurTraceContext.Set(context)
	trace.ScopeLevel.Set([]trace.BaseTraceDetail{})
	return context
}

// EntryFSchedule 创建调度中心入口
func (*traceManager) EntryFSchedule(taskGroupName string, taskId int64, data map[string]string) *trace.TraceContext {
	traceId := strconv.FormatInt(sonyflake.GenerateId(), 10)
	context := &trace.TraceContext{
		AppId:         strconv.FormatInt(core.AppId, 10),
		AppName:       core.AppName,
		AppIp:         core.AppIp,
		ParentAppName: "",
		TraceId:       traceId,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.FSchedule,
		CreateAt:      dateTime.Now(),
		TaskContext: trace.TaskContext{
			TaskName:      taskGroupName,
			TaskGroupName: taskGroupName,
			TaskId:        taskId,
			TaskData:      collections.NewDictionaryFromMap(data),
		},
	}
	trace.CurTraceContext.Set(context)
	trace.ScopeLevel.Set([]trace.BaseTraceDetail{})
	return context
}

// EntryWatchKey 创建etcd入口
func (*traceManager) EntryWatchKey(key string) *trace.TraceContext {
	traceId := strconv.FormatInt(sonyflake.GenerateId(), 10)
	context := &trace.TraceContext{
		AppId:         strconv.FormatInt(core.AppId, 10),
		AppName:       core.AppName,
		AppIp:         core.AppIp,
		ParentAppName: "",
		TraceId:       traceId,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.WatchKey,
		CreateAt:      dateTime.Now(),
		WatchKeyContext: trace.WatchKeyContext{
			WatchKey: key,
		},
	}
	trace.CurTraceContext.Set(context)
	trace.ScopeLevel.Set([]trace.BaseTraceDetail{})
	return context
}

// TraceDatabase 数据库埋点
func (*traceManager) TraceDatabase() trace.ITraceDetail {
	detail := &TraceDetailDatabase{
		BaseTraceDetail: trace.NewTraceDetail(eumCallType.Database, ""),
	}
	add(detail)
	return detail
}

// TraceDatabaseOpen 数据库埋点
func (*traceManager) TraceDatabaseOpen(dbName string, connectString string) trace.ITraceDetail {
	detail := &TraceDetailDatabase{
		BaseTraceDetail:  trace.NewTraceDetail(eumCallType.Database, ""),
		DbName:           dbName,
		ConnectionString: connectString,
	}
	add(detail)
	return detail
}

// TraceElasticsearch Elasticsearch埋点
func (*traceManager) TraceElasticsearch(method string, IndexName string, AliasesName string) trace.ITraceDetail {
	detail := &TraceDetailEs{
		BaseTraceDetail: trace.NewTraceDetail(eumCallType.Elasticsearch, method),
		IndexName:       IndexName,
		AliasesName:     AliasesName,
	}
	add(detail)
	return detail
}

// TraceEtcd etcd埋点
func (*traceManager) TraceEtcd(method string, key string, leaseID int64) trace.ITraceDetail {
	detail := &TraceDetailEtcd{
		BaseTraceDetail: trace.NewTraceDetail(eumCallType.Etcd, method),
		Key:             key,
		LeaseID:         leaseID,
	}
	add(detail)
	return detail
}

// TraceHand 手动埋点
func (*traceManager) TraceHand(name string) trace.ITraceDetail {
	detail := &trace.TraceDetailHand{
		BaseTraceDetail: trace.NewTraceDetail(eumCallType.Hand, ""),
		Name:            name,
	}
	add(detail)
	return detail
}

// TraceEventPublish 事件发布
func (*traceManager) TraceEventPublish(eventName string) trace.ITraceDetail {
	detail := &TraceDetailEventConsumer{
		BaseTraceDetail: trace.NewTraceDetail(eumCallType.EventPublish, ""),
		Name:            eventName,
	}
	add(detail)
	return detail
}

// TraceMqSend mq发送埋点
func (*traceManager) TraceMqSend(method string, server string, exchange string, routingKey string) trace.ITraceDetail {
	detail := &TraceDetailMq{
		BaseTraceDetail: trace.NewTraceDetail(eumCallType.Mq, method),
		Server:          server,
		Exchange:        exchange,
		RoutingKey:      routingKey,
	}
	add(detail)
	return detail
}

// TraceMq open、create埋点
func (*traceManager) TraceMq(method string, server string, exchange string) trace.ITraceDetail {
	detail := &TraceDetailMq{
		BaseTraceDetail: trace.NewTraceDetail(eumCallType.Mq, method),
		Server:          server,
		Exchange:        exchange,
	}
	add(detail)
	return detail
}

// TraceRedis Redis埋点
func (*traceManager) TraceRedis(method string, key string, field string) trace.ITraceDetail {
	detail := &TraceDetailRedis{
		BaseTraceDetail: trace.NewTraceDetail(eumCallType.Redis, method),
		Key:             key,
		Field:           field,
	}
	add(detail)
	return detail
}

// TraceHttp http埋点
func (*traceManager) TraceHttp(method string, url string) trace.ITraceDetail {
	detail := &TraceDetailHttp{
		BaseTraceDetail: trace.NewTraceDetail(eumCallType.Http, method),
		Method:          method,
		Url:             url,
	}
	add(detail)
	return detail
}

// TraceGrpc grpc埋点
func (*traceManager) TraceGrpc(method string, url string) trace.ITraceDetail {
	detail := &TraceDetailGrpc{
		BaseTraceDetail: trace.NewTraceDetail(eumCallType.Grpc, method),
		Method:          method,
		Url:             url,
	}
	add(detail)
	return detail
}

func add(traceDetail trace.ITraceDetail) {
	if t := trace.CurTraceContext.Get(); t != nil {
		detail := traceDetail.GetTraceDetail()
		// 时间轴：上下文入口起点时间到本次开始时间
		detail.Timeline = time.Duration(detail.StartTs-t.StartTs) * time.Microsecond
		details := t.List
		if len(details) > 0 {
			detail.UnTraceTs = time.Duration(detail.StartTs-details[len(details)-1].(trace.ITraceDetail).GetTraceDetail().EndTs) * time.Microsecond
		} else {
			detail.UnTraceTs = time.Duration(detail.StartTs-t.StartTs) * time.Microsecond
		}
		detail.TraceId, detail.AppName, detail.AppId, detail.AppIp, detail.ParentAppName = t.GetAppInfo()
		t.AddDetail(traceDetail)
	}
}

// End 结束当前链路
func (receiver *traceManager) Push(traceContext *trace.TraceContext, err error) {
	// 清空当前上下文
	defer trace.CurTraceContext.Remove()

	if traceContext.IsIgnore() {
		return
	}
	traceContext.Error(err)
	traceContext.EndTs = time.Now().UnixMicro()
	traceContext.UseTs = time.Duration(traceContext.EndTs-traceContext.StartTs) * time.Microsecond
	traceContext.UseDesc = traceContext.UseTs.String()

	var isError bool
	// 移除忽略的明细
	var newList []any
	for _, detail := range traceContext.List {
		traceDetail := detail.(trace.ITraceDetail).GetTraceDetail()
		// 打印日志，供上传到FOPS
		if traceDetail.Exception != nil {
			isError = true
			flog.Errorf("%s %s %s:%d %s", traceDetail.CallType.ToString(), traceDetail.Exception.ExceptionCallFile, traceDetail.Exception.ExceptionCallFuncName, traceDetail.Exception.ExceptionCallLine, traceDetail.Exception.ExceptionMessage)
		}
		if !traceDetail.IsIgnore() {
			newList = append(newList, detail)
		}
	}

	// 打印日志，供上传到FOPS
	// 判断是否有异常,如果有异常，就要把异常信息打印到控制台
	if !isError && traceContext.Exception != nil {
		flog.Errorf("%s %s %s:%d %s", traceContext.TraceType.ToString(), traceContext.Exception.ExceptionCallFile, traceContext.Exception.ExceptionCallFuncName, traceContext.Exception.ExceptionCallLine, traceContext.Exception.ExceptionMessage)
	}
	traceContext.List = newList
	traceContext.TraceCount = len(newList)

	// 启用了链路追踪后，把数据写入到本地队列中
	if defConfig.Enable {
		queue.Push("TraceContext", traceContext)
	}

	// 打印日志
	if defConfig.PrintLog { //  || isError
		lst := collections.NewList[string]()
		for i := 0; i < len(traceContext.List); i++ {
			traceDetail := traceContext.List[i].(trace.ITraceDetail)
			tab := strings.Repeat("\t", traceDetail.GetLevel()-1)
			detail := traceDetail.GetTraceDetail()
			log := fmt.Sprintf("%s%s (%s)：%s", tab, flog.Blue(i+1), flog.Green(detail.UnTraceTs.String()), traceDetail.ToString())
			lst.Add(log)

			if detail.Exception != nil && detail.Exception.ExceptionIsException {
				lst.Add(fmt.Sprintf("%s:%s %s 出错了：%s", detail.Exception.ExceptionCallFile, flog.Blue(detail.Exception.ExceptionCallLine), flog.Red(detail.Exception.ExceptionCallFuncName), flog.Red(detail.Exception.ExceptionMessage)))
			}
		}

		if traceContext.Exception != nil && traceContext.Exception.ExceptionIsException {
			lst.Add(fmt.Sprintf("%s%s:%s %s %s", flog.Red("【异常】"), flog.Blue(traceContext.Exception.ExceptionCallFile), flog.Blue(traceContext.Exception.ExceptionCallLine), flog.Green(traceContext.Exception.ExceptionCallFuncName), flog.Red(traceContext.Exception.ExceptionMessage)))
		}

		lst.Add("-----------------------------------------------------------------")
		logs := strings.Join(lst.ToArray(), "\n")
		switch traceContext.TraceType {
		case eumTraceType.WebApi, eumTraceType.WebSocket:
			flog.Printf("【%s链路追踪】TraceId:%s，耗时：%s，%s\n%s\n", traceContext.TraceType.ToString(), flog.Green(traceContext.TraceId), flog.Red(traceContext.UseTs.String()), traceContext.WebContext.WebPath, logs)
		case eumTraceType.MqConsumer, eumTraceType.QueueConsumer, eumTraceType.EventConsumer:
			flog.Printf("【%s链路追踪】TraceId:%s，耗时：%s，%s\n%s\n", traceContext.TraceType.ToString(), flog.Green(traceContext.TraceId), flog.Red(traceContext.UseTs.String()), traceContext.ConsumerContext.ConsumerQueueName, logs)
		case eumTraceType.Task, eumTraceType.FSchedule:
			flog.Printf("【%s链路追踪】TraceId:%s，耗时：%s，%s %s\n%s\n", traceContext.TraceType.ToString(), flog.Green(traceContext.TraceId), flog.Red(traceContext.UseTs.String()), traceContext.TaskContext.TaskName, traceContext.TaskContext.TaskGroupName, logs)
		default:
			flog.Printf("【%s链路追踪】TraceId:%s，耗时：%s\n%s\n", traceContext.TraceType.ToString(), flog.Green(traceContext.TraceId), flog.Red(traceContext.UseTs.String()), logs)
		}
	}
}
