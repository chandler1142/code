package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stdout, "you got something wrong %v", err)
			continue
		}
		go handleConn(conn)
	}
}

type client chan<- string  // send only channel

var (
	entering = make(chan client)
	leaving = make(chan client)
	messages = make(chan string)
)

func broadcaster() {
	clients := make(map[client]bool) //all connected clients
	for {
		select {
		case msg := <- messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients{
				cli <- msg
			}
		case cli := <- entering:
			clients[cli] = true
		case cli := <- leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}


func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has joined us"
	// 这里不是把数据吐出去。。而是吐出去了本身一个channel这里比较难理解。
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}

	leaving <- ch
	messages <- who + " has left "
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}