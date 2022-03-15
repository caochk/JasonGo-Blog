package main

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/beego/beego/v2/server/web/session/redis"
	_ "github.com/go-sql-driver/mysql"
	. "my_blog/models"
	_ "my_blog/routers"
	"my_blog/utils"
)

func init() {
	var mysqlurls, _ = beego.AppConfig.String("mysqlurls")
	var mysqlport, _ = beego.AppConfig.String("mysqlport")
	var mysqluser, _ = beego.AppConfig.String("mysqluser")
	var mysqlpass, _ = beego.AppConfig.String("mysqlpass")
	var mysqldb, _ = beego.AppConfig.String("mysqldb")
	var dsn = mysqluser + ":" + mysqlpass + "@tcp(" + mysqlurls + ":" + mysqlport + ")/" + mysqldb + "?charset=utf8&loc=Asia%2FShanghai"
	err := orm.RegisterDataBase("default", "mysql", dsn)
	if err != nil {
		fmt.Println(err)
	}
	orm.RegisterModel(new(Article))
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(Credit))
	orm.RegisterModel(new(Favorite))

	if err := orm.RunSyncdb("default", false, true); err != nil {
		fmt.Println("[ERROR] sync db:", err)
	}
}

func main() {
	// 注册模板函数：下一页
	if err := beego.AddFuncMap("NextPage", utils.NextPage); err != nil {
		fmt.Println("[ERROR] next page:", err)
	}
	// 注册模板函数：上一页
	if err := beego.AddFuncMap("PrevPage", utils.PrevPage); err != nil {
		fmt.Println("[ERROR] prev page:", err)
	}
	err := beego.AddFuncMap("DealTotalPage", utils.DealTotalPage)
	if err != nil {
		return
	}

	beego.Run()
}
