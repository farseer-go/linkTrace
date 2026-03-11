package linkTrace

import (
	"os"
	"time"

	"github.com/farseer-go/fs/batchFileWriter"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/queue"
)

// getLogBasePath 获取日志基础路径
// 如果 /var/ 目录存在，返回 /var/log/linkTrace/，否则返回 ./linkTrace/
func getlinkTraceBasePath() string {
	if _, err := os.Stat("/var/"); err == nil {
		return "/var/log/linkTrace/"
	}
	return "./linkTrace/"
}

// Enable 是否启用
var defConfig config

// 链路数据写入器
var writer *batchFileWriter.BatchFileWriter

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{queue.Module{}}
}

func (module Module) PreInitialize() {
	defConfig = configure.ParseConfig[config]("LinkTrace")

	// 使用了链路追踪组件，则要把空组件移除后，重新注册
	container.Remove[trace.IManager]()
	container.Register(func() trace.IManager {
		return &traceManager{}
	})

	// 启用了链路追踪后，才需要初始化ES和消费
	if defConfig.Enable {
		// 初始化链路数据写入器
		writer = batchFileWriter.NewWriter(getlinkTraceBasePath()+core.AppName+"/", "trace", "hour", 10, 0, time.Second*5, true)
	}
}

func (module Module) Shutdown() {
	// 关闭链路数据写入器
	if defConfig.Enable {
		writer.Close()
	}
}
