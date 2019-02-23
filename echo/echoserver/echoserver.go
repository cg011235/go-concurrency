package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"
)

func echoLower(c net.Conn, str string) {
	t := time.Duration(rand.Uint32() % 10)
	time.Sleep(t * time.Second)
	fmt.Fprintln(c, strings.ToLower(str))
}

func echoUpper(c net.Conn, str string) {
	t := time.Duration(rand.Uint32() % 10)
	time.Sleep(t * time.Second)
	fmt.Fprintln(c, strings.ToUpper(str))
}

func handle(c net.Conn) {
	defer c.Close()
	log.Println("Accepted connection from", c.RemoteAddr())
	input := bufio.NewScanner(c)
	for input.Scan() {
		log.Println("Received:", input.Text(), "from:", c.RemoteAddr())
		go echoLower(c, input.Text())
		go echoUpper(c, input.Text())
	}
}

func main() {
	port := flag.String("port", "9000", "port number for echo server")
	netw := flag.String("net", "tcp", "network type (tcp | udp)")
	flag.Parse()

	log.Println("Starting echo server at port", *port, "with", *netw, "network")
	listner, err := net.Listen(*netw, "localhost:" + *port)
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	go func() {
		for {
			conn, err := listner.Accept()
			if err != nil {
				log.Fatal(err)
			}
			go handle(conn)
		}
	}()

	<-stop  // Block till we receive any value

	log.Println("Stopping the echo server")
}
