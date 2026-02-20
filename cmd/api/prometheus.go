package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func prometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()

		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = c.Request.URL.Path
		}

		httpRequestDuration.WithLabelValues(
			c.Request.Method,
			endpoint,
		).Observe(duration)

		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			endpoint,
			strconv.Itoa(c.Writer.Status()),
		).Inc()
	}
}
