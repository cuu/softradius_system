package libs

import (
	"reflect"
	"github.com/pkg4go/camelcase"
	"time"
	"strings"
	"runtime"
	"fmt"
	"log"
	"strconv"
	_ "github.com/mitchellh/mapstructure"
	"github.com/cuu/softradius/libs/times"
	"math"
	mrand "math/rand"
	"crypto/rand"
	"io"
)

var _base_id int

func RemoveSuffix(str string, suf string) string {
	if len(str) == 1 {
		return str
	}
	for {
		if strings.HasSuffix(str,suf) {
			str = str[:len(str)-1]
		} else {
			break
		}
	}
	return str
}

func GetType(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}



// Type返回的可能是ptr,string,struct,slice ,比较可视化
func Type(v interface{}) string {
	t := reflect.TypeOf(v);
	k := t.Kind()
	return k.String()
}

func GetTplName(controller interface{} ) string {
	//不管是ptr还是struct本体,都返回struct的名子,不带星号
	var ret string
	if t := reflect.TypeOf(controller); t.Kind() == reflect.Ptr {
		ret =  t.Elem().Name()
	} else {
		ret = t.Name()
	}
	//XxxYyyy => xxx_yyyy
	return camelcase.Reverse(ret)+".tpl"
}

func GetCurrTimeNano()string {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	
	return fmt.Sprintf("%s",timestamp)
}
func Get_currtime() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
	
}

func Debug(format string, a ...interface{}) {
	pc := make([]uintptr,10)
	runtime.Callers(2,pc)
	f := runtime.FuncForPC(pc[0])
	file,line := f.FileLine(pc[0])
	
	info := fmt.Sprintf(format, a...)

	log.Printf("[sr] %s:%d,%s %v", file, line,f.Name(), info)
}


func get_val(val reflect.Value) string{
	
	switch val.Kind() {
	case reflect.Invalid:
		return "invalid"		
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)
	case reflect.String:
		return val.String()    
		// etc...
	case reflect.Bool:
		return strconv.FormatBool(val.Bool())
	case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Slice, reflect.Map:
		return val.Type().String() + " 0x" + strconv.FormatUint(uint64(val.Pointer()), 16)
	default: // reflect.Array, reflect.Struct, reflect.Interface
		return val.Type().String() + " value"
	}
	
}

func Getdict(V interface{}) map[string]string {
	m := make(map[string]string)
	var self = reflect.ValueOf(V).Elem()
	
	for i := 0; i < self.NumField(); i++ {
		valueField := self.Field(i)
		typeField := self.Type().Field(i)
		f := valueField.Interface()
		val := reflect.ValueOf(f)	
		m[typeField.Name] = get_val(val)
	}
	return m
}

func Getdict2(V interface{}) map[string]interface{}{
	m := make(map[string]interface{})
	var self = reflect.ValueOf(V).Elem()
	
	for i := 0; i < self.NumField(); i++ {
		valueField := self.Field(i)
		typeField := self.Type().Field(i)
		//f := valueField.Interface()
		//val := reflect.ValueOf(f)	
		m[typeField.Name] = valueField
	}
	return m
}



func SetField(obj interface{}, name string, value interface{}) error {

	structValue := reflect.ValueOf(obj).Elem()
	fieldVal := structValue.FieldByName(name)

	if !fieldVal.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !fieldVal.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	val := reflect.ValueOf(value)

	if fieldVal.Type() != val.Type() {

		if m,ok := value.(map[string]interface{}); ok {

			// if field value is struct
			if fieldVal.Kind() == reflect.Struct {
				return FillStruct(m, fieldVal.Addr().Interface())
			}

			// if field value is a pointer to struct
			if fieldVal.Kind()==reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.Struct {
				if fieldVal.IsNil() {
					fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
				}
				// fmt.Printf("recursive: %v %v\n", m,fieldVal.Interface())
				return FillStruct(m, fieldVal.Interface())
			}

		}

		return fmt.Errorf("Provided value type didn't match obj field type")
	}

	fieldVal.Set(val)
	return nil

}

