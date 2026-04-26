package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type visitor struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	rps      rate.Limit
	burst    int
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(3 * time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.LastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func NewRateLimite(rps int64, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rps:      rate.Limit(rps),
		burst:    burst,
	}

	go rl.cleanupVisitors()
	return rl
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.rps, rl.burst)
		rl.visitors[ip] = &visitor{Limiter: limiter, LastSeen: time.Now()}
		return limiter
	}
	v.LastSeen = time.Now()
	return v.Limiter
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr //нужно учитывать X-Forwarded-For и другие заголовки
		if idx := len(ip) - 1; idx >= 0 {
			for i := len(ip) - 1; i >= 0; i-- {
				if ip[i] == ':' {
					ip = ip[:i]
					break
				}
			}
		}

		if !rl.getVisitor(ip).Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
