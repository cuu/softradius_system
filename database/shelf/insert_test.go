package shelf

import r "gopkg.in/dancannon/gorethink.v2"
import . "github.com/pkg4go/assert"
import "testing"

func TestInsert(t *testing.T) {
	a := A{t}

	db, err := Open(r.ConnectOpts{
		Address: "localhost:28015",
	}, "test")

	a.Nil(err)

	db.Insert(map[string]string{
		"name": "hello",
		"desc": "world",
	}, "user")

	type User struct {
		Name string
		Desc string
	}

	db.Insert(User{
		Name: "haoxin",
		Desc: "haha",
	})
}
