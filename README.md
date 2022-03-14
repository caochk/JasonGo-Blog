# **JasonGo-Blog**

## 技术栈
框架：beego  
orm：beego自带（_未来打算换成gorm，因为个人感觉beego自带的orm操作不是太好用_）  
数据库：MySQL  
缓存中间件：Redis  
日志系统（暂未引入）

## 项目结构
-JasonGo-Blog  
&emsp;&emsp;|-conf 配置文件目录   
&emsp;&emsp;|-controllers 控制器目录     
&emsp;&emsp;|-models 数据库访问目录     
&emsp;&emsp;|-cache 操作Redis缓存目录  
&emsp;&emsp;|-utils 公共方法目录  
&emsp;&emsp;|-static 静态资源目录  
&emsp;&emsp;|-views 模板文件目录  
&emsp;&emsp;|-main.go 程序执行入口

## TODO
- [x] ~~首页文章列表展示功能~~
- [x] ~~分页浏览功能~~
- [X] ~~分类浏览功能~~
- [X] ~~搜索功能~~（较为初级，可深入优化）
- [X] ~~注册功能（含邮箱验证码的实现）~~
- [X] ~~登录功能~~（图片验证码尚未实现）~~（利用session自动登录已实现）~~
- [X] ~~文章全文展示功能~~
- [ ] 接入Redis缓存   
    - 【读】先读Redis，未查到则读MySQL，读到后往Redis存一份
    - 【写】先写MySQL，后定期更新至Redis
- [ ] 利用Redis实现自动登录
- [ ] 用户认证方式由session改为JWT
    - 同时引入Redis解决用户手动退出登录时token尚未过期问题
- [ ] 引入日志系统   
- [ ] 文章发布功能   
- [ ] 评论功能  
- [ ] 后台管理功能