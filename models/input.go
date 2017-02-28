package models

import (
	"fmt"
	"strings"
	"reflect"
	"strconv"
	"regexp"
	"github.com/cuu/softradius/libs"
	"github.com/pkg4go/camelcase"
)

type STRINGS struct{
	str []string
}

type Regex struct{
	rexp *regexp.Regexp
	Msg string
}


type Alert struct{
	Title string
	Type  string
	Msg   string
}


var Alert_type = make(map[string]string)

type Input struct {
	Id   string
	Name string
	Class string
	Value string
	Description string `input desc text`
	Type string
	Checked bool
	Required bool
	Note  string `show error text`
	Valid *Regex
	Size int
	ReadOnly bool
	Default string  //此可取Value而代之
	
}

type Select struct{
	Id string
	Name string
	Class string
	Value []string
	Type string
	Multiple bool
	Description string
	Args [][]string
	Size int
	ReadOnly bool
	Required bool
}

type Form struct{
	Title string
	Action string
	Style string `form:primary,warning,info,danger,success`
	Children []interface{}
}

func init(){
	//init some global vars
	Alert_type["Danger"] = "danger"
	Alert_type["Info"]   = "info"
	Alert_type["Warning"] = "warning"
	Alert_type["Success"] = "success"

//	fmt.Println(Alert_type)
	
}


func (this *STRINGS) Append(str string){
	this.str = append(this.str,str)
}

func (this *STRINGS) Remove(str string){
	for i,v := range this.str{
		if v  == str {
			this.str = append(this.str[:i],this.str[(i+1):]...)
			break
		}
	}
}

func (this *STRINGS) String() string{
	return strings.Join(this.str,"")
}

func (self *Alert) Render() string{
	var out STRINGS

	out.Append("<div class=\"box box-default\">")
	out.Append("<div class=\"box-header with-border\">")
	out.Append("<h3 class=\"box-title\">"+self.Title+"</h3>")
	out.Append("</div>")
	out.Append("<div class=\"box-body\">")
	out.Append(self.String())
	out.Append("</div>")
	return out.String()
}

func (self *Alert) String() string{
	if self.Type == ""{
		self.Type = "Danger"
	}
	
	type_str := Alert_type[self.Type]
	
	return fmt.Sprintf("<div class=\"alert alert-%s alert-dismissible\">"+
                "<button type=\"button\" class=\"close\" data-dismiss=\"alert\" aria-hidden=\"true\">×</button>"+
                "<h4><i class=\"icon fa fa-%s\"></i> </h4>"+
		"%s"+
                "</div>",	type_str,type_str, self.Msg)
	
}
// Override String() 
func (self *Input) String() string{
	var out STRINGS
	
	idstr := ""
	cls_str := ""
	name_str := ""
	placehold := ""
	req_str   := ""
	rdo_str   := ""
	size_str  := ""
	
	if self.Id != "" {
		idstr = "id=\""+self.Id+"\""
	}else {
		idstr = "id=\""+self.get_default_id()+"\""
	}
	
	if self.Class != ""{
		cls_str = "class=\"form-control "+self.Class+"\""
	}else{
		cls_str = "class=\"form-control\""
	}
	if self.Name != ""{
		name_str = "name=\""+self.Name+"\""
	}

	if self.Description !=""{
		placehold =self.Description
	}else{
		self.Description=""	
	}

	if self.Valid != nil {
		placehold += ","+self.Valid.Msg
	}

	if self.Value == "" && self.Default != "" {
		placehold = "placeholder=\""+placehold +",默认是"+self.Default+"\""
	}else {
		placehold = "placeholder=\""+placehold +"\""
	}
	
	if self.Required == true{
		req_str = "required=\"required\""
	}
	if self.ReadOnly == true {
		rdo_str = "readonly=\"readonly\""
	}
	
	if self.Size != 0 {
		size_str ="size=\""+strconv.Itoa(self.Size)+"\""
	}
	
	if self.Type == "button" || self.Type == "submit" {
		cls_str = "class=\""+ self.Class+"\""
		out.Append(fmt.Sprintf("<button type=\"%s\" %s %s %s %s  >%s</button>",self.Type,name_str,idstr,cls_str,placehold,self.Value))
	}else if self.Type == "checkbox"{
		cls_str = "class=\""+ self.Class+"\""
		if self.Checked == true{
		 	out.Append(fmt.Sprintf("<input type=\"%s\" value=\"%s\" %s %s %s %s checked=\"checked\" />",self.Type,self.Value,name_str,idstr,cls_str,placehold))
		}else{
		 	out.Append(fmt.Sprintf("<input type=\"%s\" value=\"%s\" %s %s %s %s />",self.Type,self.Value,name_str,idstr,cls_str,placehold))			
		}
	}else if self.Type == "textarea"{
		out.Append(fmt.Sprintf("<textarea  rows=4 %s %s %s  %s %s %s>%s</textarea>",name_str,idstr,cls_str,placehold,req_str, rdo_str,self.Value))
	}else if self.Type == "hr" {
		out.Append(fmt.Sprintf("<%s size=%d />",self.Type,self.Size))
	}else{
		out.Append(fmt.Sprintf("<input type=\"%s\" value=\"%s\"   %s %s %s %s  %s %s %s />", self.Type,self.Value,name_str,idstr,cls_str,placehold, req_str,rdo_str, size_str))
	}
	return out.String()
}

