package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

//发送信息
func Send(conn net.Conn, body string) {
	conn.Write([]byte(body))
	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)
	if err != nil {
		Log(conn.RemoteAddr().String(), "waiting server back msg error: ", err)
		return
	}
	Log(conn.RemoteAddr().String(), "receive server back msg: \n", string(buffer[:n]))
}


//日志
func Log(v ...interface{}) {
	log.Println(v...)
}

func main()  {
	server := "laptop-d5d42j5u:5222"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("connection success")

	streamStart := "<?xml version=’1.0’?><stream:stream from='thzw1142@163.com' to='localhost' version='1.0' xml:lang='en' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams'>"
	Send(conn, streamStart)

	//有问题
	//auth := "<auth mechanism='DIGEST-MD5' xmlns='urn:ietf:params:xml:ns:xmpp-sasl'></auth>"
	//Send(conn, auth)

	//mandatory的有点问题
	//requireFeatures := `
	//	<stream:features>
 	//		<starttls xmlns=’urn:ietf:params:xml:ns:xmpp-tls’>
 	//			<required/>
 	//		</starttls>
	//	</stream:features>
	//`
	//Send(conn, requireFeatures)


	bindResources := `
	<iq id="wSBRk-4" type="set">  
		<bind xmlns="urn:ietf:params:xml:ns:xmpp-bind">  
			<resource>Smack</resource>  
			<terminal>android</terminal>  
		</bind>  
	</iq> 
	`
	Send(conn, bindResources)

	//voluntary-to-negotiate
	features2 := `
		<stream:features>
 			<compression xmlns='http://jabber.org/features/compress'>
 				<method>zlib</method>
 				<method>lzw</method>
 			</compression>
		</stream:features>
	`
	Send(conn, features2)

	//empty features
	//features3 := `
	//	<stream:features/>
	//`
	//Send(conn, features3)

	streamEnd := "</stream:stream>"
	Send(conn, streamEnd)

	defer conn.Close()
}

func handleRead(conn net.Conn, done chan string) {
	buf := make([]byte, 1024)
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error to read message because of ", err)
		return
	}
	fmt.Println(string(buf[:reqLen-1]))
	//done <- "Read"
}
