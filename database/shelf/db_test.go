package shelf

import r "gopkg.in/dancannon/gorethink.v2"
import . "github.com/pkg4go/assert"
import "testing"

func TestDB(t *testing.T) {
	a := A{t}

	db, err := Open(r.ConnectOpts{
		Address: "localhost:28015",
	}, "test")

	a.Equal(err, nil)
	a.Equal(db.IsConnected(), true)

	db.Close()

	a.Equal(db.IsConnected(), false)
}
