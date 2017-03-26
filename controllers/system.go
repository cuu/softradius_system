package controllers

import (
	r "github.com/cuu/softradius/routers"
	"github.com/cuu/softradius/models"
	"github.com/cuu/softradius/libs"
	"fmt"
	"strconv"
//	"reflect"
	"github.com/astaxie/beego"

	sort "github.com/cuu/softradius/libs/sortutil"
	rdb  "github.com/cuu/softradius/database/shelf"
//	"encoding/json"
//	"github.com/pquerna/ffjson/ffjson"
	
)

var ProductPolicys = map[int]string{r.PPMonth: "预付费包月",
	r.PPTimes: "预付费时长", r.BOMonth: "买断包月",
	r.BOTimes: "买断时长", r.PPFlow: "预付费流量", r.BOFlows: "买断流量",
	r.AwesomeFee:"买断流量+时间",
	r.AwesomeFeeBoTime:"买断时间",}

type Bas struct{
	Id string `gorethink:"id,omitempty"`
	IpAddr    string
	Name      string
	Secret    string
	CoaPort   int
	TimeType  int
	VendorId  int
}

type Operators struct{
	Id   string `gorethink:"id,omitempty"`
	Name string
	Pass string
	Desc string
	Type int
	Status int
	RuleItem []string
	Products []string  `only products ids`
	Nodes    []string  `only nodes ids`
	
}

type Node struct {
	Id string `gorethink:"id,omitempty"`
	Name string
	Desc string
}

type ProductAttr struct {
	Id         string `gorethink:"id,omitempty"`
	ProductId  string
	Name   string
	Value  string
	Desc   string
}


type Products struct {
	Id string  `gorethink:"id,omitempty"`
	Name string
	Policy int
	Status int
	BindMac int
	BindVlan int
	ConcurNumber int
	FeePeriod string
	FeeMonths int
	FeeTimes  int /// unit: Seconds
	FeeFlows  int
	FeePrice  int  // unit: cents
	InputMaxLimit int
	OutputMaxLimit int
	CreateTime  string
	UpdateTime  string
	Attrs []ProductAttr
}


type SysController struct {
	BaseController
	Rules []string
}


var TimeTypeMap = map[int]string {
	0:"标准时区，北京时间",
	1:"时间与时区相同",
}


//标准永远是 >>1<< 
var BasVendorTypeMap = map[int]string {
	1: "标准",
	9: "思科",
	3041: "阿尔卡特",
	2352: "爱立信",
	2011: "华为",
	25506: "H3C",
	3902: "中兴",
	10055: "爱快",
	14988: "RouterOS",
}


//每个controller 有这样的命名规则,保证不重复
var _sys_ctl SysController

