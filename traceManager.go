package linkTrace

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/snowflake"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/linkTrace/eumCallType"
	"github.com/farseer-go/linkTrace/eumTraceType"
	"time"
)

type traceManager struct {
}

func (*traceManager) GetCurTrace() trace.ITraceContext {
	return curTraceContext.Get()
}

// EntryWebApi Webapi入口
func (*traceManager) EntryWebApi(domain string, path string, method string, contentType string, headerDictionary collections.ReadonlyDictionary[string, string], requestBody string, requestIp string) trace.ITraceContext {
	traceId := parse.ToInt64(headerDictionary.GetValue("TraceId"))
	if traceId == 0 {
		traceId = snowflake.GenerateId()
	}
	context := &TraceContext{
		AppId:         fs.AppId,
		AppName:       fs.AppName,
		AppIp:         fs.AppIp,
		ParentAppName: headerDictionary.GetValue("AppName"),
		TraceId:       traceId,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.WebApi,
		Web: WebContext{
			Domain:      domain,
			Path:        path,
			Method:      method,
			ContentType: contentType,
			Headers:     headerDictionary.ToDictionary(),
			RequestBody: requestBody,
			RequestIp:   requestIp,
		},
		List: collections.NewList[trace.ITraceDetail](),
	}
	curTraceContext.Set(context)
	trace.ScopeLevel.Set(collections.NewList[trace.BaseTraceDetail]())
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

// TraceKeyLocation 关键位置埋点
func (*traceManager) TraceKeyLocation(name string) trace.ITraceDetail {
	detail := &TraceDetailHand{
		BaseTraceDetail: newTraceDetail(eumCallType.KeyLocation, ""),
		Name:            name,
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

// EntryMqConsumer mq 消费埋点
func (*traceManager) EntryMqConsumer(server string, queueName string, routingKey string) trace.ITraceContext {
	traceId := snowflake.GenerateId()
	context := &TraceContext{
		AppId:         fs.AppId,
		AppName:       fs.AppName,
		AppIp:         fs.AppIp,
		ParentAppName: "",
		TraceId:       traceId,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.MqConsumer,
		Consumer: ConsumerContext{
			Server:     server,
			QueueName:  queueName,
			RoutingKey: routingKey,
		},
		List: collections.NewList[trace.ITraceDetail](),
	}
	curTraceContext.Set(context)
	trace.ScopeLevel.Set(collections.NewList[trace.BaseTraceDetail]())
	return context
}

// EntryQueueConsumer queue 消费埋点
func (*traceManager) EntryQueueConsumer(subscribeName string) trace.ITraceContext {
	traceId := snowflake.GenerateId()
	context := &TraceContext{
		AppId:         fs.AppId,
		AppName:       fs.AppName,
		AppIp:         fs.AppIp,
		ParentAppName: "",
		TraceId:       traceId,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.QueueConsumer,
		Consumer: ConsumerContext{
			Server:    fmt.Sprintf("%s/%s/%v", fs.AppName, fs.AppIp, fs.AppId),
			QueueName: subscribeName,
		},
		List: collections.NewList[trace.ITraceDetail](),
	}
	curTraceContext.Set(context)
	trace.ScopeLevel.Set(collections.NewList[trace.BaseTraceDetail]())
	return context
}

// EntryTask 创建本地任务入口
func (*traceManager) EntryTask(taskName string) trace.ITraceContext {
	traceId := snowflake.GenerateId()
	context := &TraceContext{
		AppId:         fs.AppId,
		AppName:       fs.AppName,
		AppIp:         fs.AppIp,
		ParentAppName: "",
		TraceId:       traceId,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.Task,
		Task: TaskContext{
			TaskName: taskName,
		},
		List: collections.NewList[trace.ITraceDetail](),
	}
	curTraceContext.Set(context)
	trace.ScopeLevel.Set(collections.NewList[trace.BaseTraceDetail]())
	return context
}

// EntryFSchedule 创建调度中心入口
func (*traceManager) EntryFSchedule(taskGroupName string, taskGroupId int64, taskId int64) trace.ITraceContext {
	traceId := snowflake.GenerateId()
	context := &TraceContext{
		AppId:         fs.AppId,
		AppName:       fs.AppName,
		AppIp:         fs.AppIp,
		ParentAppName: "",
		TraceId:       traceId,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.FSchedule,
		Task: TaskContext{
			TaskName:    taskGroupName,
			TaskGroupId: taskGroupId,
			TaskId:      taskId,
		},
		List: collections.NewList[trace.ITraceDetail](),
	}
	curTraceContext.Set(context)
	trace.ScopeLevel.Set(collections.NewList[trace.BaseTraceDetail]())
	return context
}

// EntryWatchKey 创建etcd入口
func (*traceManager) EntryWatchKey(key string) trace.ITraceContext {
	traceId := snowflake.GenerateId()
	context := &TraceContext{
		AppId:         fs.AppId,
		AppName:       fs.AppName,
		AppIp:         fs.AppIp,
		ParentAppName: "",
		TraceId:       traceId,
		StartTs:       time.Now().UnixMicro(),
		TraceType:     eumTraceType.WatchKey,
		WatchKey: WatchKeyContext{
			Key: key,
		},
		List: collections.NewList[trace.ITraceDetail](),
	}
	curTraceContext.Set(context)
	trace.ScopeLevel.Set(collections.NewList[trace.BaseTraceDetail]())
	return context
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

func newTraceDetail(callType eumCallType.Enum, callMethod string) trace.BaseTraceDetail {
	// 获取当前层级列表
	lstScope := trace.ScopeLevel.Get()
	baseTraceDetail := trace.BaseTraceDetail{
		DetailId:       snowflake.GenerateId(),
		Level:          lstScope.Count() + 1,
		ParentDetailId: lstScope.Last().DetailId,
		CallMethod:     callMethod,
		CallType:       callType,
		StartTs:        time.Now().UnixMicro(),
		EndTs:          time.Now().UnixMicro(),
	}
	// 加入到当前层级列表
	if !lstScope.IsNil() {
		lstScope.Add(baseTraceDetail)
		trace.ScopeLevel.Set(lstScope)
	}
	return baseTraceDetail
}

func add(traceDetail trace.ITraceDetail) {
	if t := curTraceContext.Get(); t != nil && defConfig.Enable {
		detail := traceDetail.GetTraceDetail()
		// 时间轴：上下文入口起点时间到本次开始时间
		detail.Timeline = time.Duration(detail.StartTs-t.GetStartTs()) * time.Microsecond
		if t.GetList().Count() > 0 {
			detail.UnTraceTs = time.Duration(detail.StartTs-t.GetList().Last().GetTraceDetail().EndTs) * time.Microsecond
		} else {
			detail.UnTraceTs = time.Duration(detail.StartTs-t.GetStartTs()) * time.Microsecond
		}
		t.AddDetail(traceDetail)
	}
}
