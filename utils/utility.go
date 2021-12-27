package utils

import (
	"crypto/md5"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
}

// AlertAndRedirect 这个函数应该就是提示之用
func (p *BaseController) AlertAndRedirect(msg string, url string) {
	if url == "" {
		p.Ctx.WriteString("<script>alert('" + msg + "');window.history.go(-1);</script>")
		p.StopRun()
	} else {
		p.Redirect(url, 302)
	}
}

// Md5 生成MD5密文
func Md5(plaintext string) string {
	plain := []byte(plaintext)
	return fmt.Sprintf("%X", md5.Sum(plain))
}