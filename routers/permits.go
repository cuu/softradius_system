package routers

import (
//	"github.com/cuu/softradius/controllers"
	"github.com/astaxie/beego"
//	"github.com/cuu/softradius/models"
	"github.com/cuu/softradius/libs"
	sort "github.com/cuu/softradius/libs/sortutil"
	
	"fmt"
)

//大写单元名, 可以被外部包调用此struct
type Route struct{
	Path string
	Name string
	Category string
	Is_menu bool
	Order float32
	Is_open bool
	Oprs []string
	Methods string `method,eg *:Index`
}

type Menu struct {
	Category string
	Items []Route
}

type Permit struct {
	routes map[string]*Route
	Menu []Menu
}

var Permits *Permit

// FEES, PP=>Pre Pay, BO => Buy out
const (
	PPMonth = iota // 0 
	PPTimes        // 1
	BOMonth        // 2
	BOTimes        // 3
	PPFlow         // 4
	BOFlows        // 5
	AwesomeFee     // 6
	AwesomeFeeBoTime  // 7
)

//ACCOUNT_STATUS
const (
	UsrPreAuth = iota
	UsrNormal
	UsrPause
	UsrCancel
	UsrExpire
)

//CARD_STATUS
const (
	CardInActive = iota
	CardActive
	CardUsed
	CardRecover
)

//CARD_TYPE
const (
	ProductCard =iota
	BalanceCard
)

var ACCEPT_TYPES = map[string]string{
	"open":"开户",
	"pause":"停机",
	"resume":"复机",
	"cancel":"销户",
	"next":"续费",
	"charge":"充值",
	"change":"变更",
}


const (
	MenuSys    = "系统管理" // 0
	MenuBus    = "营业管理"
	MenuOpt    = "维护管理"
	MenuStat   = "统计分板"
	MenuWlan   = "Wlan管理"
	MenuMpp    = "微信接入"
	MenuPlugin = "插件管理"
	MenuAgency = "代理管理" // 8
	MenuSysMenu= "系统菜单" // 9
)

/*
var MenuCats = [8]string{
	MenuSys,MenuBus,MenuOpt,
	MenuStat,MenuWlan,MenuMpp,
	MenuPlugin,MenuAgency,
}
*/
var MENU_ICONS = map[string]string{
	MenuSys:"fa fa-cog",
	MenuBus:"fa fa-user",
	MenuOpt:"fa fa-wrench",
	MenuStat:"fa fa-bar-chart",
	MenuPlugin:"fa fa-bar-chart",
	MenuAgency:"fa fa-hand-spock-o",
}


const MAX_EXPIRE_DATE = "3000-12-30"


/*
e(sh string, arr []string) []string{
	for i,v := range arr {
		if v == sh {
			return append( arr[:i],arr[i+1:]... )
		}
	}
}
*/

func init() {

	Permits = &Permit{ routes:make(map[string]*Route)}

	beego.AddFuncMap("Fen2yuan",libs.Fen2yuan)
	beego.AddFuncMap("Bps2mbps",libs.Bps2mbps)
	beego.AddFuncMap("Kb2mb",   libs.Kb2mb)
	beego.AddFuncMap("Sec2hour",libs.Sec2hour)
	
	beego.AddFuncMap("In",      libs.In)
}

func (r *Route) GetOprs() []string {
	return r.Oprs
}

func (p *Permit) Add_route( path string ,arg *Route) {
	if path == "" { return }
	arg.Path = path
	p.routes[path] = arg
}

func (p *Permit) Get_route(path string) (*Route ,bool) {
	r,ok := p.routes[path]
	return r,ok
}


func (p *Permit) Bind_super(opr string) {
	
	for k,_ := range  p.routes {
		p.routes[k].Oprs = append(p.routes[k].Oprs,opr)
	}
	
	
}


func (p *Permit) Bind_opr(opr string,path string){
	if path =="" { return }

	if _, ok := p.routes[path]; ok {
		oprs := p.routes[path].Oprs
		for _,o := range oprs {
			if o == opr {
				return
			}
		}
		p.routes[path].Oprs = append(p.routes[path].Oprs,opr)
	}
}

func (p *Permit) Unbind_opr(opr string, path ...string){
	if len(path) > 0 {
		p.routes[ path[0] ].Oprs = libs.Slice_rm(opr,p.routes[path[0]].Oprs )
		
	}else
	{
		//通杀
		for i,_ := range p.routes {
			p.routes[i].Oprs = libs.Slice_rm(opr,p.routes[i].Oprs)
		}
	}
}

func (p *Permit) Check_open(path string) bool{
	return p.routes[path].Is_open
}

func (p *Permit) Check_opr_category(opr string, category string)  bool {
	
	for _,v := range p.routes {
		if v.Category == category {
			for _,o := range v.Oprs {
				if o == opr {
					return true
				}
			}
			return false
		}
	}

	return false
}

func (p *Permit) Build_menus( order_cats []string) {
//	menu := []Menu{}
	for _,v := range order_cats {
		p.Menu = append(p.Menu,Menu{ Category:v } )
	}

	//fmt.Println(menu)
	for k,v := range p.routes {
		for i,m := range p.Menu {
			if v.Category == m.Category {
				fmt.Println("Build_menus:  ",k,v)
				p.Menu[i].Items = append(p.Menu[i].Items,*v)
				
				break
			}
		}
	}

	for idx,_  := range p.Menu {
		sort.AscByField(p.Menu[idx].Items, "Order")
	}
		
	return 
}

func (p *Permit) Match(username string,  path string ) bool {
	if _,ok := p.routes[path]; ok {
		for _,v := range p.routes[path].Oprs {
			if v == username {
				return true
			}
		}
	}
	return false
}
