package cache

import (
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"strconv"
	"time"
)

//var ctx = context.Background()
//var rdb = utils.InitRedisClient()

type CommentCacheController struct {
	beego.Controller
}

func (c CommentCacheController) AddComment2Redis(articleId string, userId int, content string, now time.Time) {
	commentMap := map[string]string{
		"articleId": articleId,
		"userId":    strconv.Itoa(userId),
		"content":   content,
		"date":      now.Format("2006-01-02 15:04:05"),
	}
	commentJson, _ := json.Marshal(commentMap)
	commentStr := string(commentJson)
	key := "article:comments:" + articleId
	rdb.LPush(ctx, key, commentStr)
	rdb.Expire(ctx, key, time.Hour*24)
}
