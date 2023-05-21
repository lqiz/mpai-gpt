package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lqiz/mpai/api/v1"
	"github.com/lqiz/mpai/response"
)

func RouteMp(parentRoute *gin.Engine) {
	// keepalive
	parentRoute.GET("/", func(c *gin.Context) {
		data := response.NewOKResponse(gin.H{"tests": "pong"})
		c.JSON(http.StatusOK, data)
	})

	v1Route := parentRoute.Group("/api")
	{
		v1.WxMsg(v1Route)
	}

}
