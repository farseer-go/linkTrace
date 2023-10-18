package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/linkTrace/eumCallType"
)

type TraceDetailMq struct {
	TraceDetail
	Server     string // redis key
	Exchange   string // hash field
	RoutingKey string // hash field
}

func (receiver *TraceDetailMq) GetTraceDetail() *TraceDetail {
	return &receiver.TraceDetail
}

func (receiver *TraceDetailMq) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s，%s Server=%s，Exchange=%s，RoutingKey=%s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.CallMethod, receiver.Server, receiver.Exchange, receiver.RoutingKey)
}

// TraceMq mq send埋点
func TraceMq(method string, server string, exchange string, routingKey string) *TraceDetailMq {
	detail := &TraceDetailMq{
		TraceDetail: newTraceDetail(eumCallType.Mq, method),
		Server:      server,
		Exchange:    exchange,
		RoutingKey:  routingKey,
	}
	add(detail)
	return detail
}
