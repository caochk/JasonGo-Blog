package cache

import (
	"context"
	"encoding/json"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-redis/redis/v8"
	"my_blog/models"
	"my_blog/utils"
	"strconv"
)

var ctx = context.Background()
var rdb = utils.InitRedisClient()

type ArticleCacheController struct {
	beego.Controller
}

// Articles2Redis 将MySQL中article表的数据以有序集合类型写入Redis【测试通过】 TODO 后续将此函数改造为定时任务，后期能否将其协程化
func (ArticleCacheController)Articles2Redis()  {
	article_model := models.Article{}
	if articles, err := article_model.FindAllArticles(); err == nil {
		for _, article := range articles {
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
				"checked": strconv.Itoa(int(article.Checked)),
				"create_time": article.CreateTime.Format("2006-01-02 15:04:05"),
				"update_time": article.UpdateTime.Format("2006-01-02 15:04:05"),
			}
			artile_json,_ := json.Marshal(article_map)
			article_str := string(artile_json)
			var article_zset = redis.Z{}
			article_zset.Score = float64(article.Id)
			article_zset.Member = article_str
			rdb.ZAdd(ctx, "article", &article_zset)
		}
	}
}