func init(){
	//在这儿独立处理映身关系
	// 事实上只有 方法是 get的 Route才应该是 Is_menu显示出来

	_ctl := &_sys_ctl

	_cate := r.MenuSys
	
	_ctl.routes = append( _ctl.routes,
	r.Route{Path:"/node",Name:"区域信息管理",Category:_cate,Is_menu :true, Order:1.0,Is_open:true, Methods:"*:SysNode"})

	_ctl.routes = append( _ctl.routes,
	r.Route{Path:"/node/add",Name:"区域信息添加",Category:_cate,Is_menu :false, Order:1.1,Is_open:true, Methods:"*:SysNodeAdd"})

	_ctl.routes = append( _ctl.routes,
	r.Route{Path:"/node/update",Name:"区域信息改变",Category:_cate,Is_menu :false, Order:1.2,Is_open:true, Methods:"*:SysNodeUpdate"})

	_ctl.routes = append( _ctl.routes,
	r.Route{Path:"/node/delete",Name:"区域信息删除",Category:_cate,Is_menu :false, Order:1.3,Is_open:true, Methods:"*:SysNodeDel"})

	
 	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/opr",Name:"操作员管理",Category:_cate,Is_menu :true, Order:1.9,Is_open:true, Methods:"*:SysOpr"})
	
 	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/opr/add",Name:"增加操作员",Category:_cate,Is_menu :false, Order:2.0,Is_open:true, Methods:"*:SysAddOpr"})

	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/opr/update",Name:"修改操作员",Category:_cate,Is_menu :false, Order:2.01,Is_open:true, Methods:"*:SysUpdateOpr"})

	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/opr/delete",Name:"删除操作员",Category:_cate,Is_menu :false, Order:2.02,Is_open:true, Methods:"*:SysDeleteOpr"})

	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/opr/changepassword",Name:"操作员改密码",Category:_cate,Is_menu :false, Order:2.03,Is_open:false, Methods:"*:SysChangePassword"})
	
 	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/product",Name:"资费信息管理",Category:_cate,Is_menu :true, Order:2.1,Is_open:true, Methods:"*:SysProduct"})

 	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/product/add",Name:"增加资费信息",Category:_cate,Is_menu: false, Order:2.2,Is_open:true, Methods:"*:SysAddProduct"})
 	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/product/update",Name:"更新资费信息",Category:_cate,Is_menu: false, Order:2.3,Is_open:true, Methods:"*:SysUpdateProduct"})	
 	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/product/delete",Name:"删除资费信息",Category:_cate,Is_menu: false, Order:2.4,Is_open:true, Methods:"*:SysDelProduct"})

	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/product/detail",Name:"资费信息详细",Category:_cate,Is_menu: false, Order:2.5,Is_open:true, Methods:"*:SysProductDetail"})
	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/product/attr/add",Name:"添加资费属性",Category:_cate,Is_menu: false, Order:2.6,Is_open:true, Methods:"*:SysProductAttrAdd"})	

	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/product/attr/update",Name:"修改自定资费属性",Category:_cate,Is_menu: false, Order:2.7,Is_open:true, Methods:"*:SysProductAttrUpdate"})	

	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/product/attr/delete",Name:"删除自定资费属性",Category:_cate,Is_menu: false, Order:2.8,Is_open:true, Methods:"*:SysProductAttrDelete"})	
	
	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/bas",Name:"BAS节点管理",Category:_cate,Is_menu: true, Order:2.9,Is_open:true, Methods:"*:SysBas"})	

	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/bas/add",Name:"BAS节点添加",Category:_cate,Is_menu: false, Order:3.0,Is_open:true, Methods:"*:SysBasAdd"})	

	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/bas/update",Name:"BAS节点更新",Category:_cate,Is_menu: false, Order:3.1,Is_open:true, Methods:"*:SysBasUpdate"})
	
	_ctl.routes = append( _ctl.routes,
		r.Route{Path:"/bas/delete",Name:"BAS节点删除",Category:_cate,Is_menu: false, Order:3.2,Is_open:true, Methods:"*:SysBasDelete"})	
	
	
// 	_ctl.routes = append( _ctl.routes,
//		r.Route{Path:"/blacklist",Name:"黑白名单管理",Category:_cate,Is_menu :false, Order:2.2,Is_open:true, Methods:"*:SysBlacklist"})

	
	_ctl.AddRoutes()

	
	
}

func (this *SysController) NodeUpdateForm( title string,action string) {
	this.Forms["node_update_form"] = models.InfoForm(title,action,
		models.TextBox(&models.Input{Name:"Name",Id:"node_name",Description:"区域名称",Required:true,Valid:models.Is_not_empty}),
		models.TextBox(&models.Input{Name:"Desc",Id:"node_desc",Description:"区域描述",Required:true,Valid:models.Is_not_empty}),
		models.Hidden(&models.Input{Name:"Id",Id:"node_id",Description:"node id"}),
		models.Submit(&models.Input{Name:"submit",Description:"",Value:"Submit",Class:"btn btn-info"}),	
	)

	
}



//把this中的routes 放到 routers.Permits 中,每个Controller写一遍
func (this *SysController) AddRoutes() {
	
	for i,v := range  this.routes {
		if v.Methods != "" {
			beego.Router(v.Path, this, v.Methods)
		}else{
			beego.Router(v.Path, this)
		}
		
		//if v.Is_menu {
			//Permits.routes is map,so no confict path key ! 
			r.Permits.Add_route(v.Path,&this.routes[i])
		//}
	}
	
}

//主要做Form的要关准备
func (this *SysController) GuuPrepare(){

	this.TplName = libs.GetTplName(this)
	this.Data["STATUS"] = Opr_status
	this.Data["ProductPolicys"]= ProductPolicys
	this.Data["msg"] = ""
}


