package controllers

import (
	//"context"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"my_blog/models"
	"my_blog/utils"
	"my_blog/utils/jwtUtils"
	"my_blog/utils/respUtils"
	"regexp"
	"strings"
	"time"
)

//var ctx = context.Background()
//var rdb = utils.InitRedisClient()

type LoginController struct {
	utils.BaseController
	beego.Controller
}

type RespondMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
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
		var resp RespondMsg
		resp.Code = -1
		resp.Msg = "[ERROR] email-invalid"
		c.Data["json"] = &resp
		c.ServeJSON()
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

// EcodeRedis 生成验证码及向邮箱发送验证码之Redis版
func (c *LoginController) EcodeRedis() {
	//var resp RespondMsg
	var resp respUtils.Resp
	to := c.GetString("email")
	if matched, _ := regexp.MatchString(".+@.+\\..+", to); !matched {
		//fmt.Println("[ERROR] email-invalid")
		//resp.Code = -1
		//resp.Msg = "[ERROR] email-invalid"
		//c.Data["json"] = &resp
		//c.ServeJSON()
		//return "email-invalid"
		c.Data["json"] = resp.NewResp(respUtils.ERROR_CODE, "[ERROR] email-invalid")
		c.ServeJSON()
	}
	// 邮箱地址合法性校验通过
	var ec utils.Email
	ecode := ec.GenEcode() // 【通过】
	fmt.Println("ecode:", ecode)
	if err := ec.SendEmail(to, ecode); err == nil {
		//err := c.SetSession("ecode", strings.ToLower(ecode))
		err := rdb.Set(ctx, to, strings.ToLower(ecode), 300*time.Second).Err() // 我规定了邮箱验证码必须在5分钟即300s内输入，不然就会过期【已通过测试】
		if err != nil {
			//fmt.Println("[ERROR] set session:", err)
			//resp.Code = -1
			//resp.Msg = "[ERROR] set session:" + err.Error()
			//c.Data["json"] = &resp
			//c.ServeJSON()
			//return "send-fail"
			c.Data["json"] = resp.NewResp(respUtils.ERROR_CODE, "[ERROR] set session at redis:"+err.Error())
			c.ServeJSON()
		} else {
			//fmt.Println("send-pass")
			//resp.Code = 200
			//resp.Msg = "send-pass"
			//c.Data["json"] = &resp
			//c.ServeJSON()
			//return "send-pass"
			c.Data["json"] = resp.NewResp(respUtils.SUCCESS_CODE, "send-pass")
			c.ServeJSON()
		}
	} else {
		//fmt.Println("[ERROR] send-fail:", err)
		//resp.Code = -1
		//resp.Msg = "[ERROR] send-fail:" + err.Error()
		//c.Data["json"] = &resp
		//c.ServeJSON()
		//return "send-fail"
		c.Data["json"] = resp.NewResp(respUtils.ERROR_CODE, "[ERROR] send-fail:"+err.Error())
		c.ServeJSON()
	}
}

