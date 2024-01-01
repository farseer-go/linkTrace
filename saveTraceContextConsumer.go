package linkTrace

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/fops"
	"github.com/farseer-go/fs/trace"
	"time"
)

// FopsServer fops地址
var FopsServer string

// SaveTraceContextConsumer 上传到FOPS中心
func SaveTraceContextConsumer(subscribeName string, lstMessage collections.ListAny, remainingCount int) {
	trace.CurTraceContext.Get().Ignore()
	lstTraceContext := collections.NewList[TraceContext]()
	lstMessage.Foreach(func(item *any) {
		// 上下文
		dto := (*item).(TraceContext)
		lstTraceContext.Add(dto)
	})
	if err := fops.UploadTrace(lstTraceContext); err != nil {
		exception.ThrowRefuseException(err.Error())
	}

	// 控制3秒执行一次
	<-time.After(3 * time.Second)
	return
}
