package http_proxy_middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go_gateway/dao"
	"github.com/go_gateway/middleware"
	"github.com/go_gateway/reverse_proxy"
)
//基于请求信息，反向代理层
func HTTPReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c,2001,errors.New("service not found"))
			c.Abort()
			return
		}
		//fmt.Printf("serviceInterface",public.Obj2Json(serviceInterface))
		serviceDetail := serviceInterface.(*dao.ServiceDetail)
		lb, err := dao.LoadBalanceHandler.GetLoadBalancer(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2002, err)
			c.Abort()
			return
		}

		tran, err := dao.TransportorHandler.GetTrans(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			c.Abort()
			return
		}

		//创建reverseProxy,参数需要连接池和负载均衡器
		//使用reverseProxy.serverHttp(c,request,response)

		proxy := reverse_proxy.NewLoadBalanceReverseProxy(c,lb,tran)
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
		return




	}
}