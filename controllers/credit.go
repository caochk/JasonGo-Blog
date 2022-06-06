package controllers

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"math/rand"
	creditCache "my_blog/cache"
	"my_blog/utils/respUtils"
	"strconv"
)

var creditCacheController creditCache.CreditCacheController

type CreditController struct {
	beego.Controller
}

// SendRedPackage 发积分红包【测试通过】
func (c CreditController) SendRedPackage() {
	totalCredit, _ := strconv.Atoi(c.GetString("totalCredit"))
	redPackageNum, _ := strconv.Atoi(c.GetString("redPackageNum"))
	// 将管理员输入的总积分平均分拆到设置好的个数的红包中
	redPackageNums := c.splitRedPackage(totalCredit, redPackageNum)
	// 以list方式存红包入Redis
	creditCacheController.RedPackage2Redis(redPackageNums)
	c.Data["json"] = resp.NewResp(respUtils.SUCCESS_CODE, "红包准备就绪")
	c.ServeJSON()
}

// GetRedPackage 抢积分红包并记录【测试通过】
func (c CreditController) GetRedPackage() {
	redPackageKey := c.GetString("redPackageKey")
	userId := c.GetString("userId")
	// 未抢过红包
	if !creditCacheController.HaveGottenRedPackage(redPackageKey, userId) {
		// 抢一个红包
		if redPackageVal, err := creditCacheController.GetRedPackage(redPackageKey); err == nil {
			// 记录该用户抢了一个红包
			creditCacheController.AddRedPackageRecord(redPackageKey, userId, redPackageVal)
			c.Data["json"] = resp.NewResp(respUtils.SUCCESS_CODE, "用户"+userId+"抢到"+redPackageVal)
			c.ServeJSON()
		} else {
			// 该红包已被抢空
			c.Data["json"] = resp.NewResp(respUtils.ERROR_CODE, "来晚了，红包被抢光了！")
			c.ServeJSON()
		}
	} else {
		c.Data["json"] = resp.NewResp(respUtils.ERROR_CODE, "用户"+userId+"抢过红包了！")
		c.ServeJSON()
	}
}

// SplitRedPackage 均值拆分积分红包：二倍均值法【测试通过】
func (c CreditController) splitRedPackage(totalCredit int, redPackageNum int) []int {
	// 已经被抢掉的积分
	var usedCredit = 0
	// 每次抢到的积分
	var redPackageNums = make([]int, 5)
	for i := 0; i < redPackageNum; i++ {
		if i == redPackageNum-1 {
			redPackageNums[i] = totalCredit - usedCredit
		} else {
			avgCredit := ((totalCredit - usedCredit) / (redPackageNum - i)) * 2
			redPackageNums[i] = 1 + rand.Intn(avgCredit-1)
		}
		usedCredit = usedCredit + redPackageNums[i]
	}
	fmt.Println(redPackageNums)
	return redPackageNums
}
