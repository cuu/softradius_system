# SoftRadius  a golang based radius accounting system
# 大宝剑计费系统 ,基于go语言,beego框架,AdminLTE UI

## Requirements
* beego 1.7+
* cuu radius libs go get -u github.com/cuu/radius
* AdminLTE
* rethinkdb
* gorethink v3


## License

MPL 2.0

## Talk
* irc.freenode.net  #softradius

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

![Operators](screenshots/operators.png?raw=true "operators")
![products](screenshots/products.png?raw=true   "Products")
![member](screenshots/member_quick.png?raw=true "Members")
![online](screenshots/online.png?raw=true  "Online")
![Bas](screenshots/bas.png?raw=true "Bas")

## Development status
- [x] Opertors, admin fully worked,now going for the privilege of normal operator and agency
- [x] Online user ,accouting worked,
- [x] Radius auth ,auth user from database, have not  implemented kickoff 
- [x] Radius acct ,accouting flow and time to User self,will add acct tickets for logging
- [ ] Agency
- [x] Members , only quick open, no searching
- [ ] Multi Vendor radius support

Still needs a lot of test ,right now just very beginning of the project

