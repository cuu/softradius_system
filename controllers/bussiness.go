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
	"math"
	//	"math/rand"
	"strings"
	"github.com/astaxie/beego"

//	sort "github.com/cuu/softradius/libs/sortutil"
	rdb  "github.com/cuu/softradius/database/shelf"
	
)

const (
	USEROK =1
	USERPAUSE = 2
	USERSOLD= 3
	USEREXPIRE = 4
)

var UserState = map[int]string{1: "正常", 2: "停机", 3: "销户", 4: "到期"}

type BusController struct {
	BaseController
}

//只允许连续生成,无随机生成的批量
type Batch struct {
	Id            string `gorethink:"id,omitempty"`
	Name          string
	NodeId        string
	ProductId     string
	Quantity      int
	UserPrefix    string
	Desc          string
	Count         int
	CreateTime    string
}

//Batch rules 
type BatRule struct {
	Id string `gorethink:"id,omitempty"`
	Name string   //Rule name
	UserPrefix string   //Batch account Prefix
	UserSuffixLen int // Length of account Suffix part
	UserSn   int     // SerialNumber of this Batch in this rule
	Count    int     // total count
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

//Member order
type OrderLog struct{
	Id        string `gorethink:"id,omitempty"`
	MemberId  string //GeneratedKeys
	ProductId string //
	OrderFee  int
	ActualFee int
	PayStatus int    // normally it's 1
	AcceptId  string // from AcceptLog
	OrderSource string //
	OrderDesc   string
	CreateTime  string // 2017-02-12 15:32:41	
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
	ActiveCode     string  // an uuid 
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
		r.Route{Path:"/member/batch", Name:"用户批量开户",Category:_cate,Is_menu:true, Order:1.21,Is_open:true, Methods:"*:MemberBatch"})

	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/member/batch/add", Name:"创建批量开户",Category:_cate,Is_menu:false, Order:1.211,Is_open:true, Methods:"*:MemberBatchAdd"})

	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/member/batch/delete", Name:"删除批量",Category:_cate,Is_menu:false, Order:1.212,Is_open:true, Methods:"*:MemberBatchDel"})
	

	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/member/batrule", Name:"批量开户规则",Category:_cate,Is_menu:true, Order:1.221,Is_open:true, Methods:"*:MemberBatchRules"})

	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/member/batrule/add", Name:"增加批量规则",Category:_cate,Is_menu:false, Order:1.222,Is_open:true, Methods:"*:MemberBatchRuleAdd"})
	
	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/member/batrule/update", Name:"更新批量规则",Category:_cate,Is_menu:false, Order:1.223,Is_open:true, Methods:"*:MemberBatchRuleUpdate"})

	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/member/batrule/delete", Name:"删除批量规则",Category:_cate,Is_menu:false, Order:1.224,Is_open:true, Methods:"*:MemberBatchRuleDelete"})

	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/member/detail", Name:"用户详细页面",Category:_cate,Is_menu:false, Order:1.3,Is_open:true, Methods:"*:MemberDetail"})

	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/member/delete", Name:"用户详细页面",Category:_cate,Is_menu:false, Order:1.31,Is_open:true, Methods:"*:MemberDelete"})
	
	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/bus/opencalc", Name:"用户开户函数",Category:_cate,Is_menu:false, Order:8.3,Is_open:false, Methods:"*:OpenCalc"})
	
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
		
		//Permits.routes is map,so no confict path key ! 
		r.Permits.Add_route(v.Path,&this.routes[i])
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


func (this *BusController) MemberCreateForm(nodes [][]string, pdus [][]string ,agencies [][]string ,user_state [][]string ) *models.Form {
	
	f := models.InfoForm("Memeber create","/member/create",
		models.Dropdown(&models.Select{Name:"NodeId",Description:"区域",Args:nodes}),
		models.Dropdown(&models.Select{Name:"ProductId",Description:"资费",Args:pdus}),
		models.TextBox(&models.Input{Name:"RealName",Description:"用户姓名",Required:true,Valid:models.Len_of(2,32),Value:"User"}),
		models.Dropdown(&models.Select{Name:"AgencyId",Description:"代理",Args:agencies}),
		models.TextBox(&models.Input{Name:"IdCard",Description:"证件号码",Valid:models.Len_of(0,32)}),
		models.TextBox(&models.Input{Name:"Mobile",Description:"用户手机号码",Valid:models.Len_of(0,32)}),
		models.TextBox(&models.Input{Name:"Address",Description:"用户地址"}),
		models.Hr(&models.Input{}),
		models.TextBox(&models.Input{Name:"Name",Description:"用户帐号",Required:true,Valid:models.Len_of(2,128)}),
		models.TextBox(&models.Input{Name:"Password",Description:"认证密码",Required:true,Valid:models.Len_of(6,32)}),
		models.TextBox(&models.Input{Name:"IpAddress",Description:"用户IP地址"}),
		models.TextBox(&models.Input{Name:"Months",Valid:models.Is_number,Description:"月数(包月有效)",Required:true}),
		models.TextBox(&models.Input{Name:"FeeValue",Valid:models.Is_rmb,Description:"缴费金额",Required:true}),
		models.TextBox(&models.Input{Name:"ExpireDate",Description:"过期日期",Required:true,ReadOnly:true,Valid:models.Is_date}),
		models.Dropdown(&models.Select{Name:"Status",Description:"用户状态",Args:user_state}),
		models.TextArea(&models.Input{Name:"Desc",Description:"用户描述"}),
		models.Submit(&models.Input{Name:"Submit",Value:"<b>提交</b>",Class:"btn btn-info"}) )
	
	return f
	
}

func (this *BusController) MemberCreate() {
	nods := this.NodeList()
	pdus := this.ProductList()
	agcs := this.AgencyList()
	var user_state [][]string
	
	allnodes := this.Items(nods,[]string{"Id","Name"})
	allproducts := this.Items(pdus,[]string{"Id","Name"})
	agcs_items := this.Items(agcs,[]string{"Id","Name"})

	var allagency [][]string
	allagency = append(allagency,[]string{"0",""})
	allagency = append(allagency,agcs_items...)
	
	for i:=1; i <= len(UserState); i++ {
		v := UserState[i]
		a := []string{strconv.Itoa(i),v}
		user_state = append(user_state,a)
	}

	f:= this.MemberCreateForm(allnodes,allproducts,allagency,user_state)
	this.TplName = "bus_open_form.html"
	
	if this.Ctx.Input.IsPost() {

		fmt.Println("in post")
		if this.Validator2(f) == false {
			this.Data["Form"] = f
			this.Render()
			return
		}
		
		agc_id := this.GetString("AgencyId")
		agc := &Agency{}
		opr := &Operators{}
		one := &Members{}
		order_log := &OrderLog{}
		accept_log := &AcceptLog{}
		feevalue := math.Ceil(this.GetStringF("FeeValue"))
		
		fmt.Println("FeeValue is :", feevalue)
		
		if len(agc_id) > 32 {  // With agency 
			err := rdb.DataBase().QuOne(agc,agc_id)
			
			if err == nil {
				if agc.Amount < libs.Yuan2fen( int(feevalue) ) {
					this.ShowTips("代理商金额不足")
					this.Render()
					return
				}

				err = rdb.DataBase().FilterOne(opr,map[string]string{"Name":agc.OperatorName})
				if err != nil {
					this.ShowTips("代理商操作员不存在")
					this.Render()
					return
				}
			}else {
				this.ShowTips("代理商信息错误")
				this.Render()
				return
			}
			
			
		}else { //没有代理商
			fmt.Println("no agency")
		}

		balance     := 0
		order_fee   := 0
		expire_date := this.GetString("ExpireDate")

		this.ParsePostToStruct(one)

		
		if libs.InSlice(one.NodeId, opr.Nodes) == false && len(agc_id) > 32 {
			this.ShowTips("代理商在此区域无权新增用户")
			this.Render()
			return
		}
		
		if libs.InSlice(one.ProductId, opr.Products) == false && len(agc_id) > 32 {
			this.ShowTips("代理商在此资费下无权新增用户")
			this.Render()
			return
		}
		
		for _,p := range pdus {
			if one.ProductId == p.Id {
				one.ConcurNumber = p.ConcurNumber
				
				if p.Policy == r.PPMonth {
					months := this.GetStringI("Months")
					order_fee = p.FeePrice *  months
					
				}else if libs.In(p.Policy, r.BOMonth,r.BOTimes){
					order_fee = p.FeePrice
				}else if libs.In(p.Policy,r.PPTimes,r.PPFlow) {
					balance = libs.Yuan2fen( int(feevalue))
					expire_date = r.MAX_EXPIRE_DATE
				}else if p.Policy == r.AwesomeFee {
					expire_date = r.MAX_EXPIRE_DATE
				}else if p.Policy ==r.AwesomeFeeBoTime {
					expire_date = r.MAX_EXPIRE_DATE
					order_fee = p.FeePrice
				}
				break
			}
		}
		
		one.CreateTime   = libs.Get_currtime()
		one.UpdateTime   = libs.Get_currtime()
		one.Balance      = balance
		one.ActiveCode,_ = libs.NewUUID()
		one.ExpireDate   = expire_date
		
		rsp,cnt := rdb.DataBase().FilterInsert(one,"Name")
		if cnt  == 0 {  // successfully fucked into
			// insert accept log and order log
			member_id := rsp[0] 
			accept_log.AcceptType   = "open"
			accept_log.AcceptSource = "console"
			accept_log.Account      = one.Name
			accept_log.AcceptTime   = one.CreateTime
			accept_log.Operator     = this.GetCookie("username")
			accept_log.AcceptDesc   = "用户新开帐号"
			
			rsp,err := rdb.DataBase().InsertQ(accept_log)
			if err != nil { fmt.Println(err) }
			
			order_log.MemberId  = member_id
			order_log.ProductId = one.ProductId
			order_log.OrderFee  = order_fee
			order_log.ActualFee = libs.Yuan2fen( int(feevalue))
			order_log.PayStatus = 1
			order_log.AcceptId  = rsp[0]
			order_log.OrderSource = "console"
			order_log.OrderDesc  = "用户新开帐号"
			order_log.CreateTime = one.CreateTime
			
			rsp,err = rdb.DataBase().InsertQ(order_log)
			
			if len(agc_id) > 32 {// Blow job, Do the Agency
				//agency order 1,open account
				//agency order 2,agency share cut
				//agency share log...
				agc_order1 := &AgencyOrder{}
				agc_order2 := &AgencyOrder{}
				agc_share  := &AgencyShare{}
				
				agc_order1.AgencyId      = agc_id
				agc_order1.MemberOrderId = rsp[0]
				agc_order1.FeeType       = "cost"
				agc_order1.FeeValue      = libs.Yuan2fen( int(feevalue) )
				agc_order1.FeeTotal      = (agc.Amount - agc_order1.FeeValue)
				agc_order1.FeeDesc       = ("代理商开户 "+one.Name)
				agc_order1.CreateTime    = libs.Get_currtime()
				agc.Amount = agc_order1.FeeTotal
				
				rsp,err = rdb.DataBase().InsertQ(agc_order1)
				
				agc_order2.AgencyId      = agc_id
				agc_order2.MemberOrderId = agc_order1.MemberOrderId
				agc_order2.FeeType       = "share"

				
				agc_order2.FeeValue      = agc_order1.FeeValue
				rate := float64(float64(agc.ShareRate)/100.00)
			//	fmt.Println(agc_order1.FeeValue, rate)
				agc_order2.FeeValue = int(float64(agc_order1.FeeValue)*rate)
			//	fmt.Println(agc_order2.FeeValue, rate)
				
				agc_order2.FeeTotal      = (agc.Amount + agc_order2.FeeValue)
				agc_order2.FeeDesc       = fmt.Sprintf("代理商分成 %s %f ",one.Name,rate)
				agc_order2.CreateTime    = libs.Get_currtime()
				agc.Amount = agc_order2.FeeTotal
				rsp,err = rdb.DataBase().InsertQ(agc_order2)
				
				agc_share.AgencyId   = agc_id
				agc_share.OrderId    = rsp[0]
				agc_share.ShareRate  = agc.ShareRate
				agc_share.ShareFee   = agc_order2.FeeValue
				agc_share.FeeValue   = agc_order1.FeeValue
				agc_share.NodeId     = one.NodeId
				agc_share.ProductId  = one.ProductId
				agc_share.CreateTime = libs.Get_currtime()
				rsp,err = rdb.DataBase().InsertQ(agc_share)

				//refresh the amount
				fmt.Println(agc.Amount)
				rdb.DataBase().Update(agc_id,agc)
				
			}
			this.Redirect("/members", 302)
		}else {
			
			this.ShowTips("用户名有重复 "+strconv.Itoa(cnt) +"个" )

			this.Render()
			return
		}

		this.Render()
		return
		
	}

	this.Data["Form"] = f
	
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
		models.Hidden(&models.Input{Name:"FeeValue",Value:"0",Description:"费用"}),
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
		balance := 0
		order_fee := 0
		expire_date := r.MAX_EXPIRE_DATE
		
		this.ParsePostToStruct(one)
		
		var _pdu = &Products{}
		for _,p := range pdus {
			if one.ProductId == p.Id {
				one.ConcurNumber = p.ConcurNumber
				_pdu = &p
				if p.Policy == r.PPMonth {
					months := this.GetStringI("Months")
					order_fee = p.FeePrice *  months
				}else if libs.In(p.Policy, r.BOMonth,r.BOTimes){
					order_fee = p.FeePrice
				}else if libs.In(p.Policy,r.PPTimes,r.PPFlow) {
					balance = libs.Yuan2fen(this.GetStringI("FeeValue"))
					expire_date = r.MAX_EXPIRE_DATE
				}else if p.Policy == r.AwesomeFee {
					expire_date = r.MAX_EXPIRE_DATE
				}else if p.Policy ==r.AwesomeFeeBoTime {
					expire_date = r.MAX_EXPIRE_DATE
					order_fee = p.FeePrice
				}				
				break
			}
		}
		
		one.CreateTime = libs.Get_currtime()
		one.UpdateTime = one.CreateTime
		one.Balance = balance
		one.TimeLength = _pdu.FeeTimes
		one.FlowLength = _pdu.FeeFlows
		one.ExpireDate = expire_date
		rsp,cnt := rdb.DataBase().FilterInsert(one,"Name")
		if cnt  == 0 {
			//accept log
			//order log

			order_log := &OrderLog{}
			accept_log := &AcceptLog{}
			// accept log
			// order log
			member_id := rsp[0]
			
			accept_log.AcceptType   = "open"
			accept_log.AcceptSource = "console"
			accept_log.Account      = one.Name
			accept_log.AcceptTime   = one.CreateTime
			accept_log.Operator     = this.GetCookie("username")
			accept_log.AcceptDesc   = "用户新开帐号"
			
			rsp,err := rdb.DataBase().InsertQ(accept_log)
			if err != nil { fmt.Println(err) }
			
			order_log.MemberId  = member_id
			order_log.ProductId = one.ProductId
			order_log.OrderFee  = order_fee
			order_log.ActualFee = libs.Yuan2fen(this.GetStringI("FeeValue"))
			order_log.PayStatus = 1
			order_log.AcceptId  = rsp[0]
			order_log.OrderSource = "console"
			order_log.OrderDesc  = "用户新开帐号"
			order_log.CreateTime = one.CreateTime
			
			rsp,err = rdb.DataBase().InsertQ(order_log)
			
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

func (this *BusController) MemberDelete() {
	id := this.GetString("member_id")

	if id == "" {
		this.Abort("403")
		return
	}
	one := &Members{}
	err:= rdb.DataBase().QuOne(one,id)
	if err == nil {
		bat := &Batch{}
		rule := &BatRule{}
		if one.BatchId != "" {	
			err = rdb.DataBase().QuOne(bat,one.BatchId)
			if err == nil {
				err = rdb.DataBase().FilterOne(rule,map[string]string{"UserPrefix":bat.UserPrefix})
				rule.Count--
				//Count减少,UserSn只增不减,直到爆破
				bat.Count--
				rdb.DataBase().Update(rule.Id,rule)
				rdb.DataBase().Update(bat.Id,bat)
			}
		}
	
		this.IdRowDelete(Members{},id,"/members")
		this.Render()
		return
	}
	return	
}

func (this *BusController) MemberBatchPageData(skip int)(int,[]Batch) {
	var items []Batch
	rdb.DataBase().SkipGet2(&items,skip,this.PerPage)
	total := rdb.DataBase().TableCount(&items)

	
	return total,items
}

func (this *BusController) MemberBatch() {
	
	this.TplName = "bus_batch_list.html"
	
	nods := this.NodeList()
	pdus := this.ProductList()

	page := this.InitPage()	
	total,results := this.MemberBatchPageData( libs.Or(page.Page,0).(int)*this.PerPage )
	page.MakePager(total)
	
	this.Data["Results"]     = results
	this.Data["Paginator"]   = page.Render()	
	this.Data["NodeMap"]     = this.ToPairMapS(nods,[]string{"Id","Name"})
	this.Data["ProductsMap"] = this.ToPairMapS(pdus,[]string{"Id","Name"})
	
	this.Render()
	
	
}

func (this *BusController) MemberBatchAddForm(nodes [][]string, pdus [][]string, rules [][]string) *models.Form {
	f := models.InfoForm("AddBatch","/member/batch/add",
		models.Dropdown(&models.Select{Name:"NodeId",Description:"用户区域",Args:nodes}),
		models.Dropdown(&models.Select{Name:"ProductId",Description:"用户资费",Args:pdus}),
		models.Dropdown(&models.Select{Name:"RuleId",Description:"创建规则",Args:rules}),
		models.TextBox(&models.Input{Name:"Name",Description:"批量名称",Valid:models.Notnull}),
		models.TextArea(&models.Input{Name:"Desc",Description:"批量描述",Valid:models.Len_of(0,100)}),
		models.TextBox(&models.Input{Name:"Count",Description:"创建用户数量,最大1000",Valid:models.Is_number}),
		models.TextBox(&models.Input{Name:"PassLen",Description:"密码长度,2-10位",Valid:models.Is_number}),
		models.Submit(&models.Input{Name:"Submit",Value:"<b>提交</b>",Class:"btn btn-info"}))
	return f
}


func (this *BusController) MemberBatchAddUserPassword(pass_len int ) string {
	//2017 03 15: only random number password for batch users now
	if pass_len - 1 <= 0 {  // pass_len must bigger than 1 
		return "NoPass"
	}
	start ,_:= strconv.Atoi("1"+ strings.Repeat("0",pass_len-1))
	end ,_  := strconv.Atoi( strings.Repeat("9",pass_len))

	n := libs.Random(start,end)
	ns := strconv.Itoa(n)
	return ns
}

func (this *BusController) MemberBatchAdd() {
	var rules []BatRule
	
	nods := this.NodeList()
	pdus := this.ProductList()
	rdb.DataBase().SkipGet2(&rules,0,1000)
	
	allnodes    := this.Items(nods,[]string{"Id","Name"})
	allproducts := this.Items(pdus,[]string{"Id","Name"})
	allrules    := this.Items(rules,[]string{"Id","Name"})

	f := this.MemberBatchAddForm(allnodes,allproducts,allrules)

	if this.InPost() {
		pdu_id  := this.POST("ProductId")
		_pdu := &Products{}
		for _,v := range pdus {
			if v.Id == pdu_id {
				_pdu = &v
				break
			}
		}
		
		rule_id := this.POST("RuleId")
		if len(rule_id)< 32 {
			this.ShowTips("必须要有合法的批量规则")
			this.Render()
			return
		}

		count := this.GetStringI("Count")
		if count <= 0 || count > 1000 {
			this.ShowTips("创建数量不正确,0-1000范围内")
			this.Render()
			return
		}
		
		passlen := this.GetStringI("PassLen")
		if passlen < 2 || passlen > 10 {
			this.ShowTips("密码长度2-10位")
			this.Render()
			return
		}
		
		rule := &BatRule{}
		rdb.DataBase().QuOne(rule,rule_id)		
		
		bat := &Batch{}
		
		this.ParsePostToStruct(bat)
		bat.UserPrefix = rule.UserPrefix
		bat.Count = 0 // bat count depends on the actual inserted
		
		agc := &Agency{}
		opera := this.GetOperator()
		if opera.Type == AGENCYOPERA {
			err := rdb.DataBase().FilterOne(agc,map[string]string{"OperatorName":opera.Name})
			if err != nil {
				fmt.Println("BaseController GetOperator error")
			}
			
			
			allfee := _pdu.FeePrice *count
			if agc.Amount < allfee {
				this.ShowTips(fmt.Sprintf("代理商您的余额不足于生成此次批量,余额为%d,此次共需%d,%d个用户,单价为%d",libs.Fen2yuan(agc.Amount),libs.Fen2yuan(allfee),count,libs.Fen2yuan(_pdu.FeePrice)))
				this.Render()
				return
		
			}
			
		}
				
		resp,cnt := rdb.DataBase().FilterInsert(bat,"Name")
		if cnt  == 0 {
			this.AddOperLog(fmt.Sprintf("新增批量%s",bat.Name))
			//Check oper type,if is agency,deal with the Amount
			//insert members
			
			total_injection := 0
			//检查位数与实际批量数字的大小,不可以超过UserSuffixeLen范围
			rule_format := fmt.Sprintf("%%0%dd",rule.UserSuffixLen)
			max_count := count
			if s,err := strconv.Atoi(strings.Repeat("9",rule.UserSuffixLen)); err == nil {
				if s < max_count {
					max_count = s
				}
			}

			order_fee := 0
			balance := 0 // 批量的用户所有的初始余额是0
			expire_date := r.MAX_EXPIRE_DATE
			
			if _pdu.Policy == r.PPMonth {
				order_fee = _pdu.FeePrice * 1 //默认1个月批量
			}else if libs.In(_pdu.Policy, r.BOMonth,r.BOTimes){
				order_fee = _pdu.FeePrice
			}else if libs.In(_pdu.Policy,r.PPTimes,r.PPFlow) {
				expire_date = r.MAX_EXPIRE_DATE
			}else if _pdu.Policy == r.AwesomeFee {
				expire_date = r.MAX_EXPIRE_DATE
			}else if _pdu.Policy ==r.AwesomeFeeBoTime {
				expire_date = r.MAX_EXPIRE_DATE
				order_fee = _pdu.FeePrice
			}

			for i:=0 ;i<max_count;i++ {
				member := &Members{}
				
				member.Name = rule.UserPrefix + fmt.Sprintf(rule_format,i+rule.UserSn)
				member.RealName = fmt.Sprintf("批量%s用户",bat.Name)
				member.BatchId = resp[0]
				member.NodeId  = this.POST("NodeId")
				member.ProductId = _pdu.Id
				member.AgencyId = agc.Id
				member.ConcurNumber = _pdu.ConcurNumber
				member.CreateTime  = libs.Get_currtime()
				member.UpdateTime  = member.CreateTime
				member.Password    = this.MemberBatchAddUserPassword(passlen)
				member.Desc      = fmt.Sprintf("批量 %s,操作员 %s",bat.Name,opera.Name)
				member.ExpireDate = expire_date
				member.Balance    = balance
				member.Status     = USEROK
				member.TimeLength = _pdu.FeeTimes
				member.FlowLength = _pdu.FeeFlows
				
				rsp,err := rdb.DataBase().Insert(member)
				if err != nil {
					//也许有重复之类的错误,不管它
					fmt.Println("Batch Insert: ",err)
				}else{
					total_injection += rsp.Inserted
					order_log := &OrderLog{}
					accept_log := &AcceptLog{}
					// accept log
					// order log
					member_id := rsp.GeneratedKeys[0]
					
					accept_log.AcceptType   = "open"
					accept_log.AcceptSource = "console"
					accept_log.Account      = member.Name
					accept_log.AcceptTime   = member.CreateTime
					accept_log.Operator     = opera.Name
					accept_log.AcceptDesc   = "用户新开帐号"
					
					rsp,err := rdb.DataBase().InsertQ(accept_log)
					if err != nil { fmt.Println(err) }
					
					order_log.MemberId  = member_id
					order_log.ProductId = member.ProductId
					order_log.OrderFee  = order_fee
					order_log.ActualFee = order_fee
					order_log.PayStatus = 1
					order_log.AcceptId  = rsp[0]
					order_log.OrderSource = "console"
					order_log.OrderDesc  = "用户新开帐号"
					order_log.CreateTime = member.CreateTime
					
					rsp,err = rdb.DataBase().InsertQ(order_log)
					
					if opera.Type == AGENCYOPERA {
						/// add agency order and share
						agc_order1 := &AgencyOrder{}
						agc_order2 := &AgencyOrder{}
						agc_share  := &AgencyShare{}
						
						agc_order1.AgencyId      = agc.Id
						agc_order1.MemberOrderId = rsp[0]
						agc_order1.FeeType       = "cost"
						agc_order1.FeeValue      = order_fee
						agc_order1.FeeTotal      = (agc.Amount - agc_order1.FeeValue)
						agc_order1.FeeDesc       = ("代理商批量开户 "+member.Name)
						agc_order1.CreateTime    = libs.Get_currtime()
						agc.Amount = agc_order1.FeeTotal
						
						rsp,err = rdb.DataBase().InsertQ(agc_order1)
						
						agc_order2.AgencyId      = agc.Id
						agc_order2.MemberOrderId = agc_order1.MemberOrderId
						agc_order2.FeeType       = "share"
						
						
						agc_order2.FeeValue      = agc_order1.FeeValue
						rate := float64(float64(agc.ShareRate)/100.00)
						//	fmt.Println(agc_order1.FeeValue, rate)
						agc_order2.FeeValue = int(float64(agc_order1.FeeValue)*rate)
						//	fmt.Println(agc_order2.FeeValue, rate)
						
						agc_order2.FeeTotal      = (agc.Amount + agc_order2.FeeValue)
						agc_order2.FeeDesc       = fmt.Sprintf("代理商分成 %s %f ",member.Name,rate)
						agc_order2.CreateTime    = libs.Get_currtime()
						agc.Amount = agc_order2.FeeTotal
						rsp,err = rdb.DataBase().InsertQ(agc_order2)
						
						agc_share.AgencyId   = agc.Id
						agc_share.OrderId    = rsp[0]
						agc_share.ShareRate  = agc.ShareRate
						agc_share.ShareFee   = agc_order2.FeeValue
						agc_share.FeeValue   = agc_order1.FeeValue
						agc_share.NodeId     = member.NodeId
						agc_share.ProductId  = member.ProductId
						agc_share.CreateTime = libs.Get_currtime()
						rsp,err = rdb.DataBase().InsertQ(agc_share)
						
						//refresh the amount
						fmt.Println(agc.Amount)
						rdb.DataBase().Update(agc.Id,agc)
						
					}
				}
			}
				
			rule.UserSn += total_injection
			rule.Count  += total_injection 
			bat.Count    = total_injection
			
			rdb.DataBase().Update(resp[0],bat)
			rdb.DataBase().Update(rule_id,rule)
			
			this.Redirect("/member/batch", 302)
			return
			
		}else {
			
			this.ShowTips("批量名称有重复 "+strconv.Itoa(cnt) +"个" )
			this.Render()
			return
		}
		
		return
	}
	
	this.Data["Form"] = f.Render()

	this.Render()
	
}

func (this *BusController) MemberBatchDel() {
	id := this.GetString("batch_id")

	if id == "" {
		this.Abort("403")
		return
	}
	bat := &Batch{}
	err := rdb.DataBase().QuOne(bat,id)
	if err == nil {
		rule := &BatRule{}
		err = rdb.DataBase().FilterOne(rule,map[string]string{"UserPrefix":bat.UserPrefix})	
		resp,err:= rdb.DataBase().FilterDel(Members{},map[string]string{"BatchId":id})
		if err == nil {
			rule.Count -= resp.Deleted
			rdb.DataBase().Update(rule.Id,rule)
		}else {
			fmt.Println(err)
		}
		
	}
	this.IdRowDelete(bat,id,"/member/batch")
	this.Render()
	return	
}


func (this *BusController) MemberBatchRules() {
	var items []BatRule
	
	this.TplName = "bus_rule_list.html"

	rdb.DataBase().SkipGet2(&items,0,1000)

	this.Data["Rules"] = items

	this.Render()
	
}


func (this *BusController) MemberBatchRuleAddForm() *models.Form {
	f := models.InfoForm("Add Rule","/member/batrule/add",
		models.TextBox(&models.Input{Name:"Name",Required:true,Valid:models.Notnull,Description:"规则名称"}),
		models.TextBox(&models.Input{Name:"UserPrefix",Required:true,Valid:models.Notnull,Description:"帐号前辍"}),
		models.TextBox(&models.Input{Name:"UserSuffixLen",Valid:models.Is_number,Description:"帐号后辍长度",Value:"6",Required:true}),
		models.Submit(&models.Input{Name:"submit",Value:"<b>提交</b>",Class:"btn btn-info"}),
	)

	return f
		
}

func (this *BusController) MemberBatchRuleAdd() {

	f:= this.MemberBatchRuleAddForm()

	if this.InPost() {
		if this.Validator2(f) == false {
			this.Data["Form"] = f.Render()
			this.Render()
			return
		}

		one := &BatRule{}
		this.ParsePostToStruct(one)
		_,cnt := rdb.DataBase().FilterInsert(one,"UserPrefix")
		if cnt  == 0 {
			this.AddOperLog(fmt.Sprintf("新增规则%s",one.Name))
			this.Redirect("/member/batrule", 302)
			return
		}else {
			
			
			this.ShowTips("用户前辍有重复 "+strconv.Itoa(cnt) +"个" )
			
		}
		
		this.Render()
		return
	}
	
	this.Data["Form"] = f.Render()

	this.Render()
	return
}

func (this *BusController) MemberBatchRuleUpdateForm() *models.Form {
	f := models.InfoForm("Update","/member/batrule/update",
		models.TextBox(&models.Input{Name:"Name",Required:true,Valid:models.Notnull,Description:"规则名称",ReadOnly:true}),
		models.TextBox(&models.Input{Name:"UserPrefix",Required:true,Valid:models.Notnull,Description:"帐号前辍"}),
		models.TextBox(&models.Input{Name:"UserSuffixLen",Valid:models.Is_number,Description:"帐号后辍长度",Value:"6",Required:true}),
		models.Hidden(&models.Input{Name:"Id",Description:"编号"}),
		models.Submit(&models.Input{Name:"submit",Value:"<b>提交</b>",Class:"btn btn-info"}),
	)
	return f
	
}

func (this *BusController) MemberBatchRuleUpdate() {
	f := this.MemberBatchRuleUpdateForm()
	if this.InPost() {
		if this.Validator2(f) == false {
			this.Data["Form"] = f.Render()
			this.Render()
			return
			
		}else {
			id := this.POST("Id")
			if len(id) < 32 {
				this.ShowTips("Form error")
				this.Render()
				return
			}
			
			one := &BatRule{}
			this.ParsePostToStruct(one)
			resp,err := rdb.DataBase().Update(id,one)
			if err == nil {
				fmt.Println("Replaced ",resp.Replaced )
			}else {
				fmt.Println(err)
			}
			
			this.Redirect("/member/batrule",302)
			return	
		}
		return
	}
	
	id := this.GET("rule_id")	
	if id == "" {
		this.Abort("403")
		return
	}
	
	one := &BatRule{}
	err := rdb.DataBase().QuOne(one,id)
	if err == nil {
		f.FillFormFromStruct(one)
		this.Data["Form"] = f.Render()
	}
	
	this.Render()
	
}

func (this *BusController) MemberBatchRuleDelete() {
	
	id := this.GET("rule_id")
	if id == "" {
		this.Abort("403")
		return
	}

	
	this.IdRowDelete(BatRule{},id,"/member/batrule")

	
	this.Render()
	return
}

