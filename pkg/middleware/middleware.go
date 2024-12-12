package middleware

import (
	"air-monolith/internal/rww"
	pkg "air-monolith/pkg/utils"
	"context"
	"net/http"
	"time"
)

func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Wrap ResponseWriter
			rw := &rww.ResponseWriterWrapper{ResponseWriter: w}

			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			// Стандартная структура - все перекинуть в pkg
			// TODO - написать тесты с помощью mockery
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
				pkg.ErrorJSON(rw, http.StatusRequestTimeout) // TODO after user was sent 408 response, upper goroutine continues working
				return
			}
		})
	}
}
