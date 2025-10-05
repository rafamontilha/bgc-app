package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type routeMetrics struct {
	Requests     int64   `json:"requests"`
	Status2xx    int64   `json:"status_2xx"`
	Status4xx    int64   `json:"status_4xx"`
	Status5xx    int64   `json:"status_5xx"`
	SumLatencyMs int64   `json:"sum_latency_ms"`
	AvgLatencyMs float64 `json:"avg_latency_ms"`
}

var (
	metricsStart   = time.Now()
	metricsMu      sync.RWMutex
	totalRequests  int64
	statusCounters = map[int]int64{}
	byRoute        = map[string]*routeMetrics{}
)

func newReqID() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return hex.EncodeToString(b)
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.Request.Header.Get("X-Request-Id")
		if rid == "" {
			rid = newReqID()
		}
		c.Set("req_id", rid)
		c.Writer.Header().Set("X-Request-Id", rid)
		c.Next()
	}
}

func MetricsAndLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latMs := time.Since(start).Milliseconds()
		status := c.Writer.Status()
		method := c.Request.Method

		route := c.FullPath()
		if route == "" {
			route = c.Request.URL.Path
		}
		key := method + " " + route

		metricsMu.Lock()
		totalRequests++
		statusCounters[status]++
		rm := byRoute[key]
		if rm == nil {
			rm = &routeMetrics{}
			byRoute[key] = rm
		}
		rm.Requests++
		rm.SumLatencyMs += latMs
		switch {
		case status >= 200 && status < 300:
			rm.Status2xx++
		case status >= 400 && status < 500:
			rm.Status4xx++
		case status >= 500:
			rm.Status5xx++
		}
		if rm.Requests > 0 {
			rm.AvgLatencyMs = float64(rm.SumLatencyMs) / float64(rm.Requests)
		}
		metricsMu.Unlock()

		rid, _ := c.Get("req_id")
		log.Printf(`{"ts":"%s","level":"info","req_id":"%v","method":"%s","path":"%s","route":"%s","status":%d,"latency_ms":%d,"ip":"%s","ua":"%s"}`,
			time.Now().Format(time.RFC3339Nano),
			rid, method, c.Request.URL.Path, route, status, latMs, c.ClientIP(), c.Request.UserAgent(),
		)
	}
}

func GetMetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		metricsMu.RLock()
		defer metricsMu.RUnlock()

		statusByString := map[string]int64{}
		for code, cnt := range statusCounters {
			statusByString[strconv.Itoa(code)] = cnt
		}
		routesCopy := map[string]routeMetrics{}
		for k, v := range byRoute {
			routesCopy[k] = *v
		}

		c.JSON(http.StatusOK, gin.H{
			"uptime_seconds":     int64(time.Since(metricsStart).Seconds()),
			"requests_total":     totalRequests,
			"requests_by_status": statusByString,
			"routes":             routesCopy,
		})
	}
}
