package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLumberjacks(t *testing.T) {
	var fixtures = []struct {
		s      string
		expect LumberjackConfig
	}{
		{
			s: "filename=abc.log,compress=true;maxsize=500",
			expect: LumberjackConfig{
				Filename:   "abc.log",
				MaxSize:    500,
				MaxBackups: DefaultLumberjackMaxBackups,
				Compress:   true,
			},
		},
		{
			s: "filename=stdout.log,MaxAge=3200,localtime=true",
			expect: LumberjackConfig{
				Filename:   "stdout.log",
				MaxSize:    DefaultLumberjackMaxSize,
				MaxBackups: DefaultLumberjackMaxBackups,
				MaxAge:     3200,
				LocalTime:  true,
			},
		},
	}
	vs := make([]string, len(fixtures))
	for i := range fixtures {
		vs[i] = fixtures[i].s
	}

	// Parse multi
	lumberjacks, err := ParseLumberjacks(vs...)
	if !assert.NoError(t, err) {
		return
	}

	assert.Len(t, lumberjacks, len(vs))
	for i, l := range lumberjacks {
		assert.Equal(t, l, fixtures[i].expect, "lumberjack config at %d not match", i)
	}
}
