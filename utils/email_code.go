package utils

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-gomail/gomail"
	"math/rand"
	"time"
)

type Email struct {
	beego.Controller
	From string
	To string
}

// SendEmail 发送邮件【测试通过】
func (c *Email) SendEmail(to string, ecode string) error {
	var m *gomail.Message
	m = gomail.NewMessage()
	// 发件人
	m.SetAddressHeader("From", "cookiemallgroup@outlook.com", "JasonGo")
	// 收件人
	m.SetAddressHeader("To", to, "")
	// 主题
	subject := "[JasonGo] Please verify your email address!"
	m.SetHeader("Subject", subject)
	// 正文
	body := "<br/>Welcome to JasonGo Blog:<br/>To complete the sign up, please enter the verification code:" +
		"<span style='color: red; font-size: 20px;'>" + ecode + "</span>, " +
		"Please do not reply to this notification, this inbox is not monitored.<br>" +
		"<span style='color: dodgerblue; font-size: 18px;'>Thanks for using JasonGo Blog!</span><br/>"
	m.SetBody("text/html", body)
	// 发送
	password, _ := beego.AppConfig.String("EmailPassword")
	d := gomail.NewDialer("smtp.office365.com", 587, "cookiemallgroup@outlook.com", password)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("[ERROR] send email:", err)
		return err
	} else {
		return nil
	}
}

// GenEcode 生成邮件验证码【测试通过】
func (c *Email) GenEcode() string {
	ecode := c.RandCode(6)
	return string(ecode)
}

func init() {
	rand.Seed(time.Now().Unix())
}

// RandCode 生成随机字符串，包含 0~9\a~z\A~Z【测试通过】
func (c *Email)RandCode(n int) []byte {
	var letters = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	if n <= 0 {
		return []byte{}
	}
	b := make([]byte, n)
	arc := uint8(0)
	if _, err := rand.Read(b[:]); err != nil {
		return []byte{}
	}
	for i, x := range b {
		arc = x & 61
		b[i] = letters[arc]
	}
	return b
}
