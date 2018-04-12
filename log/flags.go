package log

import (
	"flag"
	"net/url"
)

// AddFlags adds the flags used by this package to the given FlagSet. That's
// useful if working with a custom FlagSet. The init function of this package
// adds the flags to flag.CommandLine anyway. Thus, it's usually enough to call
// flag.Parse() to make the logging flags take effect.
func AddFlags(fs *flag.FlagSet) error {
	fs.Var(
		levelFlag("info"),
		"log.level",
		"Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, dpanic, panic, fatal]",
	)

	u, err := url.Parse(defaultLogFormatURI)
	if err != nil {
		return err
	}

	fs.Var(
		logFormatFlag(*u),
		"log.format",
		`Set the log target and format. Example: "logger:console?disableCaller=false&development=true&outputPaths=stdout&errorOutputPaths=stderr" or "logger:json?disableStacktrace=true"`,
	)

	return nil
}

type levelFlag string

// String implements flag.Value interface.
func (f levelFlag) String() string {
	return string(f)
}

// Set implements flag.Value interface.
func (f levelFlag) Set(level string) error {
	return baseLoggerLevel.UnmarshalText([]byte(level))
}

type logFormatFlag url.URL

// String implements flag.Value.
func (f logFormatFlag) String() string {
	u := url.URL(f)
	return u.String()
}

// Set implements flag.Value.
func (f logFormatFlag) Set(format string) error {
	return initGlobalLogger(format)
}
