package main

import (
	"bytes"
	"flag"
	"fmt"
	//"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

func addLines(str string) []string {
	lines := strings.Split(str, ",")
	//fmt.Printf("%#v\n", lines)
	return lines
}

func unique(lines []string) (strs []string) {
	m := make(map[string]int)
	for _, line := range lines {
		if _, ok := m[line]; ok {
			m[line] = m[line] + 1
		} else {
			m[line] = 1
		}
	}

	for k, _ := range m {
		line := string(k)
		if line != "" {
			strs = append(strs, line)
		}
	}

	return strs
}

func writeString(file *os.File, strs []string) {
	str := strings.Join(strs, "\n")
	str = fmt.Sprintf("%s\n", str)
	_, err := file.Seek(0, 0)
	if err != nil {
		log.Fatalf("Seek:%s\n", err.Error())
	}

	err = file.Truncate(0)
	if err != nil {
		log.Fatalf("WriteString:%s\n", err.Error())
	}

	_, err = file.WriteString(str)
	if err != nil {
		log.Fatalf("WriteString:%s\n", err.Error())
	}
}

func sortStrs(strs []string) []string {
	sort.Sort(sort.Reverse(sort.StringSlice(strs)))
	return strs
}

func readAll(file *os.File) (lines []string) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("ReadAll file: %s\n", err.Error())
	}
	for _, v := range bytes.SplitN(data, []byte("\n"), -1) {
		lines = append(lines, string(v))
	}

	return lines
}

var (
	filePath    *string
	file        *os.File
	uniqueFlag  *bool
	sortFlag    *bool
	addLinesStr *string
)

func init() {
	filePath = flag.String("f", "", "set file filePath")
	//unique
	uniqueFlag = flag.Bool("u", false, "unique")
	//sort
	sortFlag = flag.Bool("s", false, "sort")
	//add lines
	addLinesStr = flag.String("a", "", "add lines eg: line_1,line_2,line_3,line_4")

	flag.Parse()
	if *filePath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	file, err := os.OpenFile(*filePath, os.O_RDWR, 0777)
	if err == nil {
		defer file.Close()
	} else {
		log.Fatalf("OpenFile:%s \n", err.Error())
	}

	mainFlag := false
	oldStrs := readAll(file)
	var newStrs []string
	newStrs = oldStrs

	if *addLinesStr != "" {
		for _, line := range addLines(*addLinesStr) {
			newStrs = append(newStrs, line)
		}
		mainFlag = true
	}

	if *uniqueFlag {
		newStrs = unique(newStrs)
		mainFlag = true
	}

	if *sortFlag {
		newStrs = sortStrs(newStrs)
		mainFlag = true
	}

	if !mainFlag {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if len(newStrs) != 0 {
		writeString(file, newStrs)
	}
}
