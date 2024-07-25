package middlewear

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func LoggingRequest(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("new request: method: %s, urlpath: %s", r.Method, r.URL.Path))

		start := time.Now()

		next.ServeHTTP(w, r)

		logger.Info(fmt.Sprintf("request %s completed in %v ms", r.URL.Path, time.Since(start).Milliseconds()))
	}
}
