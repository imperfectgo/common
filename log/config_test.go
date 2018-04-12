package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func defaultConfigWith(opts ...configOption) Config {
	c := defaultConfig
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

type configOption func(c *Config)

func withEncoderType(encoderType EncoderType) configOption {
	return func(c *Config) {
		c.EncoderType = encoderType
	}
}

func withDevelopment(development bool) configOption {
	return func(c *Config) {
		c.Development = development
	}
}

func withDisableCaller(disableCaller bool) configOption {
	return func(c *Config) {
		c.DisableCaller = disableCaller
	}
}

func withOutputPaths(paths []string) configOption {
	return func(c *Config) {
		c.OutputPaths = paths
	}
}

func withErrorOutputPaths(paths []string) configOption {
	return func(c *Config) {
		c.ErrorOutputPaths = paths
	}
}

func TestParseConfigFromURI(t *testing.T) {
	fixtures := []struct {
		uri      string
		expected Config
	}{
		{
			uri: "logger:console?outputPaths=stdout&development=true",
			expected: defaultConfigWith(
				withEncoderType(ConsoleEncoder),
				withDevelopment(true),
				withOutputPaths([]string{"stdout"}),
			),
		},
		{
			uri: "logger:json?outputPaths=stdout,stderr&errorOutputPaths=stdout,stderr&disableCaller=false",
			expected: defaultConfigWith(
				withEncoderType(JSONEncoder),
				withDisableCaller(false),
				withOutputPaths([]string{"stdout", "stderr"}),
				withErrorOutputPaths([]string{"stdout", "stderr"}),
			),
		},
	}

	for i, f := range fixtures {
		c, err := ParseConfigFromURIString(f.uri)
		if !assert.NoError(t, err, "Error parsing config, at index %d for uri %s", i, f.uri) {
			return
		}

		// DeepEqual() won't work for function types
		if !assert.Equal(t, &f.expected, c, "config not match, at index %d for uri %s", i, f.uri) {
			return
		}
	}
}
