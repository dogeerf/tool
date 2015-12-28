package main

import (
	"fmt"
	//"github.com/mmitton/ldap"
	"github.com/mqu/openldap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
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

var (
	url     string
	base_dn string
	ldap    *openldap.Ldap
	err     error
)

func init() {
	initLdap()
}

func initLdap() {
	url = "ldap://ldap.domain.com:1389/"
	base_dn = "ou=People,dc=dc,dc=com"

	ldap, err = openldap.Initialize(url)
	if err != nil {
		fmt.Printf("LDAP::Initialize() : connexion error\n")
		return
	}
	ldap.SetOption(openldap.LDAP_OPT_PROTOCOL_VERSION, openldap.LDAP_VERSION3)
}

func checkUser(username string, password string) (err error) {
	who := fmt.Sprintf("uid=%s,%s", username, base_dn)
	return ldap.Bind(who, password)
}

func checkUserHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		w.Write([]byte("username or password is null"))
	}
	err := checkUser(username, password)
	if err != nil {
		w.Write([]byte(err.Error()))
		initLdap()
	}
	w.Write([]byte("success"))
}

func getUser(username string, password string) (m map[string]string, err error) {
	scope := openldap.LDAP_SCOPE_SUBTREE
	filter := fmt.Sprintf("uid=%s", username)
	attributes := []string{}

	result, err := ldap.SearchAll(base_dn, scope, filter, attributes)
	CheckErr(err, "searall")
	fmt.Printf("%#v\n", result)

	println("--------------------------------")
	fmt.Printf("Filter: %s\n", result.Entries)
	if result.Count() != 1 {
		fmt.Printf("ERROR: count ne 1\n")
	}
	entrie := result.Entries()[0]
	_ = entrie
	for _, attr := range entrie.Attributes() {
		fmt.Printf("%-10s %-20s", attr.Name(), attr.String())
	}
	println("--------------------------------")
	fmt.Println(result)
	println("--------------------------------")

	return
}

func signalHandle() {
	c := make(chan os.Signal)
	signal.Notify(c)
	signal.Notify(c, syscall.SIGHUP)
	s := <-c
	fmt.Println("get signal:", s)
}

func main() {
	defer ldap.Close()
	http.HandleFunc("/checkUser", checkUserHandler)
	getUser("user", "passwd")
	//http.ListenAndServe(":8888", nil)
}
