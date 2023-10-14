package linkTrace

import "github.com/farseer-go/collections"

type TrackContext struct {
	ParentAppName   string                                 // 上游应用
	TraceId         string                                 // 上下文ID
	StartTs         int64                                  // 调用开始时间戳
	EndTs           int64                                  // 调用结束时间戳
	UseTs           int64                                  // 总共使用时间毫秒
	LinkType        EumLinkType                            // 状态码
	Domain          string                                 // 请求域名
	Path            string                                 // 请求地址
	Method          string                                 // 请求方式
	ContentType     string                                 // 请求内容类型
	StatusCode      string                                 // 状态码
	Headers         collections.Dictionary[string, string] // 请求头部
	RequestBody     string                                 // 请求参数
	ResponseBody    string                                 // 输出参数
	RequestIp       string                                 // 客户端IP
	List            collections.List[LinkTrackDetail]      // 调用的上下文
	ExceptionDetail ExceptionDetail                        // 是否执行异常
}
