package middlewares

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	defaultMaxRequestsPerMinutes = 30
	defaultRequestTimeout        = 60 * time.Second
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string]map[string][]time.Time
}

var rateLimiter *RateLimiter
var requestTimeout time.Duration = defaultRequestTimeout
var maxRequests int = defaultMaxRequestsPerMinutes

func RateLimiterMiddleware(next http.Handler, options map[string]interface{}) http.Handler {
	if rateLimiter == nil {
		rateLimiter = &RateLimiter{
			requests: make(map[string]map[string][]time.Time),
		}
	}

	if maxRequestsOption, ok := options["maxRequests"]; ok {
		maxRequests = int(maxRequestsOption.(float64))
	}

	if requestTimeoutOption, ok := options["requestTimeout"]; ok {
		requestTimeout = time.Duration(int(requestTimeoutOption.(float64))) * time.Second
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		rateLimiter.mu.Lock()
		defer rateLimiter.mu.Unlock()

		now := time.Now()
		if rateLimiter.requests[r.Host] == nil {
			rateLimiter.requests[r.Host] = make(map[string][]time.Time)
		}

		rateLimiter.requests[r.Host][clientIP] = filterOldRequests(rateLimiter.requests[r.Host][clientIP], now)

		if len(rateLimiter.requests[r.Host][clientIP]) >= maxRequests {
			fmt.Println("Too Many Requests")
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		rateLimiter.requests[r.Host][clientIP] = append(rateLimiter.requests[r.Host][clientIP], now)
		next.ServeHTTP(w, r)
	})
}

func filterOldRequests(requests []time.Time, currentTime time.Time) []time.Time {
	var result []time.Time
	for _, reqTime := range requests {
		if currentTime.Sub(reqTime) <= requestTimeout {
			result = append(result, reqTime)
		}
	}
	return result
}
