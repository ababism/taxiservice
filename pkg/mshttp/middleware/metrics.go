package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Prometheus(g *gin.RouterGroup) {
	g.GET("/metrics/prometheus", prometheusHandler())
	//g.GET("/metrics/prometheus",  metrics.HandleFunc())
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
