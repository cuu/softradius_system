package controllers

import (
	r "github.com/cuu/softradius/routers"
	"github.com/cuu/softradius/models"
	"github.com/cuu/softradius/libs"
	"fmt"
	//	"reflect"
	"github.com/astaxie/beego"
	rdb  "github.com/cuu/softradius/database/shelf"
//	sort "github.com/cuu/softradius/libs/sortutil"
	"strconv"
	"strings"
	
)

type AgencyController struct {
	BaseController
	Rules []string
}


const (
	AGENCY_OPERATOR =2
)

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
	OperatorName string
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
	this.Data["msg"] = ""
	fmt.Println("in agenies")
	
}

//------------------------------------------------------------

func (this *AgencyController) CheckInRules(path string ) string {
	for i:=0;i<len(this.Rules);i++ {
		if this.Rules[i] == path {
			return "checked"
		}
		
	}
	return ""
}


func (this *AgencyController) Agency() {
	this.TplName = "agency_list.html"

	this.Data["AgencyList"] = this.AgencyList()
	
	this.Render()
}

func (this *AgencyController) AgencyOpenForm(nodes [][]string ,pdus [][]string) *models.Form {
	f := models.InfoForm("Open Agency","/agency/open",
		models.TextBox(&models.Input{Name:"Name",Valid:models.Len_of(2,255),Description:"代理商名称",Required:true}),
		models.TextBox(&models.Input{Name:"Contact",Valid:models.Len_of(2,255),Description:"联系人",Required:true }),
		models.TextBox(&models.Input{Name:"Mobile",Valid:models.Is_telephone,Description:"手机号", }),
		models.TextBox(&models.Input{Name:"Amount",Valid:models.Is_number,Description:"初始余额(元)",Required:true,}),
		models.TextBox(&models.Input{Name:"ShareRate",Valid:models.Is_number,Description:"分成比例1-100",Required:true }),
		models.TextBox(&models.Input{Name:"OperatorName",Valid:models.Len_of(2,255),Description:"操作员帐号",Required:true }),
		models.TextBox(&models.Input{Name:"OperatorPass",Valid:models.Len_of(2,32),Description:"操作员密码",Required:true }),
		models.GroupDropdown(&models.Select{Name:"Nodes",Args:nodes,Required:true,Size:4}),
		models.GroupDropdown(&models.Select{Name:"Products",Args:pdus,Required:true,Size:6}),
		models.TextArea(&models.Input{Name:"Desc",Description:"代理商描述",Size:4}),
		models.Submit(&models.Input{Name:"Submit",Value:"<b>提交</b>",Class:"btn btn-info"}),
	)

	return f
}

//to limit the power of agency
func (this *AgencyController) AgencyOpenMenus() []r.Menu {
	menu := r.Permits.Menu
	ms := []string{r.MenuBus}
	var ret  []r.Menu
	
	for _,v := range ms {
		for _,u := range menu {
			if u.Category == v {
				ret = append(ret,u)
				break
			}
		}
	}

	return ret
}

