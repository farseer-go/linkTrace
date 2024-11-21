package linkTrace

import (
	"fmt"
	"strconv"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/sonyflake"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/fs/trace/eumCallType"
	"github.com/farseer-go/linkTrace/eumTraceType"
)

type traceManager struct {
}

func (*traceManager) GetCurTrace() trace.ITraceContext {
	return trace.CurTraceContext.Get()
}

func (*traceManager) GetTraceId() string {
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
		return traceContext.GetTraceId()
	}
	return ""
}

// EntryWebApi Webapi入口
func (*traceManager) EntryWebApi(domain string, path string, method string, contentType string, header map[string]string, requestIp string) trace.ITraceContext {
	headerDictionary := collections.NewDictionaryFromMap(header)
	traceId := parse.ToString(headerDictionary.GetValue("Trace-Id"))
	traceLevel := parse.ToInt(headerDictionary.GetValue("Trace-Level"))
	if traceId == "" {
		traceId = strconv.FormatInt(sonyflake.GenerateId(), 10)
	} else {
		traceLevel++ // 来自上游的请求，自动+1层
	}
	context := &TraceContext{
		AppId:         strconv.FormatInt(core.AppId, 10),
		AppName:       core.AppName,
		AppIp:         core.AppIp,
		ParentAppName: headerDictionary.GetValue("Trace-App-Name"),
		TraceId:       traceId,
		TraceLevel:    traceLevel,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.WebApi,
		WebContext: WebContext{
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
func (*traceManager) EntryWebSocket(domain string, path string, header map[string]string, requestIp string) trace.ITraceContext {
	headerDictionary := collections.NewDictionaryFromMap(header)
	parentTraceId := parse.ToString(headerDictionary.GetValue("Trace-Id"))
	traceLevel := parse.ToInt(headerDictionary.GetValue("Trace-Level"))
	if parentTraceId == "" {
		parentTraceId = strconv.FormatInt(sonyflake.GenerateId(), 10)
	} else {
		traceLevel++ // 来自上游的请求，自动+1层
	}
	context := &TraceContext{
		AppId:         strconv.FormatInt(core.AppId, 10),
		AppName:       core.AppName,
		AppIp:         core.AppIp,
		ParentAppName: headerDictionary.GetValue("Trace-App-Name"),
		TraceId:       parentTraceId,
		TraceLevel:    traceLevel,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.WebSocket,
		WebContext: WebContext{
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
func (*traceManager) EntryMqConsumer(parentTraceId, parentAppName, server string, queueName string, routingKey string) trace.ITraceContext {
	// 如果来自上游，则要自动+1层，默认为0
	var traceLevel int
	if parentTraceId == "" {
		parentTraceId = strconv.FormatInt(sonyflake.GenerateId(), 10)
	} else {
		traceLevel++ // 来自上游的请求，自动+1层
	}
	context := &TraceContext{
		AppId:         strconv.FormatInt(core.AppId, 10),
		AppName:       core.AppName,
		AppIp:         core.AppIp,
		ParentAppName: parentAppName,
		TraceId:       parentTraceId,
		TraceLevel:    traceLevel,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.MqConsumer,
		ConsumerContext: ConsumerContext{
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
func (*traceManager) EntryQueueConsumer(queueName, subscribeName string) trace.ITraceContext {
	context := &TraceContext{
		AppId:         strconv.FormatInt(core.AppId, 10),
		AppName:       core.AppName,
		AppIp:         core.AppIp,
		ParentAppName: "",
		TraceId:       strconv.FormatInt(sonyflake.GenerateId(), 10),
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.QueueConsumer,
		ConsumerContext: ConsumerContext{
			ConsumerServer:    fmt.Sprintf("本地Queue/%s/%s/%v", core.AppName, core.AppIp, core.AppId),
			ConsumerQueueName: queueName + "/" + subscribeName,
		},
	}
	trace.CurTraceContext.Set(context)
	trace.ScopeLevel.Set([]trace.BaseTraceDetail{})
	return context
}

// EntryEventConsumer event 事件消费埋点
func (receiver *traceManager) EntryEventConsumer(server, eventName, subscribeName string) trace.ITraceContext {
	// 事件消费，一般是由其它入口的程序触发的，所以这里先看能不能取到之前的上下文
	var traceId string
	var traceLevel int
	var parentAppName string
	if cur := receiver.GetCurTrace(); cur != nil {
		traceId, parentAppName, _, _, _ = cur.GetAppInfo()
		traceLevel = cur.GetTraceLevel() + 1
	} else {
		traceId = strconv.FormatInt(sonyflake.GenerateId(), 10)
		parentAppName = core.AppName
	}
	context := &TraceContext{
		AppId:         strconv.FormatInt(core.AppId, 10),
		AppName:       core.AppName,
		AppIp:         core.AppIp,
		ParentAppName: parentAppName,
		TraceId:       traceId,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.EventConsumer,
		TraceLevel:    traceLevel,
		ConsumerContext: ConsumerContext{
			ConsumerServer:    server,
			ConsumerQueueName: eventName + "/" + subscribeName,
		},
	}
	trace.CurTraceContext.Set(context)
	trace.ScopeLevel.Set([]trace.BaseTraceDetail{})
	return context
}

// EntryTask 创建本地任务入口
func (*traceManager) EntryTask(taskName string) trace.ITraceContext {
	traceId := strconv.FormatInt(sonyflake.GenerateId(), 10)
	context := &TraceContext{
		AppId:         strconv.FormatInt(core.AppId, 10),
		AppName:       core.AppName,
		AppIp:         core.AppIp,
		ParentAppName: "",
		TraceId:       traceId,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.Task,
		TaskContext: TaskContext{
			TaskName: taskName,
		},
	}
	trace.CurTraceContext.Set(context)
	trace.ScopeLevel.Set([]trace.BaseTraceDetail{})
	return context
}

// EntryTaskGroup 创建本地任务入口（调度中心专用）
func (*traceManager) EntryTaskGroup(taskName string, taskGroupName string, taskId int64) trace.ITraceContext {
	traceId := strconv.FormatInt(sonyflake.GenerateId(), 10)
	context := &TraceContext{
		AppId:         strconv.FormatInt(core.AppId, 10),
		AppName:       core.AppName,
		AppIp:         core.AppIp,
		ParentAppName: "",
		TraceId:       traceId,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.Task,
		TaskContext: TaskContext{
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
func (*traceManager) EntryFSchedule(taskGroupName string, taskId int64, data map[string]string) trace.ITraceContext {
	traceId := strconv.FormatInt(sonyflake.GenerateId(), 10)
	context := &TraceContext{
		AppId:         strconv.FormatInt(core.AppId, 10),
		AppName:       core.AppName,
		AppIp:         core.AppIp,
		ParentAppName: "",
		TraceId:       traceId,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.FSchedule,
		TaskContext: TaskContext{
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
func (*traceManager) EntryWatchKey(key string) trace.ITraceContext {
	traceId := strconv.FormatInt(sonyflake.GenerateId(), 10)
	context := &TraceContext{
		AppId:         strconv.FormatInt(core.AppId, 10),
		AppName:       core.AppName,
		AppIp:         core.AppIp,
		ParentAppName: "",
		TraceId:       traceId,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.WatchKey,
		WatchKeyContext: WatchKeyContext{
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
		BaseTraceDetail: newTraceDetail(eumCallType.Database, ""),
	}
	add(detail)
	return detail
}

// TraceDatabaseOpen 数据库埋点
func (*traceManager) TraceDatabaseOpen(dbName string, connectString string) trace.ITraceDetail {
	detail := &TraceDetailDatabase{
		BaseTraceDetail:  newTraceDetail(eumCallType.Database, ""),
		DbName:           dbName,
		ConnectionString: connectString,
	}
	add(detail)
	return detail
}

// TraceElasticsearch Elasticsearch埋点
func (*traceManager) TraceElasticsearch(method string, IndexName string, AliasesName string) trace.ITraceDetail {
	detail := &TraceDetailEs{
		BaseTraceDetail: newTraceDetail(eumCallType.Elasticsearch, method),
		IndexName:       IndexName,
		AliasesName:     AliasesName,
	}
	add(detail)
	return detail
}

// TraceEtcd etcd埋点
func (*traceManager) TraceEtcd(method string, key string, leaseID int64) trace.ITraceDetail {
	detail := &TraceDetailEtcd{
		BaseTraceDetail: newTraceDetail(eumCallType.Etcd, method),
		Key:             key,
		LeaseID:         leaseID,
	}
	add(detail)
	return detail
}

// TraceHand 手动埋点
func (*traceManager) TraceHand(name string) trace.ITraceDetail {
	detail := &TraceDetailHand{
		BaseTraceDetail: newTraceDetail(eumCallType.Hand, ""),
		Name:            name,
	}
	add(detail)
	return detail
}

// TraceEventPublish 事件发布
func (*traceManager) TraceEventPublish(eventName string) trace.ITraceDetail {
	detail := &TraceDetailEventConsumer{
		BaseTraceDetail: newTraceDetail(eumCallType.EventPublish, ""),
		Name:            eventName,
	}
	add(detail)
	return detail
}

// TraceMqSend mq发送埋点
func (*traceManager) TraceMqSend(method string, server string, exchange string, routingKey string) trace.ITraceDetail {
	detail := &TraceDetailMq{
		BaseTraceDetail: newTraceDetail(eumCallType.Mq, method),
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
		BaseTraceDetail: newTraceDetail(eumCallType.Mq, method),
		Server:          server,
		Exchange:        exchange,
	}
	add(detail)
	return detail
}

// TraceRedis Redis埋点
func (*traceManager) TraceRedis(method string, key string, field string) trace.ITraceDetail {
	detail := &TraceDetailRedis{
		BaseTraceDetail: newTraceDetail(eumCallType.Redis, method),
		Key:             key,
		Field:           field,
	}
	add(detail)
	return detail
}

// TraceHttp http埋点
func (*traceManager) TraceHttp(method string, url string) trace.ITraceDetail {
	detail := &TraceDetailHttp{
		BaseTraceDetail: newTraceDetail(eumCallType.Http, method),
		Method:          method,
		Url:             url,
	}
	add(detail)
	return detail
}

// TraceGrpc grpc埋点
func (*traceManager) TraceGrpc(method string, url string) trace.ITraceDetail {
	detail := &TraceDetailGrpc{
		BaseTraceDetail: newTraceDetail(eumCallType.Grpc, method),
		Method:          method,
		Url:             url,
	}
	add(detail)
	return detail
}

func newTraceDetail(callType eumCallType.Enum, methodName string) trace.BaseTraceDetail {
	// 获取当前层级列表
	lstScope := trace.ScopeLevel.Get()
	var parentDetailId string
	if len(lstScope) > 0 {
		parentDetailId = lstScope[len(lstScope)-1].DetailId
	}
	baseTraceDetail := trace.BaseTraceDetail{
		DetailId:       strconv.FormatInt(sonyflake.GenerateId(), 10),
		Level:          len(lstScope) + 1,
		ParentDetailId: parentDetailId,
		MethodName:     methodName,
		CallType:       callType,
		StartTs:        time.Now().UnixMicro(),
		EndTs:          time.Now().UnixMicro(),
		Comment:        trace.GetComment(),
	}
	// 加入到当前层级列表
	trace.ScopeLevel.Set(append(lstScope, baseTraceDetail))
	return baseTraceDetail
}

func add(traceDetail trace.ITraceDetail) {
	if t := trace.CurTraceContext.Get(); t != nil {
		detail := traceDetail.GetTraceDetail()
		// 时间轴：上下文入口起点时间到本次开始时间
		detail.Timeline = time.Duration(detail.StartTs-t.GetStartTs()) * time.Microsecond
		details := t.GetList()
		if len(details) > 0 {
			detail.UnTraceTs = time.Duration(detail.StartTs-details[len(details)-1].(trace.ITraceDetail).GetTraceDetail().EndTs) * time.Microsecond
		} else {
			detail.UnTraceTs = time.Duration(detail.StartTs-t.GetStartTs()) * time.Microsecond
		}
		detail.TraceId, detail.AppName, detail.AppId, detail.AppIp, detail.ParentAppName = t.GetAppInfo()
		t.AddDetail(traceDetail)
	}
}
