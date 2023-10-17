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

func (receiver *TraceRedisDetail) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s，%s Key=%s，FieldKey=%s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.CallMethod, receiver.Key, receiver.Field)
}

func TraceRedis(redisMethod string, key string, field string) *TraceRedisDetail {
	detail := &TraceRedisDetail{
		TraceDetail: TraceDetail{
			//CallStackTrace: CallStackTrace{},
			CallMethod: redisMethod,
			CallType:   eumCallType.Redis,
			StartTs:    time.Now().UnixMicro(),
		},
		Key:   key,
		Field: field,
	}

	if trace := GetCurTrace(); trace != nil && defConfig.Enable {
		// 时间轴：上下文入口起点时间到本次开始时间
		detail.Timeline = time.Duration(detail.StartTs-trace.StartTs) * time.Microsecond
		if trace.List.Count() > 0 {
			detail.UnTraceTs = time.Duration(detail.StartTs-trace.List.Last().GetEndTs()) * time.Microsecond
		} else {
			detail.UnTraceTs = time.Duration(detail.StartTs-trace.StartTs) * time.Microsecond
		}
		trace.List.Add(detail)
	}
	return detail
}
