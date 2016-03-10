package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
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

func GetBytesByFilePath(filePath string) (fileBytes []byte) {
	file, err := os.Open(filePath)
	CheckErr(err, "open file:"+filePath)
	fileBytes, err = ioutil.ReadAll(file)
	CheckErr(err, "read all file:"+filePath)
	file.Close()
	return fileBytes
}

func GetKeyStatByBytes(fileBytes []byte) (keyStats []*KeyStat) {
	dec := json.NewDecoder(bytes.NewReader(fileBytes))
	for {
		var keyStat KeyStat
		if err := dec.Decode(&keyStat); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		keyStats = append(keyStats, &keyStat)
	}
	return keyStats
}

const (
	ERR_NUM = -1
)

var (
	filePath       string
	redisKeyRegexp *regexp.Regexp
	redisStats     []*RedisStat
)

func init() {
	flag.StringVar(&filePath, "f", "", "set file path")
	flag.Parse()

	if filePath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	redisKeyRegexp := regexp.MustCompile("_|-|.|:")
	_ = redisKeyRegexp
}

type KeyStat struct {
	DB         uint32 `json:"db"`
	Type       string `json:"type"`
	ExpireAt   uint64 `json:"expireat"`
	Key        string `json:"key"`
	FieldCount int    `json:"fieldCount"`
	KeySize    int    `json:"keySize"`
	ValueSize  int    `json:"valueSize"`
}

type Node struct {
	Root          bool
	Name          string
	Childrens     map[string]*Node
	KeyStats      []*KeyStat
	ChildrenCount int
	KeyStatCount  int
}

type Tree struct {
	Node
}

func NewTree() *Tree {
	tree := &Tree{}
	tree.Name = "root."
	tree.Root = true
	tree.Childrens = make(map[string]*Node)
	return tree
}

func NewNode(name string) *Node {
	return &Node{
		Name:          name,
		ChildrenCount: 0,
		Childrens:     make(map[string]*Node),
	}
}

func GetExpiredByNode(node *Node) (num int) {
	if node.ChildrenCount == 0 {
		for _, keyStat := range node.KeyStats {
			if keyStat.ExpireAt != 0 {
				num = num + 1
			}
		}
	} else {
		num = num + node.KeyStatCount
		for _, node := range node.Childrens {
			num = GetExpiredByNode(node) + num
		}
	}
	return num
}

func GetKeySizeByNode(node *Node) (size int) {
	if node.ChildrenCount == 0 {
		for _, keyStat := range node.KeyStats {
			size = size + keyStat.KeySize
		}
	} else {
		for _, node := range node.Childrens {
			size = GetKeySizeByNode(node) + size
		}
	}
	return size
}

func GetValueSizeByNode(node *Node) (size int) {
	if node.ChildrenCount == 0 {
		for _, keyStat := range node.KeyStats {
			size = size + keyStat.ValueSize
		}
	} else {
		for _, node := range node.Childrens {
			size = GetValueSizeByNode(node) + size
		}
	}
	return size
}

func GetNumByNode(node *Node) (num int) {
	if node.ChildrenCount == 0 {
		return node.KeyStatCount
	} else {
		num = num + node.KeyStatCount
		for _, node := range node.Childrens {
			num = GetNumByNode(node) + num
		}
	}
	return num
}

func (this *Node) GetNumByChildrens() (num int) {
	if this.ChildrenCount == 0 {
		num = this.KeyStatCount
	} else {
		for _, node := range this.Childrens {
			num = GetNumByNode(node)
		}
	}
	return num
}

type RedisStat struct {
	Path      string
	Childrens int
	Num       int
	Expired   int
	KeySize   int
	ValueSize int
}

func (this RedisStat) Show() {
	keySize := fmt.Sprintf("%dM", this.KeySize/1024/1024)
	valueSize := fmt.Sprintf("%dM", this.ValueSize/1024/1024)
	fmt.Printf("path=%-100s childrens=%-8d num=%-8d expired=%-8d keySize=%-3s valueSize=%-3s\n", this.Path, this.Childrens, this.Num, this.Expired, keySize, valueSize)
}

func (this *Node) Show(path string, count int) (err error) {
	path = fmt.Sprintf("%s%s", path, this.Name)
	if this.ChildrenCount >= count {
		RedisStat := &RedisStat{
			Path:      fmt.Sprintf("%s%s", path, "*"),
			Childrens: this.ChildrenCount,
			Num:       GetNumByNode(this),
			Expired:   GetExpiredByNode(this),
			KeySize:   GetKeySizeByNode(this),
			ValueSize: GetValueSizeByNode(this),
		}
		redisStats = append(redisStats, RedisStat)
	}

	if this.ChildrenCount > 0 {
		for _, node := range this.Childrens {
			node.Show(path, count)
		}
	}

	return err
}

func (this *Node) AddKeyStatToChildren(keyStat *KeyStat) (err error) {
	this.KeyStats = append(this.KeyStats, keyStat)
	this.KeyStatCount = this.KeyStatCount + 1
	return err
}

func (this *Node) AddNodeToChildren(name string) (node *Node, err error) {
	if node, find := this.Childrens[name]; find {
		return node, err
	}

	node = NewNode(name)
	this.Childrens[name] = node
	this.ChildrenCount = this.ChildrenCount + 1
	return node, err
}

func (this *Node) FindChildrenByNode(name string) (node *Node, err error) {
	var find bool
	if node, find = this.Childrens[name]; find {
		return node, err
	}

	return nil, errors.New("no find node")
}

func (this *Tree) AddNodesToTree(field []string) (node *Node, err error) {
	for deep, name := range field {
		if deep == 0 {
			node, err = this.AddNodeToChildren(name)
			CheckErr(err, "add node to tree :")
		} else {
			node, err = node.AddNodeToChildren(name)
			CheckErr(err, "add node to node :")
		}
	}
	return node, err
}

func (this *Tree) FindNodeBYTree(field []string) (node *Node, err error) {
	for deep, name := range field {
		_ = deep
		for _, node := range this.Childrens {
			node, err = node.FindChildrenByNode(name)
			if err != nil {
				return node, err
			}
		}
	}
	return node, err
}

func (this *Tree) AddKeyStatToTree(keyStat *KeyStat) (err error) {
	field := strings.SplitAfter(keyStat.Key, ".")
	node, err := this.AddNodesToTree(field)
	CheckErr(err, "add nodels to tree")
	err = node.AddKeyStatToChildren(keyStat)
	return err
}

func (this *Tree) Show(count int) {
	for _, node := range this.Childrens {
		node.Show(this.Name, count)
	}
}

func ShowTime(msg string) {
	fmt.Printf("%s: %s\n", time.Now().Format("2006-01-02 15:04:05"), msg)
}

func main() {
	tree := NewTree()

	fileBytes := GetBytesByFilePath(filePath)
	keyStats := GetKeyStatByBytes(fileBytes)

	for _, keyStat := range keyStats {
		err := tree.AddKeyStatToTree(keyStat)
		CheckErr(err, "add keyStat to tree")
	}

	tree.Show(100)

	//for _, node := range tree.Childrens {
	//	num := GetNumByNode(node)
	//	fmt.Printf("%s: %d\n", node.Name, num)
	//}

	for _, redisStat := range redisStats {
		redisStat.Show()
	}
}