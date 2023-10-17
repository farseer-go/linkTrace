package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/linkTrace/eumCallType"
)

// TraceHandDetail 手动埋点
type TraceHandDetail struct {
	TraceDetail
	Name string
}

func (receiver *TraceHandDetail) GetTraceDetail() *TraceDetail {
	return &receiver.TraceDetail
}

func (receiver *TraceHandDetail) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s， %s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.Name)
}

// TraceHand 手动埋点
func TraceHand(name string) *TraceHandDetail {
	detail := &TraceHandDetail{
		TraceDetail: newTraceDetail(eumCallType.Hand),
		Name:        name,
	}

	add(detail)
	return detail
}

// TraceKeyLocation 关键位置埋点
func TraceKeyLocation(name string) *TraceHandDetail {
	detail := &TraceHandDetail{
		TraceDetail: newTraceDetail(eumCallType.KeyLocation),
		Name:        name,
	}

	add(detail)
	return detail
}
