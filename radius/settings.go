package radius

import (
//	"fmt"
	ctl "github.com/cuu/softradius/controllers"
	rdb "github.com/cuu/softradius/database/shelf"
)

func GetClientsMap() map[string]string {
	var ret = make(map[string]string)
	var nods []ctl.Bas
	rdb.DataBase().SkipGet2(&nods,0,1000)

	if len(nods) > 0 {
		for _,v := range  nods {
			ret[v.IpAddr] = v.Secret
		}
	}

	return ret
}

