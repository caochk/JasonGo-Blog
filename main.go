package main

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
	. "my_blog/models"
	_ "my_blog/routers"
)

func init()  {
	var mysqlurls, _ = beego.AppConfig.String("mysqlurls")
	var mysqlport, _ = beego.AppConfig.String("mysqlport")
	var mysqluser, _ = beego.AppConfig.String("mysqluser")
	var mysqlpass, _ = beego.AppConfig.String("mysqlpass")
	var mysqldb, _  = beego.AppConfig.String("mysqldb")
	var dsn = mysqluser + ":" + mysqlpass + "@tcp(" + mysqlurls + ":" + mysqlport + ")/" + mysqldb + "?charset=utf8&loc=Asia%2FShanghai"
	err := orm.RegisterDataBase("default", "mysql", dsn)
	if err != nil {
		fmt.Println(err)
	}
	orm.RegisterModel(new(Article))
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(Credit))

	orm.RunSyncdb("default", false, true)
}

func main() {
	beego.Run()
}

