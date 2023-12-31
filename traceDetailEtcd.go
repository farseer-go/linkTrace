package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
)

type TraceDetailEtcd struct {
	trace.BaseTraceDetail
	Key     string // key
	LeaseID int64
}

func (receiver *TraceDetailEtcd) GetTraceDetail() *trace.BaseTraceDetail {
	return &receiver.BaseTraceDetail
}

func (receiver *TraceDetailEtcd) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s，%s Key=%s LeaseID=%v", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.MethodName, receiver.Key, receiver.LeaseID)
}

func (receiver *TraceDetailEtcd) Desc() (caption string, desc string) {
	caption = fmt.Sprintf("执行Etcd => %s %v", receiver.Key, receiver.LeaseID)
	desc = fmt.Sprintf("%s %v", receiver.Key, receiver.LeaseID)
	return
}
