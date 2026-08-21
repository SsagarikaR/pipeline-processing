package middleware

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)



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

func (rl *rateLimiter) cleanupVisitors() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	for ip, v := range rl.visitors {
		if time.Since(v.lastSeen) > 3*time.Minute {
			delete(rl.visitors, ip)
		}
	}
}

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
