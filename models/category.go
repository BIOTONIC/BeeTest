package models

import (
	"time"
	"github.com/astaxie/beego/orm"
	"strconv"
)

type Category struct {
	Id              int64
	Title           string
	Created         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	TopicTime       time.Time `orm:"index"`
	TopicCount      int64
	TopicLastUserId int64
}

func AddCategory(name string) error {
	o := orm.NewOrm()
	cate := &Category{Title: name, Created: time.Now(), Views: 0, TopicTime: time.Now()}

	qs := o.QueryTable("category")
	// 先根据title去category表查询 查询一条填充到cate里
	err := qs.Filter("title", name).One(cate)
	// 查询结果是nil 说明当前title的category已存在
	if err == nil {
		return err
	}
	// 插入新的cate
	_, err1 := o.Insert(cate)
	if err1 != nil {
		return err1
	}
	return nil
}

func GetAllCategories() ([]*Category, error) {
	o := orm.NewOrm()
	cates := make([]*Category, 0)
	qs := o.QueryTable("category")
	_, err := qs.All(&cates)
	return cates, err
}

// 根据title获得category
func GetCategory(cate string) (*Category, error) {
	o := orm.NewOrm()
	qs := o.QueryTable("category")

	category := new(Category)
	err := qs.Filter("title", cate).One(category)

	return category, err
}

func DelCategory(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	cate := &Category{Id: cid}
	_, err = o.Delete(cate)
	return err
}