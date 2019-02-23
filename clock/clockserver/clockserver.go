package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func handle(c net.Conn) {
	log.Println("Accepted connection from", c.RemoteAddr().String())
	defer c.Close()
	// Loop and send current time with the interval of 1 sec
	for {
		_, err := io.WriteString(c, time.Now().String() + "\n")
		if err != nil {
			log.Println(err)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func usage() {
	fmt.Println("clockserver [-port][-net]")
	os.Exit(2)
}

func main() {
	port := flag.String("port", "9000", "port number for server")
	netw := flag.String("net", "tcp", "network type (tcp | udp)")
	flag.Usage = usage
	flag.Parse()

	log.Println("Starting clock server at port", *port, "with", *netw, "network")

	listner, err := net.Listen(*netw, "localhost:" + *port)
	if err != nil {
		log.Fatal(err)
	}

	stop := make(chan os.Signal)
	// On sigterm or pressing ctrl-c the server will recieve singal from OS
	// This will be fed to the stop channel we just created.
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {  // anonymous goroutine
		for {
			conn, err := listner.Accept()
			if err != nil {
				log.Fatal(err)
			}
			go handle(conn)
		}
	}()

	// On receiving value on this channel, it will unblock and main will
	// exit causing other goroutines to stop as well.
	<-stop  // Block for signal

	log.Println("Recived signal so stopping server")
}
