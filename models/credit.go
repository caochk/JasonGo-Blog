package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"
)

type Credit struct {
	Id int
	//UserId   int
	Category string
	Target int
	Credit int
	CreateTime time.Time
	UpdateTime time.Time
	User *User `orm:"rel(fk)"`
}

// AddCreditDetail 积分记录表插入记录【测试通过】
func (m *Credit) AddCreditDetail(category string, target int, credit int, user_id int) error {
	//fmt.Println(user_id)
	user := &User{Id: user_id}
	now := time.Now()
	credit_m := Credit{
		User : user,
		Category: category,
		Target: target,
		Credit: credit,
		CreateTime: now,
		UpdateTime: now,
	}
	o := orm.NewOrm()
	if _, err := o.Insert(&credit_m); err == nil {
		return err
	} else {
		return err
	}
}