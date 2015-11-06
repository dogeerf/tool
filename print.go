package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	//"strconv"
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

func main() {
	fmt.Printf("==============================================================\n")
	i := 1
	for {
		fmt.Printf("%d %s \r", i, "                                        ")
		i++
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("==============================================================\n")
}
