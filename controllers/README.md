
## Controller 开发过程 :
Controller 以BaseController  为基础

并且声明一个 以 `_没有Controller的类名_ctl` 格式的全局变量

#### init() 
中写入 route映射关系
并且运行 _xxx_ctl.AddRoutes()
这个函数会在被import时就运行


#### AddRoutes 
每个controller有一个AddRoutes 函数,内空都一样,但是必须写

#### GuuPrepare 函数
用于构建Form
统一写上 
> 	this.TplName = libs.GetTplName(this)
用Controller的名子来自动指定tpl的名子
规则就是
XxxController => xxx_controller.tpl


剩下就是处理请求的函数区

### 表单Get请求 处理流程
1. this.Data["Form"] = this.MyForm.Render(),让Form显示出来


### 表单Post 处理流程
1. 使用 this.Validator() 做验证
2. 取得数据,与数据库交互

