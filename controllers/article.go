package controllers

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"my_blog/models"
	"strconv"
)

type ArticleController struct {
	beego.Controller
}

// GetArticleDetail 展示文章详情【测试通过】
func (c *ArticleController) GetArticleDetail() string {
	article_id, _ := strconv.Atoi(c.GetString("id"))
	article_model := models.Article{}
	if result, err := article_model.FindById(article_id); err != nil {
		fmt.Println("[ERROR] find by id:", err)
		return "500"
	} else {
		if result == nil {
			fmt.Println("[ERROR] cannot find any article from this id.")
			return "404"
		} else {
			fmt.Println(result.Content)
			return string(rune(result.Id))
		}
	}
}

// AddToFavorite 添加收藏【测试通过】
func (c *ArticleController) AddToFavorite() string {
	article_id, _ := strconv.Atoi(c.GetString("article_id"))
	if is_login := c.GetSession("islogin"); is_login == nil { // 未登录
		fmt.Println("未登录")
		return "not login"
	} else { // 已登录
		user_id, _ := c.GetSession("userid").(int)
		var favorite_model models.Favorite
		err := favorite_model.Like(article_id, user_id)
		if err != nil {
			fmt.Println("[ERROR] add to favorite:", err)
			return "favorite fail"
		}
		fmt.Println("成功收藏文章")
		return "favorite success"
	}
}

// CancelFavorite 取消收藏
func (c *ArticleController) CancelFavorite() string {
	article_id, _ := strconv.Atoi(c.GetString("article_id"))
	fmt.Println("article id from request:", article_id)
	if is_login := c.GetSession("islogin"); is_login == false { // 未登录
		fmt.Println("未登录")
		return "not login"
	} else { // 已登录
		user_id, _ := c.GetSession("userid").(int)
		var favorite_model models.Favorite
		err := favorite_model.Dislike(article_id, user_id)
		if err != nil {
			fmt.Println("[ERROR] cancel favorite:", err)
			return "cancel favorite fail"
		}
		fmt.Println("取消收藏文章")
		return "cancel favorite success"
	}
}
