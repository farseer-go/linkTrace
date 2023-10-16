package linkTrace

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/snowflake"
	"github.com/farseer-go/linkTrace/eumCallType"
	"github.com/farseer-go/linkTrace/eumLinkType"
	"github.com/farseer-go/queue"
	"strings"
	"time"
)

type TraceContext struct {
	TraceId         int64                                  `es:"primaryKey"` // 上下文ID
	AppId           int64                                  // 应用ID
	AppName         string                                 // 应用名称
	AppIp           string                                 // 应用IP
	ParentAppName   string                                 // 上游应用
	StartTs         int64                                  // 调用开始时间戳
	EndTs           int64                                  // 调用结束时间戳
	UseTs           time.Duration                          // 总共使用时间毫秒
	LinkType        eumLinkType.Enum                       // 状态码
	Domain          string                                 // 请求域名
	Path            string                                 `es_type:"text"` // 请求地址
	Method          string                                 // 请求方式
	ContentType     string                                 // 请求内容类型
	StatusCode      int                                    // 状态码
	Headers         collections.Dictionary[string, string] `es_type:"flattened"` // 请求头部
	RequestBody     string                                 `es_type:"text"`      // 请求参数
	ResponseBody    string                                 `es_type:"text"`      // 输出参数
	RequestIp       string                                 // 客户端IP
	List            collections.List[*TraceDetail]         `es_type:"object"` // 调用的上下文
	ExceptionDetail ExceptionDetail                        `es_type:"object"` // 是否执行异常
}

// End 结束当前链路
func (receiver *TraceContext) End() {
	receiver.EndTs = time.Now().UnixMicro()
	receiver.UseTs = time.Duration(receiver.EndTs-receiver.StartTs) * time.Microsecond

	// 启用了链路追踪后，把数据写入到本地队列中
	if defConfig.Enable {
		queue.Push("TraceContext", *receiver)
	}

	// 打印日志
	receiver.printLog()
}

// printLog 打印日志
func (receiver *TraceContext) printLog() {
	// 打印日志
	if defConfig.PrintLog {
		lst := collections.NewList[string]()
		for i := 0; i < receiver.List.Count(); i++ {
			item := *receiver.List.Index(i)
			switch item.CallType {
			case eumCallType.Database:
				tableName := parse.ToString(item.Data["TableName"])
				sql := flog.ReplaceBlues(parse.ToString(item.Data["Sql"]), "SELECT ", "UPDATE ", "DELETE ", " FROM ", " WHERE ", " LIMIT ", " SET ", " ORDER BY ", " and ", " or ")
				sql = strings.ReplaceAll(sql, tableName, flog.Red(tableName))
				lst.Add(fmt.Sprintf("%s：[%s] %s", flog.Blue(i+1), flog.Yellow(item.CallType.ToString()), sql))
			case eumCallType.Redis:
				lst.Add(fmt.Sprintf("%s：[%s]耗时：%s，%s Key=%s，FieldKey=%s", flog.Blue(i+1), flog.Yellow(item.CallType.ToString()), flog.Red(item.UseTs.String()), item.CallMethod, item.Data["Key"], item.Data["FieldKey"]))
			}
		}
		flog.Printf("【链路追踪】TraceId:%s，耗时：%s，%s，详情如下：\n%s\n", flog.Green(parse.ToString(receiver.TraceId)), flog.Red(receiver.UseTs.String()), receiver.Path, strings.Join(lst.ToArray(), "\n"))
		fmt.Println("-----------------------------------------------------------------")
	}
}

func TraceWebApi(domain string, path string, method string, contentType string, headerDictionary collections.ReadonlyDictionary[string, string], requestBody string, requestIp string) *TraceContext {
	traceId := parse.ToInt64(headerDictionary.GetValue("TraceId"))
	if traceId == 0 {
		traceId = snowflake.GenerateId()
	}
	return &TraceContext{
		AppId:         fs.AppId,
		AppName:       fs.AppName,
		AppIp:         fs.AppIp,
		ParentAppName: headerDictionary.GetValue("AppName"),
		TraceId:       traceId,
		StartTs:       time.Now().UnixMicro(),
		LinkType:      eumLinkType.WebApi,
		Domain:        domain,
		Path:          path,
		Method:        method,
		ContentType:   contentType,
		Headers:       headerDictionary.ToDictionary(),
		RequestBody:   requestBody,
		RequestIp:     requestIp,
		List:          collections.NewList[*TraceDetail](),
		//ExceptionDetail: ExceptionDetail{},
	}
}

func TraceDatabase(dbName, tableName, sql string) *TraceDetail {
	detail := &TraceDetail{
		CallType: eumCallType.Database,
		Data: map[string]any{
			"DbName":    dbName,
			"TableName": tableName,
			"Sql":       sql,
		},
		CallStackTrace: CallStackTrace{},
		CallMethod:     "",
		StartTs:        time.Now().UnixMicro(),
	}

	if trace := GetCurTrace(); trace != nil && defConfig.Enable {
		trace.List.Add(detail)
	}
	return detail
}

func TraceRedis(redisMethod string, key string, fieldKey string) *TraceDetail {
	detail := &TraceDetail{
		CallType: eumCallType.Redis,
		Data: map[string]any{
			"Key":      key,
			"FieldKey": fieldKey,
		},
		CallMethod:     redisMethod,
		StartTs:        time.Now().UnixMicro(),
		CallStackTrace: CallStackTrace{},
	}

	if trace := GetCurTrace(); trace != nil && defConfig.Enable {
		trace.List.Add(detail)
	}
	return detail
}
