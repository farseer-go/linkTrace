package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/linkTrace/eumCallType"
)

type TraceDetailEtcd struct {
	TraceDetail
	Key     string // key
	LeaseID int64
}

func (receiver *TraceDetailEtcd) GetTraceDetail() *TraceDetail {
	return &receiver.TraceDetail
}

func (receiver *TraceDetailEtcd) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s，%s Key=%s LeaseID=%v", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.CallMethod, receiver.Key, receiver.LeaseID)
}

// TraceEtcd etcd埋点
func TraceEtcd(method string, key string, leaseID int64) *TraceDetailEtcd {
	detail := &TraceDetailEtcd{
		TraceDetail: newTraceDetail(eumCallType.Etcd, method),
		Key:         key,
		LeaseID:     leaseID,
	}
	add(detail)
	return detail
}
