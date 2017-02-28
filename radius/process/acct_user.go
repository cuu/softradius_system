package process


import (
	rdb "github.com/cuu/softradius/database/shelf"
	ctl "github.com/cuu/softradius/controllers"
	"fmt"
	"github.com/cuu/softradius/libs"
//	r "github.com/cuu/softradius/routers"
//	"reflect"
)


type AcctUser struct {
	UserGot bool
	User *ctl.Members
	
	Resp ProcessResp
	Data map[string]interface{}
}

//Constructor
func NewAcctUser()  *AcctUser {
	one := &AcctUser{}
	one.ResetResp()

	return one
}



func (self *AcctUser) ResetResp() {
	self.Resp.Attrs = nil
	self.Resp.Code  = AccessAccept
	self.Resp.Attrs = make(map[string]interface{})
}

func (this *AcctUser) get_data_s(key string ) string{
	if val,ok := this.Data[key]; ok {
		return fmt.Sprintf("%v",val)
	}

	return ""
}

func (this *AcctUser) get_data_i(key string) uint32 {
	if val,ok := this.Data[key]; ok {
		return val.(uint32)
	}

	return 0
}

func (this *AcctUser) get_input_total() int {
	octets := this.get_data_i("Acct-Input-Octets")
	gigas  := this.get_data_i("Acct-Input-Gigawords")

	return int(octets/1024 + gigas*4*1024*1024)
	
}

func (this *AcctUser) get_output_total() int {
	octets := this.get_data_i("Acct-Output-Octets")
	gigas  := this.get_data_i("Acct-Output-Gigawords")

	return int(octets/1024+ gigas*4*1024*1024)
	
}

func (this *AcctUser) FillTicket() ctl.AcctTicket {
	p := ctl.AcctTicket{}
	p.MemberName   = this.get_data_s("User-Name")
	return p
	
}

func (this *AcctUser) FillOnline() ctl.AcctOnline {
	p := ctl.AcctOnline{}
	p.MemberName    = this.get_data_s("User-Name")
	p.NasAddr       = this.get_data_s("NAS-IP-Address")
	p.AcctSessionId = this.get_data_s("Acct-Session-Id")
	p.AcctStartTime = libs.Get_currtime()
	p.FramedIpAddr  = this.get_data_s("Framed-IP-Address")
	p.NasPortId     = this.get_data_s("NAS-Port")
	p.StartSource   = STATUS_TYPE_START
	return p
}

func (this *AcctUser) GetUser() bool {

	if name,ok := this.Data["User-Name"] ; ok {
		one := &ctl.Members{}
		
		err := rdb.DataBase().FilterOne(one,map[string]string{"Name":name.(string)})
		if err == nil {
			this.User  = one
			this.UserGot = true  
			return this.UserGot
		}
	
		this.User = one
		return this.UserGot
	}
	
	return false
	
}


func (self *AcctUser) AcctOnOff() {
	if status_type,ok := self.Data["Acct-Status-Type"]; ok {
		if status_type.(uint32) == STATUS_TYPE_ACCT_ON {
			fmt.Println("acct on")
		}
	}
}

func (self *AcctUser) AcctStart() {
	if status_type,ok := self.Data["Acct-Status-Type"]; ok {
		if status_type.(uint32) == STATUS_TYPE_START {
			one := self.FillOnline()
			rdb.DataBase().FilterInsert(&one,"AcctSessionId")
		}
	}	
}

func (self *AcctUser) AcctStop() {
	if status_type,ok := self.Data["Acct-Status-Type"]; ok {
		if status_type.(uint32) == STATUS_TYPE_STOP {
			one := &ctl.AcctOnline{}
			session_id := self.get_data_s("Acct-Session-Id")
			
			resp,err := rdb.DataBase().FilterDel(one,map[string]string{"AcctSessionId":session_id})
			if err == nil {
				fmt.Println("kill online: ",resp.Deleted)
			}else{
				fmt.Println(err)
				return
			}
			
			user := &ctl.Members{}
			name := self.get_data_s("User-Name")
			rdb.DataBase().FilterOne(user,map[string]string{"Name":name})
			
			session_time  := self.get_data_i("Acct-Session-Time")

			user.InFlow   += self.get_input_total()
			user.OutFlow  += self.get_output_total()
			user.UsedTime += int(session_time)

			fmt.Println(user.InFlow," ", user.OutFlow)
			resp,err = rdb.DataBase().Update(user.Id, user)
			if err == nil {
				fmt.Println("Replaced ",resp.Replaced )
			}else {
				fmt.Println(err)
			}
			
			
		}
	}	
}


func (self *AcctUser) AcctUpdate() {
	if status_type,ok := self.Data["Acct-Status-Type"]; ok {
		if status_type.(uint32) == STATUS_TYPE_UPDATE {
			fmt.Println("AcctUpdate")
			one := &ctl.AcctOnline{}
			session_id := self.get_data_s("Acct-Session-Id")
			err := rdb.DataBase().FilterOne(one,map[string]string{"AcctSessionId":session_id})
			
			if err == nil {
				one.InputTotal  += self.get_input_total()
				one.OutputTotal += self.get_output_total()
				resp,err := rdb.DataBase().Update(one.Id,one)
				if err == nil {
					fmt.Println("Online Replaced ",resp.Replaced)
				}else {
					fmt.Println(err)
				}
			}else {
				fmt.Println(one.InputTotal," ",one.OutputTotal)
				fmt.Println(session_id)
				fmt.Println(err)
			}

			user := &ctl.Members{}
			name := self.get_data_s("User-Name")
			rdb.DataBase().FilterOne(user,map[string]string{"Name":name})
			
			session_time  := self.get_data_i("Acct-Session-Time")

			user.InFlow   += self.get_input_total()
			user.OutFlow  += self.get_output_total()
			user.UsedTime += int(session_time)

			fmt.Println(user.InFlow," ", user.OutFlow)
			resp,err := rdb.DataBase().Update(user.Id, user)
			if err == nil {
				fmt.Println("Member Replaced ",resp.Replaced )
			}else {
				fmt.Println(err)
			}
				
		}
	}	
}
