package main

import (
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"net"
	"os"
	"runtime"
	"time"
)

func SetCacheByKey(key string, value string) (err error) {
	netConn, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		return err
	}
	redisConn := redis.NewConn(netConn, 200*time.Millisecond, 200*time.Millisecond)
	ok, err := redis.Bool(redisConn.Do("SET", key, value))
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("set not ok")
	}
	return err
}
func GetKeyByCache(key string) (value string, err error) {
	netConn, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		return
	}
	redisConn := redis.NewConn(netConn, 200*time.Millisecond, 200*time.Millisecond)
	value, err = redis.String(redisConn.Do("GET", key))
	if err != nil {
		return
	}
	return
}

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

const (
	OK_EXIT  = 0
	ERR_EXIT = 1
)

var (
	filePath string
)

func init() {
	flag.StringVar(&filePath, "f", "", "set file path")
	flag.Parse()

	if filePath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

}

func main() {
	fmt.Printf("==============================================================\n")

	fmt.Printf("==============================================================\n")
}
