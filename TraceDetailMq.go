package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/linkTrace/eumCallType"
)

type TraceDetailMq struct {
	TraceDetail
	Key   string // redis key
	Field string // hash field
}

func (receiver *TraceDetailMq) GetTraceDetail() *TraceDetail {
	return &receiver.TraceDetail
}

func (receiver *TraceDetailMq) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s，%s Key=%s，FieldKey=%s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.CallMethod, receiver.Key, receiver.Field)
}

// TraceMq mq send埋点
func TraceMq(redisMethod string, key string, field string) *TraceDetailMq {
	detail := &TraceDetailMq{
		TraceDetail: newTraceDetail(eumCallType.Mq),
		Key:         key,
		Field:       field,
	}
	add(detail)
	return detail
}
