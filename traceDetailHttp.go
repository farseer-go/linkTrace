package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/linkTrace/eumCallType"
)

type TraceDetailHttp struct {
	TraceDetail
	Method string // post/get/put/delete
	Url    string // url
}

func (receiver *TraceDetailHttp) GetTraceDetail() *TraceDetail {
	return &receiver.TraceDetail
}

func (receiver *TraceDetailHttp) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s，%s [%s]%s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.CallMethod, receiver.Method, receiver.Url)
}

// TraceHttp http埋点
func TraceHttp(method string, url string) *TraceDetailHttp {
	detail := &TraceDetailHttp{
		TraceDetail: newTraceDetail(eumCallType.Http, method),
		Method:      method,
		Url:         url,
	}
	add(detail)
	return detail
}
