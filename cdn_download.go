package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"errors"
	"github.com/bitly/go-simplejson"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
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

func init() {

}

func RedirectPolicy(req *http.Request, via []*http.Request) error {
	return errors.New("Stop Redirect")
}

func GetHeaderByHead(url, key string) (value string) {
	client := &http.Client{CheckRedirect: RedirectPolicy}

	for i, n := 1, 20; i <= n; i++ {
		resp, _ := client.Head(url)
		println(i, resp.StatusCode, url)
		if resp.StatusCode != 302 {
			resp.Body.Close()
			//println("err: " + err.Error())
			time.Sleep(1 * time.Second)
			continue
		}
		defer resp.Body.Close()
		value = resp.Header.Get(key)
		break
	}
	return
}

func CreateDir(dir string) error {
	//check dir exist
	if file, err := os.Stat(dir); err == nil {
		if file.IsDir() {
			return nil
		} else {
			return errors.New("find file : not dir")
		}
	}
	//create dir
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return nil
}

func showDownLoadProgress(path string) {
	var thisSize int64
	var lastSecondSize int64
	var oneSecondDownSize int64
	file, err := os.OpenFile(path, os.O_RDONLY, 755)
	defer file.Close()
	if err != nil {
		println(err.Error())
		return
	}

	for {
		stat, err := file.Stat()
		if err != nil || lastSecondSize == 0 {
			continue
		}
		thisSize = stat.Size()
		oneSecondDownSize = thisSize - lastSecondSize
		lastSecondSize = thisSize
		println(oneSecondDownSize)

		time.Sleep(1 * time.Second)
	}
}

func DownLoadUrl(url, path string) (err error) {
	//err = CreateDir(path)
	//if err != nil {
	//	return err
	//}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	//defer file.Close()
	stat, err := file.Stat()
	fileSizeStr := strconv.FormatInt(stat.Size(), 10)
	fmt.Printf("fileSizeStr:%d m\n", stat.Size()/1024/1024)

	client := &http.Client{}
	//client.Timeout = 5 * time.Second
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", "spriteapp client")
	req.Header.Add("Range", "bytes="+fileSizeStr+"-")

	//sent req
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	contentLength := resp.Header.Get("Content-Length")
	if contentLength != "" {
		n, _ := strconv.Atoi(contentLength)
		fmt.Printf("contentLength: %d m \n", n/1024/1024)
	}

	if resp.StatusCode == 416 {
		return nil
	}
	go func(file *os.File) {
		var thisSize int64
		var lastSecondSize int64
		var oneSecondDownSize int64
		tick := time.Tick(1000 * time.Millisecond)
		println("=================================")
		for {
			select {
			case <-tick:
				stat, err := file.Stat()
				if err != nil {
					continue
				}
				thisSize = stat.Size()
				oneSecondDownSize = thisSize - lastSecondSize

				fmt.Printf("thisSize: %d m \toneSecondDownSize: %d k %s \r", thisSize/1024/1024, oneSecondDownSize/1024, "               ")
				//fmt.Printf("oneSecondDownSize: %d k\r", oneSecondDownSize/1024)

				lastSecondSize = thisSize
				//fmt.Println("1000 tick.")
			}
		}
	}(file)
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func GetBodyByUrl(url string) (data []byte, err error) {
	var maxTry int = 10
	var errSleep time.Duration = 1
	for i := 1; i <= maxTry; i++ {
		fmt.Printf("try get url:%d\n", i)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			fmt.Printf("http status:%d\n", resp.StatusCode)
			time.Sleep(errSleep * time.Second)
			continue
		}

		defer resp.Body.Close()
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		} else {
			break
		}
	}
	return
}

func DownloadCdnLog(url, userName, passWord, channels string, num int) (logUrl string) {
	if num <= 0 {
		num = 1
	}

	num = num - 2*num

	dateTime := time.Now().AddDate(0, 0, num).Format("2006-01-02")
	startTime := dateTime + "-0000"
	endTime := dateTime + "-2330"

	var location string
	location = GetHeaderByHead(url, "Location")
	if location == "" {
		println("ERROR: location is null:", url)
		os.Exit(1)
	}

	url = location + "?" + "u=" + userName + "&p=" + passWord
	location = GetHeaderByHead(url, "Location")
	if location == "" {
		println("ERROR: location is null:", url)
		os.Exit(1)
	}

	url = location + "&start_time=" + startTime + "&end_time=" + endTime + "&channels=" + channels

	data, err := GetBodyByUrl(url)
	CheckErr(err, "Get url: "+url)

	//println("date: ", string(data))
	fmt.Printf("==============================================================\n")
	js, err := simplejson.NewJson(data)
	CheckErr(err, "Parse json")
	fmt.Printf("==============================================================\n")
	//fmt.Printf("%#v", js)
	fmt.Printf("==============================================================\n")
	logUrl, err = js.Get("logs").GetIndex(0).Get("files").GetIndex(0).Get("url").String()
	CheckErr(err, "get json")
	//println(logUrl)
	//fmt.Printf("==============================================================\n")

	fileDate := time.Now().AddDate(0, 0, num).Format("20060102")
	println(logUrl)
	// http://dx.wslog.chinanetcenter.com/log/sprite/wvideo.spriteapp.cn/2015-02-25-0000-2330_wvideo.spriteapp.cn.cn.log.gz?wskey=30110003578c72be00f05b464d575400096de28e160041c2
	ext := "gz"
	fileName := channels + ".log" + "-" + fileDate + "." + ext
	println(fileName)
	err = DownLoadUrl(logUrl, fileName)
	CheckErr(err, "download file")
	return
}

func Usag() {
	fmt.Printf("===========================\n")
	fmt.Printf("  Usag:\n")
	fmt.Printf("    %s domain num\n", os.Args[0])
	fmt.Printf("  eg:\n")
	fmt.Printf("    %s wvideo.spriteapp.cn 1\n", os.Args[0])
	fmt.Printf("===========================\n")
	os.Exit(1)
}

func main() {
	if len(os.Args) != 3 {
		Usag()
	}
	//channels := "wvideo.spriteapp.cn"
	channels := os.Args[1]
	num, err := strconv.Atoi(os.Args[2])
	CheckErr(err, "strconv str to num")

	url := "http://dx.wslog.chinanetcenter.com/logQuery/access"
	userName := "sprite"
	passWord := "sprite_AABB"

	logUrl := DownloadCdnLog(url, userName, passWord, channels, num)
	println(logUrl)
	fmt.Printf("==============================================================\n")
}
