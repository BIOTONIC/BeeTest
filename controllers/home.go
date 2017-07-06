package controllers

import (
	"BeeTest/models"
	"github.com/astaxie/beego"
	"github.com/beego/i18n"
)

type BaseController struct {
	beego.Controller
	i18n.Locale
}

// 语言 本地化 i18n
// Prepare runs before other methods
func (c *BaseController) Prepare() {
	lang := c.GetString("lang")
	if lang == "zh-CN" {
		c.Lang = lang
	} else {
		c.Lang = "en-US"
	}
	c.Data["Lang"] = lang
}

// 结构体类型 嵌入类型时 该类型名充当字段名
type MainController struct {
	BaseController
}

func (c *MainController) Get() {
	// c.Tr()即translate
	// c.Data存储context data
	c.Data["hi"] = c.Tr("hi")
	c.Data["bye"] = c.Tr("bye")

	c.Data["IsHome"] = true
	var err error
	// 塞入查询的所有文章
	c.Data["Topics"], err = models.GetAllTopics(
		c.Input().Get("cate"),
		c.Input().Get("label"),
		true)
	if err != nil {
		beego.Error(err)
	}
	c.TplName = "home.html"

	c.Data["Categories"], err = models.GetAllCategories()
	c.Data["IsLogin"] = checkAccount(c.Ctx)
}
