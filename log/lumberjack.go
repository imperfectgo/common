package log

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	// DefaultLumberjackMaxSize is the default maximum size in megabytes of the
	// log file before it gets rotated.
	DefaultLumberjackMaxSize = 100
	// DefaultLumberjackMaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	DefaultLumberjackMaxBackups = 50
)

// LumberjackConfig offers a declarative way to construct a lumberjack Logger with rolling.
type LumberjackConfig struct {
	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.  It uses <processname>-lumberjack.log in
	// os.TempDir() if empty.
	Filename string `json:"filename" yaml:"filename"`
	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `json:"maxsize" yaml:"maxsize"`
	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `json:"maxage" yaml:"maxage"`
	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `json:"maxbackups" yaml:"maxbackups"`
	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool `json:"localtime" yaml:"localtime"`
	// Compress determines if the rotated log files should be compressed
	// using gzip.
	Compress bool `json:"compress" yaml:"compress"`
}

func parseLumberjack(s string) (*LumberjackConfig, error) {
	c := LumberjackConfig{
		MaxSize:    DefaultLumberjackMaxSize,
		MaxBackups: DefaultLumberjackMaxBackups,
	}

	// Scan into a kv map
	m := make(map[string]string)
	for s != "" {
		key := s
		if i := strings.IndexAny(s, ",;"); i >= 0 {
			key, s = key[:i], key[i+1:]
		} else {
			s = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}

		m[key] = value
	}

	// Parse kv map to lumberjack config
	const errMessageFormat = "lumberjack: error parsing %s"
	for key, value := range m {
		switch strings.ToLower(key) {
		case "filename":
			c.Filename = value
		case "maxsize":
			v, err := strconv.Atoi(value)
			if err != nil {
				return nil, errors.WithMessage(err, fmt.Sprintf(errMessageFormat, key))
			}
			c.MaxSize = v
		case "maxage":
			v, err := strconv.Atoi(value)
			if err != nil {
				return nil, errors.WithMessage(err, fmt.Sprintf(errMessageFormat, key))
			}
			c.MaxAge = v
		case "maxbackups":
			v, err := strconv.Atoi(value)
			if err != nil {
				return nil, errors.WithMessage(err, fmt.Sprintf(errMessageFormat, key))
			}
			c.MaxBackups = v
		case "localtime":
			v, err := strconv.ParseBool(value)
			if err != nil {
				return nil, errors.WithMessage(err, fmt.Sprintf(errMessageFormat, key))
			}
			c.LocalTime = v
		case "compress":
			v, err := strconv.ParseBool(value)
			if err != nil {
				return nil, errors.WithMessage(err, fmt.Sprintf(errMessageFormat, key))
			}
			c.Compress = v
		}
	}

	return &c, nil
}

// ParseLumberjacks parses lumberjack config values into LumberjackConfig slice.
func ParseLumberjacks(vs ...string) ([]LumberjackConfig, error) {
	lumberjacks := make([]LumberjackConfig, len(vs))
	for i, v := range vs {
		l, err := parseLumberjack(v)
		if err != nil {
			return nil, err
		}
		lumberjacks[i] = *l
	}
	return lumberjacks, nil
}

func openLumberjack(configs ...LumberjackConfig) zapcore.WriteSyncer {
	writers := make([]zapcore.WriteSyncer, 0, len(configs))

	for _, config := range configs {
		l := &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize,
			MaxAge:     config.MaxAge,
			MaxBackups: config.MaxBackups,
			LocalTime:  config.LocalTime,
			Compress:   config.Compress,
		}

		writers = append(writers, zapcore.AddSync(l))
	}

	writer := zapcore.NewMultiWriteSyncer(writers...)
	return writer
}
