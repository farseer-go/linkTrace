package linkTrace

import (
	"github.com/timandy/routine"
)

// TraceId 当前请求的TraceId
var TraceId routine.ThreadLocal[int64]
