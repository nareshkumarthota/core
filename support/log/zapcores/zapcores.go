package zapcores

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Format int

const (
	EnvKeyLogFormat         = "FLOGO_LOG_FORMAT"
	DefaultLogFormat        = FormatConsole
	FormatConsole    Format = iota
	FormatJSON
)

// zapCoreMap holds logger core impl
var zapCoreMap map[string]zapcore.Core

// zapTraceCoreMap holds trace logger core impl
var zapTraceCoreMap map[string]zapcore.Core

// LogCfg Default log configuration for core creation
var LogCfg zap.Config

// TraceLogCfg Default trace log configuration for core creation
var TraceLogCfg zap.Config

// LogLvl log level variable which will control type of logs
var LogLvl *zap.AtomicLevel

// TraceLogLvl log level variable which will control type of logs
var TraceLogLvl *zap.AtomicLevel

func init() {

	createDefaultConfiguration()

	createTraceConfiguration()

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

func createDefaultConfiguration() {

	logFormat := DefaultLogFormat
	envLogFormat := strings.ToUpper(os.Getenv(EnvKeyLogFormat))
	if envLogFormat == "JSON" {
		logFormat = FormatJSON
	}

	cfg := zap.NewProductionConfig()
	cfg.DisableCaller = true

	eCfg := cfg.EncoderConfig
	eCfg.TimeKey = "timestamp"
	//eCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	eCfg.EncodeTime = zapcore.EpochNanosTimeEncoder

	if logFormat == FormatConsole {
		eCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		cfg.Encoding = "console"
		eCfg.EncodeName = nameEncoder
	}

	cfg.EncoderConfig = eCfg

	lvl := cfg.Level

	// assign single lvl instance to global level
	LogLvl = &lvl

	// assigning default configuration to global Cfg
	LogCfg = cfg
}

func createTraceConfiguration() {
	cfg := LogCfg

	if strings.Compare(cfg.Encoding, "console") == 0 {
		cfg.EncoderConfig.EncodeLevel = traceLevelEncoder
	}

	lvl := cfg.Level
	lvl.SetLevel(zapcore.DebugLevel)

	TraceLogCfg = cfg
	TraceLogLvl = &lvl
}

func nameEncoder(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + loggerName + "] -")
}

func traceLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[TRACE]")
}
