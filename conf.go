package main

import (
	"fmt"
	"github.com/revel/config"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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

type section struct {
	name       string
	configPath string
	config     *config.Config
	mKeyValue  map[string]string
}

func (this section) SetValue(key, value string) {
	this.mKeyValue[key] = value
}

func (this section) PersistenceConfig() error {
	for key, value := range this.mKeyValue {
		if _, err := this.config.String(this.name, key); err != nil {
			this.config.AddOption(this.name, key, value)
		}
	}
	return this.config.WriteFile(this.configPath, 0644, "auto add values")
}

func NewSection(name string, configPath string, c *config.Config) *section {
	return &section{
		name:       name,
		config:     c,
		configPath: configPath,
		mKeyValue:  make(map[string]string),
	}
}

func GetCurrPath() (path string, err error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err = filepath.Abs(file)
	return path, err
}

func main() {
	fmt.Printf("==============================================================\n")

	var err error
	hostname, _ := os.Hostname()
	_ = hostname
	cnfPath := "./conf.conf"
	path, _ := GetCurrPath()
	println(path)

	c, err := config.ReadDefault(cnfPath)
	CheckErr(err, "read config")

	s := NewSection("mail", cnfPath, c)
	s.SetValue("username", "gold")
	s.SetValue("password", "")
	s.SetValue("domain", "monitor.spriteapp.com")
	s.SetValue("port", "25")
	s.SetValue("to", "yunwei@spriteapp.com;")

	err = s.PersistenceConfig()
	CheckErr(err, "save config")

	c, err = config.ReadDefault(cnfPath)

	userName, err := c.String("mail", "username")
	CheckErr(err, "get value mail->username")
	println(userName)

	fmt.Printf("==============================================================\n")
}
