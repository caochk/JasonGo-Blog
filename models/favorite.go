package models

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"time"
)

type Favorite struct {
	Id         int      `orm:"pk;auto"`
	Article    *Article `orm:"rel(fk)"`
	User       *User    `orm:"rel(fk)"`
	Canceled   int8
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateTime time.Time `orm:"auto_now;type(datetime)"`
}

// Like 收藏文章【测试通过】（逻辑有问题）
func (m *Favorite) Like(article_id int, user_id int) error {
	o := orm.NewOrm()
	user := &User{Id: user_id}
	article := &Article{Id: article_id}
	favorite := Favorite{
		User:     user,
		Article:  article,
		Canceled: 0,
	}
	if created, _, err := o.ReadOrCreate(&favorite, "article_id", "user_id"); err == nil {
		if created { // 原不存在，现已创建
			fmt.Println("insert an article to favorite.")
			return nil
		} else { // 原已存在
			favorite.Canceled = 0
			if _, err := o.Update(&favorite, "canceled"); err != nil {
				fmt.Println("[ERROR] update favorite failed:", err)
				return err
			}
			fmt.Println("收藏文章更新成功")
			return nil
		}
	} else {
		fmt.Println("[ERROR] read or create error:", err)
		return err
	}
}

// Dislike 取消收藏
func (m *Favorite) Dislike(article_id int, user_id int) error {
	fmt.Println(article_id)
	fmt.Println(user_id)
	o := orm.NewOrm()
	//user := &User{Id: user_id}
	//article := &Article{Id: article_id}
	//favorite := Favorite{
	//	User: user,
	//	Article: article,
	//	Canceled: 1,
	//}
	// 问题出在根据这些过滤条件找到的行为0
	var favorites []Favorite
	if _, err := o.QueryTable("favorite").Filter("Article", 20).Filter("User", 0).RelatedSel().All(&favorites); err != nil {
		fmt.Println("[ERROR] cancel the favorite article:", err)
		return err
	} else {
		for _, favorite := range favorites {
			fmt.Println(favorite.Id)
		}
		return nil
	}
	//if num, err := o.QueryTable("favorite").Filter("user_id", user_id).Filter("article_id", article_id).
	//	Update(orm.Params{"canceled":1,}); err != nil {
	//	fmt.Println("influenced lines:", num)
	//	fmt.Println("[ERROR] cancel the favorite article:", err)
	//	return err
	//} else {
	//	fmt.Println("influenced lines:", num)
	//	fmt.Println("取消收藏成功，来自model")
	//	return nil
	//}
	//if num, err := o.Update(&favorite); err != nil {
	//	fmt.Println("influenced lines:", num)
	//	fmt.Println("[ERROR] cancel the favorite article:", err)
	//	return err
	//} else {
	//	fmt.Println("influenced lines:", num)
	//	fmt.Println("取消收藏成功，来自model")
	//	return nil
	//}
}
