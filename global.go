package linkTrace

import (
	"github.com/timandy/routine"
)

// CurTraceContext 当前请求的Trace上下文
var curTraceContext routine.ThreadLocal[*TrackContext]

// GetCurTrace 获取当前TrackContext
func GetCurTrace() *TrackContext {
	return curTraceContext.Get()
}

// SetCurTrace 设置当前TrackContext
func SetCurTrace(context *TrackContext) {
	curTraceContext.Set(context)
}
