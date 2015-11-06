package main

import (
	"fmt"
	"strings"
)

func GetFileExtByUrl(url string) (fileName string) {
	tmpArr := strings.SplitN(url, "/", -1)
	tmpArr = strings.SplitN(tmpArr[len(tmpArr)-1], "?", -1)
	tmpArr = strings.SplitN(tmpArr[0], ".", -1)
	fileName = tmpArr[len(tmpArr)-1]
	return
}

func GetFileNameByUrl(url string) (fileName string) {
	tmpArr := strings.SplitN(url, "/", -1)
	tmpArr = strings.SplitN(tmpArr[len(tmpArr)-1], "?", -1)
	tmpArr = strings.SplitN(tmpArr[0], ".", -1)
	fileName = tmpArr[0]
	return
}

func main() {
	//fmt.Println("On Unix:", filepath.SplitList("/a/b/c:/usr/bin"))
	url := "http://dx.wslog.chinanetcenter.com/log/sprite/wvideo.spriteapp.cn/2015-02-25-0000-2330_wvideo.spriteapp.cn.cn.log.gz?wskey=30110003578c72be00f05b464d575400096de28e160041c2&aaa=2?sad"
	fileName := GetFileExtByUrl(url)
	fmt.Printf("fileName:%s\n", fileName)
}
