package linkTrace

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/snowflake"
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
	List            collections.List[ITraceDetail]         `es_type:"object"` // 调用的上下文
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
			lst.Add(receiver.List.Index(i).ToString(i + 1))
		}
		flog.Printf("【链路追踪】TraceId:%s，耗时：%s，%s：\n%s\n", flog.Green(parse.ToString(receiver.TraceId)), flog.Red(receiver.UseTs.String()), receiver.Path, strings.Join(lst.ToArray(), "\n"))
		fmt.Println("-----------------------------------------------------------------")
	}
}

// TraceWebApi Webapi入口
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
		List:          collections.NewList[ITraceDetail](),
		//ExceptionDetail: ExceptionDetail{},
	}
}
