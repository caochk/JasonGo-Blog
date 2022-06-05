package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-redis/redis/v8"
	"math"
	"my_blog/models"
	"my_blog/utils"
	"strconv"
	"strings"
	"time"
)

//var ctx = context.Background()
// 新建上下文
//var ctx = context.Background()
// 在初始上下文的基础上新建带取消功能的上下文
var ctx, cancel = context.WithCancel(context.Background())
var rdb = utils.InitRedisClient()

type HomeController struct {
	utils.BaseController
	beego.Controller
}

// Home 展示首页文章列表【已测试通过】
func (c *HomeController) Home() {
	c.TplName = "index.html"
	// 展示首页文章
	article := models.Article{}
	page_size := 10
	if result, err := article.FindPaginatedArticles(0, page_size); err == nil {
		if total_article_num, err := article.GetTotalArticleNum(); err == nil {
			total_page_num := math.Ceil(float64(int(total_article_num)) / float64(page_size))
			fmt.Println("总页数：", total_page_num)
			fmt.Println("获取到的结果：", result)
			c.Data["result"] = result
			c.Data["page"] = 1
			c.Data["total_page_num"] = total_page_num
		} else {
			fmt.Println("获取文章总数失败")
		}
	}
}

// HomeRedis 从缓存读取首页文章列表【测试通过】
func (c HomeController) HomeRedis() {
	// 展示首页文章
	var page_size int64 = 10
	if articles, err := rdb.ZRangeByScore(ctx, "article", &redis.ZRangeBy{Min: "1", Max: "10"}).Result(); len(articles) == 10 { // 在缓存中查找到了10条记录
		if total_article_num, err := rdb.ZCard(ctx, "article").Result(); err == nil { // 从缓存的有序集合中读取文章总数
			total_page_num := math.Ceil(float64(int(total_article_num)) / float64(page_size))
			fmt.Println("从缓存中查到总页数：", total_page_num)
			//fmt.Println("获取到的结果：", articles)
			c.Data["result"] = articles
			c.Data["page"] = 1
			c.Data["total_page_num"] = total_page_num
		}
	} else if err == redis.Nil || len(articles) < 10 { // 缓存中为空，得去MySQL查询了
		article := models.Article{}
		page_size := 10
		if articles, err := article.FindPaginatedArticles(0, page_size); err != nil { // 查询MySQL时出错
			fmt.Println("获取文章总数失败")
			return
		} else {
			if articles == nil { // MySQL中也查不到
				fmt.Println("MySQL中也未查到")
			} else { // MySQL中查到了
				if total_article_num, err := article.GetTotalArticleNum(); err == nil {
					total_page_num := math.Ceil(float64(int(total_article_num)) / float64(page_size))
					// 返回结果
					fmt.Println("从MySQL中查到总页数：", total_page_num)
					//fmt.Println("获取到的结果：", articles)
					c.Data["result"] = articles
					c.Data["page"] = 1
					c.Data["total_page_num"] = total_page_num
					// 并将结果写入缓存
					for _, article := range articles {
						fmt.Println(article.Id)
						article_map := map[string]string{
							"user_id":     strconv.Itoa(article.User.Id),
							"category":    strconv.Itoa(int(article.Category)),
							"headline":    article.Headline,
							"content":     article.Content,
							"thumbnail":   article.Thumbnail,
							"credit":      strconv.Itoa(article.Credit),
							"readcount":   strconv.Itoa(article.Readcount),
							"hide":        strconv.Itoa(int(article.Hide)),
							"drafted":     strconv.Itoa(int(article.Drafted)),
							"checked":     strconv.Itoa(int(article.Checked)),
							"create_time": article.CreateTime.Format("2006-01-02 15:04:05"),
							"update_time": article.UpdateTime.Format("2006-01-02 15:04:05"),
						}
						artile_json, _ := json.Marshal(article_map)
						article_str := string(artile_json)
						var article_zset = redis.Z{}
						article_zset.Score = float64(article.Id)
						article_zset.Member = article_str
						if rel, err := rdb.ZAdd(ctx, "article", &article_zset).Result(); err == nil {
							fmt.Println("写入Redis成功：", rel)
						} else {
							fmt.Println("写入Redis失败：", err)
						}
					}
				} else {
					fmt.Println("文章列表是查到了，结果获取文章总数时出错")
				}
			}
		}
	} else { // 缓存查询出错
		fmt.Println("从缓存中获取指定数量文章失败：", err)
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

// PaginateRedis 从缓存读取首页之后页面的文章列表【未测试】
func (c HomeController) PaginateRedis() {
	page, _ := strconv.Atoi(c.Ctx.Input.Param(":page"))
	var page_size = 10
	start := (page - 1) * page_size
	if articles, err := rdb.ZRangeByScore(ctx, "article", &redis.ZRangeBy{Min: strconv.Itoa(start), Max: "10"}).Result(); len(articles) == 10 { // 在缓存中查找到了10条记录
		if total_article_num, err := rdb.ZCard(ctx, "article").Result(); err == nil { // 从缓存的有序集合中读取文章总数
			total_page_num := math.Ceil(float64(int(total_article_num)) / float64(page_size))
			fmt.Println("从缓存中查到总页数：", total_page_num)
			//fmt.Println("获取到的结果：", articles)
			c.Data["result"] = articles
			c.Data["page"] = 1
			c.Data["total_page_num"] = total_page_num
		}
	} else if err == redis.Nil || len(articles) < 10 { // 缓存中为空，得去MySQL查询了
		article := models.Article{}
		page_size := 10
		if articles, err := article.FindPaginatedArticles(0, page_size); err != nil { // 查询MySQL时出错
			fmt.Println("获取文章总数失败")
			return
		} else {
			if articles == nil { // MySQL中也查不到
				fmt.Println("MySQL中也未查到")
			} else { // MySQL中查到了
				if total_article_num, err := article.GetTotalArticleNum(); err == nil {
					total_page_num := math.Ceil(float64(int(total_article_num)) / float64(page_size))
					// 返回结果
					fmt.Println("从MySQL中查到总页数：", total_page_num)
					//fmt.Println("获取到的结果：", articles)
					c.Data["result"] = articles
					c.Data["page"] = 1
					c.Data["total_page_num"] = total_page_num
					// 并将结果写入缓存
					for _, article := range articles {
						fmt.Println(article.Id)
						article_map := map[string]string{
							"user_id":     strconv.Itoa(article.User.Id),
							"category":    strconv.Itoa(int(article.Category)),
							"headline":    article.Headline,
							"content":     article.Content,
							"thumbnail":   article.Thumbnail,
							"credit":      strconv.Itoa(article.Credit),
							"readcount":   strconv.Itoa(article.Readcount),
							"hide":        strconv.Itoa(int(article.Hide)),
							"drafted":     strconv.Itoa(int(article.Drafted)),
							"checked":     strconv.Itoa(int(article.Checked)),
							"create_time": article.CreateTime.Format("2006-01-02 15:04:05"),
							"update_time": article.UpdateTime.Format("2006-01-02 15:04:05"),
						}
						artile_json, _ := json.Marshal(article_map)
						article_str := string(artile_json)
						var article_zset = redis.Z{}
						article_zset.Score = float64(article.Id)
						article_zset.Member = article_str
						if rel, err := rdb.ZAdd(ctx, "article", &article_zset).Result(); err == nil {
							fmt.Println("写入Redis成功：", rel)
						} else {
							fmt.Println("写入Redis失败：", err)
						}
					}
				} else {
					fmt.Println("文章列表是查到了，结果获取文章总数时出错")
				}
			}
		}
	} else { // 缓存查询出错
		fmt.Println("从缓存中获取指定数量文章失败：", err)
	}
}

// Classify 展示分类页面【测试已通过】TODO 待接入Redis
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

//func (c *HomeController) SearchExecutor() {
//	var resp respUtils.Resp
//	fmt.Println("进入搜索执行者")
//	//var ctx context.Context
//	keyword := c.GetString("keyword")
//	page, _ := strconv.Atoi(c.GetString("page"))
//	select {
//	case <-c.Search(keyword, page):
//		fmt.Println("进入通道一")
//		//res := <-c.Search()
//		fmt.Println("Completed")
//	case <-ctx.Done():
//		fmt.Println("进入通道二")
//		fmt.Println("Canceled")
//		c.Data["json"] = resp.NewResp(respUtils.ERROR_CODE, "Canceled")
//		c.ServeJSON()
//	}
//}

// Search 搜索功能 【测试已通过】
func (c *HomeController) Search() {
	//res := make(chan int, 3)
	time.Sleep(10 * time.Second)

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

			//res <- v.Id
			//fmt.Println("写入通道完毕")
			//return res
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

	//res <- -1
	//return res
}

// Close 标签页关闭时调用一下该函数（临时）（后续打算应用至网站的方方面面以实现通过context通知作业协程退出）TODO 提供一种思路，通过SearchExecutor函数进行了尝试，暂时还有漏洞。
//func (c HomeController) Close() {
//	cancel()
//}
