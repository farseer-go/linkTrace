package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/linkTrace/eumCallType"
)

type TraceDetailEtcd struct {
	TraceDetail
	Key   string // redis key
	Field string // hash field
}

func (receiver *TraceDetailEtcd) GetTraceDetail() *TraceDetail {
	return &receiver.TraceDetail
}

func (receiver *TraceDetailEtcd) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s，%s Key=%s，FieldKey=%s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.CallMethod, receiver.Key, receiver.Field)
}

// TraceEtcd etcd埋点
func TraceEtcd(redisMethod string, key string, field string) *TraceDetailEtcd {
	detail := &TraceDetailEtcd{
		TraceDetail: newTraceDetail(eumCallType.Mq),
		Key:         key,
		Field:       field,
	}
	add(detail)
	return detail
}
