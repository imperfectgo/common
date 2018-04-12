package log

import (
	"sync"

	"go.uber.org/zap"
)

var (
	globalMu     sync.RWMutex
	globalL      *zap.Logger
	globalQuickL *zap.Logger
	globalS      *zap.SugaredLogger
)

const (
	quickLDepth = 1
)

func init() {
	err := initGlobalLogger(defaultLogFormatURI)
	if err != nil {
		panic(err)
	}
}

func initGlobalLogger(format string) error {
	c, err := ParseConfigFromURIString(format)
	if err != nil {
		return err
	}

	l, err := c.Build()
	if err != nil {
		return err
	}
	ReplaceGlobals(l)
	return nil
}

// L returns the global zap.Logger.
//
// It's safe for concurrent use.
func L() *zap.Logger {
	globalMu.RLock()
	l := globalL
	globalMu.RUnlock()
	return l
}

// S returns the global zap.SugaredLogger.
// ReplaceGlobals.
//
// It's safe for concurrent use.
func S() *zap.SugaredLogger {
	globalMu.RLock()
	s := globalS
	globalMu.RUnlock()
	return s
}

// ReplaceGlobals replaces the global zap.Logger and the zap.SugaredLogger, and returns
// a function to restore the original values.
//
// It's safe for concurrent use.
func ReplaceGlobals(logger *zap.Logger) func() {
	globalMu.Lock()
	prev := globalL
	globalL = logger
	globalQuickL = logger.WithOptions(zap.AddCallerSkip(quickLDepth))
	globalS = logger.Sugar()
	globalMu.Unlock()

	// Replace zap's global logger as well
	restoreZapLogger := zap.ReplaceGlobals(logger)

	return func() {
		restoreZapLogger()
		ReplaceGlobals(prev)
	}
}
