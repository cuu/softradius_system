package controllers

import (
	"time"
	"strconv"
	r "github.com/cuu/softradius/routers"
	"github.com/cuu/softradius/models"
	"github.com/cuu/softradius/libs"
	"github.com/cuu/softradius/libs/times"
	"fmt"
	"reflect"
	"github.com/astaxie/beego"

//	sort "github.com/cuu/softradius/libs/sortutil"
	rdb  "github.com/cuu/softradius/database/shelf"
	
)

type BusController struct {
	BaseController
	PerPage int
}

type AcceptLog struct {
	Id           string `gorethink:"id,omitempty"`
	AcceptType   string  //open
	AcceptDesc   string  //用户新开户：(0000)user
	Account      string  //认证帐号,可以是数字,也可能是email
	Operator     string  //admin
	AcceptSource string  //console
	AcceptTime   string  // 2017-02-12 15:32:41
}

type OrderLog struct{
	Id        string `gorethink:"id,omitempty"`
	MemberId  string
	ProductId string
	OrderFee  int
	ActualFee int
	PayStatus int
	AcceptId  string // from AcceptLog
	OrderSource string
	OrderDesc   string
	CreateTime  string
}

	
type Members struct {
	Id             string  `gorethink:"id,omitempty"`
	NodeId         string  // from Node table
	AgencyId       string //
	Name           string
	Password       string
	RealName       string
	IdCard         string
	Sex            int
	Age            int
	Status         int
	Email          string
	EmailActive    bool
	Mobile         string
	MobileActive   bool
	Address        string
	Desc           string
	BatchId        string
	ProductId      string
	InstallAddress string
	Balance        int    //余额, 分
	TimeLength     int    // 时间 秒
	FlowLength     int    // 流量 kb
	UsedTime       int    //被计费时间
	InFlow         int    //被计费入流量
	OutFlow        int    //被计费出流量
	BindMac        int     
	BindVlan       int
	ConcurNumber   int 
	MacAddr        string
	IpAddress      string
	VlanId         int
	VlanId2        int
	LastPause      string
	ExpireDate     string //格式 Y-m-d ,与 MAX_EXPIRE_DATE 一致	
	CreateTime     string
	UpdateTime     string

}

var _bus_ctl BusController
func init(){
	_ctl := &_bus_ctl

	_cate := r.MenuBus
	
 	_ctl.routes = append( _ctl.routes,
r.Route{Path:"/members",Name:"用户信息管理",Category:_cate,Is_menu :true, Order:1.0,Is_open:true, Methods:"*:Members"})

 	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/member/create",Name:"用户正常开户",Category:_cate,Is_menu :true, Order:1.1,Is_open:true, Methods:"*:MemberCreate"})
	
	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/member/quick", Name:"用户快速开户",Category:_cate,Is_menu:true, Order:1.2,Is_open:true, Methods:"*:MemberQuick"})


	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/member/detail", Name:"用户详细页面",Category:_cate,Is_menu:false, Order:1.3,Is_open:true, Methods:"*:MemberDetail"})
	
	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/bus/opencalc", Name:"用户开户函数",Category:_cate,Is_menu:false, Order:4.3,Is_open:false, Methods:"*:OpenCalc"})
	
	_ctl.AddRoutes()
	
}

