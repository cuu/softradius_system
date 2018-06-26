package shelf

import  (
	"github.com/pkg4go/camelcase"
	"github.com/pkg4go/convert"
	"reflect"
//	"fmt"
)

//返回的是 变量的名称,如struct的名子
//接受的参数有 struct,&struct,&[]struct, &[]map[xx]yy
func TypeName( v interface{} ) string {
	t:= reflect.TypeOf(v)
	k:= t.Kind()
	//fmt.Println(t)
	if k == reflect.Ptr {
		//fmt.Print( t.Elem().String()," ")
		iv := reflect.Indirect( reflect.ValueOf(v) )
		vk := iv.Type().Kind()
		//fmt.Println(vk)
		if vk == reflect.Slice {
			//fmt.Println("-->[]", iv.Type().Elem().Name())
			if iv.Type().Elem().Name() == "" {
				return iv.Type().Elem().String()
			}else {
				return iv.Type().Elem().Name()
			}
		}else if vk == reflect.Struct {
			//fmt.Println("->*" ,iv.Type().Name())
			return iv.Type().Name()
		}
	}else if k == reflect.Slice {
		if t.Elem().Name() == "" {
			return t.Elem().String()
		}else {
			return t.Elem().Name()
		}
	}else {
		return t.Name()
	}
	return ""
}

// Type返回的可能是ptr,string,struct,slice ,比较可视化
func Type(v interface{}) string {
	t := reflect.TypeOf(v);
	k := t.Kind()
	return k.String()
}

func getTableName(args ...interface{}) string {
	var name string
	if Type( args[0] ) == "slice" {
		panic("getTableName on slice")
	}
	
	if TypeName(args[0]) == "string" {
		name = convert.String(args[0])
	}else{
		name = camelcase.Reverse(TypeName(args[0]))
	}
	
	if len(args) == 1 {
		return name
	}

	if len(args) == 2 {
		if n := convert.String(args[1]); n != "" {
			return n
		}
	}

	return name
}
