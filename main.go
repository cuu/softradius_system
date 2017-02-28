package main

import (
	r "github.com/cuu/softradius/routers"
	. "github.com/cuu/softradius/controllers"
//	"github.com/cuu/softradius/libs"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"fmt"
	"flag"
	"errors"
	"runtime"
	//	"os"
 	re "gopkg.in/gorethink/gorethink.v3"
	rdb "github.com/cuu/softradius/database/shelf"
	rad "github.com/cuu/softradius/radius"
	
)


/*
func page_not_found(rw http.ResponseWriter, r *http.Request){
	t,_:= template.New("404.html").ParseFiles(beego.ViewsPath+"/404.html")
	data :=make(map[string]interface{})
	data["content"] = "page not found"
	t.Execute(rw, data)
}
*/

func GuuRecoverPanic(ctx *context.Context) {
	ErrAbort := errors.New("User stop run")
	if err := recover(); err != nil {
		if err == ErrAbort {
			return
		}

		var stack string
		logs.Critical("the request url is ", ctx.Input.URL())
		logs.Critical("Handler crashed with error", err)
		stack += fmt.Sprintln("the request url is ", ctx.Input.URL() )
		stack += fmt.Sprintln("Handler crashed with error: ", err)
		
		for i := 1; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			logs.Critical(fmt.Sprintf("%s:%d", file, line))
			stack = stack + fmt.Sprintln(fmt.Sprintf("%s:%d", file, line))
		} 
		
		if ctx.Output.Status != 0 {
			ctx.ResponseWriter.WriteHeader(ctx.Output.Status)
		} else {
			ctx.ResponseWriter.WriteHeader(500)
		}
		stack += "SoftRadius "
		ctx.WriteString(stack)
	}
	
}

func before_run_beego() {
	//after all init() called
	var menus = []string{ r.MenuSys,r.MenuBus,r.MenuOpt,r.MenuPlugin,r.MenuAgency, }
	r.Permits.Build_menus( menus )
	
	//还要  bind_super ,super是从 数据库得到 operator_type == 0 的用户们
	
	
	err := rdb.Register(re.ConnectOpts{
		Address: "localhost:28015",
	}, "SoftRadius")
	if err != nil {
		panic("Db connect failed...")
	}

	opera := &Operators{}
	err = rdb.DataBase().FilterOne(opera,map[string]int{"Type":0})
	if err == nil {
		r.Permits.Bind_super(opera.Name)
	}else{
		panic("no super admin ")
	}

	beego.ErrorController(&ErrorController{})

	
}


func run_beego(){

	beego.BConfig.WebConfig.DirectoryIndex = true
	/*
	beego.BConfig.Listen.AdminEnable = true
	beego.BConfig.Listen.AdminAddr = "localhost"
	beego.BConfig.Listen.AdminPort = 8088
	*/
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "guugosessionID"
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600

	beego.BConfig.RecoverFunc = GuuRecoverPanic
	beego.SetStaticPath("/AdminLTE", "static/AdminLTE")
//	fmt.Println(beego.AppConfig.DefaultString("DEFAULT::Secret","NULL"))
	//beego.ErrorHandler("404", page_not_found)
	before_run_beego()
	beego.Run()	
}



func main() {
	
	admin := flag.Bool("admin",false,"Run admin interface")
	radius := flag.Bool("radius",false, "Run radius server")
	radacct := flag.Bool("radacct",false,"Run radius acct server")
	
	flag.Parse()
		
	if *admin == true {
		run_beego()
	}else {
		if *radius == true {
		//	fmt.Println("Run radius server....")
			rad.BeAuthServer("testing123")
		}else if *radacct == true {
			rad.BeAcctServer("testing123")
		}
		
	}
	
}

