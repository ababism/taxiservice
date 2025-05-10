package endpoints

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Ping(g *gin.RouterGroup) {
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
