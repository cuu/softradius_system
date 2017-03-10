package models

import(
	"regexp"
	"fmt"
)

func Is_alphanum(x int) *Regex{
	rxp_str := "^[A-Za-z0-9]{%d}$"
	err_str := "必须是长度为%d的数字字母组合"
	var valid = regexp.MustCompile(fmt.Sprintf(rxp_str,x))
	return &Regex{rexp:valid, Msg: fmt.Sprintf(err_str,x)}
}

func Is_alphanum2(x int, y int) *Regex{
	rxp_str := "^[A-Za-z0-9]{%d,%d}$"
	err_str := "必须是长度为%d到%d的数字字母组合"
	var valid = regexp.MustCompile(fmt.Sprintf(rxp_str,x,y))
	return &Regex{rexp:valid, Msg: fmt.Sprintf(err_str,x,y)}
}

func Is_alphanum3(x int, y int) *Regex{
	rxp_str := "^[A-Za-z0-9\\_\\-]{%d,%d}$"
	err_str := "必须是长度为%d到%d的数字字母与下划线组合"
	var valid = regexp.MustCompile(fmt.Sprintf(rxp_str,x,y))
	return &Regex{rexp:valid, Msg: fmt.Sprintf(err_str,x,y)}
}

func Len_of( x int ,y int) *Regex{
	rxp_str := "[\\s\\S]{%d,%d}$"
	err_str := "长度必须为%d到%d"
	var valid = regexp.MustCompile(fmt.Sprintf(rxp_str,x,y))
	return &Regex{rexp:valid, Msg: fmt.Sprintf(err_str,x,y)}	
}

var Notnull = NewRegex(".*\\S+.*","不许为空")

var Is_not_empty = NewRegex(`.+`, "不允许为空")

var Is_date = NewRegex(`(\d{4})-(\d{2}-(\d\d))`, "日期格式:yyyy-MM-dd")
var Is_email = NewRegex(`[\w-]+(\.[\w-]+)*@[\w-]+(\.[\w-]+)+$`, "email格式,比如name@domain.com")
var Is_chars = NewRegex(`^[A-Za-z]+$`, "必须是英文字符串")
var Is_number = NewRegex(`^[0-9]*$`, "必须是数字")
var Is_number2 = NewRegex(`^[1-9]\d*$`,"必须是大于0的正整数")
var Is_number3 = NewRegex(`^(([1-9]\d*)|0)(\.\d{1,3})?$`, "支持包含(最大3位)小数点 xx.xxxxx")
var Is_numberOboveZore = NewRegex(`^\d+$`,"必须为大于等于0的整数")
var Is_cn = NewRegex("^[\u4e00-\u9fa5],{0,}$", "必须是汉字")
var Is_url = NewRegex(`[a-zA-z]+://[^\s]*`, "url格式 xxxx://xxx")
var Is_phone = NewRegex(`^(\(\d{3,4}\)|\d{3,4}-)?\d{7,8}$`, "固定电话号码格式：0000-00000000")
var Is_idcard = NewRegex(`^\d{15}$|^\d{18}$|^\d{17}[Xx]$`, "身份证号码格式")
var Is_ip = NewRegex(`(^$)|\d+\.\d+\.\d+\.\d+`, "ip格式：xxx.xxx.xxx.xxx")
var Is_rmb = NewRegex(`^(([1-9]\d*)|0)(\.\d{1,2})?$`, "人民币金额 xx.xx")

var Is_period = NewRegex(`(^$)|^([01][0-9]|2[0-3]):[0-5][0-9]-([01][0-9]|2[0-3]):[0-5][0-9]$`,"时间段，hh:mm-hh:mm,支持跨天，如 19:00-09:20")
var Is_telephone = NewRegex(`^1[0-9]{10}$`, "必须是手机号码")
var Is_time = NewRegex(`(\d{4})-(\d{2}-(\d\d))\s([01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]`, "时间格式:yyyy-MM-dd hh:mm:ss")
var Is_time_hm = NewRegex(`^([01][0-9]|2[0-3]):[0-5][0-9]$`, "时间格式: hh:mm")


var css_style = make(map[string]map[string]string)

func init(){
	css_style["input_style"] = map[string]string{"class":"form-control"}
	css_style["button_style"] = map[string]string{"class": "btn btn-primary"}
	css_style["button_style_block"] = map[string]string{"class": "btn btn-block"}
	css_style["div_style"] = map[string]string{"class": "block"}

}

