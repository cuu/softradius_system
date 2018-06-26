package shelf

import (
	r "gopkg.in/gorethink/gorethink.v3"
	"log"
  "fmt"
)

type DB struct {
	Name    string
	Session *r.Session
	Options r.ConnectOpts
}


var _Db *DB

func DataBase() *DB{
	return _Db
}

func Open(opts r.ConnectOpts, dbName string) (*DB, error) {
	session, err := r.Connect(opts)

	if err != nil {
		return nil, err
	}

	db := DB{
		Name:    dbName,
		Session: session,
		Options: opts,
	}

	return &db, nil
}

func (db *DB) Close(opts ...r.CloseOpts) error {
	return db.Session.Close(opts...)
}

func (db *DB) IsConnected() bool {
	return db.Session.IsConnected()
}

func (db *DB) Reconnect(opts ...r.CloseOpts) error {
	return db.Session.Reconnect(opts...)
}

func (db *DB) CreateTable(args ...interface{}) {
	name := getTableName(args...)
	fmt.Println("CreateTable ", name)
	r.DB(db.Name).TableDrop(name).Run(db.Session)
	response, err := r.DB(db.Name).TableCreate(name).RunWrite(db.Session)
	if err != nil {
		log.Fatalf("Error creating table: %s", err)
	}
	fmt.Println(response)
}

func (db *DB) TableExisted(name string){
	
}



func (db *DB) Table(name string) r.Term {
	return r.DB(db.Name).Table(name)
}


func Register( opts r.ConnectOpts, dbName string) error {
	db,err := Open(opts,dbName)
	if err != nil {
		panic( "DB Open Failed")
	}

	_Db = db
	return err
}

