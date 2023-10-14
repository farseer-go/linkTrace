package linkTrace

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/snowflake"
	"github.com/farseer-go/linkTrace/eumLinkType"
	"time"
)

type TrackContext struct {
	ParentAppName   string                                 // 上游应用
	TraceId         int64                                  // 上下文ID
	StartTs         int64                                  // 调用开始时间戳
	EndTs           int64                                  // 调用结束时间戳
	UseTs           int64                                  // 总共使用时间毫秒
	LinkType        eumLinkType.Enum                       // 状态码
	Domain          string                                 // 请求域名
	Path            string                                 // 请求地址
	Method          string                                 // 请求方式
	ContentType     string                                 // 请求内容类型
	StatusCode      int                                    // 状态码
	Headers         collections.Dictionary[string, string] // 请求头部
	RequestBody     string                                 // 请求参数
	ResponseBody    string                                 // 输出参数
	RequestIp       string                                 // 客户端IP
	List            collections.List[LinkTraceDetail]      // 调用的上下文
	ExceptionDetail ExceptionDetail                        // 是否执行异常
}

func NewWebApi(domain string, path string, method string, contentType string, headerDictionary collections.ReadonlyDictionary[string, string], requestBody string, requestIp string) *TrackContext {
	traceId := parse.ToInt64(headerDictionary.GetValue("TraceId"))
	if traceId == 0 {
		traceId = snowflake.GenerateId()
	}
	return &TrackContext{
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
		//List:            collections.List[LinkTraceDetail]{},
		//ExceptionDetail: ExceptionDetail{},
	}
}

// End 结束当前链路
func (receiver *TrackContext) End() {
	receiver.EndTs = time.Now().UnixMicro()
	receiver.UseTs = receiver.EndTs - receiver.StartTs
}

type LinkTraceDetail struct {
}

type ExceptionDetail struct {
}
