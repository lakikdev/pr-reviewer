package middlewares

import (
	"net/http"
	"sync"
	"time"

	"pr-reviewer/internal/api/utils"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

// Create a custom visitor struct which holds the rate limiter for each
// visitor and the last time that the visitor was seen.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Change the the map to hold values of the type visitor.
var visitors = make(map[string]*visitor)
var mu sync.Mutex

// Run a background goroutine to remove old entries from the visitors map.
func init() {
	go cleanupVisitors()
}

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(5, 10)
		// Include the current time when creating a new visitor.
		visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	// Update the last seen time for the visitor.
	v.lastSeen = time.Now()
	return v.limiter
}

// Every minute check the map for visitors that haven't been seen for
// more than 3 minutes and delete the entries.
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

func LimitVisitor(next echo.HandlerFunc) echo.HandlerFunc {
	return (func(c echo.Context) error {
		r := c.Request()

		xForwardedFor := r.Header["X-Forwarded-For"]
		if len(xForwardedFor) == 0 {
			return next(c)
		}

		ip := xForwardedFor[0]

		limiter := getVisitor(ip)
		if !limiter.Allow() {
			return utils.NewHttpError(nil, http.StatusTooManyRequests, "Err500x002", "Too Many Requests")
		}

		return next(c)
	})
}
