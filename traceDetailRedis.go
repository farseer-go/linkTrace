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
	return fmt.Sprintf("[%s]耗时：%s，%s Key=%s，Field=%s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.MethodName, receiver.Key, receiver.Field)
}

func (receiver *TraceDetailRedis) Desc() (caption string, desc string) {
	caption = fmt.Sprintf("执行Redis => %s %s %s", receiver.MethodName, receiver.Key, receiver.Field)
	desc = fmt.Sprintf("%s %s", receiver.Key, receiver.Field)
	return
}