/*
result := &MyStruct{}
err := FillStruct(myData,result)
*/
func FillStruct(m map[string]interface{}, s interface{}) error {
	for k, v := range m {
		err := SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func Str(in interface{} ) string {
	if in == nil {
		in = ""
	}
	t := Type(in)
	switch t {
	case "string":
		tv := reflect.ValueOf(in).String()
		return tv
	case "int":
		tv := reflect.ValueOf(in).Int()
		str := strconv.Itoa(int(tv))
		return str
	}
	
	return ""	
}

func Or(in interface{}, other interface{} ) interface{} {
	if in == nil {
		in = ""
	}
	t := Type(in)
	switch t {
	case "string":
		tv := reflect.ValueOf(in).String()
		if tv == "" {
			return other
		}
	case "int":
		tv := reflect.ValueOf(in).Int()
		if tv <= 0 {
			return other
		}
	}
	
	return ""
}

func InSlice(in interface{}, list interface{} ) bool {
	ret := false
	if in == nil {
		in = ""
	}
	t := Type(in)
	switch t {
	case "string":
		tv := reflect.ValueOf(in).String()
		for _,l := range list.([]string) {
			v:= reflect.ValueOf(l)
			if tv == v.String() {
				ret = true
				break
			}
		}
	
	case "int":
		tv := reflect.ValueOf(in).Int()
		for _,l := range list.([]int) {
			v := reflect.ValueOf(l)
			if tv == v.Int() {
				ret = true
				break
			}
		}
	}
	
	return ret	
}

func In( in interface{},  list ...interface{}) bool {
	fmt.Print("In: ",in ," ")
	fmt.Println("List: ",list)
	ret := false
	if in == nil {
		in = ""
	}
	t := Type(in)
	switch t {
	case "string":
		tv := reflect.ValueOf(in).String()
		for _,l := range list {
			v:= reflect.ValueOf(l)
			if tv == v.String() {
				ret = true
				break
			}
		}
	
	case "int":
		tv := reflect.ValueOf(in).Int()
		for _,l := range list {
			v := reflect.ValueOf(l)
			if tv == v.Int() {
				ret = true
				break
			}
		}
	}
	
	return ret
}

func ToInt(str string ) int{
	s,_ := strconv.Atoi(str)
	return s
}

//一般,乘法是int,除法是string
func Fen2yuan(fen int ) string {
	
	v := float64(float64(fen)/100.00)
	return fmt.Sprintf("%.2f",v)
}

func Yuan2fen(yuan int ) int {
	return yuan * 100
}

func Bb2mb( bb int ) string {
	v:= float64(bb/1024.0/1024.0)
	return fmt.Sprintf("%.2f",v)
	
}

func Bbgb2mb(bb int,gb int) string {
	bl:= float64(bb/1024.0/1024.0)
	gl:= float64(gb*4*1024.0*1024.0*1024.0)
	tl:= bl + gl
	return fmt.Sprintf("%.2f",tl)
}

func StrKb2mb (kb string ) string {
	//	fmt.Println(kb)
	_kb,_ := strconv.Atoi(kb)
	v:= float64(_kb/1024.0)
	return fmt.Sprintf("%.2f",v)
}
func Kb2mb (kb int ) string {
//	fmt.Println(kb)
	v:= float64(kb/1024.0)
	return fmt.Sprintf("%.2f",v)
}

func Mb2kb( mb int ) int {
	return mb*1024
}

func Hour2sec (hor int) int {
	return hor*3600
}

func Sec2hour (sec int ) string {
	v:= float64(sec /3600.0)
	return fmt.Sprintf("%.2f",v)
}

func Mbps2bps( mbps int ) int {
	return mbps *1024*1024
}

func Bps2mbps(bps int) string {
	_mbps := float64(bps/1024.0/1024.0)
	return fmt.Sprintf("%.3f",_mbps)
}

func AddMonths( t time.Time, months int ) time.Time {
	
	return t.AddDate(0,months,0)
	
}

func IsExpire(dstr string ) bool {
	if dstr == "" {
		return false
	}

	t := times.StrToLocalTime(dstr)
	now := time.Now()
	
	return now.After(t)
}


// Now和AcctStartTime 做比较
func FmtOnlineTime(acct_start_time string) string {
	t := times.StrToLocalTime(acct_start_time)
	now := time.Now()

	dt := now.Sub(t)
	times := dt.Seconds() // float64
	d := times / (3600 * 24)
	h := math.Mod(times, (3600 * 24) ) / 3600
	m := math.Mod(times, (3600 * 24) )
	m  = math.Mod(m,3600) / 60

	if int(d) > 0 {
		return fmt.Sprintf("%d天%d小时%d分钟" , int(d), int(h), int(m))
	}else if int(d) > 0 && int(h) > 0 {
		return fmt.Sprintf("%d小时%d分钟" ,int(h), int(m))
	}else {
		return fmt.Sprintf("%d分钟",int(m))
	}
	
}

func TimeBetween(frame time.Time, begin time.Time, end time.Time ) bool {
	if frame.After(begin) && frame.Before(end) {
		return true
	}

	return false 
}
func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}


func GenOrderId() string {

	if _base_id >= 9999 {
		_base_id = 0
	}
	_base_id +=1
	now := time.Now()
	strTime := times.Format("YmdHis", now)
	return fmt.Sprintf("%s%02d",strTime,_base_id)
}

func Random(min,max int) int {
	return mrand.Intn(max-min)+min
}
