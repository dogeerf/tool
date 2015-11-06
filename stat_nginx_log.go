package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strconv"
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
	if funcName != "main.main" {
		log.Printf("### OK: |%s| %s success.\n", funcName, operating)
	}
}

func StrToInt(v_str string) (v_int int) {
	if v_str == "" {
		v_int = 0
	} else {
		v_int, err := strconv.Atoi(v_str)
		if err != nil {
			v_int = 0
		}
		return v_int
	}
	return v_int
}

func StrToFloat64(v_str string) (v_float64 float64) {
	if v_str == "" {
		v_float64 = 0
	} else {
		v_float64, err := strconv.ParseFloat(v_str, 32)
		if err != nil {
			v_float64 = 0
		}
		return v_float64
	}
	return v_float64
}

func parse(str string) {
	lineArr := reLine.FindStringSubmatch(str)
	if len(lineArr) != 17 {
		println("waring: reLine error")
		return
	}
	requesHost := lineArr[1]
	//requesClientIp := lineArr[2]
	//requesTime, err := time.Parse("[02/Jan/2006:15:04:05 +0800]", lineArr[5])
	//requesTime = time.Now()
	//CheckErr(err, "Parse nginx time")

	requestAll := strings.Split(strings.Replace(lineArr[6], "\"", "", -1), " ")
	if lineArr[6] == "-" || len(requestAll) != 3 {
		println("waring: req error")
		return
	}
	requestMethod := requestAll[0]
	requesUrl := requestAll[1]
	//requesHttpVersion := requestAll[2]
	responseStatus := lineArr[7]
	//responseSize := lineArr[8]
	//responseUrl := lineArr[9]
	//requesUserAgent := lineArr[10]
	responseUseTimeSec := lineArr[16]

	responseUseTimeSec_flalt64 := StrToFloat64(responseUseTimeSec)
	//responseSize_int := StrToInt(responseSize)
	responseStatus_int := StrToInt(responseStatus)

	urlLine := "http://" + requesHost + requesUrl
	//println("urlLine: ", urlLine)
	u, err := url.Parse(urlLine)
	if err != nil {
		CheckErr(err, "Parse url")
	}
	//fmt.Printf("r:\n%#v\n", u)
	v := u.Query()
	//fmt.Printf("v:\n%#v\n", v)
	rv := v.Get("r")
	flag := fmt.Sprintf("%-4s %-3s  %s %s\n", requestMethod, responseStatus, u.Path, rv)

	if responseStatus_int >= 500 {
		ErrStatistics[flag] = ErrStatistics[flag] + 1
	} else if responseUseTimeSec_flalt64 > 1 {
		TimeOutStatistics[flag] = TimeOutStatistics[flag] + 1
	}

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

func SendMail(to, subject, msg string) {
	host := "192.168.133.187"
	port := "4000"

	to = base64UrlEncode(to)
	subject = base64UrlEncode(subject)
	msg = base64UrlEncode(msg)

	url := "http://" + host + ":" + port + "/?" + "to=" + to + "&subject=" + subject + "&msg=" + msg
	fmt.Println(url)
	resp, err := http.Get(url)
	CheckErr(err, "send mial")
	data, err := ioutil.ReadAll(resp.Body)
	CheckErr(err, "= =")
	_ = data
	fmt.Printf(string(data))
}

type Item struct {
	Key string
	Val int
}

func SortMapByValue(m map[string]int) (l []Item) {
	for k, v := range m {
		l = append(l, Item{Key: k, Val: v})
	}

	for i := 0; i+1 < len(l); i++ {
		if i+1 == len(l) {
			break
		}

		for j := i; j+1 < len(l); j++ {
			if l[i].Val < l[j+1].Val {
				l[i], l[j+1] = l[j+1], l[i]
			}
		}

	}
	return l
}

var reLine *regexp.Regexp
var err error
var ErrStatistics map[string]int
var TimeOutStatistics map[string]int
var nginxLogPath string

func init() {
	reLine, err = regexp.Compile(`(?U)(^.*) (.*) (.*) (.*) (\[.*\]) (\".*\") (.*) (.*) (\".*\") (\".*\") (\".*\") (.*) (.*) (.*) (.*) (.*)$`)
	ErrStatistics = make(map[string]int)
	TimeOutStatistics = make(map[string]int)
	nginxLogPath = "/root/tmp/safe_spriteapp_com_access.log-20150518"
}

func main() {
	fmt.Printf("==============================================================\n")

	//golang
	nginxLogPath := flag.String("f", "./safe_spriteapp_com_access.log", "defaule filePath.")
	flag.Parse()

	date, err := ioutil.ReadFile(*nginxLogPath)
	CheckErr(err, "ReadFile")
	astrings := strings.SplitN(string(date), "\n", -1)

	for _, line := range astrings {
		line = strings.TrimRight(line, "\n")
		parse(line)
	}

	subject := "spriteapp_admin_5XX_stat"
	to := "server@spriteapp.com"
	msg := fmt.Sprintf("AllCount:%d\n\n", len(astrings))

	msg = msg + "ERROR:\n"
	items := SortMapByValue(ErrStatistics)
	for _, v := range items {
		m := fmt.Sprintf("  %-5d %s", v.Val, v.Key)
		msg = msg + m
	}

	msg = msg + "TimeOut:(gt 1 sec)\n"
	items = SortMapByValue(TimeOutStatistics)
	for _, v := range items {
		m := fmt.Sprintf("  %-5d %s", v.Val, v.Key)
		msg = msg + m
	}

	SendMail(to, subject, msg)
	fmt.Printf(msg)

	fmt.Printf("==============================================================\n")
}
