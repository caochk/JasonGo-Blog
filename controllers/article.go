package controllers

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"strconv"
)
import "my_blog/models"

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