//把this中的routes 放到 routers.Permits 中,每个Controller写一遍
func (this *BusController) AddRoutes() {
	
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
func (this *BusController) GuuPrepare(){
	
	this.TplName = libs.GetTplName(this)
	
	this.PerPage = 100
	
}

//------------------------------------------------------------
func (this *BusController) OpenCalc() {
	months,_ := strconv.Atoi(this.GetString("months"))
	product_id := this.GetString("product_id")
	old_expire := this.GetString("old_expire")

	
	product := &Products{}
	err := rdb.DataBase().QuOne(product,product_id)
	if err != nil {
		panic("OpenCalc Get Product")
	}

	type data struct {
		Policy int         `json:"policy"`
		Fee_value string   `json:"fee_value"`
		Expire_date string `json:"expire_date"`
	}
	
	type Ret struct {
		Code int  `json:"code"`
		Data data `json:"data"`
	}
	//预付费时长,预付费流量
	if product.Policy == r.PPTimes ||  product.Policy == r.PPFlow {
		fee_value := libs.Fen2yuan(product.FeePrice)
		
		this.Data["json"] = &Ret{Code:0,Data:data{Policy:product.Policy,Fee_value:fee_value,Expire_date:r.MAX_EXPIRE_DATE}}
		

	//买断时长,买断流量
	}else if product.Policy == r.BOTimes || product.Policy == r.BOFlows {
		fee_value :=libs.Fen2yuan(product.FeePrice)
		this.Data["json"] = &Ret{Code:0,Data:data{Policy:product.Policy,Fee_value:fee_value,Expire_date:r.MAX_EXPIRE_DATE}}		

	// 预付费包月
	}else if product.Policy == r.PPMonth {
		fee := months * product.FeePrice
		fee_value := libs.Fen2yuan(fee)
		start_expire := time.Now()
		if len(old_expire) > 0 {
			start_expire  = times.StrToLocalTime(old_expire)
		}
		
		expire_date_t := libs.AddMonths(start_expire, months)
		expire_date   := times.Format("Y-m-d",expire_date_t)
		this.Data["json"] = &Ret{Code:0,Data:data{Policy:product.Policy,Fee_value:fee_value,Expire_date:expire_date}}
			
		
	//买断包月
	}else if product.Policy == r.BOMonth {
		start_expire := time.Now()
		if old_expire != "" {
			start_expire = times.StrToLocalTime(old_expire)
		}
		fee_value := libs.Fen2yuan(product.FeePrice)
		expire_date_t := libs.AddMonths(start_expire,product.FeeMonths)
		expire_date   := times.Format("Y-m-d",expire_date_t)
		this.Data["json"] = &Ret{Code:0,Data:data{Policy:product.Policy,Fee_value:fee_value,Expire_date:expire_date}}
		
	}else if product.Policy == r.AwesomeFee {
		fee_value := libs.Fen2yuan(product.FeePrice)
		this.Data["json"] = &Ret{Code:0,Data:data{Policy:product.Policy,Fee_value:fee_value,Expire_date:r.MAX_EXPIRE_DATE}}		
	}else if product.Policy == r.AwesomeFeeBoTime {
		fee_value := libs.Fen2yuan(product.FeePrice)
		this.Data["json"] = &Ret{Code:0,Data:data{Policy:product.Policy,Fee_value:fee_value,Expire_date:r.MAX_EXPIRE_DATE}}		
	}

	this.ServeJSON()
}

func (this *BusController) member_list( skip int) []Members {
	var nods []Members
	
	rdb.DataBase().SkipGet2(&nods,skip*this.PerPage,this.PerPage)

	return nods
}

func (this *BaseController) ToPairMapS( v interface{},list []string ) map[string]string {
	var ret = make( map[string]string)
	
	if libs.Type(v) != "slice" {
		panic("Only slice to maps")
		return ret
	}

	if len(list) != 2 {
		fmt.Println("Only Len of 2 list supported")
		return ret
	}
	
	s := reflect.ValueOf(v)
	
	for i:=0;i<s.Len();i++ {
		iv := s.Index(i)
		//now it's Node
		ret[ iv.FieldByName(list[0]).String() ] = iv.FieldByName(list[1]).String()
	}
	
	return ret
}

func (this *BaseController) ToPairMapI( v interface{},list []string ) map[string]int64 {
	var ret = make( map[string]int64)
	
	if libs.Type(v) != "slice" {
		panic("Only slice to maps")
		return ret
	}

	if len(list) != 2 {
		fmt.Println("Only Len of 2 list supported")
		return ret
	}
	
	s := reflect.ValueOf(v)
	
	for i:=0;i<s.Len();i++ {
		iv := s.Index(i)
		//now it's Node
		ret[ iv.FieldByName(list[0]).String() ] = iv.FieldByName(list[1]).Int()
	}
	
	return ret
}

func (this *BusController) Members () {
	nods := this.NodeList()
	pdus := this.ProductList()
	
	mbms := this.member_list(0)


	this.Data["MemberList"] = mbms
	this.Data["ProductMap"] = this.ToPairMapS(pdus,[]string{"Id","Name"})
	this.Data["NodeMap"]    = this.ToPairMapS(nods,[]string{"Id","Name"})
	this.Data["IsExpire"]  = libs.IsExpire
	
	this.TplName ="bus_member_list.html"
	
	this.Render()
}


func (this *BusController) MemberCreateForm(nodes [][]string, pdus [][]string  ) *models.Form {
	f := models.InfoForm("Memeber create","/memeber/create")
	return f
	
}

func (this *BusController) MemberCreate() {

	/*
	nods := this.NodeList()
	pdus := this.ProductList()
	allnodes := this.Items(nods,[]string{"Id","Name"})
	allproducts := this.Items(pdus,[]string{"Id","Name"})

	*/
	
	this.Render()
}

func (this *BusController) MemberQuickForm(nodes [][]string, pdus [][]string ) *models.Form {
	
	f := models.InfoForm("Member Quick","/member/quick",
		models.Dropdown(&models.Select{Name:"NodeId",Description:"区域",Args:nodes}),
		models.Dropdown(&models.Select{Name:"ProductId",Description:"资费",Args:pdus}),
		models.TextBox(&models.Input{Name:"Name",Description:"用户帐号",Required:true,Valid:models.Len_of(2,128)}),
		models.TextBox(&models.Input{Name:"Password",Description:"认证密码",Required:true,Valid:models.Len_of(6,32)}),
		models.TextBox(&models.Input{Name:"Months",Description:"月数(包月有效)",Required:true}),
		models.TextBox(&models.Input{Name:"ExpireDate",Description:"过期日期",Required:true,ReadOnly:true,Valid:models.Is_date}),
		models.Hidden(&models.Input{Name:"Status",Value:"1",Description:"用户状态"}),
		models.TextArea(&models.Input{Name:"Desc",Description:"用户描述"}),
		models.Submit(&models.Input{Name:"Submit",Value:"<b>提交</b>"}) )
		
	return f 
	
}
func (this *BusController) MemberQuick() {

	this.TplName ="bus_quick_form.html"
	nods := this.NodeList()
	pdus := this.ProductList()
	allnodes := this.Items(nods,[]string{"Id","Name"})
	allproducts := this.Items(pdus,[]string{"Id","Name"})

	f := this.MemberQuickForm(allnodes,allproducts)
	
	if this.Ctx.Input.IsPost() {
		if this.Validator2(f) == false {
			this.Data["Form"] = f.Render()
			this.Render()
			return
		}
		one := &Members{}
		this.ParsePostToStruct(one)

		for _,p := range pdus {
			if one.ProductId == p.Id {
				one.ConcurNumber = p.ConcurNumber
				break
			}
		}
		one.CreateTime = libs.Get_currtime()
		one.UpdateTime = libs.Get_currtime()
		
		_,cnt := rdb.DataBase().FilterInsert(one,"Name")
		if cnt  == 0 {
			
			this.Redirect("/members", 302)
		}else {
			
			this.ShowTips("用户有重复 "+strconv.Itoa(cnt) +"个" )
		}

		this.Render()
		return
	}

	
	//类似TR1 member/open
	
	this.Data["Form"] = f.Render()
	
	this.Render()
}


func (this *BusController) MemberDetail() {
	member_id := this.GetString("member_id")
	nods := this.NodeList()
	pdus := this.ProductList()
	agcs := this.AgencyList()
	
	//allproducts := this.Items(pdus,[]string{"Id","Name"})
	
	if this.Ctx.Input.IsPost() {


		
	}


	
	if member_id == ""{
		this.Abort("403")
		return
	}
	one := &Members{}
	err := rdb.DataBase().QuOne(one,member_id)
	if err == nil {

		var orderlogs []OrderLog
		var acceptlogs []AcceptLog
		
		this.Data["User"] = one
		this.Data["ProductsMap"] = this.ToPairMapS(pdus,[]string{"Id","Name"})
		this.Data["AgencyMap"  ] = this.ToPairMapS(agcs,[]string{"Id","Name"})
		this.Data["PolicyMap"] = this.ToPairMapI(pdus,[]string{"Id","Policy"})
		this.Data["NodeMap"]    = this.ToPairMapS(nods,[]string{"Id","Name"})
		this.Data["OrderLogs" ] = orderlogs
		this.Data["AcceptLogs"] = acceptlogs
		this.Data["IsExpire"]   = libs.IsExpire
		this.Data["YESNO"] = YESNO
		
		
	}
	
	this.TplName = "bus_member_detail.html"
	
	this.Render()
}
