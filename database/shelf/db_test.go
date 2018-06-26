package shelf

import r "gopkg.in/gorethink/gorethink.v3"
import . "github.com/pkg4go/assert"
import "testing"
import "fmt"
func TestDB(t *testing.T) {
	a := A{t}

	db, err := Open(r.ConnectOpts{
		Address: "localhost:28015",
	}, "test")

	a.Equal(err, nil)
	a.Equal(db.IsConnected(), true)

	db.Close()

	a.Equal(db.IsConnected(), false)
	
	fmt.Println("testing..")
  
  fmt.Println(db.Table("user"))
}
