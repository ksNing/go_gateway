package controller

import (
	"encoding/json"
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go_gateway/dao"
	"github.com/go_gateway/dto"
	"github.com/go_gateway/middleware"
	"github.com/go_gateway/public"
)

type AdminController struct {}

func AdminRegister(group *gin.RouterGroup) {
	adminLogin := &AdminController{}
	group.GET("/admin_info",adminLogin.AdminInfo)
	group.POST("/change_pwd",adminLogin.ChangePwd)


}
// AdminInfo godoc
// @Summary 管理员信息
// @Description 获取管理员信息接口
// @Tags 管理员
// @ID /admin/admin_info
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin/admin_info [get]
func (admin *AdminController) AdminInfo(c *gin.Context) {

	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)),adminSessionInfo); err != nil {
		middleware.ResponseError(c,2000,err)
		return
	}

	//1.读取sessionKey对应的json转换为结构体
	//2.读取数据然后封装输出结构体


	out := &dto.AdminInfoOutput{
		ID: adminSessionInfo.ID,
		Name: adminSessionInfo.UserName,
		LoginTime: adminSessionInfo.LoginTime,
		Avatar: "",
		Introduction: "",
		Roles:[]string{"admin"},
	}
	middleware.ResponseSuccess(c,out)

}

// ChangePwd godoc
// @Summary 修改密码
// @Description 修改密码
// @Tags 管理员
// @ID /admin/change_pwd
// @Accept  json
// @Produce  json
// @Param body body dto.ChangePwdInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (admin *AdminController) ChangePwd(c *gin.Context) {
	params := &dto.ChangePwdInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c,2000,err)
	}

	//1.session读取当前用户信息到结构体中sessInfo
	//2.sessInfo.username 读取数据库信息中的adminInfo
	//3.params.password + adminInfo.salt 生成加盐的新密码 saltPassword
	//4.将saltPassword存入到adminInfo中，并入库

	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)),adminSessionInfo); err != nil {
		middleware.ResponseError(c,2000,err)
		return
	}

	//连接数据库
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c,2001,err)
		return
	}
	adminInfo := &dao.Admin{}
	adminInfo, err = adminInfo.Find(c,tx,&dao.Admin{UserName:adminSessionInfo.UserName})
	if err != nil {
		middleware.ResponseError(c,2002,err)
		return
	}
	saltPwd := public.GenSaltPassword(adminInfo.Salt,params.Password)
	adminInfo.Password = saltPwd
	if err := adminInfo.Save(c,tx); err != nil {
		middleware.ResponseError(c,2003,err)
		return
	}
	middleware.ResponseSuccess(c,"")

}

