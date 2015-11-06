package main

import (
	"fmt"
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

type IFunc func() string

type T struct {
	createTime time.Time
	startTime  time.Time
	endTime    time.Time
	publicfun  IFunc
}

func NewT(i IFunc) *T {
	return &T{
		publicfun: i,
	}
}

func (this *T) t_rb() string {
	println("i am t_rb")
	return this.publicfun
}

func f_1() string {
	println("i am f1")
	return "ok"
}

func main() {
	fmt.Printf("==============================================================\n")
	t := NewT(f_1)
	t.t_rb()
	fmt.Printf("==============================================================\n")
}
