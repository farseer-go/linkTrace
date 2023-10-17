package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/linkTrace/eumCallType"
	"time"
)

// TraceHandDetail 手动埋点
type TraceHandDetail struct {
	TraceDetail
	Name string
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

// TraceKeyLocation 关键位置埋点
func TraceKeyLocation(name string) *TraceHandDetail {
	detail := &TraceHandDetail{
		TraceDetail: newTraceDetail(eumCallType.KeyLocation),
		Name:        name,
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
