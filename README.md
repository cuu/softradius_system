# SoftRadius  a golang based radius accounting system
# 大宝剑计费系统 ,基于go语言,beego框架,专注流量与时间

## Requirements
* beego 1.7+
* cuu radius libs go get -u github.com/cuu/radius
* AdminLTE
* rethinkdb
* gorethink v3


## License

MPL 2.0

## Installation
* go get github.com/astaxie/beego
* go get -u github.com/cuu/radius
* go get gopkg.in/gorethink/gorethink.v3
* go get -u github.com/cuu/softradius_system
* mv $GOCODE/src/github.com/cuu/softradius_system $GOCODE/src/github.com/cuu/softradius
* cd $GOCODE/src/github.com/cuu/softradius && go build
* Install rethinkdb and run it anywhere you like,I compiled it from source,latest version 2.3.5
* cd dbfiles && ./import.sh && cd ..
* ./softradius -admin to start web ui
* ./softradius -radius to start radius auth server
* ./softradius -radacct to start radius acct server
* go to http://127.0.0.1:8081 to see the interface,default password would be "admin/admin"

![Dashboard](screenshots/dashboard.png?raw=true "bashboard")
![Operators](screenshots/operators.png?raw=true "operators")
![products](screenshots/products.png?raw=true   "Products")
![member](screenshots/member_quick.png?raw=true "Members")
![agency](screenshots/agency.png?raw=true "Agency")
![online](screenshots/online.png?raw=true  "Online")
![Bas](screenshots/bas.png?raw=true "Bas")
![CreateBatch](screenshots/createbatch?raw=true "Batch")
![BatchRule](screenshots/batchrule.png?raw=true "BatchRule")

## Development status
- [x] Opertors, admin fully working,now going for the privilege of normal operators and agency
- [x] Online user ,accouting works,
- [x] Radius auth ,auth user from database 
- [x] Radius acct ,accouting flow and time to User self,will add acct tickets for logging
- [x] Agency
- [x] Members , only quick open
- [x] Sidbar search
- [ ] Multi Vendor radius support
- [ ] MSCHAP 
- [ ] Portal CMCC V1/V2 HuaWei V1/V2
- [ ] API


