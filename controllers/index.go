package controllers

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"my_blog/models"
)

type HomeController struct {
	beego.Controller
}

func (c *HomeController) Home() {
	fmt.Println("111")
	if c.Ctx.Request.Method == "GET" {
		fmt.Println("222")
		// 展示首页文章
		article := models.Article{}
		result, err := article.Find_paginated_articles(10, 1)
		//result, err := article.Find_all()
		//fmt.Println(err)
		if err == nil {
			c.Data["result"] = result
			for _, v := range result{
				fmt.Println(v.Id)
			}
		}
		//c.TplName = "base.html"
	}

}