package models

import (
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"math/rand"
	"my_blog/utils"
	"strconv"
	"strings"
	"time"
)

type User struct {
	//beego.Controller  // 写了会报错：panic: reflect: call of reflect.Value.Interface on zero Value
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
var c beego.Controller

// FindByUsername 查询用户名（用于注册时判断用户名是否已注册、还用于登录校验）【测试通过】
func (u *User) FindByUsername(username string) ([]*User, error) {
	o := orm.NewOrm()
	var users []*User
	_, err := o.QueryTable("user").Filter("username", username).All(&users)  // 用.One()时若未找到任何记录，会返回空
	if err == nil {
		return users, err
	} else {  // 测试时要注意一点，就是返回的结果会有0条这种情形吗，当没有一个用户名相匹配时。要是实在不行，就改回用切片配all
		return users, err
	}
}

// UpdateCredit 增减用户积分（阅读收费文章会扣除相应积分）
func (m *User) UpdateCredit(credit int) {
	o := orm.NewOrm()
	// 先读出来
	user_id_from_session, _ := strconv.Atoi(c.GetSession("userid").(string))  // 【？】.(string)是什么意思，可以直接写.(int)吗
	user := User{Id: user_id_from_session}
	if err := o.Read(&user, "id"); err == nil {
		// 再更新
		user.Credit += credit
		if _, err := o.Update(&user); err == nil {
			fmt.Println("积分更新成功")
		} else {
			fmt.Println("[ERROR] orm update:", err)
		}
	} else {
		fmt.Println("[ERROR] orm read:", err)
	}
}

// 图片验证码
func (m *User) Vcode() {
	var imageCode utils.ImageCode
	imageCode.GetCode()
}

// Signup 注册工作【测试通过】
func (m *User) Signup(username string, password string) (int, error) {
	now := time.Now()
	nickname := strings.Split(username, "@")[0]
	avatar := "(" + strconv.Itoa(rand.Intn(14)) + ")" + ".svg"
	user := User{
		Username: username,
		Password: password,
		Role: "user",
		Credit: 50,
		Nickname: nickname,
		Avatar: avatar,
		CreateTime: now,
		UpdateTime: now,
	}
	o := orm.NewOrm()
	if id_of_inserted_value, err := o.Insert(&user); err == nil {
		return int(id_of_inserted_value), err
	} else {
		fmt.Println("[ERROR] orm insert:", err)
		return -1, err
	}
}

// FindById 根据id查找用户
func (m *User) FindById(id int) (User, error) {
	user := User{
		Id: id,
	}
	o := orm.NewOrm()
	if err := o.Read(&user); err == nil {
		return user, err
	} else {
		return user, err
	}
}

// FindAllUsers 查询user表中的所有数据
func (User) FindAllUsers() ([]*User, error) {
	o := orm.NewOrm()
	var users []*User
	_, err := o.QueryTable("user").All(&users)
	if err == nil {
		return users, err
	}
	return nil, err
}