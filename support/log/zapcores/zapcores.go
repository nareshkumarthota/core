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

// logCfg Default log configuration for core creation
var logCfg zap.Config

// traceLogCfg Default trace log configuration for core creation
var traceLogCfg zap.Config

// logLvl log level variable which will control type of logs
var logLvl *zap.AtomicLevel

// traceLogLvl log level variable which will control type of logs
var traceLogLvl *zap.AtomicLevel

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

func GetDefaultLogConfig() zap.Config {
	return logCfg
}

func GetDefaultLogLevl() *zap.AtomicLevel {
	return logLvl
}

func GetDefaultTraceLogConfig() zap.Config {
	return traceLogCfg
}

func GetDefaultTraceLogLevl() *zap.AtomicLevel {
	return traceLogLvl
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
	logLvl = &lvl

	// assigning default configuration to global Cfg
	logCfg = cfg
}

func createTraceConfiguration() {
	cfg := logCfg

	if strings.Compare(cfg.Encoding, "console") == 0 {
		cfg.EncoderConfig.EncodeLevel = traceLevelEncoder
	}

	lvl := cfg.Level
	lvl.SetLevel(zapcore.DebugLevel)

	traceLogCfg = cfg
	traceLogLvl = &lvl
}

func nameEncoder(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + loggerName + "] -")
}

func traceLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[TRACE]")
}