//---------------------------------------------------------------------------
func (this *SysController) OprList() []Operators {
	var opr []Operators
	rdb.DataBase().SkipGet2(&opr,0,1000)
	
	return opr
}

func (this *SysController) SysOpr () {
	this.TplName = "sys_opr_list.tpl"

	this.Data["Opr_list"]  = this.OprList()
	this.Data["AllOprTypes"] = Opr_type
	this.Data["AllOprStatus"] = Opr_status

	this.Render()
}


func (this *SysController) SysAddOprForm(title string, action string) *models.Form {
	nodes := this.NodeList()
	pdus  := this.ProductList()
	var opr_status [][]string
	
	allnodes := this.Items(nodes,[]string{"Id","Name"})
	allproducts := this.Items(pdus,[]string{"Id","Name"})

	for i:=0; i < len(Opr_status); i++ {
		v := Opr_status[i]
		a := []string{strconv.Itoa(i),v}
		
		opr_status = append(opr_status,a)
	}
	
	form:= models.InfoForm(title,action,
		models.TextBox(&models.Input{Name:"Name",Description:"操作员名称",Required:true,Valid:models.Len_of(2,32)}),
		models.Password(&models.Input{Name:"Pass",Description:"操作员密码",Required:true,Valid:models.Len_of(6,128)}),
		models.Dropdown(&models.Select{Name:"Status",Description:"操作员状态",Args:opr_status}),
		models.GroupDropdown(&models.Select{Name:"Nodes",Description:"关联区域(多选)",Size:4,Args:allnodes}),
		models.GroupDropdown(&models.Select{Name:"Products",Description:"关联资费(多选)",Size:6,Args:allproducts}),
		models.TextArea(&models.Input{Name:"Desc",Description:"操作员描述" }),
		models.Hidden(&models.Input{Name:"Id",Id:"opr_id",Description:"operator id"}),
		models.Submit(&models.Input{Name:"Submit",Value:"提交"}))

	return form
}

func (this *SysController) CheckInRules(path string ) string {
	for i:=0;i<len(this.Rules);i++ {
		if this.Rules[i] == path {
			return "checked"
		}
		
	}
	return ""
}


func (this *SysController) SysAddOpr() {

	f := this.SysAddOprForm("Add","/opr/add")
	this.TplName = "sys_opr_form.html"
	this.Data["AllMenus"] = r.Permits.Menu
	this.Data["CheckInRules"] = this.CheckInRules
	this.Data["CheckOpen"]  = r.Permits.Check_open

	
	if this.Ctx.Input.IsPost() {

		if this.Validator2(f) == false {
			this.Data["Form"] = f
			this.Render()
			return
		}
		opera := &Operators{}
		this.ParsePostToStruct(opera)
		fmt.Println(opera)
		opera.Type = 1 // 只允许一个超级管理员
		
		_,cnt := rdb.DataBase().FilterInsert(opera,"Name")
		if cnt  == 0 {
			
			this.Redirect("/opr", 302)
		}else {
			
			this.ShowTips("名称有重复 "+strconv.Itoa(cnt) +"个" )
		}

		this.Render()
		return
		//fmt.Println(this.GetStrings("Products"))
		
	}


	this.Data["Form"] = f
	this.Render()
	return
	
}

func (this *SysController) SysUpdateOpr() {
	
	opr_id := this.GetString("opr_id")
	
	this.TplName = "sys_opr_form.html"
	this.Data["AllMenus"] = r.Permits.Menu
	this.Data["CheckInRules"] = this.CheckInRules
	this.Data["CheckOpen"]  = r.Permits.Check_open
	
	f := this.SysAddOprForm("Update","/opr/update")
	
	if this.Ctx.Input.IsPost(){
		
		if this.Validator2(f) == false {
			this.Data["Form"] = f
			this.Render()
			return
		}
		
		fmt.Println("Products:", this.GetStrings("Products"))
		one := &Operators{}
		this.ParsePostToStruct(one)
		one.Type = 1
		fmt.Println(one.RuleItem)
		resp,err := rdb.DataBase().Update(this.GetString("Id"),one)
		if err == nil {
			fmt.Println("Replaced ",resp.Replaced )
		}else {
			fmt.Println(err)
		}

		this.Redirect("/opr",302)		
		return
	}
	
	if opr_id == "" {
		this.Abort("403")
		return
	}

	
	one := &Operators{}	
	err := rdb.DataBase().QuOne(one,opr_id)
	if err == nil {
		
		this.Rules = one.RuleItem		
		f.FillFormFromStruct(one)
		this.Data["Form"]  = f
		
	}

	
	this.Render()

	return
}

