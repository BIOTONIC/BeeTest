package controllers

import (
	"BeeTest/models"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"reflect"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	// 如果是退出登录，清空 Cookie 和 Session
	if c.GetString("exist") != "" {
		c.Ctx.SetCookie("uname", "", -1, "/")
		c.DelSession("uname")
	}
	// 用login.html的模版
	c.TplName = "login.html"
}

func (ctrl *LoginController) Post() {
	// Input()携带了表单提交项
	uname := ctrl.Input().Get("uname")
	psw := ctrl.Input().Get("psw")
	// ==优先级大于:= 所以autoLogin是true or false
	autoLogin := ctrl.Input().Get("autoLogin") == "on"

	user, err := models.GetUser(uname)
	if err != nil {
		// 报错
		// 302 该资源原本确实存在，但已经被临时改变了位置
		beego.Error(err)
		ctrl.Redirect("/login", 302)
		return
	}

	// 密码正确 明文的 没有加密
	if user.PassWord == psw {
		maxAge := 0
		// 如果autoLogin maxAge设置到最大
		if autoLogin {
			maxAge = 1 << 31 - 1
			// TODO maxAge应该是int32类型吧
			fmt.Println(reflect.TypeOf(maxAge))
		}
		// 登录成功，放入cookie和session
		ctrl.Ctx.SetCookie("uname", uname, maxAge, "/")
		ctrl.SetSession("uname", uname)
		// 301 该资源已经被永久改变了位置 发送HTTP Location来重定向到正确的新位置
		ctrl.Redirect("/", 301)
	} else {
		ctrl.Redirect("/login", 301)
	}
	return
}

// 检查账户
func checkAccount(ctx *context.Context) bool {
	if ctx.Input.CruSession.Get("uname") != "" &&
		ctx.Input.CruSession.Get("uname") != nil {
		return true
	}
	// Cookie中没有uname 当前是未登录状态
	ck, err := ctx.Request.Cookie("uname")
	if err != nil {
		return false
	}
	// 根据Cookie中存在的uname给session加上uname 达到autoLogin的目的
	uname := ck.Value
	if beego.AppConfig.String("uname") == uname {
		ctx.Input.CruSession.Set("uname", uname)
		return true
	}
	return false
}