package shelf

import (
	r "gopkg.in/gorethink/gorethink.v3"
	"adminlte/libs"
	"fmt"
)

/*
var hero map[string]interface{}
err = res.One(&hero)
*/
// precise update

//根据id取得一条数据
func (db *DB) QuOne(one interface{},id string) error  {

	name := getTableName(one)
	var err error
	if Type(one) != "ptr" {
		err = r.DB(db.Name).Table(name).Get(id).ReadOne(&one,db.Session)
	}else {
		err = r.DB(db.Name).Table(name).Get(id).ReadOne(one,db.Session)
	}
	
	if err == nil {
		return err
	}else { libs.Debug("%s",err) }

	return err
}


func (db *DB) Get(key interface{}, args ...interface{}) (*r.Cursor, error){
	
	if len(args) < 1 {
		return nil,fmt.Errorf("1 arguments least")
	}
	
	name := getTableName(key)
	
	return r.DB(db.Name).Table(name).Get(args[0]).Run(db.Session)
}

// 取得所有结果 
func (db *DB) QuAll(args ...interface{}) error {
	//args[0] 应该是pointer
	if len(args) < 2 {
		panic("QuAll needs two arguments")
		return fmt.Errorf("QuAll needs two arguments")
	}
	
	name := getTableName(args[0])
	fmt.Println("QueryAll ",name)
	var rsp *r.Cursor
	var err error
	if len(args) > 2 {
		//假定 第3个参数是string或是 []string
		
		rsp,err = r.DB(db.Name).Table(name).Filter(args[1]).Pluck(args[2]).Run(db.Session)
	}else {
		rsp,err = r.DB(db.Name).Table(name).Filter(args[1]).Run(db.Session)	
	}
	
	if err == nil {
		rsp.All(args[0])
	}else {
		fmt.Println("QuAll ",err)
	}
	return err
}


func (db *DB) GetAll(key interface{}, args ...interface{}) (*r.Cursor, error){
	name := getTableName(key)
	fmt.Println("GetAll ",name)
	return r.DB(db.Name).Table(name).GetAll(args...).Run(db.Session)	
}

func (db *DB) FilterCount(tab interface{},args ...interface{} ) int {
	if len(args) < 1 {
		return -1
	}
	name := getTableName(tab)
	fmt.Println("FilterCount table name:", name,args[0])
	count := -1
	rsp,err := r.DB(db.Name).Table(name).Filter(args[0]).Count().Run(db.Session)
	if err ==nil {
		err = rsp.One(&count)
	}
	return count
	
}

func (db *DB) FilterOne(store interface{},args ...interface{}) error {

	if len(args) < 1 {
		return fmt.Errorf("FilterOne an arguments required")
	}
	
	name := getTableName(store)
	fmt.Println("FilterOne ",name)
	rsp,err := r.DB(db.Name).Table(name).Filter(args[0]).Run(db.Session)
	if err == nil {
		err = rsp.One(store)
	}else {
		fmt.Println("FilterOne: ",err)
	}
	return err	
}

func (db *DB) FilterGet(table interface{},args ...interface{}) (*r.Cursor, error)  {

	if len(args) < 1 {
		return nil,fmt.Errorf("1 arguments least")
	}
	
	name := getTableName(table)
	fmt.Println("FilterGet Table name:", name,args)

	return r.DB(db.Name).Table(name).Filter(args[0]).Run(db.Session)
}



func (db *DB) SkipGet(table interface{},skip int, limit int ,args ...interface{}) (*r.Cursor, error){
	name := getTableName(table)
	fmt.Println("SkipGet ",skip," ", limit, " ",name)
	if len(args) > 0{
		return r.DB(db.Name).Table(name).OrderBy(args...).Skip(skip).Limit(limit).Run(db.Session)
	}else
	{
		return r.DB(db.Name).Table(name).Skip(skip).Limit(limit).Run(db.Session)
	}
}

func (db *DB) SkipGet2(store interface{},skip int, limit int ,args ...interface{}) error {
	name := getTableName(store)
	fmt.Println("SkipGet2 ",skip," ", limit, " ",name)
	
	var rsp *r.Cursor
	var err error
	
	if len(args) > 0 {
		rsp,err = r.DB(db.Name).Table(name).OrderBy(args...).Skip(skip).Limit(limit).Run(db.Session)
		
	}else
	{
		rsp,err =  r.DB(db.Name).Table(name).Skip(skip).Limit(limit).Run(db.Session)
	}

	if err == nil {
		rsp.All(store)
	}else {
		fmt.Println("SkipGet2 ",err)
	}

	return err
}
