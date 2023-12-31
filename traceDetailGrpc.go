package linkTrace

import (
	"fmt"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/trace"
)

type TraceDetailGrpc struct {
	trace.BaseTraceDetail
	Method       string                                 // post/get/put/delete
	Url          string                                 // url
	Headers      collections.Dictionary[string, string] // 头部
	RequestBody  string                                 // 入参
	ResponseBody string                                 // 出参
	StatusCode   int                                    // 状态码
}

func (receiver *TraceDetailGrpc) GetTraceDetail() *trace.BaseTraceDetail {
	return &receiver.BaseTraceDetail
}

func (receiver *TraceDetailGrpc) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s，%s [%s]%s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.MethodName, receiver.Method, receiver.Url)
}

func (receiver *TraceDetailGrpc) SetHttpRequest(url string, head map[string]any, requestBody string, responseBody string, statusCode int) {
	receiver.Url = url
	receiver.Headers = collections.NewDictionary[string, string]()
	receiver.RequestBody = requestBody
	receiver.ResponseBody = responseBody
	receiver.StatusCode = statusCode
	for k, v := range head {
		receiver.Headers.Add(k, parse.ToString(v))
	}
}

func (receiver *TraceDetailGrpc) Desc() (caption string, desc string) {
	caption = fmt.Sprintf("调用http => %v %s %s", receiver.StatusCode, receiver.Method, receiver.Url)
	lstHeader := collections.NewList[string]()
	for k, v := range receiver.Headers.ToMap() {
		lstHeader.Add(fmt.Sprintf("%s=%v", k, v))
	}
	desc = fmt.Sprintf("头部：%s 入参：%s 出参：%s", lstHeader.ToString(","), receiver.RequestBody, receiver.ResponseBody)
	return
}
