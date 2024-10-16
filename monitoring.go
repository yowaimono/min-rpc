package minrpc

import (
    "log"
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// LoggingMiddleware 是一个日志中间件
type LoggingMiddleware struct{}

func (m *LoggingMiddleware) Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
    })
}

// PrometheusMiddleware 是一个Prometheus监控中间件
type PrometheusMiddleware struct {
    requestCounter *prometheus.CounterVec
    requestLatency *prometheus.HistogramVec
}

func NewPrometheusMiddleware() *PrometheusMiddleware {
    m := &PrometheusMiddleware{
        requestCounter: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "http_requests_total",
                Help: "Total number of HTTP requests.",
            },
            []string{"method", "path"},
        ),
        requestLatency: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "http_request_duration_seconds",
                Help: "HTTP request latency distribution.",
            },
            []string{"method", "path"},
        ),
    }

    prometheus.MustRegister(m.requestCounter)
    prometheus.MustRegister(m.requestLatency)

    return m
}

func (m *PrometheusMiddleware) Handle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        m.requestCounter.WithLabelValues(r.Method, r.URL.Path).Inc()
        m.requestLatency.WithLabelValues(r.Method, r.URL.Path).Observe(time.Since(start).Seconds())
    })
}

// StartPrometheusServer 启动Prometheus监控服务器
func StartPrometheusServer(addr string) {
    http.Handle("/metrics", promhttp.Handler())
    log.Printf("Starting Prometheus server on %s", addr)
    go http.ListenAndServe(addr, nil)
}