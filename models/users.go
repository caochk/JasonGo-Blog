package models

import "time"

type Users struct {
	Id int `orm:"pk"`
	Username string
	Password string
	Nickname string
	Avatar string
	Qq string
	Role string
	Credit int
	CreateTime time.Time
	UpdateTime time.Time
}