package log

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/timonwong/zap-syslog"
	"github.com/timonwong/zap-syslog/syslog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultConfig = Config{
		Level:                   baseLoggerLevel,
		DisableCaller:           true,
		Lumberjacks:             nil,
		ErrorLumberjacks:        nil,
		defaultOutputPaths:      []string{"stderr"},
		defaultErrorOutputPaths: []string{"stderr"},
	}
)

// EncoderType represents encoder type.
type EncoderType int

const (
	// JSONEncoder represents a json encoder type.
	JSONEncoder EncoderType = iota
	// ConsoleEncoder represents a console encoder type, which is mainly for human friendly output.
	ConsoleEncoder
	// SyslogEncoder represents a Syslog (RFC5425) encoder type.
	SyslogEncoder
)

// Config offers a declarative way to construct a logger. It doesn't do
// anything that can't be done with New, Options, and the various
// zapcore.WriteSyncer and zapcore.Core wrappers, but it's a simpler way to
// toggle common options.
type Config struct {
	EncoderType EncoderType `json:"encoderType" yaml:"encoderType"`
	// Level is the minimum enabled logging level. Note that this is a dynamic
	// level, so calling Config.Level.SetLevel will atomically change the log
	// level of all loggers descended from this config.
	Level zap.AtomicLevel `json:"level" yaml:"level"`
	// Development puts the logger in development mode, which changes the
	// behavior of DPanicLevel and takes stacktraces more liberally.
	Development bool `json:"development" yaml:"development"`
	// DisableCaller stops annotating logs with the calling function's file
	// name and line number. By default, all logs are annotated.
	DisableCaller bool `json:"disableCaller" yaml:"disableCaller"`
	// DisableStacktrace completely disables automatic stacktrace capturing. By
	// default, stacktraces are captured for WarnLevel and above logs in
	// development and ErrorLevel and above in production.
	DisableStacktrace bool     `json:"disableStacktrace" yaml:"disableStacktrace"`
	OutputPaths       []string `json:"outputPaths" yaml:"outputPaths"`
	// ErrorOutputPaths is a list of paths to write internal logger errors to.
	// The default is standard error.
	//
	// Note that this setting only affects internal errors; for sample code that
	// sends error-level logs to a different location from info- and debug-level
	// logs, see the package-level AdvancedConfiguration example.
	ErrorOutputPaths []string           `json:"errorOutputPaths" yaml:"errorOutputPaths"`
	Lumberjacks      []LumberjackConfig `json:"lumberjacks" json:"lumberjacks"`
	ErrorLumberjacks []LumberjackConfig `json:"errorLumberjacks" json:"errorLumberjacks"`

	OutputAddresses []string `json:"outputAddresses" yaml:"outputAddresses"`
	// Syslog related config
	Framing  zapsyslog.Framing `json:"framing" yaml:"framing"`
	Facility syslog.Priority   `json:"facility" yaml:"facility"`
	Hostname string            `json:"hostname" yaml:"hostname"`
	PID      int               `json:"pid" yaml:"pid"`
	App      string            `json:"app" yaml:"app"`

	defaultOutputPaths      []string
	defaultErrorOutputPaths []string
}

