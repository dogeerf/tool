package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"net"
	"os"
	"runtime"
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

func setColour(text string, colour string) (colourText string) {
	m := make(map[string]int)
	m["black"] = 30
	m["red"] = 31
	m["green"] = 32
	m["yellow"] = 33
	m["blue"] = 34
	m["purple"] = 34
	m["dark green"] = 36
	m["dark_green"] = 36
	m["white"] = 37
	if textColourNum, exist := m[colour]; exist {
		colourText = fmt.Sprintf("\033[%d;1m%s\033[0m", textColourNum, text)
	} else {
		colourText = text
	}
	return
}

func main() {
	fmt.Printf("==============================================================\n")
	redisAddr := "localost:7003"
	conn, err := net.Dial("tcp", redisAddr)
	if err != nil {
		fmt.Printf("%s: conn redis :%s\n", setColour("WARNING", "red"), redisAddr)
	}
	c := redis.NewConn(conn, 2e7, 2e7)
	// get all keys *
	// dump value save to hd
	c.Do()
	_ = c
	fmt.Printf("==============================================================\n")
}
