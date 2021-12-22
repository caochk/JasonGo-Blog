package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"
)

type Article struct {
	Id       int `orm:"pk"`
	UserId   int
	Category int8
	Headline string
	Content string
	Thumbnail string
	Credit int
	Readcount int
	Hide int8
	Drafted int8
	Checked int8
	CreateTime time.Time
	UpdateTime time.Time
	Users *Users `orm:"rel(fk)"`
}

// 查询article表中的所有数据
func (a *Article) Find_all() (*Article, error) {
	o := orm.NewOrm()
	article := &Article{}
	_, err := o.QueryTable("article").All(article)
	if err == nil {
		return article, err
	}
	return nil, err
}

// 根据id在article表中找到唯一对应数据
func (a *Article) Find_by_id(article_id int) (*Article, error) {
	o := orm.NewOrm()
	article := &Article{Id: article_id}
	err := o.Read(article, "articleid")
	if err == nil {
		return article, err
	}
	return nil, err
}

// article表与users表进行连接查询，返回10条记录。
// 返回10条记录的原因是博客系统首页中每页肯定只能展示一部分文章，在此定为每页10篇，然后分页。
func (a * Article) Find_paginated_articles(start int, count int) ([]*Article, error) {
	o := orm.NewOrm()
	var article []*Article
	_, err := o.QueryTable("article").Filter("hide", 0).
		Filter("drafted", 0).Filter("checked", 1).OrderBy("-id").Limit(count, start).All(&article)
	if err == nil {
		return article, err
	}
	return nil, err
}