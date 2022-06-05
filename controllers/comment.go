package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	commentCache "my_blog/cache"
	"my_blog/utils/respUtils"
	"time"
)

type CommentController struct {
	beego.Controller
}

//var resp respUtils.Resp
var commentCacheController commentCache.CommentCacheController

// AddComment 添加评论【测试通过】（还差MySQL部分）
func (c CommentController) AddComment() {
	articleId := c.GetString("article_id")
	content := c.GetString("content")
	userId := 2
	now := time.Now()
	// 写入数据库
	// 写入缓存
	commentCacheController.AddComment2Redis(articleId, userId, content, now)
	c.Data["json"] = resp.NewResp(respUtils.SUCCESS_CODE, "评论添加成功")
	c.ServeJSON()
}
