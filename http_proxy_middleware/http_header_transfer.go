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
func HTTPHeaderTransferMiddleware() gin.HandlerFunc {
	    return func(c *gin.Context) {
			serviceInterface, ok := c.Get("service")
			if !ok {
				middleware.ResponseError(c,2001,errors.New("service not found"))
				c.Abort()
				return
			}
			serviceDetail := serviceInterface.(*dao.ServiceDetail)
			fmt.Println("header1",public.Obj2Json(serviceDetail.HTTPRule.HeaderTransfor))
			for _,item := range strings.Split(serviceDetail.HTTPRule.HeaderTransfor, ",") {
				items := strings.Split(item, " ")
				if len(items) != 3 {
					continue
				}
				if items[0] == "add" || items[0] == "edit" {
					c.Request.Header.Set(items[1],items[2])
				}
				if items[0] == "del" {
					c.Request.Header.Del(items[1])
				}
			}
			fmt.Println("header2",serviceDetail.HTTPRule.HeaderTransfor)
		    c.Next()


	}
}