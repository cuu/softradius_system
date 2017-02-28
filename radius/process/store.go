package process

import (
//	"fmt"
	ctl "github.com/cuu/softradius/controllers"
//	rdb "github.com/cuu/softradius/database/shelf"
	r "github.com/cuu/softradius/routers"
)

//封装一下radius的特定操作

func UpdateUserExpire(user *ctl.Members) {
	if user.ExpireDate == r.MAX_EXPIRE_DATE {
		
	}
}