func (this *AgencyController) AgencyOpen() {
	nodes := this.NodeList()
	pdus  := this.ProductList()

	allnodes    := this.Items(nodes,[]string{"Id","Name"})
	allproducts := this.Items(pdus,[]string{"Id","Name"})
	
	f := this.AgencyOpenForm(allnodes,allproducts)
	this.Data["AllMenus"]     = this.AgencyOpenMenus()
	this.Data["CheckInRules"] = this.CheckInRules
	this.Data["CheckOpen"]    = r.Permits.Check_open
	this.TplName = "agency_form.html"

	if this.Ctx.Input.IsPost() {
		if this.Validator2(f) == false {
			this.Data["Form"] = f
			this.Render()
			return
		}
		
		op := &Operators{}
		opname := this.GetString("OperatorName")
		cnt := rdb.DataBase().FilterCount(op, map[string]string{"Name":opname})
		if cnt > 0 {
			
			this.ShowTips("管理员有重复,请选择另一个管理员名称")
			this.Render()
			return
		}

		agc := &Agency{}
		this.ParsePostToStruct(agc)
		fmt.Println("agency: ",agc)
		agc.Amount = libs.Yuan2fen(agc.Amount)
		if agc.ShareRate > 100 {
			agc.ShareRate = 100
		}else if agc.ShareRate < 0 {
			agc.ShareRate = 0
		}
		
		_,cnt = rdb.DataBase().FilterInsert(agc,"Name")
		if cnt  == 0 {
			this.Redirect("/agency", 302)
		}else {
			
			this.ShowTips("代理名称有重复 "+strconv.Itoa(cnt) +"个" )
			this.Render()
			return
		}

		this.ParsePostToStruct(op)
		op.Name = opname
		op.Pass = this.GetString("OperatorPass")
		op.Type = AGENCY_OPERATOR
		op.Desc = "代理商"
		
		fmt.Println("op: ",op)
		_,cnt = rdb.DataBase().FilterInsert(op,"Name")
		if cnt  == 0 {
			this.Redirect("/agency", 302)
		}else {
			
			this.ShowTips("管理员名称有重复 "+strconv.Itoa(cnt) +"个" )
			this.Render()
			return
		}
		
		this.Render()
		return
	}
	
	this.Data["Form"] = f
	
	this.Render()
}

func (this *AgencyController) AgencyUpdate() {
	
}


func (this *AgencyController) AgencyRechargeForm(name string, agency_id string) *models.Form {
	f := models.InfoForm("Recharge","/agency/recharge",
		models.TextBox(&models.Input{Name:"Name",Value:name, Valid:models.Len_of(2,255),Description:"代理商名称",Required:true,ReadOnly:true }),
		models.TextBox(&models.Input{Name:"FeeValue",Valid:models.Is_number,Description:"充值余额(元),格式xx.xx", Required:true}),
		models.Submit(&models.Input{Name:"Submit",Value:"<b>提交</b>",Class:"btn btn-info"}),
		models.Hidden(&models.Input{Name:"Id",Value:agency_id,Description:"编号"}),
	)
	return f
}

func (this *AgencyController) AgencyRecharge() {
	var id string

	if this.Ctx.Input.IsGet() {
		id = this.GetString("agency_id")
		if strings.Trim(id, " ") == "" {
			this.Abort("403")
			return
		}
	}else if this.Ctx.Input.IsPost() {
		id = this.GetString("Id")
	}
	
	agc := &Agency{}
	err := rdb.DataBase().QuOne(agc,id)
	if err == nil {
		f:= this.AgencyRechargeForm(agc.Name,agc.Id)
		if this.Ctx.Input.IsPost() {
			if this.Validator2(f) == false {
				this.Data["Form"] = f.Render()
				this.Render()
				return
			}

			fee_value := this.GetStringI("FeeValue")
			agc.Amount += libs.Yuan2fen(fee_value)
			resp,err := rdb.DataBase().Update(id,agc)
			if err == nil {
				fmt.Println("Replaced ",resp.Replaced )
			}else {
				fmt.Println(err)
			}
			
			this.Redirect("/agency",302)
			return
			
		}
		
		this.Data["Form"] = f.Render()
		this.Render()
		
	}else {
		this.Abort("403")
	}
	
}

func (this *AgencyController) AgencyDelete() {
	id := this.GetString("agency_id")
	
	if id == "" || strings.Trim(id," ")  == "" {
		this.Abort("403")
		return
	}

	agc := &Agency{}
	err := rdb.DataBase().QuOne(agc,id)
	if err == nil {
		_,err = rdb.DataBase().FilterDel("operators",map[string]string{"Name":agc.OperatorName})
		if err == nil {
			
		}else {
			fmt.Println("代理删除时出错: ",err)
		}
		
		_,err = rdb.DataBase().Del("agency",id)
		if err == nil {
			this.Redirect("/agency",302)
		}else {
			this.ShowTips( err )
		}
	}
	this.Render()			
}

func (this *AgencyController) AgencyOrders() {
	
}

func (this *AgencyController) AgencyShares() {

}


