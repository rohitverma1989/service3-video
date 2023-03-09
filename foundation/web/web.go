// Package we contains the small web based extensions
package web

import (
	"os"
	"syscall"

	"github.com/ardanlabs/service/vendor/github.com/dimfeld/httptreemux/v5"
)

// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
	*httptreemux.ContextMux
	shutdown chan os.Signal
}

func NewApp(shutdown chan os.Signal) *App {
	return &App{
		ContextMux: httptreemux.NewContextMux(),
		shutdown:   shutdown,
	}
}

func (a *App) SignalShutDown() {
	a.shutdown <- syscall.SIGTERM
}
