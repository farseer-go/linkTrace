package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/linkTrace/eumCallType"
)

type TraceDetailRedis struct {
	TraceDetail
	Key   string // redis key
	Field string // hash field
}

func (receiver *TraceDetailRedis) GetTraceDetail() *TraceDetail {
	return &receiver.TraceDetail
}

func (receiver *TraceDetailRedis) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s，%s Key=%s，Field=%s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.CallMethod, receiver.Key, receiver.Field)
}

// TraceRedis Redis埋点
func TraceRedis(method string, key string, field string) *TraceDetailRedis {
	detail := &TraceDetailRedis{
		TraceDetail: newTraceDetail(eumCallType.Redis, method),
		Key:         key,
		Field:       field,
	}
	add(detail)
	return detail
}
