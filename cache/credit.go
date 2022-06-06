package cache

import (
	"errors"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/google/uuid"
)

type CreditCacheController struct {
	beego.Controller
}

// RedPackage2Redis 按管理员设置生成一个红包【测试通过】
func (c CreditCacheController) RedPackage2Redis(redPackageNums []int) {
	key := "RED_PACKAGE_KEY:" + uuid.New().String()
	for _, redPackageNum := range redPackageNums {
		fmt.Println(redPackageNum)
		rdb.LPush(ctx, key, redPackageNum)
	}
}

// HaveGottenRedPackage 判断某用户是否已抢过该红包【测试通过】
func (c CreditCacheController) HaveGottenRedPackage(redPackageKey string, userId string) bool {
	key := "RED_PACKAGE_RECORD:" + redPackageKey
	if err := rdb.HGet(ctx, key, userId).Err(); err != nil {
		return false // key不存在即这是抢该红包的第一人或这人没抢过该红包
	} else { // 说明该人在哈希类型中有记录，即已抢过红包
		return true
	}
	//fmt.Println("从哈希类型返回的结果：", res, err)
	//return true
}

// GetRedPackage 抢一个红包【测试通过】
func (c CreditCacheController) GetRedPackage(redPackageKey string) (string, error) {
	key := "RED_PACKAGE_KEY:" + redPackageKey
	if res, err := rdb.LPop(ctx, key).Result(); err == nil {
		fmt.Println("抢到一个红包：", res)
		return res, nil
	} else {
		return "", errors.New("empty")
	}

}

// AddRedPackageRecord 记录某用户抢到一个红包【测试通过】
func (c CreditCacheController) AddRedPackageRecord(redPackageKey string, userId string, val string) {
	key := "RED_PACKAGE_RECORD:" + redPackageKey
	rdb.HSet(ctx, key, userId, val)
}
