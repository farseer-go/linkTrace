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

func TraceHand(name string) *TraceHandDetail {
	detail := &TraceHandDetail{
		TraceDetail: TraceDetail{
			//CallStackTrace: CallStackTrace{},
			CallMethod: "",
			CallType:   eumCallType.Hand,
			StartTs:    time.Now().UnixMicro(),
			EndTs:      time.Now().UnixMicro(),
		},
		Name: name,
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
