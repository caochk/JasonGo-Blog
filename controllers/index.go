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
	//c.TplName = "index.html"
	// 展示首页文章
	article := models.Article{}
	page_size := 10
	if result, err := article.FindPaginatedArticles(0, 10); err == nil {
		if total_article_num, err := article.GetTotalArticleNum(); err == nil {
			total_page_num := math.Ceil(float64(int(total_article_num)) / float64(page_size))
			fmt.Println("总页数：", total_page_num)
			c.Data["result"] = result
			fmt.Println("1")
			c.Data["page"] = 1
			fmt.Println("2")
			c.Data["total_page_num"] = total_page_num
			fmt.Println("3")
		} else {
			fmt.Println("获取文章总数失败")
			//c.AlertAndRedirect("", "")
		}
	}
}

// Paginate 用于展示首页之后的页面中文章列表（第二页、第三页等）【已测试通过】
func (c *HomeController) Paginate() {
	page, _ := strconv.Atoi(c.Ctx.Input.Param(":page"))
	page_size := 10
	start := (page - 1) * page_size
	article := models.Article{}
	if result, err := article.FindPaginatedArticles(start, page_size); err == nil {
		//for _, v := range result{
		//	fmt.Println(v.Id)
		//	fmt.Println("分页：", v.User)
		//}
		if total_article_num, err := article.GetTotalArticleNum(); err == nil {
			total_page_num := math.Ceil(float64(int(total_article_num)) / float64(page_size))
			fmt.Println("总页数：", total_page_num)
			c.Data["result"] = result
			c.Data["page"] = page
			c.Data["total_page_num"] = total_page_num
			//c.TplName = "index.html"
		} else {
			fmt.Println("获取文章总数失败")
			//c.AlertAndRedirect("", "")
		}
	} else {
		fmt.Println("error getting result:", err)
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
		for _, v := range result {
			fmt.Println("分类分页：", v.Id)
		}
	} else {
		fmt.Println("error getting result:", err)
	}
	// 获取用于展示分页的总页数
	if total_article_num, err := article.GetTotalArticleNumByCategory(category); err == nil {
		total_page_num := math.Ceil(float64(int(total_article_num)) / float64(page_size)) // golang的除法在计算小数/大数时=0，利用float64可避此坑
		fmt.Println("分类总页数：", total_article_num, total_page_num)
	} else {
		c.AlertAndRedirect("获取文章总数失败", "")
	}
}

// Search 搜索功能 【测试已通过】
func (c *HomeController) Search() {
	// 后端校验过滤无效搜索内容
	keyword := c.GetString("keyword")
	//var keyword string  // 获取前端参数方法2之bind
	//c.Ctx.Input.Bind(&keyword, "keyword")
	if len(keyword) == 0 || len(keyword) > 10 || strings.Contains(keyword, "%") {
		fmt.Println("404")
	}

	page, _ := strconv.Atoi(c.GetString("page"))
	page_size := 10
	start := (page - 1) * page_size
	var article models.Article
	// 获取搜索结果
	if result, err := article.FindByHeadline(keyword, start, page_size); err == nil {
		for _, v := range result {
			fmt.Println("搜索结果：", v.Id)
		}
	} else {
		fmt.Println("error getting result:", err)
	}
	// 搜索结果对应页数
	if total_article_num, err := article.GetTotalArticleNumByKeyword(keyword); err == nil {
		total_page_num := math.Ceil(float64(int(total_article_num)) / float64(page_size)) // golang的除法在计算小数/大数时=0，利用float64可避此坑
		fmt.Println("搜索总页数：", total_article_num, total_page_num)
	} else {
		c.AlertAndRedirect("获取文章总数失败", "")
	}
}
