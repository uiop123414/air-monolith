package main

import (
	"context"
	"net/http"
	"time"
)

func (app *application) TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r = r.WithContext(ctx)

			done := make(chan struct{})

			go func() {
				next.ServeHTTP(w, r)
				close(done)
			}()

			select {
			case <-done:
			case <-ctx.Done():
				http.Error(w, "Request Timeout", http.StatusRequestTimeout) // TODO after user was sent 408 response, upper goroutine continues working
				return
			}
		})
	}
}
