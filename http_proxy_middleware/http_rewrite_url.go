package http_proxy_middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go_gateway/dao"
	"github.com/go_gateway/middleware"
	regexp2 "regexp"
	"strings"
)
//基于请求信息，匹配接入方式
func HTTPReWriteUrlMiddleware() gin.HandlerFunc {
	    return func(c *gin.Context) {
			serviceInterface, ok := c.Get("service")
			if !ok {
				middleware.ResponseError(c,2001,errors.New("service not found"))
				c.Abort()
				return
			}
			serviceDetail := serviceInterface.(*dao.ServiceDetail)
			for _, item := range strings.Split(serviceDetail.HTTPRule.UrlRewrite, ",") {
				items := strings.Split(item, " ")
				if len(items) != 2 {
					continue
				}

				regexp, err := regexp2.Compile(items[0])
				if err != nil {
					fmt.Println("err",err)
					continue
				}
				replacePath := regexp.ReplaceAll([]byte(c.Request.URL.Path), []byte(items[1]))
				c.Request.URL.Path = string(replacePath)

			}

		    c.Next()


	}
}