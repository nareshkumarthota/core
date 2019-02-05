package zapcores

import (
	"go.uber.org/zap/zapcore"
)

// zapCoreMap holds logger core impl
var zapCoreMap map[string]zapcore.Core

// zapTraceCoreMap holds trace logger core impl
var zapTraceCoreMap map[string]zapcore.Core

func init() {
	zapCoreMap = make(map[string]zapcore.Core)
	zapTraceCoreMap = make(map[string]zapcore.Core)
}

// RegisterLogCore adds core to zapcoremap
func RegisterLogCore(name string, core zapcore.Core) {
	zapCoreMap[name] = core
}

// RegisterTraceLogCore adds trace core to zapTraceCoreMap
func RegisterTraceLogCore(name string, core zapcore.Core) {
	zapTraceCoreMap[name] = core
}

// GetZapCoreMap returns complete log core map
func GetZapCoreMap() map[string]zapcore.Core {
	return zapCoreMap
}

// GetZapTraceCoreMap returns complete trace log core map
func GetZapTraceCoreMap() map[string]zapcore.Core {
	return zapTraceCoreMap
}
