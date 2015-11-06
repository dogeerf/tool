package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strings"
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

func main() {
	fmt.Printf("==============================================================\n")

	//var name []string

	ip := "192.168.133.100"
	name := ip
	names, err := net.LookupAddr(ip)
	if err == nil || len(names) >= 1 {
		name = names[0]
		name = strings.Split(name, ".")[0]
	}

	fmt.Printf("%v\n", name)
	fmt.Printf("==============================================================\n")
}
