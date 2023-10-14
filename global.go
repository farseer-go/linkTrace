package linkTrace

import (
	"github.com/timandy/routine"
)

// CurTraceContext 当前请求的Trace上下文
var curTraceContext = routine.NewInheritableThreadLocal[*TraceContext]()

// GetCurTrace 获取当前TrackContext
func GetCurTrace() *TraceContext {
	return curTraceContext.Get()
}

// SetCurTrace 设置当前TrackContext
func SetCurTrace(context *TraceContext) {
	curTraceContext.Set(context)
}
