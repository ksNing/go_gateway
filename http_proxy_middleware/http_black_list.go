package http_proxy_middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go_gateway/dao"
	"github.com/go_gateway/middleware"
	"github.com/go_gateway/public"
	"strings"
)
//基于请求信息，匹配接入方式
func HTTPBlackListMiddleware() gin.HandlerFunc {
	    return func(c *gin.Context) {
			serviceInterface, ok := c.Get("service")
			if !ok {
				middleware.ResponseError(c,2001,errors.New("service not found"))
				c.Abort()
				return
			}
			serviceDetail := serviceInterface.(*dao.ServiceDetail)
			whiteIplist := []string{}
			if serviceDetail.AccessControl.WhiteList != "" {
				whiteIplist = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
			}

			blackIplist := []string{}
			if serviceDetail.AccessControl.BlackList != "" {
				blackIplist = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
			}

			if serviceDetail.AccessControl.OpenAuth == 1 && len(whiteIplist) == 0 && len(blackIplist) > 0 {
				if public.InStringSlice(blackIplist,c.ClientIP()) {
					middleware.ResponseError(c,2002,errors.New(fmt.Sprintln("%s in black ip list",c.ClientIP())))
					c.Abort()
					return
				}
			}
			
		    c.Next()


	}
}