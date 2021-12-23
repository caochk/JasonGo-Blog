package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"my_blog/controllers"
)

func init() {
    beego.Router("/", &controllers.HomeController{}, "get:Home")
	beego.Router("/page/:page([0-9]+)", &controllers.HomeController{}, "get:Paginate")
	beego.Router("/category/:category_page", &controllers.HomeController{}, "get:Classify")
}
