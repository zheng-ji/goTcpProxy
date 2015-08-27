// zheng-ji.info

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	src          = flag.String("src", "127.0.0.1:8999", "proxy server's host.")
	dest         = flag.String("dest", "172.16.1.250:3306", "where proxy server forward requests to.")
	maxConn      = flag.Int("max_avail_conn", 25, "max active connection.")
	maxWaitQueue = flag.Int("max_wait_conn", 10000, "max connections in the queue wait for servers.")
	expire       = flag.Int("ttl", 20, "timeout of read and write")
)

func main() {

	flag.Parse()
	fmt.Printf("Proxying %s->%s.\n", *src, *dest)

	server, err := net.Listen("tcp", *src)
	if err != nil {
		log.Fatal(err)
	}

	waitQueue := make(chan net.Conn, *maxWaitQueue)
	availPools := make(chan bool, *maxConn)
	for i := 0; i < *maxConn; i++ {
		availPools <- true
	}

	go loop(waitQueue, availPools)
	go onExitSignal()

	for {
		connection, err := server.Accept()
		if err != nil {
			log.Print(err)
		} else {
			log.Printf("Received connection from %s.\n", connection.RemoteAddr())
			waitQueue <- connection
		}
	}
}

func onExitSignal() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGUSR1, syscall.SIGTERM, syscall.SIGINT)
L:
	for {
		sig := <-sigChan
		switch sig {
		case syscall.SIGUSR1:
			log.Fatal("Reopen log file")
		case syscall.SIGTERM, syscall.SIGINT:
			log.Fatal("Catch SIGTERM singal, exit.")
			break L
		}
	}
}

func loop(waitQueue chan net.Conn, availPools chan bool) {
	for connection := range waitQueue {
		<-availPools
		go func(connection net.Conn) {
			handleConnection(connection)
			availPools <- true
			log.Printf("Closed connection from %s.\n", connection.RemoteAddr())
		}(connection)
	}
}

func handleConnection(connection net.Conn) {
	defer connection.Close()

	remote, err := net.Dial("tcp", *dest)
	defer remote.Close()

	if err != nil {
		log.Print(err)
		return
	}

	//等待双向连接完成
	complete := make(chan bool, 2)
	one_side := make(chan bool, 1)
	other_side := make(chan bool, 1)
	go pass(connection, remote, complete, one_side, other_side)
	go pass(remote, connection, complete, other_side, one_side)
	<-complete
	<-complete
}

// copy Content two-way
func pass(from net.Conn, to net.Conn, complete chan bool, one_side chan bool, other_side chan bool) {
	var err error = nil
	var bytes []byte = make([]byte, 256)
	var read int = 0

	for {
		select {

		case <-other_side:
			complete <- true
			return

		default:

			from.SetReadDeadline(time.Now().Add(time.Duration(*expire) * time.Second))
			read, err = from.Read(bytes)
			if err != nil {
				complete <- true
				one_side <- true
				return
			}

			to.SetWriteDeadline(time.Now().Add(time.Duration(*expire) * time.Second))
			_, err = to.Write(bytes[:read])
			if err != nil {
				complete <- true
				one_side <- true
				return
			}
		}
	}
}
