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

func (receiver *TraceDetailHand) Desc() (caption string, desc string) {
	caption = fmt.Sprintf("手动埋点 => %s", receiver.Name)
	desc = fmt.Sprintf("%s", receiver.Name)
	return
}
