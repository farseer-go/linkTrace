package linkTrace

type CallStackTrace struct {
	CallMethod     string            // 调用方法
	FileName       string            // 执行文件名称
	FileLineNumber int               // 方法执行行数
	ReturnType     string            // 方法返回类型
	MethodParams   map[string]string // 方法入参
}

type ExceptionDetail struct {
}