func (this *SysController) SysDeleteOpr() {

	opr_id := this.GetString("opr_id")

	if opr_id == "" {
		this.Abort("403")
		return
	}

	this.IdRowDelete(Operators{},opr_id,"/opr")
	
	this.Render()
	return
}


func (this *SysController) SysChangePasswordForm() *models.Form {
	
	f := models.InfoForm("改密","/opr/changepassword",
		models.TextBox(&models.Input{Name:"operator_name",Description:"管理员名",Value:this.GetCookie("username"),Size:32,ReadOnly:true,Required:true}),
		models.Password(&models.Input{Name:"operator_pass",Description:"管理员新密码",Value:"",Required:true,Size:32,Valid:models.Len_of(6,32)}),
		models.Password(&models.Input{Name:"operator_pass_chk",Description:"确认管理员新密码",Required:true,Size:32,Valid:models.Len_of(6,32)}),
		models.Submit(&models.Input{Name:"Submit",Value:"<b>提交</b>"}))
	return f
}

func (this *SysController) SysChangePassword() {
	f := this.SysChangePasswordForm()
	this.Data["msg"] = ""
	if this.Ctx.Input.IsPost() {
		if this.Validator2(f) == false {
			this.Data["Form"] = f.Render()
			this.Render()
			return
		}

		pass1 := this.GetString("operator_pass")
		pass2 := this.GetString("operator_pass_chk")
		if pass1 != pass2 {
			this.Data["msg"] = "密码不一致"
			this.Data["Form"] = f.Render()
			this.Render()
			return
		}

		opera := &Operators{}
		err := rdb.DataBase().FilterOne(opera,map[string]string{"Name":this.GetCookie("username")})
		opera.Pass = pass1
		
		resp,err := rdb.DataBase().Update(opera.Id,opera)
		if err == nil {
			fmt.Println("Replaced ",resp.Replaced )
			
		}else {
			fmt.Println(err)
		}

		this.Redirect("/",302)
		
	}
	
	this.Data["Form"] = f.Render()

	this.Render()
}

func (this *SysController) SysNode() {
	//区域的删除时要检查是否有用户在此区域下,否则不能删除区域
	this.TplName = "sys_node_list.tpl"
	this.Data["Nodes"] = this.NodeList()

	this.Render()
}

func (this *SysController) SysNodeAdd(){

	this.NodeUpdateForm("Add","/node/add")	
	if this.Ctx.Input.IsPost() {

		if this.Validator("node_update_form") == false {
			this.Data["Form"] = this.Forms["product_add_form"].Render()
			this.Render()
			return
		}
		
		//_,cnt := rdb.DataBase().FilterInsert(Node{},map[string]string{"Name":this.GetString("Name"),"Desc":this.GetString("Desc")},"Name")
//		_,cnt := rdb.DataBase().FilterInsert(Node{},&Node{Name:this.GetString("Name"),Desc:this.GetString("Desc")},"Name")
		_,cnt := rdb.DataBase().FilterInsert(&Node{Name:this.GetString("Name"),Desc:this.GetString("Desc")},"Name")		
		if cnt  == 0 {
			
			this.Redirect("/node", 302)
		}else {

			this.ShowTips("名称重复 "+strconv.Itoa(cnt) +"个" )
		}

		
		/*
		var shval = map[string]interface{}{"Name":this.GetString("Name")}
		
		cnt := rdb.DataBase().FilterCount(Node{},shval)
		fmt.Println(cnt)
		if cnt == 0 {
			resp,err := rdb.DataBase().InsertQ(Node{Name:this.GetString("Name"),Desc:this.GetString("Desc")})
			if err == nil {
				fmt.Println(resp.GeneratedKeys)
			}
			
			this.Data["Form"] = this.Toast("","OK")
		}else {
			this.Data["Form"] = "名称重复 "+strconv.Itoa(cnt) +"个"
	        }
*/
		this.Render()
		return
	}
	
	
	this.Data["Form"] = this.Forms["node_update_form"].Render()
	this.Render()
}

