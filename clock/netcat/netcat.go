package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"fmt"
)

func handle(c net.Conn) {
	_, err := io.Copy(os.Stdout, c)
	if err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Println("netcat [-host][-port][-net]")
	os.Exit(2)
}

func main() {
	host := flag.String("host", "localhost", "host to connect with")
	port := flag.String("port", "9000", "port number for server")
	netw := flag.String("net", "tcp", "network type (tcp | udp)")
	flag.Usage = usage
	flag.Parse()

	conn, err := net.Dial(*netw, *host + ":" + *port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	handle(conn)
}
