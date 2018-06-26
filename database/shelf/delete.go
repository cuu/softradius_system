package shelf

import (
	r "gopkg.in/gorethink/gorethink.v3"
	//"adminlte/libs"
	"fmt"
)


func (db *DB) Del(table interface{}, args ...interface{})  (r.WriteResponse, error) {
	var ret r.WriteResponse
	if len(args) < 1 {
		return ret,fmt.Errorf("1 arguments least")
	}
	
	name := getTableName(table)
	
	return r.DB(db.Name).Table(name).Get(args[0]).Delete().RunWrite(db.Session)
}

func (db *DB) FilterDel(table interface{}, args interface{} )  (r.WriteResponse, error)  {

	name := getTableName(table)
	
	return r.DB(db.Name).Table(name).Filter(args).Delete().RunWrite(db.Session)	
	
}
