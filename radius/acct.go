package radius

import (
//	"flag"
//	"fmt"
	"log"
//	"os"
//	"os/exec"
//	"strings"
	//	"unicode"
	//	"layeh.com/radius"
	"github.com/cuu/radius"
	"github.com/cuu/softradius/radius/process"
	
)


func BeAcctServer( secret string ) {

	log.Println("rad acct server starting")
	
	entry := &process.Entry{}
	acct_server := radius.Server{
		Handler:	 radius.HandlerFunc(entry.AcctHandler),
		Secret:		[]byte(secret),
		Dictionary: radius.Builtin,
		Addr:		":1813",
	}

	clm := GetClientsMap()
	log.Println(clm)
	acct_server.AddClientsMap( clm )
	
	if err := acct_server.ListenAndServe(); err != nil{
		log.Fatal(err)
	}

}
