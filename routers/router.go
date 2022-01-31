package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"my_blog/controllers"
)

func init() {
	beego.Router("/", &controllers.HomeController{}, "get:Home")
	beego.Router("/page/:page([0-9]+)", &controllers.HomeController{}, "get:Paginate")
	beego.Router("/category/:category_page", &controllers.HomeController{}, "get:Classify")
	beego.Router("/search", &controllers.HomeController{}, "get:Search")
	beego.Router("/login", &controllers.LoginController{}, "post:Login")
	beego.Router("/ecode", &controllers.LoginController{}, "post:EcodeRedis")
	beego.Router("/signup", &controllers.LoginController{}, "post:SignupRedis")
	beego.Router("/article", &controllers.ArticleController{}, "get:GetArticleDetail")
	beego.Router("/favorite", &controllers.ArticleController{}, "post:AddToFavorite")
	beego.Router("/favorite", &controllers.ArticleController{}, "delete:CancelFavorite")
}
