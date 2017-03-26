package process

import (
	"time"
	"fmt"
//	ctl "github.com/cuu/softradius/controllers"
	rdb "github.com/cuu/softradius/database/shelf"
	r "github.com/cuu/softradius/routers"
	"github.com/cuu/softradius/libs"
	"github.com/cuu/softradius/libs/times"
)

//封装一下radius的特定操作
func UpdateUserExpire(au *AuthUser) {
	if au.User.ExpireDate == r.MAX_EXPIRE_DATE {
		if au.Product.Policy == r.AwesomeFee || au.Product.Policy == r.AwesomeFeeBoTime {
			sec := au.Product.FeeTimes
			once_created := times.StrToLocalTime(au.User.CreateTime) // panic if there is no CreateTime
			new_date := libs.AddDuration(once_created,fmt.Sprintf("%ds",sec))
			strTime := times.Format("Y-m-d H:i:s", new_date)
			au.User.ExpireDate = strTime
			
			resp,err := rdb.DataBase().Update(au.User.Id,au.User)
			if err == nil {
				fmt.Println("Replaced ",resp.Replaced )
			}else {
				fmt.Println(err)
			}
		}
	}
}
