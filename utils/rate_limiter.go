package utils

import (
	"net/http"

	"golang.org/x/time/rate"
)

// NewRateLimiter builds a limiter with given rps and burst.
func NewRateLimiter(rps float64, burst int) *rate.Limiter {
	return rate.NewLimiter(rate.Limit(rps), burst)
}

// RateLimitMiddleware applies request limiting to handlers.
func RateLimitMiddleware(limiter *rate.Limiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
