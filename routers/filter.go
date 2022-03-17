package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"my_blog/utils/jwtUtils"
	"my_blog/utils/respUtils"
	"strings"
)

func Filter(ctx *context.Context) {
	var c *beego.Controller
	// 拦截以进行登录认证
	auth(ctx, c)
}

// 用户JWT认证
func auth(ctx *context.Context, c *beego.Controller) {
	var req = ctx.Request
	var resp respUtils.Resp

	var authString = ctx.Input.Header("Authorization")
	var kv = strings.Split(authString, " ")
	var token string
	if len(kv) != 2 || kv[0] != "Bearer" {
		c.Data["json"] = resp.NewResp(respUtils.ERROR_CODE, "AuthString无效")
		c.ServeJSON()
	} else {
		token = kv[1]
	}
	// 访问非登录页面时拦截检查登录状态【未完】
	if req.RequestURI != "/login" {
		if _, err := jwtUtils.ParseToken(token); err != nil {
			c.Data["json"] = resp.NewResp(respUtils.ERROR_CODE, "令牌解析失败")
			c.ServeJSON()
		}
	}
}
