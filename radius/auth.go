package radius

import (
//	"flag"
//	"fmt"
	"log"
//	"os"
//	"os/exec"
//	"strings"
//	"unicode"

	//	"github.com/blind-oracle/go-radius"
	//	"layeh.com/radius"
	"github.com/cuu/radius"
 //	re "gopkg.in/gorethink/gorethink.v3"	
//	rdb "github.com/cuu/softradius/database/shelf"
//	ctl "github.com/cuu/softradius/controllers"

	"github.com/cuu/softradius/radius/process"
)

/* 处理逻辑
首先取得用户名
其次取得套餐信息
根据套餐信息,比对用户的属性
判断返回accept或是reject

*/

func BeAuthServer(secret string ) {

	entry := &process.Entry{}
	
	server := radius.Server{
		Handler:    radius.HandlerFunc(entry.AuthHandler),
		Secret:     []byte(secret),
		Dictionary: radius.Builtin,
		Addr:		":1812",
	}

	clm := GetClientsMap()
	log.Println(clm)
	server.AddClientsMap( clm )
	
	log.Println("Radauth server starting ", secret)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}


}
