package linkTrace

import (
	"fmt"

	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
)

// TraceDetailHand 手动埋点
type TraceDetailHand struct {
	trace.BaseTraceDetail
	Name string
}

func (receiver *TraceDetailHand) GetTraceDetail() *trace.BaseTraceDetail {
	return &receiver.BaseTraceDetail
}

func (receiver *TraceDetailHand) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s， %s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.Name)
}

func (receiver *TraceDetailHand) SetName(name string) {
	receiver.Name = name
}
