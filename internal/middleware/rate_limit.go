package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	requests map[uint][]time.Time
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[uint][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKeyID := GetAPIKeyID(r.Context())
		if apiKeyID == 0 {
			next.ServeHTTP(w, r)
			return
		}

		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()
		// Cleanup old requests
		threshold := now.Add(-rl.window)
		var valid []time.Time
		for _, t := range rl.requests[apiKeyID] {
			if t.After(threshold) {
				valid = append(valid, t)
			}
		}

		if len(valid) >= rl.limit {
			http.Error(w, `{"error":"Too Many Requests: rate limit exceeded"}`, http.StatusTooManyRequests)
			return
		}

		valid = append(valid, now)
		rl.requests[apiKeyID] = valid

		next.ServeHTTP(w, r)
	})
}
