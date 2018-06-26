package shelf

import (
	 r "gopkg.in/gorethink/gorethink.v3"
	"fmt"
	"reflect"
	//"adminlte/libs"
)


func (db *DB) Insert(args ...interface{}) (r.WriteResponse, error)  {
	v := args[0]

	name := getTableName(args...)
	fmt.Println("Insert Table name:", name)
	
	return r.DB(db.Name).Table(name).Insert(v ).RunWrite(db.Session)
}

//适合一次性抽插时,对某个key产生的重复做个判断
// FilterInsert(Node{},Node{xxx:"xxx",yy:"yy"}, "xxx" )
// check table Node{}'s xxx key, if there is a value "xxx" ,or insert the new record
// Node{xxx:"xxx",yy:"yy"}
///目前最好是map[string]string 类型 ,

func (db *DB) FilterInsert(arg interface{}, dupkey string ) (ret []string,count int ) {

	name := getTableName(arg)
	kind := reflect.TypeOf(arg).Kind()
	if kind != reflect.Struct && kind != reflect.Map  && kind != reflect.Ptr {
		panic("only struct,map,ptr are supported")
		return 
	}
	
	var search = make(map[string]interface{})
	
	var shval interface{}
	
	if kind == reflect.Ptr {
		//	val := reflect.ValueOf(arg).Elem().FieldByName(dupkey)
		ss := reflect.Indirect( reflect.ValueOf(arg)).FieldByName(dupkey)
		if reflect.TypeOf(ss).Kind() == reflect.Struct {
			shval = ss.String()
		}else
		{
			shval = ss
		}
	}
	
	if kind == reflect.Map {
		v := reflect.ValueOf(arg)
		ss := v.MapIndex( reflect.ValueOf(dupkey))
		if reflect.TypeOf(ss).Kind() == reflect.Struct {
			shval = ss.String()
		}else
		{
			shval = ss
		}
	}

//	fmt.Println(shval,reflect.TypeOf(shval).Kind())
	search[dupkey] = shval
//	fmt.Println(search)
	count = db.FilterCount(name,search)
	
	/*
	rsp,err := r.DB(db.Name).Table(name).Filter(search).Count().Run(db.Session)
	if err == nil {
		rsp.One(&count)
	}
*/
	if  count == 0 {
		
		resp,err := db.Insert(arg,name)
		if err == nil {
			fmt.Println(resp.GeneratedKeys)
			ret = resp.GeneratedKeys
		}else {
			fmt.Println(err)
			
		}
//		fmt.Println("No Old record")
	}else {
		fmt.Println("FilterInsert ,count ",count)
	}

	return
	
}
