package http_proxy_router

import (
	"github.com/gin-gonic/gin"
	"github.com/go_gateway/http_proxy_middleware"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Use(http_proxy_middleware.HTTPAccessModeMiddleware(),
		http_proxy_middleware.HTTPFlowCountMiddleware(),
		http_proxy_middleware.HTTPWhiteListMiddleware(),
		http_proxy_middleware.HTTPBlackListMiddleware(),
		http_proxy_middleware.HTTPHeaderTransferMiddleware(),
		http_proxy_middleware.HTTPStripUrlMiddleware(),
		http_proxy_middleware.HTTPReWriteUrlMiddleware(),
		http_proxy_middleware.HTTPReverseProxyMiddleware())
	return router
}
