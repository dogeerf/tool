package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"runtime"
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
	log.Printf("### OK: |%s| %s success.\n", funcName, operating)
}

//func SendMail(host, user, password, from, to, subject, body, mailtype string) error {
//	hp := strings.Split(host, ":")
//	auth := smtp.PlainAuth("", user, password, hp[0])
//	var content_type string
//	if mailtype == "html" {
//		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
//	} else {
//		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
//	}
//	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
//	send_to := strings.Split(to, ";")
//	err := smtp.SendMail(host, auth, user, from, msg)
//	return err
//}

func sendmail_1() {
	host := "192.168.133.187:25"
	user := "gold"
	from := "web5@monitor.domain.com"
	to := "user@domain.com"
	body := "This is the email body"
	content_type := "html"
	subject := "subject"

	msg := "To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body

	c, err := smtp.Dial(host)
	CheckErr(err, "conn smtp host.")

	err = c.Mail(from)
	CheckErr(err, "set form user:")

	err = c.Rcpt(to)
	CheckErr(err, "set to user:")

	wc, err := c.Data()
	CheckErr(err, "open body io")

	_, err = fmt.Fprintf(wc, msg)
	CheckErr(err, "copy body")

	err = wc.Close()
	CheckErr(err, "close wc")

	err = c.Quit()
	CheckErr(err, "quit conn")

}

func SendMail(host, user, password, from, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	println("1")
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}
	msg := []byte("To: " + to + "\r\nFrom: " + from + "<" + from + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	println("2")
	err := smtp.SendMail(host, auth, from, send_to, msg)
	println("3")
	return err
}

func main() {
	fmt.Printf("==============================================================\n")
	//sendmail_1()
	//host := "192.168.133.187:25"
	//user := ""
	//password := ""

	host := "smtp.exmail.qq.com:465"
	user := "user@domain.com"
	password := "password"

	from := "user@domain.com"
	to := "user@domain.com;user@domain.com;"

	subject := "this is subject!"
	body := "This is body!"
	mailtype := "html"

	err := SendMail(host, user, password, from, to, subject, body, mailtype)
	CheckErr(err, "send mail")
	fmt.Printf("==============================================================\n")
}
