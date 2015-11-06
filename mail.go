package main

import (
	"crypto/tls"
	"fmt"
	"log"
	//"net"
	"errors"
	"net/smtp"
	"os"
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

type Mailer struct {
	host     string
	port     int
	addr     string
	username string
	password string
	to       []string
	header   map[string]string
}

func NewMailer(host string, port int, username, password string, to []string) *Mailer {
	addr := fmt.Sprintf("%s:%d", host, port)
	header := make(map[string]string)
	header["From"] = username + "<" + username + ">"
	header["Content-Type"] = "text/html; charset=UTF-8"
	//header["Content-Type"] = "text/plain; charset=UTF-8"
	return &Mailer{
		host:     host,
		port:     port,
		addr:     addr,
		username: username,
		password: password,
		to:       to,
		header:   header,
	}
}

func (this Mailer) SetContentType(t string) error {
	if t == "html" {
		this.header["Content-Type"] = "text/html; charset=UTF-8"
	} else if t == "plain" || t == "txt" {
		this.header["Content-Type"] = "text/plain; charset=UTF-8"
	} else {
		return errors.New("mail SetContentType type not find")
	}
	return nil
}

func (this Mailer) Send(subject, body string) (err error) {
	if this.port == 25 {
		return this.SendMail(subject, body)
	} else {
		if this.SendMailUsingTLS(subject, body) != nil {
			return this.SendMail(subject, body)
		}
	}
	return nil
}

func (this Mailer) SendMail(subject, body string) (err error) {
	c, err := smtp.Dial(this.addr)
	if err != nil {
		return err
	}
	defer c.Close()

	//create auth
	auth := smtp.PlainAuth("", this.username, this.password, this.host)
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}

	message := ""
	//set subject
	this.header["Subject"] = subject
	for k, v := range this.header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	for _, u := range this.to {
		message += fmt.Sprintf("To: %s\r\n", u)
	}

	message += "\r\n\r\n" + body
	msg := []byte(message)
	//set from
	if err = c.Mail(this.username); err != nil {
		return err
	}
	//set addr
	for _, addr := range this.to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	//send mail
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()

}

func (this Mailer) SendMailUsingTLS(subject, body string) (err error) {
	//create smtp client
	conn, err := tls.Dial("tcp", this.addr, nil)
	if err != nil {
		return err
	}
	c, err := smtp.NewClient(conn, this.host)
	if err != nil {
		return err
	}
	defer c.Close()

	//create auth
	auth := smtp.PlainAuth("", this.username, this.password, this.host)
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}

	message := ""
	//set subject
	this.header["Subject"] = subject
	for k, v := range this.header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	for _, u := range this.to {
		message += fmt.Sprintf("To: %s\r\n", u)
	}

	message += "\r\n\r\n" + body
	msg := []byte(message)
	//set from
	if err = c.Mail(this.username); err != nil {
		return err
	}
	//set addr
	for _, addr := range this.to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	//send mail
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}

func main() {
	users := []string{"user@domain.com", "user@domain.com", "user@domain.com"}
	//mailer := NewMailer("smtp.exmail.qq.com", 465, "user@domain.com", "test1234", users)
	hostName, _ := os.Hostname()
	domain := "@monitor.domain.com"
	user := hostName + domain

	mailer := NewMailer("gold", 25, user, "", users)
	mailer.SetContentType("plain")
	msg := ` ~,~ `
	//err := mailer.SendMailUsingTLS("我就是标题", msg)
	//err := mailer.SendMail("我就是标题", msg)
	err := mailer.Send("我就是标题", msg)
	CheckErr(err, "send mail")
	//os.Exit(0)
}
