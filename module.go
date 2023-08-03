package linkTrack

import (
	"farseer-go/fs/configure"
	"farseer-go/fs/flog"
	"github.com/farseer-go/fs/modules"
)

type Module struct {
}

var linkTr linkTrack

func (module Module) DependsModule() []modules.FarseerModule {
	return nil
}

func (module Module) Initialize() {
	linkTrackConfigs := configure.GetSubNodes("LinkTrack")
	for _, configString := range linkTrackConfigs {
		config := configure.ParseString[linkTrackConfig](configString.(string))
		if config.TimeInterval == 0 {
			_ = flog.Error("LinkTrack配置缺少TimeInterval")
			continue
		}
	}
	linkTr.Init()
}

func (module Module) PreInitialize() {
}