// Build builds options into zap.Logger.
func (cfg Config) Build(opts ...zap.Option) (*zap.Logger, error) {
	var enc zapcore.Encoder
	switch cfg.EncoderType {
	case JSONEncoder:
		enc = zapcore.NewJSONEncoder(defaultJSONEncoderConfig)
	case ConsoleEncoder:
		enc = zapcore.NewConsoleEncoder(defaultConsoleEncoderConfig)
	case SyslogEncoder:
		encoderCfg := defaultSyslogEncoderConfig
		encoderCfg.Framing = cfg.Framing
		encoderCfg.Facility = cfg.Facility
		encoderCfg.Hostname = cfg.Hostname
		encoderCfg.PID = cfg.PID
		encoderCfg.App = cfg.App
		enc = zapsyslog.NewSyslogEncoder(defaultSyslogEncoderConfig)
	default:
		return nil, fmt.Errorf("unknown encoder type: %d", int(cfg.EncoderType))
	}

	var err error
	var sink zapcore.WriteSyncer
	var errSink zapcore.WriteSyncer

	switch cfg.EncoderType {
	case JSONEncoder, ConsoleEncoder:
		sink, errSink, err = cfg.openStandardSinks()
		if err != nil {
			return nil, err
		}

		lumberSink, errLumberSink := cfg.openLumberjackSinks()
		sink = zapcore.NewMultiWriteSyncer(sink, lumberSink)
		errSink = zapcore.NewMultiWriteSyncer(errSink, errLumberSink)

	case SyslogEncoder:
		sink, errSink, err = cfg.openSyslogSinks()
		if err != nil {
			return nil, err
		}
	}

	cfgOpts := cfg.buildOptions(errSink)
	zapOpts := make([]zap.Option, 0, len(cfgOpts)+len(opts))
	zapOpts = append(zapOpts, cfgOpts...)
	zapOpts = append(zapOpts, opts...)

	l := zap.New(
		zapcore.NewCore(enc, sink, cfg.Level),
		zapOpts...,
	)
	return l, nil
}

func (cfg Config) openStandardSinks() (zapcore.WriteSyncer, zapcore.WriteSyncer, error) {
	var outputPaths = cfg.OutputPaths
	var errorOutputPaths = cfg.ErrorOutputPaths
	if outputPaths == nil {
		outputPaths = cfg.defaultOutputPaths
	}
	if errorOutputPaths == nil {
		errorOutputPaths = cfg.defaultErrorOutputPaths
	}

	sink, closeOut, err := zap.Open(outputPaths...)
	if err != nil {
		return nil, nil, err
	}

	errSink, _, err := zap.Open(errorOutputPaths...)
	if err != nil {
		closeOut()
		return nil, nil, err
	}

	return sink, errSink, nil
}

func (cfg Config) openLumberjackSinks() (zapcore.WriteSyncer, zapcore.WriteSyncer) {
	sink := openLumberjack(cfg.Lumberjacks...)
	errSink := openLumberjack(cfg.ErrorLumberjacks...)
	return sink, errSink
}

func (cfg Config) openSyslogSinks() (zapcore.WriteSyncer, zapcore.WriteSyncer, error) {
	if len(cfg.OutputAddresses) == 0 {
		return nil, nil, fmt.Errorf("")
	}

	errSink, closeErr, err := zap.Open("stderr")
	if err != nil {
		return nil, nil, err
	}

	writeSyncers := make([]zapcore.WriteSyncer, 0, len(cfg.OutputAddresses))
	for _, addr := range cfg.OutputAddresses {
		networkAddr := strings.SplitN(addr, ":", 2)

		var address string
		network := "tcp"
		if len(networkAddr) == 1 {
			address = networkAddr[0]
		} else {
			network, address = networkAddr[0], networkAddr[1]
		}

		s, err := zapsyslog.NewConnSyncer(network, address)
		if err != nil {
			closeErr()
			return nil, nil, err
		}

		writeSyncers = append(writeSyncers, s)
	}

	sink := zapcore.NewMultiWriteSyncer(writeSyncers...)
	return sink, errSink, nil
}

func (cfg Config) buildOptions(errSink zapcore.WriteSyncer) []zap.Option {
	opts := []zap.Option{zap.ErrorOutput(errSink)}

	if cfg.Development {
		opts = append(opts, zap.Development())
	}

	if !cfg.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	stackLevel := zap.ErrorLevel
	if cfg.Development {
		stackLevel = zap.WarnLevel
	}
	if !cfg.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	return opts
}

func appendStringsFromStrings(output []string, vs []string) []string {
	output = append(output, vs...)
	return output
}

func appendStringsFromCommaSeparatedStrings(output []string, vs []string) []string {
	for _, s := range vs {
		output = append(output, strings.Split(s, ",")...)
	}
	return output
}

