package logger

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"handler", "method", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of the HTTP request duration.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"handler", "method"},
	)
)

func Setup(app *echo.Echo) {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	app.Use(PrometheusMidleware)
	app.GET("debug/console", GetLogger)
	app.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}

func PrometheusMidleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().URL.Path == "/metrics" {
			promhttp.Handler().ServeHTTP(c.Response().Writer, c.Request())
			return nil
		}
		start := time.Now()
		httpRequestsTotal.WithLabelValues(c.Path(), c.Request().Method, "200").Inc()
		httpRequestDuration.WithLabelValues(c.Path(), c.Request().Method).Observe(time.Since(start).Seconds())
		return next(c)
	}
}
