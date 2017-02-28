package controllers

import (
	"github.com/astaxie/beego"
	"github.com/cuu/softradius/models"
	r "github.com/cuu/softradius/routers"
	"fmt"
	"reflect"
	"strings"
	"strconv"
	"github.com/cuu/softradius/libs"
	rdb "github.com/cuu/softradius/database/shelf"
)

var YESNO = map[int]string{ 0:"否", 1:"是", }
///0级管理员只有一个,1级有很多,2级是给代理商的管理员,2级管理员没有系统管理权限
var Opr_type = map[int]string{ 0: "系统管理员", 1: "普通操作员",2:"代理商操作员"}

var Opr_status = map[int]string{0: "正常", 1: "停用"}


type GuuPreparer interface {
        GuuPrepare()
}

type GuuRender interface {
	GuuRender()
}

type BaseController struct {
	beego.Controller
	Forms  map[string]*models.Form
	routes []r.Route
	Secret string /// secret for secure cookie
	
	
}


func (this *BaseController) Inactive(m r.Menu) string {

	var url = this.Ctx.Input.URL()
	if strings.HasSuffix(url,"/"){
		url = url[:len(url)-1]
	}
	
	for _,v := range m.Items {
		if url == v.Path {
			return "active"
		}
	}
	
	return ""
}

func (this *BaseController) AClass(path string ) string {

	var url = this.Ctx.Input.URL()
	if strings.HasSuffix(url,"/"){
		url = url[:len(url)-1]
	}
	
	if url == path {
		return "active"
	}
	
	return ""
}

func (this *BaseController) GetURLName(){
	url := libs.RemoveSuffix(this.Ctx.Input.URL(),"/")
	r,ok :=  r.Permits.Get_route(url)
	if ok {
		
		this.Data["URLName"] = r.Name
	}
}


func (this *BaseController) Prepare(){
	this.Data["adminlte"] = "/AdminLTE"
	
	this.Data["Menu"] = r.Permits.Menu
	this.Data["MenuIcon"] = r.MENU_ICONS
	this.Data["Inactive"] = this.Inactive
	this.Data["AClass"]   = this.AClass
	this.Data["URLName"] = "..."

	this.Data["Match"] = this.MatchOpr
	this.Data["GetCookie"] = this.GetCookie
	this.Data["CheckOprCate"] = this.CheckOprCategory

	this.Data["YesOrNo"] = YESNO
	
	this.Layout = "layout_base.tpl"     
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["Sidebar"] = "sidebar.tpl"
	this.LayoutSections["Header"]  = "header.tpl"
	this.LayoutSections["Footer"]  = "footer.tpl"
	this.LayoutSections["ContentHeader"] = "content-header.tpl"
	this.LayoutSections["HeadCss"] = "head_css.tpl"

	this.Forms = make(map[string]*models.Form)
	
        if app, ok := this.AppController.(GuuPreparer); ok {
                app.GuuPrepare()
        }
	
	this.Secret = beego.AppConfig.DefaultString("DEFAULT::Secret","FUCK")
	
	this.AuthOpr()

	this.GetURLName()
	
}

func (c *BaseController) Render() error {
	
	if !c.EnableRender {
		return nil
	}
	
	rb, err := c.RenderBytes()
	if err != nil {
		return err
	}
	
	c.Ctx.Output.Header("Content-Type", "text/html; charset=utf-8")
	return c.Ctx.Output.Body(rb)
}


