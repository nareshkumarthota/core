package zapcores

import (
	"go.uber.org/zap/zapcore"
)

// zapCoreMap holds core impl
var zapCoreMap map[string]zapcore.Core

func init() {
	zapCoreMap = make(map[string]zapcore.Core)
}

// RegisterCore adds core to zapcoremap
func RegisterCore(name string, core zapcore.Core) {
	zapCoreMap[name] = core
}

// GetZapCoreMap returns complete map
func GetZapCoreMap() map[string]zapcore.Core {
	return zapCoreMap
}
