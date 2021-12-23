package controllers

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"math"
	"my_blog/models"
	"my_blog/utils"
	"strconv"
	"strings"
)

type HomeController struct {
	utils.BaseController
	beego.Controller
}

// Home 展示首页文章列表【已测试通过】
func (c *HomeController) Home() {
	if c.Ctx.Request.Method == "GET" {
		// 展示首页文章
		article := models.Article{}
		result, err := article.FindPaginatedArticles(0, 10)
		if err == nil {
			c.Data["result"] = result
			for _, v := range result{
				fmt.Println(v.Id)
			}
		}
		//c.TplName = "base.html"
	}
}

// Paginate 用于展示首页之后的页面中文章列表（第二页、第三页等）【已测试通过】
func (c *HomeController) Paginate() {
	page, _ := strconv.Atoi(c.Ctx.Input.Param(":page"))
	page_size := 10
	start := (page - 1) * page_size
	article := models.Article{}
	if result, err := article.FindPaginatedArticles(start, page_size); err == nil {
		for _, v := range result{
			fmt.Println("分页：", v.Id)
		}
	} else {
		fmt.Println("error getting result:", err)
	}

	 if total_article_num, err := article.GetTotalArticleNum(); err == nil {
		 total_page_num := math.Ceil(float64(int(total_article_num)) / float64(page_size))
		 fmt.Println("总页数：", total_page_num)
	 } else {
		 c.AlertAndRedirect("获取文章总数失败", "")
	}
}

// Classify 展示分类页面【测试已通过】
func (c *HomeController) Classify() {
	category_page := c.Ctx.Input.Param(":category_page")
	category, _ := strconv.Atoi(strings.Split(category_page, "-")[0])
	page, _ := strconv.Atoi(strings.Split(category_page, "-")[1])

	page_size := 10
	start := (page - 1) * page_size
	article := models.Article{}
	// 获取分类文章结果
	if result, err := article.FindByCategory(category, start, page_size); err == nil {
		for _, v := range result{
			fmt.Println("分类分页：", v.Id)
		}
	} else {
		fmt.Println("error getting result:", err)
	}
	// 获取用于展示分页的总页数
	if total_article_num, err := article.GetTotalArticleNumByCategory(category); err == nil {
		total_page_num := math.Ceil(float64(int(total_article_num)) / float64(page_size))  // golang的除法在计算小数/大数时=0，利用float64可避此坑
		fmt.Println("分类总页数：", total_article_num, total_page_num)
	} else {
		c.AlertAndRedirect("获取文章总数失败", "")
	}
}