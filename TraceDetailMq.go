package linkTrace

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
)

type TraceDetailMq struct {
	trace.BaseTraceDetail
	Server     string // MQ服务器地址
	Exchange   string // 交换器名称
	RoutingKey string // 路由key
}

func (receiver *TraceDetailMq) GetTraceDetail() *trace.BaseTraceDetail {
	return &receiver.BaseTraceDetail
}

func (receiver *TraceDetailMq) ToString() string {
	return fmt.Sprintf("[%s]耗时：%s，%s Server=%s，Exchange=%s，RoutingKey=%s", flog.Yellow(receiver.CallType.ToString()), flog.Red(receiver.UseTs.String()), receiver.MethodName, receiver.Server, receiver.Exchange, receiver.RoutingKey)
}

func (receiver *TraceDetailMq) Desc() (caption string, desc string) {
	caption = fmt.Sprintf("发送MQ消息 => %s %s %s", receiver.Server, receiver.Exchange, receiver.RoutingKey)
	desc = fmt.Sprintf("%s %s %s", receiver.Server, receiver.Exchange, receiver.RoutingKey)
	return
}
