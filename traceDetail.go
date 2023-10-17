package linkTrace

import (
	"github.com/farseer-go/linkTrace/eumCallType"
	"time"
)

type ITraceDetail interface {
	ToString() string
	//GetEndTs() int64
	//GetUnTraceTs() time.Duration
	//GetTimeline() time.Duration
	GetTraceDetail() *TraceDetail
}

// TraceDetail 埋点明细
type TraceDetail struct {
	//CallStackTrace   CallStackTrace   // 调用栈
	CallMethod       string           // 调用方法
	CallType         eumCallType.Enum // 调用类型
	Timeline         time.Duration    // 从入口开始统计
	UnTraceTs        time.Duration    // 上一次结束到现在开始之间未Trace的时间
	StartTs          int64            // 调用开始时间戳
	EndTs            int64            // 调用停止时间戳
	UseTs            time.Duration    // 总共使用时间毫秒
	IsException      bool             // 是否执行异常
	ExceptionMessage string           // 异常信息
}

type CallStackTrace struct {
	CallMethod     string            // 调用方法
	FileName       string            // 执行文件名称
	FileLineNumber int               // 方法执行行数
	ReturnType     string            // 方法返回类型
	MethodParams   map[string]string // 方法入参
}

type ExceptionDetail struct {
}

func newTraceDetail(callType eumCallType.Enum) TraceDetail {
	return TraceDetail{
		//CallStackTrace: CallStackTrace{},
		CallMethod: "",
		CallType:   callType,
		StartTs:    time.Now().UnixMicro(),
		EndTs:      time.Now().UnixMicro(),
	}
}

// End 链路明细执行完后，统计用时
func (receiver *TraceDetail) End(err error) {
	receiver.EndTs = time.Now().UnixMicro()
	receiver.UseTs = time.Duration(receiver.EndTs-receiver.StartTs) * time.Microsecond

	if err != nil {
		receiver.IsException = true
		receiver.ExceptionMessage = err.Error()
	}
}

func add(traceDetail ITraceDetail) {
	if trace := GetCurTrace(); trace != nil && defConfig.Enable {
		detail := traceDetail.GetTraceDetail()
		// 时间轴：上下文入口起点时间到本次开始时间
		detail.Timeline = time.Duration(detail.StartTs-trace.StartTs) * time.Microsecond
		if trace.List.Count() > 0 {
			detail.UnTraceTs = time.Duration(detail.StartTs-trace.List.Last().GetTraceDetail().EndTs) * time.Microsecond
		} else {
			detail.UnTraceTs = time.Duration(detail.StartTs-trace.StartTs) * time.Microsecond
		}
		trace.List.Add(traceDetail)
	}
}
