package main

import (
	"encoding/base64"
	"fmt"
	"log"
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

func base64UrlEncode(src string) (dst string) {
	return base64.URLEncoding.EncodeToString([]byte(src))
}

func base64UrlDecode(src string) (dst string, err error) {
	data, err := base64.URLEncoding.DecodeString(src)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func main() {
	var src, dst, op string
	if len(os.Args) < 2 {
		dst = ""
	} else if len(os.Args) == 2 {
		src = os.Args[1]
		dst = base64UrlEncode(src)
	} else if len(os.Args) >= 3 {
		op = os.Args[1]
		src = os.Args[2]
		if op == "-d" || op == "-D" || op == "d" || op == "D" {
			src = os.Args[2]
			dst, _ = base64UrlDecode(src)
		} else if op == "-e" || op == "-E" || op == "e" || op == "E" {
			src = os.Args[2]
			dst = base64UrlEncode(src)
		} else {
			dst = ""
		}
	}

	fmt.Printf("%s", dst)
}
