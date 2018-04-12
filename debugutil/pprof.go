package debugutil

import (
	"net/http"
	"net/http/pprof"
	"runtime"
)

// PProfHandlers returns a map of pprof handlers keyed by the HTTP path.
func PProfHandlers(prefix string) map[string]http.Handler {
	// set only when there's no existing setting
	if runtime.SetMutexProfileFraction(-1) == 0 {
		// 1 out of 5 mutex events are reported, on average
		runtime.SetMutexProfileFraction(5)
	}

	m := make(map[string]http.Handler)

	m[prefix+"/"] = http.HandlerFunc(pprof.Index)
	m[prefix+"/profile"] = http.HandlerFunc(pprof.Profile)
	m[prefix+"/symbol"] = http.HandlerFunc(pprof.Symbol)
	m[prefix+"/cmdline"] = http.HandlerFunc(pprof.Cmdline)
	m[prefix+"/trace "] = http.HandlerFunc(pprof.Trace)
	m[prefix+"/heap"] = pprof.Handler("heap")
	m[prefix+"/goroutine"] = pprof.Handler("goroutine")
	m[prefix+"/threadcreate"] = pprof.Handler("threadcreate")
	m[prefix+"/block"] = pprof.Handler("block")
	m[prefix+"/mutex"] = pprof.Handler("mutex")
	return m
}
