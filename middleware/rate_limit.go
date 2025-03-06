package middleware

import (
	"fullcycle-goexpert-desafio-rate-limiter/limiter"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limiter_requests_total",
			Help: "Total number of requests handled by rate limiter",
		},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(requestCounter)
}

type RateLimitMiddleware struct {
	limiter limiter.RateLimiterInterface
}

func NewRateLimitMiddleware(limiter limiter.RateLimiterInterface) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiter: limiter,
	}
}

func (m *RateLimitMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		token := r.Header.Get("API_KEY")

		allowed, err := m.limiter.IsAllowed(ip, token)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !allowed {
			requestCounter.WithLabelValues("blocked").Inc()
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
			return
		} else {
			requestCounter.WithLabelValues("allowed").Inc()
		}
	}
}
