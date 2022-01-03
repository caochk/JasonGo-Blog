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

// PrevPage 前往上一页，用于模板
func PrevPage(page int) int {
	page--
	return page
}

// NextPage 前往下一页，用于模板
func NextPage(page int) int {
	page++
	return page
}

// DealTotalPage 将total_page_num为数组类型
func DealTotalPage(total_page_num float64) []float64 {
	var total_page_slice = make([]float64, int(total_page_num))
	var count = float64(1)
	for true {
		if count > total_page_num {
			break
		}
		total_page_slice = append(total_page_slice, count)
	}
	return total_page_slice
}