func (this *BaseController) Validator(key string ) bool{
	//valids self form
	is_valid := true
	for _, cld := range this.Forms[key].Children{
//	for _, cld := range form.Children{		
		_vcld := reflect.ValueOf(cld)
		_cld := _vcld.Elem()
		m := libs.Getdict(cld) // struct to map
		//fmt.Println(_cld)
		//fmt.Println(reflect.ValueOf(cld))
		_form_val := this.GetString( m["Name"] )
		fmt.Println("Validator form value: ", _form_val)
		if len(_form_val) > 0{
			_vcld.MethodByName("SetValue").Call([]reflect.Value{reflect.ValueOf(_form_val) })
		}
		
		_cvf := _cld.FieldByName("Valid")		
		if _cvf.IsValid() == true{

			
			
			if  _cvf.Elem().IsValid() &&  len(_cvf.Elem().FieldByName("Msg").String()) > 0{
				_ret :=_cvf.MethodByName("MatchString").Call([]reflect.Value{reflect.ValueOf(this.GetString(m["Name"]))})
				//fmt.Println(_cvf.Elem().FieldByName("msg"))
				//fmt.Println(this.GetString(m["Name"]))
				if _ret[0].Bool() != true {
					_note := _cld.FieldByName("Note")
					if _note.CanSet(){
						_note.SetString(_cvf.Elem().FieldByName("Msg").String())
					}
					is_valid = _ret[0].Bool()
				}
				//break
			}
			
		}
	}

	return is_valid
}

func (this *BaseController) Validator2(form *models.Form ) bool{
	//valids self form
	is_valid := true
//	for _, cld := range this.Forms[key].Children{
	for _, cld := range form.Children{		
		_vcld := reflect.ValueOf(cld)
		_cld := _vcld.Elem()
		m := libs.Getdict(cld) // struct to map
		//fmt.Println(_cld)
		//fmt.Println(reflect.ValueOf(cld))
		_form_val := this.GetString( m["Name"] )
		fmt.Println("Validator form value: ", _form_val)
		if len(_form_val) > 0{
			_vcld.MethodByName("SetValue").Call([]reflect.Value{reflect.ValueOf(_form_val) })
		}
		
		_cvf := _cld.FieldByName("Valid")		
		if _cvf.IsValid() == true{

			
			
			if  _cvf.Elem().IsValid() &&  len(_cvf.Elem().FieldByName("Msg").String()) > 0{
				_ret :=_cvf.MethodByName("MatchString").Call([]reflect.Value{reflect.ValueOf(this.GetString(m["Name"]))})
				//fmt.Println(_cvf.Elem().FieldByName("msg"))
				//fmt.Println(this.GetString(m["Name"]))
				if _ret[0].Bool() != true {
					_note := _cld.FieldByName("Note")
					if _note.CanSet(){
						_note.SetString(_cvf.Elem().FieldByName("Msg").String())
					}
					is_valid = _ret[0].Bool()
				}
				//break
			}
			
		}
	}

	return is_valid
}



func (this *BaseController) GetCookie(key string) string {
//	fmt.Println("GetCookie ",this.Ctx.Input.Cookie(key))
	//	return this.Ctx.GetCookie(key)
	return this.Ctx.Input.Cookie(key)
}

func (this *BaseController) SetCookie( name string, value string, others ...interface{}) {
	this.Ctx.Output.Cookie(name,value,others...)
}

//Tr用的是 SecureCookie
func (this *BaseController) GetSecCookie( key string) (string ,bool){
	return this.GetSecureCookie(this.Secret,key)
}

func (this *BaseController) SetSecCookie(name,value string, others ...interface{}){
	this.SetSecureCookie(this.Secret,name,value,others...)
}


func (this *BaseController) MatchOpr(path string ) bool{
	return r.Permits.Match( this.GetCookie("username"), path )
}

func (this *BaseController) CheckOprCategory(cat string) bool {
	return r.Permits.Check_opr_category(this.GetCookie("username"),cat)
}

//把this中的routes 放到 routers.Permits 中


