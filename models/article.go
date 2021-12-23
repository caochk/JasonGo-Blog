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

// FindAll 查询article表中的所有数据
func (a *Article) FindAll() (*Article, error) {
	o := orm.NewOrm()
	article := &Article{}
	_, err := o.QueryTable("article").All(article)
	if err == nil {
		return article, err
	}
	return nil, err
}

// FindById 根据id在article表中找到唯一对应数据
func (a *Article) FindById(article_id int) (*Article, error) {
	o := orm.NewOrm()
	article := &Article{Id: article_id}
	err := o.Read(article, "articleid")
	if err == nil {
		return article, err
	}
	return nil, err
}

// FindPaginatedArticles article表与users表进行连接查询，返回10条记录。
// 返回10条记录的原因是博客系统首页中每页肯定只能展示一部分文章，在此定为每页10篇，然后分页。
func (a * Article) FindPaginatedArticles(start int, count int) ([]*Article, error) {
	o := orm.NewOrm()
	var articles []*Article
	_, err := o.QueryTable("article").Filter("hide", 0).
		Filter("drafted", 0).Filter("checked", 1).OrderBy("-id").Limit(count, start).All(&articles)
	if err == nil {
		return articles, err
	}
	return nil, err
}

// GetTotalArticleNum 获取文章（未隐藏、非草稿、已审核）总数量
func (a *Article) GetTotalArticleNum() (int64, error) {
	o := orm.NewOrm()
	total_article_num, err := o.QueryTable("article").Filter("hide", 0).
		Filter("drafted", 0).Filter("checked", 1).Count()
	if err == nil {
		return total_article_num, err
	}
	return -1, err
}

// FindByCategory 按照文章类型获取文章
func (a *Article) FindByCategory(category int, start int, count int) ([]*Article, error) {
	o := orm.NewOrm()
	var articles []*Article
	_, err := o.QueryTable("article").Filter("hide", 0).Filter("drafted", 0).Filter("checked", 1).
		Filter("category", category).OrderBy("-id").Limit(count, start).All(&articles)
	if err == nil {
		return articles, err
	}
	return nil, err
}

// GetTotalArticleNumByCategory 根据文章类型来获取文章总数量
func (a *Article) GetTotalArticleNumByCategory(category int) (int64, error) {
	o := orm.NewOrm()
	total_article_num, err := o.QueryTable("article").Filter("hide", 0).
		Filter("drafted", 0).Filter("checked", 1).Filter("category", category).Count()
	if err == nil {
		return total_article_num, err
	}
	return -1, err
}