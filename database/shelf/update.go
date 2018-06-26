package shelf

import (
	 r "gopkg.in/gorethink/gorethink.v3"
	"fmt"
)

// precise update
// update("tablename",{toget},{toupdate})
func (db *DB) Update(args ...interface{}) (r.WriteResponse, error)  {

	if len(args) < 2 {
		fmt.Println("Update needs 2 arguments")
		return r.WriteResponse{},fmt.Errorf("Update needs 2 arguments")
	}
	
	_search := args[0]
	_replace:= args[1]  //一般是一个struct的 pointer
	name := getTableName(_replace)
	fmt.Println("Update Table name:", name,_search,_replace)

	return r.DB(db.Name).Table(name).Get(_search).Update(_replace ).RunWrite(db.Session)
}

// bunch update
func (db *DB) FilterUpdate(key interface{},args ...interface{}) (r.WriteResponse, error)  {
	var response r.WriteResponse
	if len(args) < 2{
		return response,fmt.Errorf("2 arguments least")
	}
	
	_search := args[0]
	_replace := args[1]
	name := getTableName(key)
	fmt.Println("FilterUpdate Table name:", name,_search,_replace)

	return r.DB(db.Name).Table(name).Filter(_search).Update(_replace ).RunWrite(db.Session)
}

func (db *DB) FilterDelete(key interface{},args ...interface{}) (r.WriteResponse, error)  {

	_search  := args[0]
	name := getTableName(key)
	fmt.Println("Delete Table name:", name)

	return r.DB(db.Name).Table(name).Filter(_search).Delete().RunWrite(db.Session)
}

func (db *DB) UpdateAll(key interface{}, args ...interface{}) (r.WriteResponse, error)  {

	_replace := args[0]
	
	name := getTableName(key)
	fmt.Println("UpdateAll Table name:", name)

	return r.DB(db.Name).Table(name).Update(_replace).RunWrite(db.Session)
}


