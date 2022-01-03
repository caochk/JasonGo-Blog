package controllers

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"my_blog/models"
	"my_blog/utils"
	"regexp"
	"strings"
)

type LoginController struct {
	utils.BaseController
	beego.Controller
}

// Login 登录函数
func (c *LoginController) Login() {
	username := strings.Trim(c.GetString("username"), " ")
	password := strings.Trim(c.GetString("password"), " ")
	vcode := strings.ToLower(strings.Trim(c.GetString("vcode"), " "))

	var user models.User
	// 先校验图形验证码是否正确
	if vcode != c.GetSession("vcode") {
		fmt.Println("vcode-error")
	} else {
		password = utils.Md5(password) // 用户输入密码加密，之后用于验证密码是否匹配
		if result, err := user.FindByUsername(username); err == nil {
			if len(result) == 1 && result[0].Password == password { // 密码验证通过，登录成功
				// 在session中保存当前登录用户的一系列信息
				if err := c.SetSession("islogin", true); err != nil {
					fmt.Println("[ERROR] set session of islogin:", err)
				}
				if err := c.SetSession("userid", result[0].Id); err != nil {
					fmt.Println("[ERROR] set session of userid:", err)
				}
				if err := c.SetSession("username", result[0].Username); err != nil {
					fmt.Println("[ERROR] set session of userName:", err)
				}
				if err := c.SetSession("nickname", result[0].Nickname); err != nil {
					fmt.Println("[ERROR] set session of nickname:", err)
				}
				if err := c.SetSession("role", result[0].Role); err != nil {
					fmt.Println("[ERROR] set session of role:", err)
				}
				// 积分更新模块，每次登录可获得1积分【未完全】
				user.UpdateCredit(1)
				// 设置cookie实现一个月自动登录【待补】
			} else {
				fmt.Println("login-fail")
			}
		} else {
			fmt.Println("[ERROR] find by username:", err)
		}
	}
}

// Signup 注册函数【测试通过】
func (c *LoginController) Signup() string {
	username := strings.Trim(c.GetString("username"), " ")
	password := strings.Trim(c.GetString("password"), " ")
	ecode := strings.ToLower(strings.Trim(c.GetString("ecode"), " "))

	var user models.User
	// 验证邮箱验证码
	//fmt.Println("从session中取出ecode:", c.GetSession("ecode"))
	if ecode != c.GetSession("ecode") {
		fmt.Println("ecode-error")
		return "ecode-error"
	} else if matched, _ := regexp.MatchString(".+@.+\\..+", username); !matched || len(password) < 5 {
		fmt.Println("invalid username or password") // html文件处要改，那里应该写的是up-invalid
		return "up-invalid"
	} else if result, _ := user.FindByUsername(username); len(result) > 0 { // 检测用户邮箱是否已注册
		fmt.Println("user-repeated")
		//fmt.Println("find by username:", result[0].Id)
		return "user-repeated"
	} else {
		password = utils.Md5(password)
		if id_of_inserted_value, err := user.Signup(username, password); err == nil {
			if result, err := user.FindById(id_of_inserted_value); err == nil {
				// 在session中保存当前登录用户的一系列信息
				if err := c.SetSession("islogin", true); err != nil {
					fmt.Println("[ERROR] set session of islogin:", err)
					return "reg-fail"
				}
				if err := c.SetSession("userid", result.Id); err != nil {
					fmt.Println("[ERROR] set session of userid:", err)
					return "reg-fail"
				}
				if err := c.SetSession("username", result.Username); err != nil {
					fmt.Println("[ERROR] set session of userName:", err)
					return "reg-fail"
				}
				if err := c.SetSession("nickname", result.Nickname); err != nil {
					fmt.Println("[ERROR] set session of nickname:", err)
					return "reg-fail"
				}
				if err := c.SetSession("role", result.Role); err != nil {
					fmt.Println("[ERROR] set session of role:", err)
					return "reg-fail"
				}
				// 更新积分记录表
				var credit models.Credit
				err := credit.AddCreditDetail("用户注册", 0, 50, c.GetSession("userid").(int))
				if err != nil {
					fmt.Println("[ERROR] add credit detail:", err)
				}
				return "reg-pass"
			} else {
				fmt.Println("[ERROR] find by id:", err)
				return "reg-fail"
			}
		} else {
			fmt.Println("[ERROR] sign up:", err)
			return "reg-fail"
		}
	}
}

// Ecode 邮箱验证码【测试通过】
func (c *LoginController) Ecode() string {
	to := c.GetString("email")
	if matched, _ := regexp.MatchString(".+@.+\\..+", to); !matched {
		fmt.Println("[ERROR] email-invalid")
		return "email-invalid"
	}
	// 邮箱地址合法性校验通过
	var ec utils.Email
	ecode := ec.GenEcode() // 【通过】
	fmt.Println("ecode:", ecode)
	if err := ec.SendEmail(to, ecode); err == nil {
		err := c.SetSession("ecode", strings.ToLower(ecode))
		if err != nil {
			fmt.Println("[ERROR] set session:", err)
			return "send-fail"
		} else {
			fmt.Println("send-pass")
			return "send-pass"
		}
	} else {
		fmt.Println("[ERROR] send-fail:", err)
		return "send-fail"
	}
}
