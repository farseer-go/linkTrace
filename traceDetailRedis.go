package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
)

type TraceDetailRedis struct {
	trace.BaseTraceDetail
	Key   string // redis key
	Field string // hash field
}

func (receiver *TraceDetailRedis) GetTraceDetail() *trace.BaseTraceDetail {
	return &receiver.BaseTraceDetail
}

func (receiver *TraceDetailRedis) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s，%s Key=%s，Field=%s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.CallMethod, receiver.Key, receiver.Field)
}
