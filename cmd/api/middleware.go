package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func (app *application) requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		app.logger.PrintInfo("request", map[string]string{
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"status":  strconv.Itoa(c.Writer.Status()),
			"latency": time.Since(start).String(),
		},
		)
	}
}

func (app *application) recoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				app.logger.PrintError(fmt.Errorf("%v", rec), map[string]string{"trace": "panic recovered"})
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
