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
- [ ] 接入Redis缓存   
- [ ] MySQL数据定时同步至Redis   
- [ ] 引入日志系统   
- [ ] 文章发布功能   
- [ ] 评论功能  
- [ ] 后台管理功能