func (self *Input) get_default_id () string {
	self.Id = camelcase.Reverse(self.Name)
	return self.Id
}

func (self *Input) SetValue(val string){
	self.Value = val
}


func (self *Select) render_option() string {
	var out STRINGS
	
	var selected = false
	for _,v := range self.Args {
		if self.Multiple == true && len(self.Value) > 0 {
			
			for _,v2 := range self.Value {
				if v2 == v[0] {
					selected = true
					break
				}
			}
			
		}else if self.Multiple != true && len(self.Value) > 0 {
			if self.Value[0] == v[0] {
				selected = true
			}
		}
		if len(v) > 1 {
			if selected == true {
				out.Append(fmt.Sprintf("<option selected=\"selected\" value=%s>%s</option>",v[0],v[1]))	
			}else {
				out.Append(fmt.Sprintf("<option value=%s>%s</option>",v[0],v[1]))
			}
		}else {
			if selected == true {
				out.Append(fmt.Sprintf("<option selected=\"selected\" value=%s>%s</option>",v[0],v[0]))	
			}else {
				out.Append(fmt.Sprintf("<option value=%s>%s</option>",v[0],v[0]))						
			}
		}
		selected = false
	}
	return out.String()	
}

func (self *Select) String() string {
	var out STRINGS;
	idstr := ""
	cls_str := ""
	name_str := ""
	//	type_str := ""
	size_str := ""
	mul_str := ""
	rdo_str :=""
	if self.Id != "" {
		idstr = "id=\""+self.Id+"\""
	}else {
		idstr = "id=\""+self.get_default_id()+"\""
	}
	
	if self.Class != ""{
		cls_str = "class=\"form-control "+self.Class+"\""
	}else{
		cls_str = "class=\"form-control\""
	}
	if self.Name != ""{
		name_str = "name=\""+self.Name+"\""
	}
	if self.Size != 0 {
		size_str ="size=\""+strconv.Itoa(self.Size)+"\""
	}
	if self.Multiple == true {
		mul_str = "multiple=\"multiple\""
	}
	
	if self.ReadOnly == true {
		rdo_str = "readonly=\"readonly\""
	}
	
	out.Append(fmt.Sprintf("<%s %s %s %s %s %s %s>",self.Type,name_str,idstr,cls_str,size_str,mul_str,rdo_str))
	out.Append( self.render_option() )
	out.Append( fmt.Sprintf("</%s>",self.Type) )
	return out.String()
} 

func (self *Select) get_default_id () string {
	self.Id = camelcase.Reverse(self.Name)
	return self.Id
}

func (self *Select) SetValue(val ...string) {
	//数组的参数传进来时,在后面加上... 转换
	//
	self.Value = val	
}

func Render( V interface{}) string{
	elem := fmt.Sprintf("%v",V)
	var m = libs.Getdict(V)
	var out STRINGS;
	out.Append("<div class=\"form-group\">")

	if m["Type"] != "hidden" &&  m["Type"] != "hr"  {
		out.Append(fmt.Sprintf("<label class=\"col-sm-2 control-label\" id=\"lab_%s\" for=\"%s\">%s</label>",m["Id"],m["Id"],m["Description"] ))
		
		out.Append("<div class=\"col-sm-6\">")
		out.Append( elem )
	
		out.Append("</div>")
		if len(m["Note"])>0 {
			out.Append(fmt.Sprintf("<span class=\"wrong\">%s</span>",m["Note"]))
		}
	}else {
		out.Append(fmt.Sprintf("%v",V))		
	}
	out.Append("</div>")
	return out.String()		
}

// from database to fill form,simple version
func (self *Form) Fill(res interface{} ) {

	res_map := libs.Getdict(res)
	fmt.Println(res_map)
	for _,child := range self.Children {
		dict := libs.Getdict(child)
		if val,ok := res_map[dict["Name"]]; ok {
			
			v:= reflect.ValueOf(child).Elem().FieldByName("Value")
			if v.IsValid() {
				
				v.SetString(val)
				/*
				switch tv := val.(type) {
				case string:
					v.SetString( val.(string) )
				case int32,int64:
					v.SetString(  strconv.FormatInt( val.(int64),10 ) )
				default:
					fmt.Println("Fill error",tv)
					panic("Fill error")
				}
*/
			}
		}
	}

}