func (this *SysController) SysNodeUpdate() {

	if this.Ctx.Input.IsPost() {

		node := &Node{}
		this.ParsePostToStruct(node)
		resp,err := rdb.DataBase().Update(this.GetString("Id"),node)
		if err == nil {
			fmt.Println("Replaced ",resp.Replaced )
		}else {
			fmt.Println(err)
		}

		this.Redirect("/node",302)
//		this.Data["Form"] = "Updated"
//		this.Render()
		return
	}

	
	id := this.GetString("node_id")
	one := &Node{}
	err := rdb.DataBase().QuOne(one,id)
	if err == nil {
		this.NodeUpdateForm("Update","/node/update")
		fmt.Println(one)
		this.Forms["node_update_form"].Fill(one)		
	}
	
	this.Data["Form"] = this.Forms["node_update_form"].Render()
	
	this.Render()
}

func (this *SysController) SysNodeDel() {
	//删除节点区域之前要先查一下此区域下是否有用户
	node_id := this.GetString("node_id")
	if node_id == "" {
		this.Abort("403")
		return
	}
	
	_,err := rdb.DataBase().Del(Node{},this.GetString("node_id"))
	if err == nil {
		this.Redirect("/node",302)
	}else {
		this.ShowTips( err )
	}


	this.Render()
}


func (this *SysController) product_add_form(title string, action string) {
	var policy_arr  [][]string
	for i:=0; i < len(ProductPolicys); i++ {
		v := ProductPolicys[i]
		a := []string{strconv.Itoa(i),v}
		fmt.Println(a)
		policy_arr = append(policy_arr,a)
	}
	
	var boolean_items [][]string
	for i,v := range YESNO{
		a := []string{ strconv.Itoa(i),v}
		boolean_items = append(boolean_items,a)
	}
	
	this.Forms["product_add_form"] = models.InfoForm(title,action,
		models.TextBox(&models.Input{Name:"Name",Description:"资费名称",Required:true,Valid:models.Len_of(4,64)}),
		models.Dropdown(&models.Select{Name:"Policy",Description:"计费策略",Args:policy_arr}),
		models.TextBox(&models.Input{Name:"FeeMonths",Description:"买断授权月数"}),
		models.TextBox(&models.Input{Name:"FeeTimes",Description:"买断时长(小时)",Valid:models.Is_number3,Value:"0"}),
		models.TextBox(&models.Input{Name:"FeeFlows",Description:"买断流量(MB)",Valid:models.Is_number3,Value:"0"}),
		
		models.TextBox(&models.Input{Name:"FeePrice",Description:"资费价格(元)",Required:true,Valid:models.Is_rmb}),
		models.TextBox(&models.Input{Name:"FeePeriod",Description:"开放认证时段",Valid:models.Is_period}),
		models.TextBox(&models.Input{Name:"ConcurNumber",Description:"并发数控制(0表示不限制)",Value:"0",Valid:models.Is_numberOboveZore}),
		models.Dropdown(&models.Select{Name:"BindMac",Description:"是否绑定Mac",Args:boolean_items}),
		models.Dropdown(&models.Select{Name:"BindVlan",Description:"是否绑定VLAN",Args:boolean_items}),
		models.TextBox(&models.Input{Name:"InputMaxLimit",Description:"最大上行速率(Mbps)",Required:true,Valid:models.Is_number3}),		
		models.TextBox(&models.Input{Name:"OutputMaxLimit",Description:"最大下行速率(Mbps)",Required:true,Valid:models.Is_number3}),
		models.Hidden(&models.Input{Name:"Id",Description:"pdu uuid "}),
		models.Submit(&models.Input{Name:"submit",Id:"pdu_submit",Description:"",Value:"提交",Class:"btn btn-info"})    )
		
}

func (this *SysController) SysProduct() {

	
	if this.Ctx.Input.IsPost() {

		this.Render()
		return
	}
	this.TplName = "sys_product_list.html"
	this.Data["Products"] = this.ProductList()
	this.Render()
}


