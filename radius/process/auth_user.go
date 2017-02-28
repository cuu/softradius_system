package process


import (
//	re "gopkg.in/gorethink/gorethink.v3"	
	rdb "github.com/cuu/softradius/database/shelf"
	ctl "github.com/cuu/softradius/controllers"
	"fmt"

	"github.com/cuu/softradius/libs"
	r "github.com/cuu/softradius/routers"
	//	"github.com/cuu/softradius/radius"
)

type AuthUser struct {
	UserGot bool
	User *ctl.Members
	ProductGot bool
	Product *ctl.Products

	Resp ProcessResp
}


//Constructor
func NewAuthUser()  *AuthUser {
	one := &AuthUser{}
	one.ResetResp()

	return one
}


func (self *AuthUser) ResetResp() {
	self.Resp.Attrs = nil
	self.Resp.Code  = AccessAccept
	self.Resp.Attrs = make(map[string]interface{})
}

func (self *AuthUser) ErrorAuth(reply string) bool {
	self.ResetResp()

	self.Resp.Code = AccessReject //#packet.AccessReject
	self.Resp.Attrs["Reply-Message"] = reply

	return false
}

func (p *AuthUser) GetUser(username string ,password string) bool {
	one := &ctl.Members{}

	err := rdb.DataBase().FilterOne(one,map[string]string{"Name":username})
	if err == nil {
		if one.Password == password {
			p.User  = one
			p.UserGot = true  
			return true
		}
	}

	p.ErrorAuth("User pass error")
	p.User = one
	return false
	
}


func (p *AuthUser) GetProduct() bool {
	
	one := &ctl.Products{}

	if p.UserGot != true {
		return p.UserGot
	}
	err := rdb.DataBase().QuOne(one,p.User.ProductId)
	if err == nil {
		p.Product = one
		p.ProductGot = true
		return true
	}else {
		fmt.Println(err)
	}
	p.ErrorAuth("Product error")
	p.Product = one
	return false
}

func (p *AuthUser) Billing() bool {
	if p.UserGot != true {
		return p.ErrorAuth("User error")
	}

	if p.ProductGot != true {
		return p.ErrorAuth("Product error")
	}

	
	if libs.IsExpire( p.User.ExpireDate) {
		return p.ErrorAuth("User Expired")
	}

	acct_policy := p.Product.Policy
	if libs.In(acct_policy ,r.PPMonth,r.BOMonth) {
		if libs.IsExpire( p.User.ExpireDate) {
			p.Resp.Attrs["Framed-Pool"] = ""
		}
	}else if libs.In(acct_policy,r.PPTimes,r.PPFlow) {
		if p.User.Balance <= 0 {
			return p.ErrorAuth("User Balance lack")
		}
		
	}else if acct_policy == r.BOTimes {
		if p.User.TimeLength <= 0 {
			return p.ErrorAuth("User Time lack")
		}
	}else if acct_policy ==r.AwesomeFee {
		/// 更新用户的过期时间,因为,初创时,都设成了3000年
		///
		UpdateUserExpire(p.User)
		if p.User.TimeLength <= 0{
			return p.ErrorAuth("User Time lack")
		}
		if p.User.FlowLength <= 0{
			return p.ErrorAuth("User Flow lack")
		}
		if libs.IsExpire(p.User.ExpireDate) {
			p.Resp.Attrs["Framed-Pool"] = ""
		}
		
	}else if acct_policy == r.AwesomeFeeBoTime {
		/// 更新用户过期先
		UpdateUserExpire(p.User)
		if p.User.TimeLength <=0 {
			return p.ErrorAuth("User Time lack")
		}
		
	}
	
	/// check concur number

	
	
	return true
}
