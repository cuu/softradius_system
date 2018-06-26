package shelf

import r "gopkg.in/dancannon/gorethink.v2"

type DB struct {
	Name    string
	Session *r.Session
	Options r.ConnectOpts
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

func (db *DB) Table(name string) r.Term {
	return r.DB(db.Name).Table(name)
}
