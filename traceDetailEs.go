package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
)

type TraceDetailEs struct {
	trace.BaseTraceDetail
	IndexName   string // 索引名称
	AliasesName string // 别名
}

func (receiver *TraceDetailEs) GetTraceDetail() *trace.BaseTraceDetail {
	return &receiver.BaseTraceDetail
}

func (receiver *TraceDetailEs) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s，%s IndexName=%s，AliasesName=%s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.MethodName, receiver.IndexName, receiver.AliasesName)
}
