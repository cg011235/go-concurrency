package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

type stats struct {
	mu			sync.Mutex		// mutex
	counter		uint64			// number of requests so far
	dataFlow	int64			// bytes per seconds
	respTime	time.Duration	// average response time
}

type server struct {
	hostAddr	string
	c			net.Conn
	connected	bool
	stat		stats
}

func initServers(netw string, serverList string) []server {
	nodes := strings.Split(serverList, ",")
	n := len(nodes)
	hostlist := make([]server, n)
	for i, node := range nodes {
		log.Println("Dialing host:", node)
		hostlist[i].hostAddr = strings.TrimSpace(node)

		c, err := net.Dial(netw, hostlist[i].hostAddr)
		if err != nil {
			hostlist[i].connected = false
			log.Println(err)
		} else {
			hostlist[i].c = c
			hostlist[i].connected = true
		}
		hostlist[i].stat.mu.Lock()
		hostlist[i].stat.counter = 0
		hostlist[i].stat.dataFlow = 0
		hostlist[i].stat.respTime = time.Duration(0)
		hostlist[i].stat.mu.Unlock()
	}
	return hostlist
}

type backend struct {
	serverList	[]server
	total		int
	nextServer	int
	mu 			sync.RWMutex
}

var b backend

func getNextServer() server {
	b.mu.Lock()
	defer b.mu.Unlock()
	for {
		if b.serverList[b.nextServer].connected {
			ret := b.serverList[b.nextServer]
			b.nextServer++
			if b.nextServer >= b.total {
				b.nextServer = 0
			}
			return  ret
		} else {
			b.nextServer++
			if b.nextServer >= b.total {
				b.nextServer = 0
			}
		}
	}
}

func handle(c net.Conn) {
	defer c.Close()
	log.Println("Accepted connection from", c.RemoteAddr())
	s := getNextServer()
	input := bufio.NewScanner(c)
	output := bufio.NewScanner(s.c)
	for input.Scan() {
		log.Println("Recieved:", input.Text())
		s.c.Write(input.Bytes())
		for output.Scan() {
			c.Write(output.Bytes())
		}
	}
	log.Println("Closing connection to", c.RemoteAddr())
}

func main() {
	port := flag.String("port", "12345", "port number for load balancer")
	netw := flag.String("net", "tcp", "transport layer protocol")
	host := flag.String("host", "0.0.0.0", "host to listen at")
	servers := flag.String("servers", "localhost:9000,localhost:9001", "comma separated list of servers")

	flag.Parse()

	listner, err := net.Listen(*netw, *host + ":" + *port)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Started load balancer server at port", *port)

	b.mu.Lock()
	b.serverList	= initServers(*netw, *servers)
	b.nextServer	= 0
	b.total			= len(b.serverList)
	b.mu.Unlock()

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

	<-stop
	
	log.Println("Stopping the server now")
}
