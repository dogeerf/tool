package main

import (
	"errors"
	"fmt"
	"os"
	//"os/exec"
	"strings"
)

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
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	return nil
}

func GetDirPath(path string) string {
	arr := strings.Split(path, "/")
	dirPath := ""
	for i, v := range arr {
		if i == len(arr)-1 {
			break
		}

		if i == 0 {
			dirPath = v
		} else {
			dirPath = dirPath + "/" + v
		}
	}
	return dirPath
}

func CreateFile(path string) error {
	if file, err := os.Stat(path); err == nil {
		if !file.IsDir() {
			return nil
		}
	}

	if dirPath := GetDirPath(path); dirPath != "./" && dirPath != "" {
		println("dirPath", dirPath)
		err := CreateDir(dirPath)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(path)
	defer f.Close()
	return err
}

func main() {
	var err error
	path := "./tmp/a/v/bbb/a.txt"
	err = CreateFile(path)
	if err != nil {
		fmt.Println(err.Error())
	}

}

