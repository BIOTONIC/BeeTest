package models

import (
	"time"
	"strconv"
	"github.com/astaxie/beego/orm"
)

type Reply struct {
	Id      int64
	Tid     int64
	Name    string
	Content string    `orm:"size(1000)"`
	Created time.Time `orm:"index"`
}

// 添加评论
func AddReply(tid, nickname, content string) error {
	id, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	reply := &Reply{
		Tid:     id,
		Content: content,
		Name:    nickname,
		Created: time.Now(),
	}

	_, err1 := o.Insert(reply)
	if err1 != nil {
		return err1
	}

	// 添加评论数
	topic, err := GetTopic(tid, false)
	if topic != nil {
		topic.ReplyCount++
		// TODO 最佳评论后文章的更新时间就变了？
		topic.Updated = time.Now()
		topic.ReplyTime = time.Now()
		_, err1 = o.Update(topic)
		if err1 != nil {
			return err1
		}
	}

	return nil
}

// 删除评论 tid是topic id id是replay id
func DelReply(tid, id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	reply := &Reply{Id: cid}
	_, err = o.Delete(reply)

	// 减少评论数
	topic, err := GetTopic(tid, false)
	if topic != nil {
		topic.ReplyCount--
		o.Update(topic)
	}
	return err
}

// 获得所有评论
func GetAllReplies(tid string) ([]*Reply, error) {
	id, err := strconv.ParseInt(tid, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	// use make to create a slice
	replies := make([]*Reply, 0)
	qs := o.QueryTable("reply")
	// OrderBy 默认是ASC column名前面加上-表示DESC
	_, err = qs.Filter("tid", id).OrderBy("-created").All(&replies)
	return replies, err
}
