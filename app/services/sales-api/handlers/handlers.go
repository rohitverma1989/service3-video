package handlers

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/ardanlabs/service/app/services/sales-api/handlers/debug/checkgrp"
	"github.com/ardanlabs/service/app/services/sales-api/handlers/v1/testgrp"
	"github.com/dimfeld/httptreemux/v5"
	"go.uber.org/zap"
)

func DebugStandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()
	// register all the standard library debug endpoints
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())
	return mux
}

func DebugMux(build string, log *zap.SugaredLogger) http.Handler {
	mux := DebugStandardLibraryMux()

	// register debug check endpoints
	cgh := checkgrp.Handlers{
		Build: build,
		Log:   log,
	}
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)
	return mux
}

// APIMuxConfig contains all the mandatory systems required by handlers.
type APIMuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

// APIMux constructs a http.Handler with all application routes defined.
func APIMux(cfg APIMuxConfig) *httptreemux.ContextMux {
	mux := httptreemux.NewContextMux()
	tgh := testgrp.Handlers{
		Log: cfg.Log,
	}
	mux.Handle(http.MethodGet, "/v1/test", tgh.Test)
	return mux
}
