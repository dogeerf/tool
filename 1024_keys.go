package main

import (
	"fmt"
	//	"github.com/garyburd/redigo/redis"
	"flag"
	"hash/crc32"
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

func GetKeys() (keys map[uint32]string) {
	keys = make(map[uint32]string)
	for i := 0; ; i++ {
		value := fmt.Sprintf("%s_%d", "op_test", i)
		hash32 := crc32.ChecksumIEEE([]byte(value))
		key := hash32 % 1024
		if _, find := keys[key]; !find {
			keys[key] = value
			if len(keys) == 1024 {
				return keys
			}
		}
	}
	return keys

}

func SetKeys(keys map[uint32]string) {
	for key, value := range keys {
		fmt.Printf("key=%-4d value=%s\n", key, value)
	}
}

func init() {

}

func main() {
	keys := GetKeys()
	SetKeys(keys)
}