func (cfg *Config) populateCommonFromQS(values url.Values) error {
	for k, vs := range values {
		switch k {
		case "development":
			development, err := strconv.ParseBool(vs[0])
			if err != nil {
				return errors.WithMessage(err, "config: error parsing development")
			}
			cfg.Development = development
		case "disableCaller":
			disableCaller, err := strconv.ParseBool(vs[0])
			if err != nil {
				return errors.WithMessage(err, "config: error parsing disableCaller")
			}
			cfg.DisableCaller = disableCaller
		case "disableStacktrace":
			disableStacktrace, err := strconv.ParseBool(vs[0])
			if err != nil {
				return errors.WithMessage(err, "config: error parsing disableStacktrace")
			}
			cfg.DisableStacktrace = disableStacktrace
		}
	}

	return nil
}

func (cfg *Config) populateStandardEncoderFromQS(values url.Values) error {
	var outputPaths []string
	var errorOutputPaths []string

	for k, vs := range values {
		if len(vs) == 0 {
			continue
		}

		switch k {
		case "outputPath":
			outputPaths = appendStringsFromStrings(outputPaths, vs)
		case "outputPaths":
			outputPaths = appendStringsFromCommaSeparatedStrings(outputPaths, vs)
		case "errorOutputPath":
			errorOutputPaths = appendStringsFromStrings(errorOutputPaths, vs)
		case "errorOutputPaths":
			errorOutputPaths = appendStringsFromCommaSeparatedStrings(errorOutputPaths, vs)
		case "lumberjack":
			lumberjacks, err := ParseLumberjacks(vs...)
			if err != nil {
				return err
			}
			cfg.Lumberjacks = lumberjacks
		case "errorLumberjack":
			lumberjacks, err := ParseLumberjacks(vs...)
			if err != nil {
				return err
			}
			cfg.ErrorLumberjacks = lumberjacks
		}
	}

	cfg.OutputPaths = outputPaths
	cfg.ErrorOutputPaths = errorOutputPaths
	return nil
}

func (cfg *Config) populateSyslogEncoderFromQS(values url.Values) error {
	var outputAddresses []string

	for k, vs := range values {
		if len(vs) == 0 {
			continue
		}

		switch k {
		case "outputAddress":
			outputAddresses = appendStringsFromStrings(outputAddresses, vs)
		case "outputAddresses":
			outputAddresses = appendStringsFromCommaSeparatedStrings(outputAddresses, vs)
		case "framing":
			v := strings.ToLower(vs[0])
			switch v {
			case "", "0", "non-transparent", "non-transparent-framing":
				cfg.Framing = zapsyslog.NonTransparentFraming
			case "1", "octet-counting", "octet-counting-framing":
				cfg.Framing = zapsyslog.OctetCountingFraming
			default:
				return fmt.Errorf("config: unknown framing: %s", v)
			}
		case "facility":
			facility, err := syslog.FacilityPriority(vs[0])
			if err != nil {
				return errors.WithMessage(err, "config: error parsing facility")
			}
			cfg.Facility = facility
		case "hostname":
			cfg.Hostname = vs[0]
		case "pid":
			pid, err := strconv.Atoi(vs[0])
			if err != nil {
				return errors.WithMessage(err, "config: error parsing pid")
			}
			cfg.PID = pid
		case "app":
			cfg.App = vs[0]
		}
	}

	cfg.OutputAddresses = outputAddresses
	return nil
}

// ParseConfigFromURI parses config from a config uri.
func ParseConfigFromURI(u *url.URL) (*Config, error) {
	if u.Scheme != "logger" {
		return nil, fmt.Errorf("invalid scheme %s", u.Scheme)
	}

	values := u.Query()
	config := defaultConfig
	err := config.populateCommonFromQS(values)
	if err != nil {
		return nil, err
	}

	switch u.Opaque {
	case "console":
		config.EncoderType = ConsoleEncoder
		err := config.populateStandardEncoderFromQS(values)
		if err != nil {
			return nil, err
		}
	case "json":
		config.EncoderType = JSONEncoder
		err := config.populateStandardEncoderFromQS(values)
		if err != nil {
			return nil, err
		}
	case "syslog":
		config.EncoderType = SyslogEncoder
		err := config.populateSyslogEncoderFromQS(values)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported logger %q", u.Opaque)
	}

	return &config, nil
}

// ParseConfigFromURIString parses config from a config uri string.
func ParseConfigFromURIString(uri string) (*Config, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	return ParseConfigFromURI(u)
}