/// 略等于tr的auth_opr,验证当前的管理员状态
func (this *BaseController) AuthOpr() (result bool) {
	//this.Ctx.Input.IsPost(),URL
	// $GOCODE/src/github.com/axxx/beego/context/input.go
	if this.GetCookie("username") == "" {
		if strings.HasPrefix(this.Ctx.Input.URL(), "/login") != true {
			
			this.Redirect("/login",302)
		}else { fmt.Println("on login page now") }
		
	}else{
		opr := this.GetCookie("username")
		fmt.Println(this.Ctx.Input.URL())
		rule,ok := r.Permits.Get_route(this.Ctx.Input.URL())
		if ok {
			for _,v := range rule.GetOprs() {
				if v == opr {
					result = true
					break
				}
			}
			
			if result != true {
				this.Abort("403")
			}
			result = false
		}
	}
	return
}

func (this *BaseController) Reset_cookie() {
	this.SetCookie("username","")
	this.SetCookie("opr_type","")
	this.SetCookie("login_time",libs.Get_currtime())
	this.SetCookie("login_ip", this.Get_clientip())
}

func (this *BaseController) Get_clientip() string {
	s := strings.Split(this.Ctx.Request.RemoteAddr, ":")
	return s[0]
}


func (this *BaseController) Toast(title string, msg string) string {
	return models.AlertBox(&models.Alert{Title:title,Type:"Info",Msg:fmt.Sprintf("%s",msg)});
}

func (this *BaseController) ShowTips( msg interface{} ){
	this.Layout = "tips.tpl"
	this.Data["Form"] = fmt.Sprintf("%v",msg)
}


func (this *BaseController) ParsePostToStruct (s interface{} ) {
	var val = reflect.ValueOf(s)
	_val := val.Elem()
	t := _val.Type()
	for i:=0; i < t.NumField();i++ {
		field := t.Field(i)

		//遇到数组,就更简单了,直接是 []string [][]string,声明是啥样子,得到的就是啥样子
		//fmt.Println(  field.Name,field.Type.String()  )
		
		
		thefield := _val.FieldByName( field.Name )
		if(thefield.CanSet()) {
			if( this.GetString(field.Name) != ""){
				fmt.Println(field.Name, field.Type.String())
				if field.Type.String() == "string"{
					thefield.SetString(this.GetString( field.Name))
				}else if field.Type.String() == "int" {
					intval,_ := strconv.Atoi(this.GetString(field.Name) )
					thefield.Set( reflect.ValueOf(intval))
				}else if field.Type.String() == "[]string" {
					fmt.Println(field.Name,":",this.GetStrings(field.Name))
					thefield.Set( reflect.ValueOf( this.GetStrings(field.Name)))
				}
			}
		}
	
		
	}
	
}


//统一删除某数据库表中的对应的Id的数据row

func (this *BaseController) IdRowDelete(table interface{},id string,redirect string) {

	_,err := rdb.DataBase().Del(table,id )
	if err == nil {
		this.Redirect(redirect,302)
	}else {
		this.ShowTips( err )
	}
	
}

func (this *BaseController) BasList() []Bas {
	var nods []Bas
	rdb.DataBase().SkipGet2(&nods,0,1000)
	return nods
}

func (this *BaseController) NodeList() []Node {
	
	var nods []Node
	
	rdb.DataBase().SkipGet2(&nods,0,1000)

	return nods
}


func (this *BaseController) ProductList() []Products {
	var nods []Products
	rdb.DataBase().SkipGet2(&nods,0,1000)
	
	return nods	
}

func (this *BaseController) AgencyList() []Agency {
	var nods []Agency
	rdb.DataBase().SkipGet2(&nods,0,1000)
	
	return nods	
}


func (this *BaseController) Items(v interface{},list []string ) [][]string {

	var ret [][]string
	
	if libs.Type(v) != "slice" {
		fmt.Println("Items for slice")
		return ret
	}
	
	s := reflect.ValueOf(v)
	for i:=0;i<s.Len();i++ {
		iv := s.Index(i)
		//now it's Node
		var a []string
		a = append(a,iv.FieldByName(list[0]).String())
		a = append(a,iv.FieldByName(list[1]).String())

		ret = append(ret,a)
	}

	return ret
}
