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

func adder() func(int) int {
	sum := 0
	innerfunc := func(x int) int {
		sum += x
		return sum
	}
	return innerfunc
}

func main() {
	fmt.Printf("==============================================================\n")
	pos, neg := adder(), adder()
	for i := 0; i < 10; i++ {
		//fmt.Println(pos(i), neg(-2*i))
		fmt.Println(pos(i), neg(-i))
	}
	fmt.Printf("==============================================================\n")
}
