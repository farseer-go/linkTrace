package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
)

// TraceDetailEventConsumer 事件埋点
type TraceDetailEventConsumer struct {
	trace.BaseTraceDetail
	Name string
}

func (receiver *TraceDetailEventConsumer) GetTraceDetail() *trace.BaseTraceDetail {
	return &receiver.BaseTraceDetail
}

func (receiver *TraceDetailEventConsumer) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s， %s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.Name)
}
