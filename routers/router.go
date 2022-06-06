package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"my_blog/cache"
	"my_blog/controllers"
)

func init() {
	beego.Router("/", &controllers.HomeController{}, "get:HomeRedis")
	beego.Router("/page/:page([0-9]+)", &controllers.HomeController{}, "get:Paginate")
	beego.Router("/category/:category_page", &controllers.HomeController{}, "get:Classify")
	beego.Router("/search", &controllers.HomeController{}, "get:Search")
	beego.Router("/login", &controllers.LoginController{}, "post:LoginJWT")
	beego.Router("/ecode", &controllers.LoginController{}, "post:EcodeRedis")
	beego.Router("/signup", &controllers.LoginController{}, "post:SignupJWT")
	beego.Router("/article", &controllers.ArticleController{}, "get:GetArticleDetail")
	beego.Router("/favorite", &controllers.ArticleController{}, "post:AddToFavorite")
	beego.Router("/favorite", &controllers.ArticleController{}, "delete:CancelFavorite")
	beego.Router("/redis2home", &controllers.HomeController{}, "get:HomeRedis")
	beego.Router("/redis", &cache.ArticleCacheController{}, "get:Articles2Redis")
	beego.Router("/addcomment", &controllers.CommentController{}, "post:AddComment")
	beego.Router("/sendredpackage", &controllers.CreditController{}, "post:SendRedPackage")
	beego.Router("/getredpackage", &controllers.CreditController{}, "post:GetRedPackage")
	//beego.Router("/close", &controllers.HomeController{}, "post:Close")

	// 过滤器
	//beego.InsertFilter("/*", beego.BeforeRouter, Filter) // 寻找路由之前进行过滤拦截
}
