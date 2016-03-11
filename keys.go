package main

import (
	"flag"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"hash/crc32"
	"log"
	"net"
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
	//log.Printf("### OK: |%s| %s success.\n", funcName, operating)
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

func ShowKeys(keys map[uint32]string) {
	for value := 0; value < 1024; value++ {
		key := keys[uint32(value)]
		fmt.Printf("slot=%-4d value=%-4d key=%s\n", value, value, key)
	}
}

func SetKeys(keys map[uint32]string) {
	//for key, value := range keys {
	if len(keys) != 1024 {
		fmt.Printf("ERROR: len(keys) != 1024\n")
		return
	}

	for value := 0; value < 1024; value++ {
		key := keys[uint32(value)]

		exist, err := redis.Bool(redisConn.Do("EXISTS", key))
		CheckErr(err, "check key exist")
		if exist {
			fmt.Printf("FIND KEY: slot=%-4d value=%-4d key=%s\n", value, value, key)
			continue
		}

		ok, err := redis.String(redisConn.Do("SET", value, key))
		CheckErr(err, "set key")
		if ok == "OK" {
			fmt.Printf("SET KEY: slot=%-4d value=%-4d key=%s\n", value, value, key)
		} else {
			fmt.Printf("SET KEY ERR: slot=%-4d value=%-4d key=%s\n", value, value, key)
		}
	}
}

func GetRedisConn(redisAddr string) (redisConn redis.Conn) {
	conn, err := net.Dial("tcp", redisAddr)
	CheckErr(err, "conn redis error")
	redisConn = redis.NewConn(conn, 200*time.Millisecond, 200*time.Millisecond)
	return redisConn
}

const (
	OK_EXIT  = 0
	ERR_EXIT = 1
)

var (
	redisHost string
	redisPort int
	redisAddr string
	showKeys  bool
	redisConn redis.Conn
)

func init() {
	flag.StringVar(&redisAddr, "addr", "", "set redis host:port")
	flag.StringVar(&redisHost, "h", "localhost", "set redis host")
	flag.IntVar(&redisPort, "p", 3000, "set redis port")
	flag.BoolVar(&showKeys, "s", false, "set redis host:port")
	flag.Parse()

	if showKeys {
		keys := GetKeys()
		ShowKeys(keys)
		os.Exit(0)
	}

	if redisAddr == "" && redisHost != "" && redisPort != 0 {
		redisAddr = fmt.Sprintf("%s:%d", redisHost, redisPort)
	}

	if redisAddr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	redisConn = GetRedisConn(redisAddr)
}

func main() {
	keys := GetKeys()
	SetKeys(keys)

	defer redisConn.Close()
}
