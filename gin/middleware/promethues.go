package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "service"

/*
	service_http_request_size_bytes_count 总的请求次数
	service_http_request_size_bytes_sum 总的数据

	Histogram会自动创建3个指标，分别为：
		service_http_request_duration_seconds_count 事件发生的总次数
		service_http_request_duration_seconds_bucket 事件产生的值分布在bucket中的次数
		service_http_request_duration_seconds_sum  所有事件产生值的大小的总和
*/

var (
	labels = []string{"status", "endpoint", "method"}

	uptime = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "uptime",
			Help:      "HTTP service uptime.",
		}, nil,
	)

	reqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_request_count_total",
			Help:      "Total number of HTTP requests made.",
		}, labels,
	)

	reqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request latencies in seconds.",
		}, labels,
	)

	reqSizeBytes = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "http_request_size_bytes",
			Help:      "HTTP request sizes in bytes.",
		}, labels,
	)

	respSizeBytes = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: namespace,
			Name:      "http_response_size_bytes",
			Help:      "HTTP request sizes in bytes.",
		}, labels,
	)
)

func init() {
	prometheus.MustRegister(uptime, reqCount, reqDuration, reqSizeBytes, respSizeBytes)
	go recordUptime()
}

func recordUptime() {
	for range time.Tick(time.Second) {
		uptime.WithLabelValues().Inc()
	}
}

func calcRequestSize(r *http.Request) float64 {
	size := 0
	if r.URL != nil {
		size = len(r.URL.String())
	}

	size += len(r.Method)
	size += len(r.Proto)

	for name, values := range r.Header {
		size += len(name)
		for _, value := range values {
			size += len(value)
		}
	}
	size += len(r.Host)

	if r.ContentLength != -1 {
		size += int(r.ContentLength)
	}
	return float64(size)
}

type PromOpts struct {
	ExcludeRegexStatus   string
	ExcludeRegexEndpoint string
	ExcludeRegexMethod   string
}

var defaultPromOpts = &PromOpts{}

func (po *PromOpts) checkLabel(label, pattern string) bool {
	if pattern == "" {
		return true
	}

	matched, err := regexp.MatchString(pattern, label)
	if err != nil {
		return true
	}
	return !matched
}

func PromMiddleware(promOpts *PromOpts) gin.HandlerFunc {
	if promOpts == nil {
		promOpts = defaultPromOpts
	}

	return func(c *gin.Context) {

		start := time.Now()
		c.Next()

		status := fmt.Sprintf("%d", c.Writer.Status())
		endpoint := c.Request.URL.Path
		method := c.Request.Method

		lvs := []string{status, endpoint, method}

		isOk := promOpts.checkLabel(status, promOpts.ExcludeRegexStatus) &&
			promOpts.checkLabel(endpoint, promOpts.ExcludeRegexEndpoint) &&
			promOpts.checkLabel(method, promOpts.ExcludeRegexMethod)

		if !isOk {
			return
		}

		reqCount.WithLabelValues(lvs...).Inc()
		reqDuration.WithLabelValues(lvs...).Observe(time.Since(start).Seconds())
		reqSizeBytes.WithLabelValues(lvs...).Observe(calcRequestSize(c.Request))
		respSizeBytes.WithLabelValues(lvs...).Observe(float64(c.Writer.Size()))
	}
}

func PromHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
