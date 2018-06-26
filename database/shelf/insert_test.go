package shelf

import r "gopkg.in/gorethink/gorethink.v3"
import . "github.com/pkg4go/assert"
import "testing"
import "fmt"

func TestInsert(t *testing.T) {
	a := A{t}

	err := Register(r.ConnectOpts{
		Address: "localhost:28015",
	}, "test")

	a.Nil(err)
	
	DataBase().CreateTable("user")

	resp,_ := DataBase().Insert(map[string]string{
		"Name": "hello",
		"Desc": "world",
	}, "user")
	
	_key := resp.GeneratedKeys
	
	fmt.Println(_key)

	
	rsp,_ := _Db.Get("user",_key[0])
	var auser map[string]interface{}
	err = rsp.One(&auser)
	fmt.Println(auser)

	//这儿有个BUG,必须要加 struct tag 才能让Id 对应上rethinkdb中的id,否则就会被认作是两个字段
	//字段也全部用tag才能变成小写,要不然就统一的大写第一个开头
	type User struct {
		Id string `gorethink:"id,omitempty"`
		Name string 
		Desc string
		Test int  `gorethink:"Test,omitempty"`
	}

	DataBase().Insert(map[string]string{
		"Name": "guu",
		"Desc": "admin",
	
	}, "user")
	
//	_Db.CreateTable(UserTest{})
	
	DataBase().Insert(User{
		Name: "haoxin",
		Desc: "haha",
		Test: 3000,
	})
	
	/*
	DataBase().Insert(User{
		Name: "aueea@1843.com",
		Desc: "nenennea",
	})

	DataBase().Insert(User{
	//	Id: 1,
		Name:"guu",
		Desc:"come on",
	})
	*/
	/*
	resp, err = DataBase().FilterUpdate("user",map[string]string{"name":"helo"},map[string]string{"Desc":"not the fuck world"})
	a.Nil(err)
	fmt.Println("Update ",resp.Replaced)
*/
	//rsp,_ = DataBase().SkipGet("user",1,1)
	//	rsp,_ = DataBase().SkipGet("user",1,10, r.OrderByOpts{Index: "name"} )
	//rsp,_ = DataBase().SkipGet("user",1,10, r.Desc("name"))
	rsp,_ = DataBase().SkipGet("user",0,10)	
//	var users []map[string]interface{}
	var users []User  // it works
	err = rsp.All(&users)
	fmt.Println("user list:", users)

	
}
