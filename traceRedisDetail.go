package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/linkTrace/eumCallType"
)

type TraceRedisDetail struct {
	TraceDetail
	Key   string // redis key
	Field string // hash field
}

func (receiver *TraceRedisDetail) GetTraceDetail() *TraceDetail {
	return &receiver.TraceDetail
}

func (receiver *TraceRedisDetail) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s，%s Key=%s，FieldKey=%s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.CallMethod, receiver.Key, receiver.Field)
}

// TraceRedis Redis埋点
func TraceRedis(redisMethod string, key string, field string) *TraceRedisDetail {
	detail := &TraceRedisDetail{
		TraceDetail: newTraceDetail(eumCallType.Redis),
		Key:         key,
		Field:       field,
	}

	add(detail)
	return detail
}