func (this *SysController) SysAddProduct () {

	this.product_add_form("Add","/product/add")
	this.TplName = "sys_product_form.html"
	
	if this.Ctx.Input.IsPost() {
		if this.Validator("product_add_form") == false {
			//this.FillFromPost()
			this.Data["Form"] = this.Forms["product_add_form"].Render()
			
		}else{
			product := &Products{}
			this.ParsePostToStruct(product)
			fmt.Println(product)

			product.FeeTimes = libs.Hour2sec(product.FeeTimes)
			product.FeeFlows = libs.Mb2kb(product.FeeFlows)
			product.FeePrice = libs.Yuan2fen(product.FeePrice)
			product.InputMaxLimit = libs.Mbps2bps(product.InputMaxLimit)
			product.OutputMaxLimit = libs.Mbps2bps(product.OutputMaxLimit)
			_datetime := libs.Get_currtime()
			product.CreateTime = _datetime
			product.UpdateTime = _datetime
			
			_,cnt := rdb.DataBase().FilterInsert(product,"Name")
			if cnt  == 0 {
				
				this.Redirect("/product", 302)
			}else {
				
				this.ShowTips("名称有重复 "+strconv.Itoa(cnt) +"个" )
			}
			
		}
		
		this.Render()
		return
	}
	
	this.Data["Form"]  = this.Forms["product_add_form"].Render()	
	this.Render()
}

func (this *SysController) SysUpdateProduct() {
	
}

func (this *SysController) SysDelProduct() {
	id := this.GetString("product_id")
	if id == "" {
		this.Abort("403")
		return
	}
	
	_,err := rdb.DataBase().Del("products",id)
	if err == nil {
		this.Redirect("/product",302)
	}else {
		this.ShowTips( err )
	}


	this.Render()	
}


func (this *SysController) SysProductDetail() {
	product_id := this.GetString("product_id")
	this.TplName = "sys_product_detail.html"

	pdus := this.ProductList()
	
	this.Data["YESNO"] = YESNO
	this.Data["OprStatus"] = Opr_status
	
	product := &Products{}
	var product_attrs []ProductAttr
	
	err := rdb.DataBase().QuOne(product,product_id)
	if err == nil {
		this.Data["Product"]      = product
		this.Data["ProductAttrs"] = product_attrs
		this.Data["PolicyMap"]    = this.ToPairMapI(pdus,[]string{"Id","Policy"})
	}

	this.Render()
}

func (this *SysController) SysProductAttrAdd() {
}

func (this *SysController) SysProductAttrUpdate() {
	
}

func (this *SysController) SysProductAttrDelete() {
	
}


func (this *SysController) SysBas() {
	this.TplName = "sys_bas_list.html"
	this.Data["TimeTypeMap"] = TimeTypeMap
	this.Data["BasVendorType"] = BasVendorTypeMap
	
	this.Data["BasList"] = this.BasList()
	this.Render()
}

func (this *SysController) SysBasAddForm(bastypes [][]string,timetypes [][]string  ) *models.Form {
	f := models.InfoForm("增加BAS","/bas/add",
		models.TextBox(&models.Input{Name:"IpAddr",Valid:models.Is_ip,Description:"BAS地址",Required:true}),
		models.TextBox(&models.Input{Name:"Name",Valid:models.Len_of(2,64),Description:"BAS名称",Required:true}),
		models.TextBox(&models.Input{Name:"Secret",Valid:models.Is_alphanum2(4,32),Description:"共享秘钥",Required:true}),
		models.Dropdown(&models.Select{Name:"VendorId",Description:"BAS类型",Args:bastypes,Required:true}),
		models.TextBox(&models.Input{Name:"CoaPort",Description:"CoA端口",Default:"3799"}),
		models.Dropdown(&models.Select{Name:"TimeType",Description:"时间类型",Args:timetypes,Required:true}),
		models.Submit(&models.Input{Name:"submit",Value:"<b>提交</b>"}))

	return f
}

