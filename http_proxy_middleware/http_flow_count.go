package http_proxy_middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go_gateway/dao"
	"github.com/go_gateway/middleware"
	"github.com/go_gateway/public"
	"time"
)

//基于请求信息，匹配接入方式
func HTTPFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serviceInterface.(*dao.ServiceDetail)
		//1统计全站
		//2统计服务
		//3统计租户
		totalCounter, err := public.FlowCounterHandler.GetCounter(public.Flowtotal)
		if err != nil {
			middleware.ResponseError(c, 2002, errors.New("flowtotal not found"))
			c.Abort()
			return
		}
		totalCounter.Increase()

		dayCount, _ := totalCounter.GetDayData(time.Now())
		fmt.Printf("totalCounter qps:%v,dayCount:%v", totalCounter.QPS, dayCount)

		serviceCounter, err := public.FlowCounterHandler.GetCounter(public.FlowServiceCountPrefix + serviceDetail.Info.ServiceName)
		if err != nil {
			middleware.ResponseError(c, 2003, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceCounter.Increase()
		dayServiceCount, _ := serviceCounter.GetDayData(time.Now())
		fmt.Printf("serviceCounter qps:%v,dayCount:%v", serviceCounter.QPS, dayServiceCount)

		c.Next()

	}
}
