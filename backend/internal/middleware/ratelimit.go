package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"swiftly/backend/internal/pkg/response"
)

// visitor stores the rate limiter and the last seen time for each client IP
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

func init() {
	// Cleanup old visitors every minute to prevent memory leak
	go cleanupVisitors()
}

func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

func getVisitor(ip string, r rate.Limit, b int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(r, b)
		visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

// RateLimitMiddleware limits the number of requests per IP address
// r: number of tokens added to the bucket per second
// b: maximum number of tokens in the bucket (burst)
func RateLimitMiddleware(r rate.Limit, b int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ip, _, err := net.SplitHostPort(req.RemoteAddr)
			if err != nil {
				// Fallback to RemoteAddr if SplitHostPort fails
				ip = req.RemoteAddr
			}

			limiter := getVisitor(ip, r, b)
			if !limiter.Allow() {
				response.Error(w, http.StatusTooManyRequests, "Too many requests. Please slow down.", nil)
				return
			}

			next.ServeHTTP(w, req)
		})
	}
}
