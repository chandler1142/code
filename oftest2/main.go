package main

import (
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func getConnection()(net.Conn) {
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
	return conn
}

func handleWrite(conn net.Conn, content string) {
	for i := 10; i > 0; i-- {
		_, e := conn.Write([]byte(content))
		if e != nil {
			fmt.Println("Error to send message because of ", e.Error())
			break
		}
	}
	//done <- "Sent"
}

type Stream struct {
	StreamId xml.Name `xml:"id"`
}

func handleRead(conn net.Conn) {

	buf := make([]byte, 4096)
	streamId := ""
	for {
		reqLen, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error to read message because of ", err)
			return
		}
		receiveStr := string(buf[:reqLen])
		println("receive: " + receiveStr)
		if strings.Contains(receiveStr, "<stream:stream") {
			start := strings.Index(receiveStr, "id=\"")
			end := strings.Index(receiveStr, "\" xml:lang=")
			streamId = receiveStr[start+5:end]
		}
		if strings.Contains(receiveStr, "starttls") {
			//may not essential
			handleWrite(conn, "<starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'/> ")
			//handleWrite(conn, "<stream:stream />")
			//handleWrite(conn, "<?xml version='1.0'?><stream:stream from='test@localhost' to='test2@localhost' version='1.0' xml:lang='en' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams'>")
		}
		if strings.Contains(receiveStr, "proceed") {
			newStream := "<stream:stream from='test@localhost' to='localhost' version='1.0' xml:lang='en' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams'>"
			handleWrite(conn, newStream)
		}
		if strings.Contains(receiveStr, "mechanisms") {
			//没有TLS，直接回复AUTH
			newStream := `
				  <auth mechanism="PLAIN" xmlns="urn:ietf:params:xml:ns:xmpp-sasl"/>
			`
			handleWrite(conn, newStream)
		}
		//bind resources
		if strings.Contains(receiveStr, "success") {
			newStream := `<stream:stream from='test@localhost' to='localhost' version='1.0'
     xml:lang='en'
     xmlns='jabber:client'
     xmlns:stream='http://etherx.jabber.org/streams'>`
			handleWrite(conn, newStream)
			newStream = "<stream:stream xmlns='jabber:client' to='localhost' xmlns:stream='http://etherx.jabber.org/streams' version='1.0' from='test@localhost' id='"+streamId+"' xml:lang='en'>"
			handleWrite(conn, newStream)
		}
		if strings.Contains(receiveStr, "<bind") {

			bindIQ := "<iq id='5q15b-3' type='set'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'><resource>Android</resource></bind></iq>"
			handleWrite(conn, bindIQ)
		}
		if strings.Contains(receiveStr,"jid") {
			//iqId := ""
			//start := strings.Index(receiveStr, "id=\"")
			//end := strings.Index(receiveStr, "\" to=")
			//iqId = receiveStr[start+4:end]
			//print(iqId)
			//messageStr := "<message from='test@localhost/Android' id='"+iqId+"' to='test2@localhost' type='chat' xml:lang='en'><body>Art thou not Romeo, and a Montague?</body></message>"
			//handleWrite(conn, messageStr)
		}
	}

	//done <- "Read"
}

func main()  {
	conn := getConnection()
	defer conn.Close()

	println("connection success")

	streamStart := "<?xml version='1.0'?><stream:stream from='test@localhost' to='localhost' version='1.0' xml:lang='en' xmlns='jabber:client' xmlns:stream='http://etherx.jabber.org/streams'>"
	//streamEnd := "</stream:stream>"
	//messageBody := `<message to='test2@localhost'>
	//<body>Wherefore art thou?</body>
	//</message>`

	go handleRead(conn)
	handleWrite(conn, streamStart)
	//handleWrite(conn, messageBody)

	//handleWrite(conn, streamEnd)

	time.Sleep(30*time.Second)

	//println(<-done)
	//println(<-done)


}
