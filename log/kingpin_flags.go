package log

import (
	"net/url"

	"gopkg.in/alecthomas/kingpin.v2"
)

// AddKingpinFlags adds the flags used by this package to the Kingpin application.
// To use the default Kingpin application, call AddFlags(kingpin.CommandLine)
func AddKingpinFlags(a *kingpin.Application) {
	s := loggerSettings{}
	a.Flag("log.level", "Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, dpanic, panic, fatal]").
		Default(baseLoggerLevel.String()).
		StringVar(&s.level)
	a.Flag("log.format", `Set the log target and format. Example: "logger:console?disableCaller=false&development=true&outputPaths=stdout&errorOutputPaths=stderr" or "logger:json?disableStacktrace=true"`).
		Default(defaultLogFormatURI).
		URLVar(&s.format)
	a.Action(s.apply)
}

type loggerSettings struct {
	level  string
	format *url.URL
}

func (s *loggerSettings) apply(ctx *kingpin.ParseContext) (err error) {
	c, err := ParseConfigFromURI(s.format)
	if err != nil {
		return err
	}

	l, err := c.Build()
	if err != nil {
		return err
	}

	restore := ReplaceGlobals(l)
	defer func() {
		if err != nil {
			restore()
		}
	}()

	return baseLoggerLevel.UnmarshalText([]byte(s.level))
}
