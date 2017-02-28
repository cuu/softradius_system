package controllers

import (
	r "github.com/cuu/softradius/routers"
	//"github.com/cuu/softradius/models"
	"github.com/cuu/softradius/libs"
	"fmt"
	//	"reflect"
	"github.com/astaxie/beego"

//	sort "github.com/cuu/softradius/libs/sortutil"
	
)

type AgencyController struct {
	BaseController
}



type AgencyShare struct {
	Id string `gorethink:"id,omitempty"`
	OrderId  string
	AgencyId string
	ShareRate int
	ShareFee  int
	CreateTime string
}

type AgencyOrder struct{
	Id string `gorethink:"id,omitempty"`
	AgencyId string
	MemberOrderId string
	FeeType  string
	FeeValue int
	FeeTotal int
	FeeDesc  string
	CreateTime string
}


type Agency struct {
	Id string `gorethink:"id,omitempty"`
	Name string
	Operator string
	Contact string
	Mobile string
	Amount  int
	ShareRate int
	Desc string
	CreateTime string
	UpdateTime string
	Orders []AgencyOrder
	Shares []AgencyShare
}



//每个controller 有这样的命名规则,保证不重复
var _agency_ctl AgencyController

func init(){
	//在这儿独立处理映身关系
	// 事实上只有 方法是 get的 Route才应该是 Is_menu显示出来
	_ctl := &_agency_ctl
	_cate := r.MenuAgency
	
 	_ctl.routes = append( _ctl.routes,
r.Route{Path:"/agency",Name:"代理列表",Category:_cate,Is_menu :true, Order:1.0,Is_open:true, Methods:"*:Agency"})

 	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/agency/open",Name:"代理开户",Category:_cate,Is_menu :true, Order:1.1,Is_open:true, Methods:"*:AgencyOpen"})

 	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/agency/update",Name:"代理商修改",Category:_cate,Is_menu :false, Order:1.2,Is_open:true, Methods:"*:AgencyUpdate"})

 	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/agency/recharge",Name:"代理商充值",Category:_cate,Is_menu :false, Order:1.3,Is_open:true, Methods:"*:AgencyRecharge"})

 	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/agency/delete",Name:"代理商删除",Category:_cate,Is_menu :false, Order:1.4,Is_open:true, Methods:"*:AgencyDelete"})
	
	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/agency/order", Name:"代理交易查询",Category:_cate,Is_menu:true, Order:1.8,Is_open:true, Methods:"*:AgencyOrders"})

	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/agency/share", Name:"代理分成",Category:_cate,Is_menu:true, Order:3.0,Is_open:true, Methods:"*:AgencyShares"})
	
	
	_ctl.AddRoutes()
	
}

//把this中的routes 放到 routers.Permits 中,每个Controller写一遍
func (this *AgencyController) AddRoutes() {
	
	for i,v := range  this.routes {
		if v.Methods != "" {
			beego.Router(v.Path, this, v.Methods)
		}else{
			beego.Router(v.Path, this)
		}
		
		if v.Is_menu {
			//Permits.routes is map,so no confict path key ! 
			r.Permits.Add_route(v.Path,&this.routes[i])
		}
	}
	
}

//主要做Form的要关准备
func (this *AgencyController) GuuPrepare(){
	
	this.TplName = libs.GetTplName(this)
	fmt.Println("in agenies")
	
}

//------------------------------------------------------------

func (this *AgencyController) Agency() {

	
}

func (this *AgencyController) AgencyOpen() {
	
}

func (this *AgencyController) AgencyUpdate() {
	
}

func (this *AgencyController) AgencyRecharge() {
	
}

func (this *AgencyController) AgencyDelete() {
	
}

func (this *AgencyController) AgencyOrders() {
	
}

func (this *AgencyController) AgencyShares() {

}


