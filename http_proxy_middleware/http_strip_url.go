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
func HTTPStripUrlMiddleware() gin.HandlerFunc {
	    return func(c *gin.Context) {
			serviceInterface, ok := c.Get("service")
			if !ok {
				middleware.ResponseError(c,2001,errors.New("service not found"))
				c.Abort()
				return
			}
			serviceDetail := serviceInterface.(*dao.ServiceDetail)
			fmt.Println("oldRule",c.Request.URL.Path)
			if serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL && serviceDetail.HTTPRule.NeedStripUri == 1 {
				c.Request.URL.Path = strings.Replace(c.Request.URL.Path, serviceDetail.HTTPRule.Rule,"", 1)

			}
			fmt.Println("newRule",c.Request.URL.Path)
		    c.Next()


	}
}