func (self *Form) FillValue(field_key string, field_val reflect.Value) {
	for _,child := range self.Children {
		iv := reflect.ValueOf(child)
		val := reflect.Indirect(iv)

		name_field := val.FieldByName("Name")
		if name_field.String() == field_key {
			val_field := val.FieldByName("Value")
			if val_field.IsValid() && val_field.CanSet() {
				fmt.Println("field_val type:",field_val.Kind())
				fmt.Println("val_field type:",val_field.Kind())
				switch val_field.Kind() {
					//这儿可能出现目标Form的元素是一个Select
					//但是结构体存储的是一个单int
					//比如状态选择
					//这样就要强制转换成一个[]string的Slice
					//一维数组only
				case reflect.Slice:
					if field_val.Kind() != reflect.Slice {
						val_field.Set(reflect.ValueOf([]string{fmt.Sprintf("%v",field_val)}))
					}else {
						val_field.Set(field_val)
					}
				case reflect.String:
					val_field.SetString(fmt.Sprintf("%v",field_val))
				default:
					val_field.Set(field_val)
				//fmt.Println("Bingo")
				}
			}
			break
		}
	}
}

// from database
func (self *Form) FillFormFromStruct(res interface{} ) {

	iv  :=reflect.ValueOf(res)
	val := reflect.Indirect(iv)

	for i:=0;i<val.NumField();i++ {
		field_name := val.Type().Field(i).Name
		
		fmt.Print( field_name, " ",val.Field(i) , " ")
		fmt.Println( val.Field(i).Kind() )
		switch val.Field(i).Kind() {
		case reflect.Int:
			
			self.FillValue(field_name,val.Field(i))
		case reflect.String:
			self.FillValue(field_name,val.Field(i))
		case reflect.Slice:
			//fmt.Println("Slice")
			self.FillValue(field_name, val.Field(i))
		default:
			fmt.Println("unknow ", val.Field(i).Kind())
			panic("Uknow Kind in FillFormFromStruct")
		}
		
		
	}
	
	//panic("haha")
	

}

func (self *Form) RenderCss() string{
	var ret string
	for _,child := range self.Children {
		ret += Render(child)
	}
	return ret
}

func (self *Form) Render() string{
	var out STRINGS;
	var ret string;
	for _, child := range self.Children{
		
		ret += Render(child)
	}
	
	out.Append("<div class=\"box box-"+self.Style+"\">")
	out.Append("<div class=\"box-header with-border\"><h3 class=\"box-title\">"+self.Title+"</h3></div>")
	out.Append("<form method=\"POST\" class=\"form-horizontal\" action=\""+self.Action+ "\">")
	out.Append("<div class=\"box-body\">")
	out.Append(ret+"</div></form>")
	out.Append("</div>")
	
	return out.String()
	
}

func Password(ipt *Input) *Input{
	ipt.Type= "password"
	return ipt
}

func TextBox(ipt *Input) *Input{
	ipt.Type= "text"
	return ipt
}

func TextArea(ipt *Input) *Input{
	ipt.Type= "textarea"
	return ipt
}

func Hidden(ipt *Input) *Input{
	ipt.Type = "hidden"
	return ipt
}

func Dropdown(sel *Select) *Select{
	//tst4 := [][]string{{"a","b","c"},{"d","e"}}
	sel.Type = "select"
	return sel
	
}

func Submit(ipt *Input) *Input{
	ipt.Type = "submit"
	return ipt
}

func Button(ipt *Input) *Input{
	ipt.Type = "button"
	return ipt
}

func GroupDropdown(sel *Select) *Select{
	sel.Type = "select"
	sel.Multiple = true
	return sel
}

func CheckBox(ipt *Input) *Input{
	ipt.Type = "checkbox"
	return ipt
}

func Hr(ipt* Input) *Input{
	ipt.Type = "hr"
	ipt.Size = 1
	return ipt
}

func AlertBox(alt *Alert) string{
	return alt.Render()
}

func InfoForm(title string,action string, clds ...interface{}) *Form{
	aform := &Form{Title:title,Action:action,Style:"info",Children:clds}
	return aform
}

/*
func WarnForm(title string,action string, clds ...string) string {
	aform := Form{Title:title,Action:action,Style:"warning",Children:clds}
	return aform.Render()
}

func DangerForm(title string,action string, clds ...string) string {
	aform := Form{Title:title,Action:action,Style:"danger",Children:clds}
	return aform.Render()
}
*/

func NewRegex(rxp_str string, err_str string ) *Regex{
	var valid = regexp.MustCompile(rxp_str)
	return &Regex{rexp:valid, Msg: err_str}
}

func (r *Regex) MatchString(str string) bool{
	return r.rexp.MatchString(str)
}
