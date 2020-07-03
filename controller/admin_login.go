package controller

import (
	"encoding/json"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go_gateway/dao"
	"github.com/go_gateway/dto"
	"github.com/go_gateway/middleware"
	"github.com/go_gateway/public"
	"time"
)

type AdminLoginController struct {}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login",adminLogin.AdminLogin)
	group.GET("/login_out",adminLogin.AdminLoginOut)

}

// AdminLogin godoc
// @Summary 管理员登陆
// @Description 管理员接口
// @Tags 管理员
// @ID /admin_login/login
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
func (adminLogin *AdminLoginController) AdminLogin(c *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c,2000,err)
		return
	}

	//1.params.UserName 取得管理员的信息adminInfo
	//2.将adminInfo中的salt字段与params中的password进行和操作
	//3.将第二步和操作的结果与adminInfo中的password对比

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c,2001,err)
		return
	}
	admin := &dao.Admin{}
	admin, err = admin.LoginCheck(c,tx,params)
	if err != nil {
		middleware.ResponseError(c,2002,err)
	}

	//设置session
	sessInfo := &dto.AdminSessionInfo{
		ID:  admin.Id,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}
	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		middleware.ResponseError(c,2003,err)
		return
	}
	sess := sessions.Default(c)
	sess.Set(public.AdminSessionInfoKey, string(sessBts))
	sess.Save()

	out := &dto.AdminLoginOutput{Token:admin.UserName}
	middleware.ResponseSuccess(c,out)

}


// AdminLogin godoc
// @Summary 管理员退出
// @Description 管理员接口
// @Tags 管理员
// @ID /admin_login/login_out
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin_login/login_out [get]
func (adminLogin *AdminLoginController) AdminLoginOut(c *gin.Context) {

	sess := sessions.Default(c)
	sess.Delete(public.AdminSessionInfoKey)
	sess.Save()
	middleware.ResponseSuccess(c,"")
}
