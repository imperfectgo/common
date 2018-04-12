package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/alecthomas/kingpin.v2"
)

func newTestApp() *kingpin.Application {
	return kingpin.New("test", "").Terminate(nil)
}

func TestAddKingpinFlags(t *testing.T) {
	app := newTestApp()
	AddKingpinFlags(app)

	_, err := app.Parse([]string{"--log.level=debug", "--log.format", "logger:console?outputPaths=stderr"})
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "debug", baseLoggerLevel.String())
}
