package linkTrack

import (
	"farseer-go/fs/configure"
	"farseer-go/fs/flog"
)

type Module struct {
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
}

func (module Module) PreInitialize() {
	linkTrack := linkTrack{}
	linkTrack.Init()
}
