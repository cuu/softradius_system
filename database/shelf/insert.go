package shelf

import r "gopkg.in/dancannon/gorethink.v2"

func (db *DB) Insert(args ...interface{}) error {
	v := args[0]

	name := getTableName(args...)

	return r.DB(db.Name).Table(name).Insert(v).Exec(db.Session)
}
