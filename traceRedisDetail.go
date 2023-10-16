package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/linkTrace/eumCallType"
	"time"
)

type TraceRedisDetail struct {
	TraceDetail
	Key   string // redis key
	Field string // hash field
}

func (receiver *TraceRedisDetail) ToString(index int) string {
	return fmt.Sprintf("%s：[%s]耗时：%s，%s Key=%s，FieldKey=%s", flog.Blue(index), flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.CallMethod, receiver.Key, receiver.Field)
}

func TraceRedis(redisMethod string, key string, field string) *TraceRedisDetail {
	detail := &TraceRedisDetail{
		TraceDetail: TraceDetail{
			CallStackTrace: CallStackTrace{},
			CallMethod:     "",
			CallType:       eumCallType.Redis,
			StartTs:        time.Now().UnixMicro(),
		},
		Key:   key,
		Field: field,
	}

	if trace := GetCurTrace(); trace != nil && defConfig.Enable {
		trace.List.Add(detail)
	}
	return detail
}
