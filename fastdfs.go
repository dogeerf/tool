package main

import (
	"fmt"
	fdfs "github.com/weilaihui/fdfs_client"
	"os"
	"runtime"
	//"strconv"
	//"strings"
	//"time"
)

func CheckErr(err error, operating string) {
	//pc,file,line,ok = runtime.Caller(1)
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	if err != nil {
		fmt.Printf("@@@ ERROR: |%s| %s failed.\n", funcName, operating)
		fmt.Printf("  %s\n", err.Error())
		//panic(err)
		os.Exit(-1)
	}
	fmt.Printf("### OK: |%s| %s success.\n", funcName, operating)
}

func fdfs_check(host string, port int, filePath string) {
	var hostList = []string{host}
	tracker := &fdfs.Tracker{
		HostList: hostList,
		Port:     port,
	}
	fdfsClient, err := fdfs.NewFdfsClientByTracker(tracker)
	CheckErr(err, "create fdfs Client")

	uploadResponse, err := fdfsClient.UploadByFilename(filePath)
	CheckErr(err, "upload file "+filePath)

	//fmt.Println(uploadResponse.GroupName)
	fmt.Println("uploadResponse.RemoteFileId: ", uploadResponse.RemoteFileId)

	fdfsClient.DeleteFile(uploadResponse.RemoteFileId)
	CheckErr(err, "delete file "+filePath)
}

func main() {
	fdfs_check("duke.jie.c", 33133, "/root/work/go/tool/file")
}
