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
func HTTPWhiteListMiddleware() gin.HandlerFunc {
	    return func(c *gin.Context) {
			serviceInterface, ok := c.Get("service")
			if !ok {
				middleware.ResponseError(c,2001,errors.New("service not found"))
				c.Abort()
				return
			}
			serviceDetail := serviceInterface.(*dao.ServiceDetail)
			iplist := []string{}
			if serviceDetail.AccessControl.WhiteList != "" {
				iplist = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
			}
			if serviceDetail.AccessControl.OpenAuth == 1 && len(iplist) != 0 {
				if !public.InStringSlice(iplist,c.ClientIP()) {
					middleware.ResponseError(c,2002,errors.New(fmt.Sprintln("%s not in white ip list",c.ClientIP())))
					c.Abort()
					return
				}
			}

		    c.Next()


	}
}