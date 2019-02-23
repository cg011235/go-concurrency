package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
)

func main()  {
	host := flag.String("host", "localhost", "host to connect to")
	port := flag.String("port", "9000", "port number to connect to")
	netw := flag.String("net", "tcp", "network layer (tcp | udp)")
	flag.Parse()

	conn, err := net.Dial(*netw, *host + ":" + *port)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn)
		done <- struct{}{}
	}()
	io.Copy(conn, os.Stdin)
	<-done
}
