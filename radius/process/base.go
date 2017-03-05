package process

import (
	re "gopkg.in/gorethink/gorethink.v3"	
	rdb "github.com/cuu/softradius/database/shelf"
	//	ctl "github.com/cuu/softradius/controllers"
	//	"fmt"
	"github.com/cuu/radius"
	"github.com/cuu/softradius/libs"
	"log"
	"fmt"
	
)

const (
	STATUS_TYPE_NONE         = 0       
	STATUS_TYPE_START        = 1        
	STATUS_TYPE_STOP         = 2        
	STATUS_TYPE_UPDATE       = 3    
	STATUS_TYPE_UNLOCK       = 4      
	STATUS_TYPE_CHECK_ONLINE = 5
	STATUS_TYPE_ACCT_ON      = 7
	STATUS_TYPE_ACCT_OFF     = 8
)


const (
//	# Packet codes
	AccessNone    = 0
	AccessRequest = 1
	AccessAccept = 2
	AccessReject = 3
	AccountingRequest = 4
	AccountingResponse = 5
	AccessChallenge = 11
	StatusServer = 12
	StatusClient = 13
	DisconnectRequest = 40
	DisconnectACK = 41
	DisconnectNAK = 42
	CoARequest = 43
	CoAACK = 44
	CoANAK = 45

)


type ProcessResp struct {
	Attrs map[string]interface{}
	Code int
	
}

type ProcessEntry interface{
	AuthProcess(user string, pass string) bool
	AuthHandler(w radius.ResponseWriter, p *radius.Packet)
	AcctProcess(user string) bool
	AcctHandler(w radius.ResponseWriter, p *radius.Packet)
	
}


type Entry struct{
	Rw   *radius.ResponseWriter
	Pkt  *radius.Packet
	//map[string]interface{} 
}


func (self *Entry) OkAuth() []*radius.Attribute {
	var attributes []*radius.Attribute
	attributes = []*radius.Attribute{
		self.Pkt.Dictionary.MustAttr("Acct-Interim-Interval",uint32(60)),
		self.Pkt.Dictionary.MustAttr("Session-Timeout",uint32(30)),
		self.Pkt.Dictionary.MustAttr("Reply-Message","success"),
	}

	return attributes
}

func (self *Entry)AcctHandler(w radius.ResponseWriter, p *radius.Packet) {
	
	self.Rw  = &w
	self.Pkt = p


	self.AcctProcess()
	
	/*
	for _, attr := range p.Attributes {
		name, ok := p.Dictionary.Name(attr.Type)
		if !ok{
			continue
		}
		
		
		value  := fmt.Sprint(attr.Value)
		log.Printf("%s %s",name,value)
	}

	
	
	var attributes []*radius.Attribute
	attributes = []*radius.Attribute{
		p.Dictionary.MustAttr("Reply-Message", "Done"),
	}

	w.AccountingResponse(attributes...)
*/
}

func (self *Entry) AuthHandler(w radius.ResponseWriter, p *radius.Packet) {
	username, password, ok := p.PAP() //PAP MSCHAPV1 V2 Needs decryption
	if !ok {
		w.AccessReject()
		return
	}

	fmt.Println( string(p.Secret))

	
	self.Rw = &w
	self.Pkt = p
	log.Printf("%s with %s requesting access (%s #%d)\n", username,password, w.RemoteAddr(), p.Identifier)

	env := make(map[string]string)

	for _, attr := range p.Attributes {
		fmt.Println(attr.Type)
		name, ok := p.Dictionary.Name(attr.Type)
		if !ok {
			continue
		}
		value := fmt.Sprint(attr.Value)
		env[name] = value
	}
	
	fmt.Println(env)
		//env = append(env, "RADIUS_USERNAME="+username, "RADIUS_PASSWORD="+password)

	self.AuthProcess(username,password)
		
}


func (self *Entry) AcctProcess() bool {
	m := make(map[string]interface{})
	
	for _,attr := range self.Pkt.Attributes {
		name, ok := self.Pkt.Dictionary.Name(attr.Type)
		if !ok {
			continue
		}
		
		value := fmt.Sprint(attr.Value)
		//codec_type := self.Pkt.Dictionary.Codec(attr.Type).String()
		//fmt.Println(name, " ", attr)
		log.Printf("%s %s",name,value)
		m[name] = attr.Value
		fmt.Println(libs.Type(attr.Value))//=>string,slice,uint32
		
	}

	ac := NewAcctUser()
	ac.Data = m

	ac.GetUser()
	ac.AcctOnOff()
	ac.AcctStart()
	ac.AcctStop()
	ac.AcctUpdate()
	
	var attributes []*radius.Attribute
	attributes = []*radius.Attribute{
		self.Pkt.Dictionary.MustAttr("Reply-Message", "Done"),
	}

	(*self.Rw).AccountingResponse(attributes...)
	return true
}

func (self *Entry) AuthProcess(user string,pass string) bool {
	au := NewAuthUser()
	au.GetUser(user,pass) 
	au.GetProduct()
	au.Billing()

	
	if au.Resp.Code == 3 {
		log.Printf("%s rejected (%s #%d)\n", user, (*self.Rw).RemoteAddr(), self.Pkt.Identifier)
		var attributes []*radius.Attribute
		for i,v := range au.Resp.Attrs {
			switch libs.Type(v) {
			case "string":
				attributes = append(attributes,self.Pkt.Dictionary.MustAttr(i,v))
			case "int":
				attributes = append(attributes,self.Pkt.Dictionary.MustAttr(i,uint32(v.(uint32))))
			}
		}
		
		(*self.Rw).AccessReject(attributes...)
		
	}else if au.Resp.Code == 2 {
		attrs := self.OkAuth()		
		log.Printf("%s accepted (%s #%d)\n", user, (*self.Rw).RemoteAddr(), self.Pkt.Identifier)
		(*self.Rw).AccessAccept(attrs...)		
	}

	return true
}



//--------------------------------------------------------------------------
func init_database() {
	err := rdb.Register(re.ConnectOpts{
		Address: "localhost:28015",
	}, "SoftRadius")
	if err != nil {
		panic("Db connect failed...")
	}
	
}

func init(){
	init_database()
}
