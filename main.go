package main

import (
	"github.com/astaxie/beego/orm"
	"github.com/beego/i18n"
	"github.com/astaxie/beego"
	"BeeTest/controllers"
	"BeeTest/models"
)

func init() {
	models.RegisterDB()
}

func main() {
	i18n.SetMessage("en-US", "conf/locale_en-US.ini")
	i18n.SetMessage("zh-CN", "conf/locale_zh-CN.ini")

	// 注册多语言模板函数
	beego.AddFuncMap("i18n", i18n.Tr)

	orm.Debug = true

	// 路由写在main里面 其实写在routers里面应该更好
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/register", &controllers.RegisterController{})
	beego.Router("/category", &controllers.CategoryController{})
	beego.Router("/topic", &controllers.TopicController{})
	beego.AutoRouter(&controllers.TopicController{})
	beego.Router("/reply", &controllers.ReplyController{})
	beego.AutoRouter(&controllers.ReplyController{})
	beego.Router("/attachment/:all", &controllers.AttachController{})
	// 启动 beego
	beego.Run()
}