func (this *SysController) SysBasAdd() {
	
	keys := []int{}
	for k,_ := range BasVendorTypeMap {
		keys = append(keys,k)
	}
	sort.Asc(keys)
	
	var bastypes [][]string
	for _,k := range keys {
		
		a := []string{strconv.Itoa(k), BasVendorTypeMap[k] }
		bastypes = append(bastypes,a)
	}
	
	var timetypes [][]string
	
	for i,v := range TimeTypeMap {
		a := []string{strconv.Itoa(i),v}
		timetypes = append(timetypes,a)
	}
	
	f := this.SysBasAddForm(bastypes,timetypes)

	if this.Ctx.Input.IsPost() {

		if this.Validator2(f) == false {
			//this.FillFromPost()
			this.Data["Form"] = f.Render()
			this.Render()
			return
			
		}else {
			one := &Bas{}
			this.ParsePostToStruct(one)
			
			if one.CoaPort == 0 {
				one.CoaPort = 3799
			}
			
			_,cnt := rdb.DataBase().FilterInsert(one,"IpAddr")
			if cnt  == 0 {
				this.Redirect("/bas", 302)
			}else {
				
				this.ShowTips("Bas地址有重复 "+strconv.Itoa(cnt) +"个" )
				return
			}
			
		}
			
	}

	
	this.Data["Form"]  = f.Render()
	
	this.Render()
}

func (this *SysController) SysBasUpdateForm( bastypes [][]string,timetypes [][]string ) *models.Form {
	
	f := models.InfoForm("更新BAS","/bas/update",
		models.TextBox(&models.Input{Name:"IpAddr",Valid:models.Is_ip,Description:"BAS地址",Required:true}),
		models.TextBox(&models.Input{Name:"Name",Valid:models.Len_of(2,64),Description:"BAS名称",Required:true}),
		models.TextBox(&models.Input{Name:"Secret",Valid:models.Is_alphanum2(4,32),Description:"共享秘钥",Required:true}),
		models.Dropdown(&models.Select{Name:"VendorId",Description:"BAS类型",Args:bastypes,Required:true}),
		models.TextBox(&models.Input{Name:"CoaPort",Description:"CoA端口",Default:"3799"}),
		models.Dropdown(&models.Select{Name:"TimeType",Description:"时间类型",Args:timetypes,Required:true}),
		models.Hidden(&models.Input{Name:"Id",Description:"编号"}),
		models.Submit(&models.Input{Name:"submit",Value:"<b>提交</b>"}))

	return f
}

func (this *SysController) SysBasUpdate( ) {

	var id string
	if this.Ctx.Input.IsPost() {
		id = this.GetString("Id")
		
	}else {
		id = this.GetString("bas_id")
	}
	if id == "" {
		this.Abort("403")
		return
	}
	keys := []int{}
	for k,_ := range BasVendorTypeMap {
		keys = append(keys,k)
	}
	sort.Asc(keys)
	
	var bastypes [][]string
	for _,k := range keys {
		
		a := []string{strconv.Itoa(k), BasVendorTypeMap[k] }
		bastypes = append(bastypes,a)
	}
	
	var timetypes [][]string
	
	for i,v := range TimeTypeMap {
		a := []string{strconv.Itoa(i),v}
		timetypes = append(timetypes,a)
	}

	f:=this.SysBasUpdateForm(bastypes,timetypes )
	
	if this.Ctx.Input.IsPost() {
		if this.Validator2(f) == false {
			this.Data["Form"] = f.Render()
			this.Render()
			return
		}else {
			one := &Bas{}
			this.ParsePostToStruct(one)
			resp,err := rdb.DataBase().Update(id,one)
			if err == nil {
				fmt.Println("Replaced ",resp.Replaced )
			}else {
				fmt.Println(err)
			}
			
			this.Redirect("/bas",302)
			return
		}
		
	}
	
	one := &Bas{}
	err := rdb.DataBase().QuOne(one,id)
	if err == nil {
		f.FillFormFromStruct(one)
		this.Data["Form"] = f.Render()
	}
	
	this.Render()
	
}

func (this *SysController) SysBasDelete() {
	id := this.GetString("bas_id")
	if id == "" {
		this.Abort("403")
		return
	}
	
	_,err := rdb.DataBase().Del("bas",id)
	if err == nil {
		this.Redirect("/bas",302)
	}else {
		this.ShowTips( err )
	}


	this.Render()		
}
