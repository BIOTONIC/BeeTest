package models

import (
	"time"
	"strconv"
	"github.com/astaxie/beego/orm"
	"strings"
	"github.com/astaxie/beego"
)

type Topic struct {
	Id              int64
	Uid             int64
	Title           string
	Content         string `orm:"size(5000)"`
	Category        string
	Labels          string
	Attachment      string
	Created         time.Time `orm:"index"`
	Updated         time.Time `orm:"index"`
	Views           int64     `orm:"index"`
	Author          string
	ReplyTime       time.Time
	ReplyCount      int64
	ReplyLastUserId int64
}

func AddTopic(category, title, label, content, attachment string) error {
	// 处理标签
	// 标签以一个空格分隔
	// 存储时以$开头 中间用#分隔
	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"

	o := orm.NewOrm()
	topic := &Topic{
		Title:      title,
		Content:    content,
		Category:   category,
		Labels:     label,
		Created:    time.Now(),
		Updated:    time.Now(),
		ReplyTime:  time.Now(),
		Attachment: attachment,
	}

	_, err1 := o.Insert(topic)
	if err1 != nil {
		return err1
	}

	// 分类文章数 + 1
	// 同一个package下的GetCategory方法直接调用
	cate, err1 := GetCategory(category)
	if cate != nil {
		cate.TopicCount++
		o.Update(cate)
	}

	return nil
}

func ModifyTopic(id, category, title, label, content, attachment string) error {
	// 处理标签
	label = "$" + strings.Join(strings.Split(label, " "), "#$") + "#"

	// id传进来是string 转成int64
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	topic := &Topic{Id: cid}

	topic, err = GetTopic(id, false)

	// 如果文章分类改变了
	// EqualFold 判断两个utf-8编码字符串，大小写不敏感
	if !strings.EqualFold(topic.Category, category) {
		// 分类文章数 + 1
		cate, err1 := GetCategory(category)
		if err1 != nil {
			// 输出 级别是Error
			// Trace Debug Info Warn Error Critical
			beego.Error(err1)
		}
		if cate != nil {
			cate.TopicCount++
			o.Update(cate)
		}

		// 修改前的分类数 - 1
		cate_before, err2 := GetCategory(topic.Category)
		if err2 != nil {
			beego.Error(err2)
		}
		if cate_before != nil {
			cate_before.TopicCount--
			o.Update(cate_before)
		}
	}

	// if can not find one
	// create a new one
	if o.Read(topic) == nil {
		topic.Title = title
		topic.Content = content
		topic.Category = category
		topic.Labels = label
		topic.Updated = time.Now()
		topic.Attachment = attachment
		_, err2 := o.Update(topic)
		if err2 != nil {
			return err2
		}
	}

	return nil
}

// 获取所有的文章
func GetAllTopics(cate, label string, isDesc bool) ([]*Topic, error) {
	o := orm.NewOrm()
	// get a origin slice of type *Topic which length is 0
	topics := make([]*Topic, 0)
	qs := o.QueryTable("topic")
	// variant err, type is error
	var err error
	if isDesc {
		// 降序
		// OrderBy 默认是ASC column名前面加上-表示DESC
		if len(cate) > 0 {
			// 按照cate来限定
			_, err = qs.Filter("category", cate).OrderBy("-created").All(&topics)
		} else if len(label) > 0 {
			// 按照label来限定 __两个下划线作用类似于. 这里就是labels字符串中包含label的意思
			_, err = qs.Filter("labels__icontains", label).OrderBy("-created").All(&topics)
		} else {
			// 没有限定
			_, err = qs.OrderBy("-created").All(&topics)

		}
	} else {
		// 升序 默认的
		// 升序时没有cate和label限定
		_, err = qs.All(&topics)
	}
	return topics, err
}

// 根据 ID 获得文章
func GetTopic(id string, count bool) (*Topic, error) {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}
	o := orm.NewOrm()
	qs := o.QueryTable("topic")

	topic := new(Topic)
	err = qs.Filter("id", cid).One(topic)

	// 是否增加浏览数
	if count {
		topic.Views++
		_, err = o.Update(topic)
	}

	// labels拆分
	// 先#变空格 -1代表无数量限制
	// 再删去所有的$符号
	topic.Labels = strings.Replace(strings.Replace(topic.Labels, "#", " ", -1), "$", "", -1)

	return topic, err
}

// 删除文章 并分类文章数-1
func DelTopic(id string) error {
	cid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	o := orm.NewOrm()
	topic := &Topic{Id: cid}
	_, err = o.Delete(topic)

	topic, err = GetTopic(id, false)

	// 分类文章数 - 1
	cate, err1 := GetCategory(topic.Category)
	if err1 != nil {
		beego.Error(err1)
	}
	if cate != nil {
		cate.TopicCount--
		o.Update(cate)
	}

	return err
}
