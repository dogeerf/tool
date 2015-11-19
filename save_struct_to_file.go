package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"runtime"
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

type User struct {
	Id         int
	Name       string
	CreateTime time.Time
	IsMan      bool
}

func (this *User) Show() {
	fmt.Printf("-----------------------------------\n")
	fmt.Printf("Id:		    %v\n", this.Id)
	fmt.Printf("Name:	    %v\n", this.Name)
	fmt.Printf("CreateTime: %v\n", this.CreateTime)
	fmt.Printf("-----------------------------------\n")
}

func main() {
	fmt.Printf("==============================================================\n")
	var err error
	filePath := "./data"
	structNum := 10 * 10000

	file, err := os.Create(filePath)
	defer file.Close()

	var users []User

	for i := 1; i <= structNum; i++ {
		user := User{
			Id:         i,
			Name:       fmt.Sprintf("user_%d", i),
			CreateTime: time.Now(),
			IsMan:      true,
		}
		users = append(users, user)
	}

	// Create an encoder and save data
	enc := gob.NewEncoder(file)
	err = enc.Encode(users)
	if err != nil {
		log.Fatal("encode:", err)
	}

	file.Seek(0, 0)
	dec := gob.NewDecoder(file)
	err = dec.Decode(&users)
	if err != nil {
		log.Fatal("Decode:", err)
	}
	for i, user := range users {
		if i%10000 == 0 {
			user.Show()
			//fmt.Printf("ID:%d User:%s\n", i, user.Name)
		}

	}
	fmt.Printf("==============================================================\n")
}
