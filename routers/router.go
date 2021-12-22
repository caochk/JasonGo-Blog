package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"my_blog/controllers"
)

func init() {
    beego.Router("/", &controllers.HomeController{}, "get:Home")
}
