package config

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

type defConfigImpl struct {
	logConfig      zap.Config
	logLevel       *zap.AtomicLevel
	traceLogConfig zap.Config
	traceLogLevel  *zap.AtomicLevel
}

// DefConfig returns default configuration values
type DefConfig interface {
	GetDefaultLogConfig() zap.Config
	GetDefaultLogLevl() *zap.AtomicLevel
	GetDefaultTraceLogConfig() zap.Config
	GetDefaultTraceLogLevl() *zap.AtomicLevel
}

var defCfg DefConfig

func init() {
	defCfg = createDefaultConfiguration()
}

// GetDefConfig returns default configuration
func GetDefConfig() DefConfig {
	return defCfg
}

func (d *defConfigImpl) GetDefaultLogConfig() zap.Config {
	return d.logConfig
}

func (d *defConfigImpl) GetDefaultLogLevl() *zap.AtomicLevel {
	return d.logLevel
}

func (d *defConfigImpl) GetDefaultTraceLogConfig() zap.Config {
	return d.traceLogConfig
}

func (d *defConfigImpl) GetDefaultTraceLogLevl() *zap.AtomicLevel {
	return d.traceLogLevel
}

func createDefaultConfiguration() DefConfig {

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

	// trace log configuration
	tcfg := cfg

	if strings.Compare(tcfg.Encoding, "console") == 0 {
		tcfg.EncoderConfig.EncodeLevel = traceLevelEncoder
	}

	tlvl := tcfg.Level
	tlvl.SetLevel(zapcore.DebugLevel)

	defCfg := &defConfigImpl{
		logConfig:      cfg,
		logLevel:       &lvl,
		traceLogConfig: tcfg,
		traceLogLevel:  &tlvl,
	}

	return defCfg
}

func nameEncoder(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + loggerName + "] -")
}

func traceLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[TRACE]")
}
