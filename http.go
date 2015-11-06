package main

import (
	"fmt"
	"net/http"
	"strings"
)

type Hello struct{}

func (h Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	to := r.Form.Get("to")
	form := r.Form.Get("form")
	msg := r.Form.Get("msg")
	if to == "" || form == "" || msg == "" {
		fmt.Fprint(w, "error")
		return
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]
	fmt.Fprint(w, "Client ip:"+ip+"\n")

	host := "dnsP(ip)"

	fmt.Fprint(w, "to:"+to+"\n")
	fmt.Fprint(w, "form:"+form+"\n")
	fmt.Fprint(w, "msg:"+msg+"\n")
	fmt.Fprint(w, "host:"+host+"\n")
}

func main() {
	var h Hello
	http.ListenAndServe("0.0.0.0:4000", h)
}

