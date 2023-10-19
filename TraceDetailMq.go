package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
)

type TraceDetailMq struct {
	trace.BaseTraceDetail
	Server     string // redis key
	Exchange   string // hash field
	RoutingKey string // hash field
}

func (receiver *TraceDetailMq) GetTraceDetail() *trace.BaseTraceDetail {
	return &receiver.BaseTraceDetail
}

func (receiver *TraceDetailMq) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s，%s Server=%s，Exchange=%s，RoutingKey=%s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.CallMethod, receiver.Server, receiver.Exchange, receiver.RoutingKey)
}
