package log

import (
	"github.com/imperfectgo/zap-syslog"
	"github.com/imperfectgo/zap-syslog/syslog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultLogFormatURI = "logger:json?outputPaths=stderr"
)

var (
	baseLoggerLevel             = zap.NewAtomicLevel()
	defaultConsoleEncoderConfig = zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	defaultJSONEncoderConfig = zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	defaultSyslogEncoderConfig = zapsyslog.SyslogEncoderConfig{
		EncoderConfig: defaultJSONEncoderConfig,
		Framing:       zapsyslog.DefaultFraming,
		Facility:      syslog.LOG_LOCAL0,
	}
)
