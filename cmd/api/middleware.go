package main

import (
	"air-monolith/internal/rww"
	"context"
	"net/http"
	"time"
)

func (app *application) TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Wrap ResponseWriter
			rw := &rww.ResponseWriterWrapper{ResponseWriter: w}

			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r = r.WithContext(ctx)

			done := make(chan struct{})

			go func() {
				next.ServeHTTP(rw, r)
				close(done)
			}()

			select {
			case <-done:
			case <-ctx.Done():
				app.errorJSON(rw, http.StatusRequestTimeout) // TODO after user was sent 408 response, upper goroutine continues working
				return
			}
		})
	}
}
