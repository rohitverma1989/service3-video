package mid

import (
	"context"
	"net/http"

	"github.com/ardanlabs/service/foundation/web"
)

// Logger...
func Logger(handler web.Handler) web.Handler {
	h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		err := handler(ctx, w, r)
		if err != nil {
			return err
		}
		return nil
	}
	return h
}
