package panicretry

import (
	"net/http"

	"github.com/acoshift/middleware"
)

// Config is the panicretry middleware config
type Config struct {
	MaxAttempts int
	Skipper     middleware.Skipper
}

const defaultAttempts = 3

// New creates new panicretry middleware
func New(c Config) middleware.Middleware {
	if c.MaxAttempts <= 0 {
		c.MaxAttempts = defaultAttempts
	}
	if c.Skipper == nil {
		c.Skipper = middleware.DefaultSkipper
	}

	return func(h http.Handler) http.Handler {
		try := func(w http.ResponseWriter, r *http.Request) (err interface{}) {
			defer func() {
				err = recover()
			}()
			h.ServeHTTP(w, r)
			return
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if c.Skipper(r) {
				h.ServeHTTP(w, r)
				return
			}

			var err interface{}
			for i := 0; i < c.MaxAttempts; i++ {
				if err = try(w, r); err == nil {
					return
				}
			}
			panic(err)
		})
	}
}
