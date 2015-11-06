package main

import (
	"fmt"
	"github.com/mmitton/ldap"
	"log"
	"os"
	"runtime"
	"time"
)

func CheckErr(err error, operating string) {
	//pc,file,line,ok = runtime.Caller(1)
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	if err != nil {
		log.Printf("@@@ ERROR: |%s| %s failed.\n", funcName, operating)
		log.Printf("  %s\n", err.Error())
		//panic(err)
		os.Exit(-1)
	}
	log.Printf("### OK: |%s| %s success.\n", funcName, operating)
}

type T struct {
	createTime time.Time
	startTime  time.Time
	endTime    time.Time
}

func NewT() *T {
	return &T{
		createTime: time.Now(),
	}
}
func (this *T) Show() {
	fmt.Printf("createTime:%v\n", this.createTime)
}

func init() {

}

func main() {
	fmt.Printf("==============================================================\n")
	//		base_dn,
	//		ScopeWholeSubtree, DerefAlways, 0, 0, false,
	//		filter[0],
	//		attributes,
	//		nil
	base_db := "dc=xiaoju,dc=com"

	search_request := ldap.NewSearchRequest()
	fmt.Printf("==============================================================\n")
}
