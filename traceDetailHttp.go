package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/linkTrace/eumCallType"
)

type TraceDetailHttp struct {
	TraceDetail
	Key   string // redis key
	Field string // hash field
}

func (receiver *TraceDetailHttp) GetTraceDetail() *TraceDetail {
	return &receiver.TraceDetail
}

func (receiver *TraceDetailHttp) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s，%s Key=%s，FieldKey=%s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.CallMethod, receiver.Key, receiver.Field)
}

// TraceHttp http埋点
func TraceHttp(method string, key string, field string) *TraceDetailHttp {
	detail := &TraceDetailHttp{
		TraceDetail: newTraceDetail(eumCallType.Http, method),
		Key:         key,
		Field:       field,
	}
	add(detail)
	return detail
}
