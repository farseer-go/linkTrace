package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/linkTrace/eumCallType"
)

// TraceDetailHand 手动埋点
type TraceDetailHand struct {
	TraceDetail
	Name string
}

func (receiver *TraceDetailHand) GetTraceDetail() *TraceDetail {
	return &receiver.TraceDetail
}

func (receiver *TraceDetailHand) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s， %s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.Name)
}

// TraceHand 手动埋点
func TraceHand(name string) *TraceDetailHand {
	detail := &TraceDetailHand{
		TraceDetail: newTraceDetail(eumCallType.Hand),
		Name:        name,
	}
	add(detail)
	return detail
}

// TraceKeyLocation 关键位置埋点
func TraceKeyLocation(name string) *TraceDetailHand {
	detail := &TraceDetailHand{
		TraceDetail: newTraceDetail(eumCallType.KeyLocation),
		Name:        name,
	}
	add(detail)
	return detail
}
