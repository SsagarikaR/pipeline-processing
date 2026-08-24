package middleware

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// NewRateLimiter creates a per-IP rate limiter allowing r requests per
// second with a burst of b, and starts a background goroutine that
// forgets IPs that haven't been seen in a while so the map doesn't grow
// forever.
func NewRateLimiter(r rate.Limit, b int) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     r,
		burst:    b,
	}

	go func() {
		for {
			time.Sleep(time.Minute)
			rl.cleanupVisitors()
		}
	}()

	return rl
}

// cleanupVisitors drops any visitor that hasn't made a request in the
// last 3 minutes, so idle IPs don't sit in memory indefinitely.
func (rl *rateLimiter) cleanupVisitors() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	for ip, v := range rl.visitors {
		if time.Since(v.lastSeen) > 3*time.Minute {
			delete(rl.visitors, ip)
		}
	}
}

// getVisitor returns the rate limiter for an IP, creating one on first
// sight. It double-checks under the write lock so two concurrent
// requests from a brand-new IP can't create two separate limiters.
func (rl *rateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.RLock()
	v, exists := rl.visitors[ip]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		defer rl.mu.Unlock()
		// Double check
		v, exists = rl.visitors[ip]
		if !exists {
			limiter := rate.NewLimiter(rl.rate, rl.burst)
			rl.visitors[ip] = &visitor{limiter, time.Now()}
			return limiter
		}
	}

	rl.mu.Lock()
	v.lastSeen = time.Now()
	rl.mu.Unlock()

	return v.limiter
}

// Middleware rejects requests from an IP once it exceeds its rate limit,
// responding with 429 Too Many Requests instead of calling next.
func (rl *rateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.Header.Get("X-Real-IP")
		}
		if ip == "" {
			ip = r.RemoteAddr
		}

		limiter := rl.getVisitor(ip)

		if !limiter.Allow() {
			http.Error(w, "429 Too Many Requests - rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
