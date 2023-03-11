// Package we contains the small web based extensions
package web

import (
	"context"
	"net/http"
	"os"
	"syscall"

	"github.com/dimfeld/httptreemux/v5"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
	mw       []Middleware
}

func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
		mw:         mw,
	}
}

func (a *App) SignalShutDown() {
	a.shutdown <- syscall.SIGTERM
}

func (a *App) Handle(method string, group string, path string, handler Handler, mw ...Middleware) {
	// First add the handler specific middlewares handlers
	handler = wrapMiddleware(mw, handler)

	// Secondly add the application specific middlewares handlers
	handler = wrapMiddleware(a.mw, handler)

	// Function to execute for each request
	h := func(w http.ResponseWriter, r *http.Request) {

		// Calling wrapped handler function
		err := handler(r.Context(), w, r)
		if err != nil {
			// Error handling code
		}
	}

	// POST CODE PROCESSING CAN BE IMPLEMENTED HERE
	// <CODE FOR POST PROCESSING>
	finalPath := path
	if group != "" {
		finalPath = group + path
	}
	a.ContextMux.Handle(method, finalPath, h)
}