// SignupRedis 注册函数的Redis版本
func (c *LoginController) SignupRedis() string {
	username := strings.Trim(c.GetString("username"), " ")
	password := strings.Trim(c.GetString("password"), " ")
	ecode := strings.ToLower(strings.Trim(c.GetString("ecode"), " "))

	var user models.User
	// 验证邮箱验证码
	//fmt.Println("从session中取出ecode:", c.GetSession("ecode"))
	if ecode_from_redis, err := rdb.Get(ctx, username).Result(); err == redis.Nil { // redis中不存在该key
		fmt.Println("username does not exist")
		return "username does not exist"
	} else if err != nil { // redis报告其他错误
		fmt.Println("redis error")
		return "redis error"
	} else if ecode != ecode_from_redis { // redis中ecode与用户输入的不匹配
		fmt.Println("ecode:", ecode)
		fmt.Println("ecode from redis:", ecode_from_redis)
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
				fmt.Println("reg-pass")
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

// 引入JWT、Redis后的注册函数
func (c *LoginController) SignupJWT() {
	username := strings.Trim(c.GetString("username"), " ")
	password := strings.Trim(c.GetString("password"), " ")
	ecode := strings.ToLower(strings.Trim(c.GetString("ecode"), " "))

	var user models.User
	var resp respUtils.Resp
	// 验证邮箱验证码
	//fmt.Println("从session中取出ecode:", c.GetSession("ecode"))
	if ecode_from_redis, err := rdb.Get(ctx, username).Result(); err == redis.Nil { // redis中不存在该key
		//fmt.Println("username does not exist")
		//jsonData := make(map[string]interface{}, 2)
		//jsonData["code"] = -1
		//jsonData["msg"] = "username does not exist"
		//c.Data["json"] = &jsonData
		//respUtils := respUtils.Resp{
		//	Code: -1,
		//	Message: "username does not exist",
		//}

		c.Data["json"] = resp.NewResp(respUtils.ERROR_CODE, "username does not exist")
		c.ServeJSON()
		//return "username does not exist"
	} else if err != nil { // redis报告其他错误
		//fmt.Println("redis error")
		//return "redis error"
		c.Data["json"] = resp.NewResp(respUtils.ERROR_CODE, "redis error")
		c.ServeJSON()
	} else if ecode != ecode_from_redis { // redis中ecode与用户输入的不匹配
		//fmt.Println("ecode:", ecode)
		//fmt.Println("ecode from redis:", ecode_from_redis)
		//fmt.Println("ecode-error")
		//return "ecode-error"
		c.Data["json"] = resp.NewResp(respUtils.ERROR_CODE, "ecode-error")
		c.ServeJSON()
	} else if matched, _ := regexp.MatchString(".+@.+\\..+", username); !matched || len(password) < 5 {
		//fmt.Println("invalid username or password") // html文件处要改，那里应该写的是up-invalid
		//return "up-invalid"
		c.Data["json"] = resp.NewResp(respUtils.ERROR_CODE, "up-invalid")
		c.ServeJSON()
	} else if result, _ := user.FindByUsername(username); len(result) > 0 { // 检测用户邮箱是否已注册
		//fmt.Println("user-repeated")
		//fmt.Println("find by username:", result[0].Id)
		//return "user-repeated"
		c.Data["json"] = resp.NewResp(respUtils.ERROR_CODE, "user-repeated")
		c.ServeJSON()
	} else {
		password = utils.Md5(password)
		if id_of_inserted_value, err := user.Signup(username, password); err == nil {
			if result, err := user.FindById(id_of_inserted_value); err == nil {
				// 在session中保存当前登录用户的一系列信息【由session改至JWT认证】
				//if err := c.SetSession("islogin", true); err != nil {
				//	fmt.Println("[ERROR] set session of islogin:", err)
				//	return "reg-fail"
				//}
				//if err := c.SetSession("userid", result.Id); err != nil {
				//	fmt.Println("[ERROR] set session of userid:", err)
				//	return "reg-fail"
				//}
				//if err := c.SetSession("username", result.Username); err != nil {
				//	fmt.Println("[ERROR] set session of userName:", err)
				//	return "reg-fail"
				//}
				//if err := c.SetSession("nickname", result.Nickname); err != nil {
				//	fmt.Println("[ERROR] set session of nickname:", err)
				//	return "reg-fail"
				//}
				//if err := c.SetSession("role", result.Role); err != nil {
				//	fmt.Println("[ERROR] set session of role:", err)
				//	return "reg-fail"
				//}

				// 生成令牌
				claims := make(jwt.MapClaims)
				claims["islogin"] = true
				claims["userid"] = result.Id
				claims["username"] = result.Username
				claims["nickname"] = result.Nickname
				claims["role"] = result.Role
				claims["exp"] = time.Now().Add(time.Hour * 2).Unix() // 设置该token在当前时间的两小时后过期，届时需重新输密码
				fmt.Println("声明：", claims)
				signedToken := jwtUtils.BuildToken(claims)
				fmt.Println("签名后的令牌：", signedToken)

				// 更新积分记录表
				var credit models.Credit
				err := credit.AddCreditDetail("用户注册", 0, 50, result.Id)
				if err != nil {
					//fmt.Println("[ERROR] add credit detail:", err)
					c.Data["json"] = resp.NewResp(respUtils.ERROR_CODE, "[ERROR] add credit detail:"+err.Error())
					c.ServeJSON()
				}
				//fmt.Println("reg-pass")
				//return "reg-pass"
				//c.Data["json"] = resp.NewResp(respUtils.SUCCESS_CODE, )
				//c.ServeJSON()
				respond := resp.NewRespWithData(respUtils.SUCCESS_CODE, "reg-pass", signedToken)
				fmt.Println("转换为字节流的token：", respond.ToBytes())
				c.Ctx.Output.Body(respond.ToBytes())
			} else {
				//fmt.Println("[ERROR] find by id:", err)
				//return "reg-fail"
				c.Data["json"] = resp.NewResp(respUtils.ERROR_CODE, "[ERROR] find by id:"+err.Error())
				c.ServeJSON()
			}
		} else {
			//fmt.Println("[ERROR] sign up:", err)
			//return "reg-fail"
			c.Data["json"] = resp.NewResp(respUtils.ERROR_CODE, "[ERROR] sign up:"+err.Error())
			c.ServeJSON()
		}
	}
}
