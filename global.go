package linkTrace

import (
	"github.com/farseer-go/fs/trace"
	"github.com/timandy/routine"
)

// CurTraceContext 当前请求的Trace上下文
var curTraceContext = routine.NewInheritableThreadLocal[trace.ITraceContext]()

// getCurTrace 获取当前TrackContext
func getCurTrace() trace.ITraceContext {
	return curTraceContext.Get()
}

// setCurTrace 设置当前TrackContext
func setCurTrace(context trace.ITraceContext) {
	curTraceContext.Set(context)
}
