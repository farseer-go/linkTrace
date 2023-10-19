package linkTrace

import (
	"github.com/farseer-go/fs/asyncLocal"
	"github.com/farseer-go/fs/trace"
)

// CurTraceContext 当前请求的Trace上下文
var curTraceContext = asyncLocal.New[trace.ITraceContext]()
