package models

import (
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)


// init will be called before main
// one package or file can contain more than one init function
// init will be called automatically
// it can not be called explicitly

// register model driver & database and sync it
func RegisterDB() {
	// get params from config
	mysqluser := beego.AppConfig.String("mysqluser")
	mysqlpass := beego.AppConfig.String("mysqlpass")
	mysqlurls := beego.AppConfig.String("mysqlurls")
	mysqldb := beego.AppConfig.String("mysqldb")

	orm.RegisterModel(new(User),new(Category), new(Topic), new(Reply))
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", mysqluser + ":" + mysqlpass + "@tcp(" + mysqlurls + ")/" + mysqldb + "?charset=utf8&loc=Asia%2FShanghai")
	orm.RunSyncdb("default", false, true)
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
}
