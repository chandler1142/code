package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

//http://www.hmxtstudy.com:8080/www/hmxtwz/fcy/dz2.jsp?id=121
func main() {

	OKAddr := "6.1.1.6" // local IP address to use

	OKAddress, _ := net.ResolveTCPAddr("tcp", OKAddr)

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			LocalAddr: OKAddress}).Dial, TLSHandshakeTimeout: 10 * time.Second}

	client := &http.Client{
		Transport: transport,
	}

	resp, err := client.Get("http://www.hmxtstudy.com:8080/www/hmxtwz/fcy/dz2.jsp?id=121")

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(html))

	// fmt.Println(os.Stdout, string(html))
	// or save the zip - see https://www.socketloop.com/tutorials/golang-download-file-example
}
