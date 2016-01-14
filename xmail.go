package main

import (
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"math/rand"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
	"xiaoju/sms"
)

func TransformPhones(phones string) (res []int64, errno int) {
	p1 := regexp.MustCompile(`[,\s]+`)
	p2 := regexp.MustCompile(`^[0-9]{11}$`) // 校验手机号的位数是否为11位
	p_str_slice := p1.Split(phones, -1)

	p_slice := make([]int64, 0, len(p_str_slice))
	for _, i := range p_str_slice {
		if p2.MatchString(i) {
			j, _ := strconv.ParseInt(i, 10, 64) // string2int64
			p_slice = append(p_slice, j)
		} else {
			fmt.Printf("Warning phones: %s\n", i)
		}
	}

	if len(p_slice) < 1 {
		fmt.Printf("Error input:%s, err: no phone numbers\n", phones)
		return p_slice, 1
	} else {
		return p_slice, 0
	}
}

func GenRandn(k int) (ret int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	ret = r.Intn(k)
	return ret
}

func SendSms(msgSend string, phones_slice []int64) {
	// random choice host
	hosts := []string{"10.231.144.36", "10.231.144.37"}
	rk := GenRandn(len(hosts))
	host := hosts[rk]
	port := "9090"

	transport, err := thrift.NewTSocket(net.JoinHostPort(host, port))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		os.Exit(1)
	}
	transportFactory := thrift.NewTBufferedTransportFactory(10240)
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	useTransport := transportFactory.GetTransport(transport)
	client := sms.NewMessageServiceClientFactory(useTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to ", host, ":", port, " ", err)
		os.Exit(1)
	}
	defer transport.Close()

	msg := sms.NewMessage()
	msg.Phones = phones_slice
	msg.BusinessId = 200100000
	msg.Message = msgSend
	client.SendMessage(msg)
}

func main() {
	var msgSend, phones string
	flag.StringVar(&msgSend, "m", "hello, go!!!", "The message want to be send, eg: hello, go")
	flag.StringVar(&phones, "p", "15210239684", "Who receives the message (at least one, separated by:comma or whitespace)")
	flag.Parse()

	phones_slice, err_p := TransformPhones(phones)
	if err_p != 0 {
		os.Exit(1)
	}
	_ = phones_slice
	SendSms(msgSend, phones_slice)
	fmt.Printf("send sms(%s) to %s\n", msgSend, phones)
}
