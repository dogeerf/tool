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

func del() {
	fmt.Printf("123")
}

func add() {

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
	filePath   *string
	file       *os.File
	uniqueFlag *bool
	sortFlag   *bool
)

func init() {
	filePath = flag.String("f", "", "set file filePath")
	//unique
	uniqueFlag = flag.Bool("u", false, "unique")
	//sort
	sortFlag = flag.Bool("s", false, "sort")

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

	oldStrs := readAll(file)
	var newStrs []string

	if *uniqueFlag {
		newStrs = unique(oldStrs)
	} else {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *sortFlag {
		newStrs = sortStrs(newStrs)
	}

	//	for _, str := range newStrs {
	//		println(str)
	//	}
	if len(newStrs) != 0 {
		writeString(file, newStrs)
	}
}
