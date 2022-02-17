package cache

import (
	"encoding/json"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"my_blog/models"
	"strconv"
)

type UserCacheController struct {
	beego.Controller
}

// Users2Redis 将MySQL中user表的数据以哈希类型写入Redis【测试通过】 TODO 后续将此函数改造为定时任务，后期能否将其协程化
func (c UserCacheController) Users2Redis()  {
	user_model := models.User{}
	if users, err := user_model.FindAllUsers(); err == nil {
		for _, user := range users {
			var user_map = map[string]interface{}{
				"user_id": strconv.Itoa(user.Id),
				"password": user.Password,
				"nickname": user.Nickname,
				"avatar": user.Avatar,
				"qq": user.Qq,
				"role": user.Role,
				"credit": strconv.Itoa(user.Credit),
				"create_time": user.CreateTime.Format("2006-01-02 15:04:05"),
				"update_time": user.UpdateTime.Format("2006-01-02 15:04:05"),
			}

			user_json, _ := json.Marshal(user_map)
			user_str := string(user_json)

			var user_slice = []string{
				user.Username,
				user_str,
				}

			if rel, err := rdb.HMSet(ctx, "user", user_slice).Result(); err == nil {
				fmt.Println("user表数据写入Redis成功：", rel)
			} else {
				fmt.Println("user表数据写入Redis失败：", err)
			}
		}
	}
}