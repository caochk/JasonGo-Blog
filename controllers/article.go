package controllers

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-redis/redis/v8"
	"my_blog/models"
	"strconv"
)

type ArticleController struct {
	beego.Controller
}

// GetArticleDetail 展示文章详情【测试通过】
func (c *ArticleController) GetArticleDetail() string {
	article_id, _ := strconv.Atoi(c.GetString("id"))
	// 先查询缓存
	if article_content, err := rdb.Get(ctx, strconv.Itoa(article_id)).Result(); err == nil {
		fmt.Println("Redis中查到的文章详情：", article_content)
		return "200"
	} else if err == redis.Nil {  // 缓存未查到，再查询MySQL TODO 优化方向：多个客户端同时查询该key，但是后台依然只产生一个线程（协程？）去MySQL进行查询（这个除了用锁，用消息队列可行吗？）
		article_model := models.Article{}
		if result, err := article_model.FindById(article_id); err != nil {
			fmt.Println("[ERROR] find by id:", err)  // 查询MySQL时出错
			return "500"
		} else {
			if result == nil {  // MySQL中也查不到
				fmt.Println("[ERROR] cannot find any article from this id.")
				return "404"
			} else {  // MySQL中查询到了
				// 将MySQL中查询到的文章详情以哈希类型存入缓存
				fmt.Println(result.Content)
				var article_content = []string{
					strconv.Itoa(article_id),
					result.Content,
				}
				if _, err := rdb.HSet(ctx, "article_content", article_content).Result(); err == nil {
					fmt.Println("将MySQL中查询结果写入缓存成功")
					return "200"
				} else {
					fmt.Println("将MySQL中查询结果写入缓存失败：", err)
					return "404"
				}
				//if _, err := rdb.Set(ctx, strconv.Itoa(article_id), result.Content, 24 * time.Hour).Result(); err == nil {
				//	fmt.Println("将MySQL中查询结果写入缓存成功")
				//	return "200"
				//} else {
				//	fmt.Println("将MySQL中查询结果写入缓存失败：", err)
				//	return "404"
				//}
			}
		}
	} else {
		fmt.Println("[ERROR] 查询缓存时报错：", err)
		return "404"
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
