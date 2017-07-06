package models

import (
	"time"
	"github.com/astaxie/beego/orm"
)

// 驼峰命名法会自动在数据库中转换成下划线分隔法
// UserName -> user_name
type User struct {
	Id       int64
	UserName string
	PassWord string
	Created  time.Time `orm:"index"`
}

// 多个变量都是string类型的
func AddUser(user_name, password string) error {
	o := orm.NewOrm()

	// user is a pointer
	user := &User{
		UserName: user_name,
		PassWord: password,
		Created:  time.Now(),
	}

	// user must be a pointer
	_, err1 := o.Insert(user)
	if err1 != nil {
		return err1
	}

	return nil
}

func GetUser(user_name string) (*User, error) {
	o := orm.NewOrm()
	// return a QuerySeter for table operations
	qs := o.QueryTable("user")
	user := new(User)
	// add condition expression to QuerySeter
	err := qs.Filter("user_name", user_name).One(user)
	return user, err
}
