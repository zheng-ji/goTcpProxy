// zheng-ji.info

package main

import (
	"net"
	"time"
)

func initProxy() {

	pLog.Infof("Proxying %s -> %s\n", pConfig.Bind, pConfig.Backend)

	server, err := net.Listen("tcp", pConfig.Bind)
	if err != nil {
		pLog.Fatal(err)
	}

	waitQueue := make(chan net.Conn, pConfig.WaitQueueLen)
	availPools := make(chan bool, pConfig.MaxConn)
	for i := 0; i < pConfig.MaxConn; i++ {
		availPools <- true
	}

	go loop(waitQueue, availPools)

	for {
		connection, err := server.Accept()
		if err != nil {
			pLog.Error(err)
		} else {
			pLog.Infof("Received connection from %s.\n", connection.RemoteAddr())
			waitQueue <- connection
		}
	}
}

func loop(waitQueue chan net.Conn, availPools chan bool) {
	for connection := range waitQueue {
		<-availPools
		go func(connection net.Conn) {
			handleConnection(connection)
			availPools <- true
			pLog.Infof("Closed connection from %s.\n", connection.RemoteAddr())
		}(connection)
	}
}

func handleConnection(connection net.Conn) {
	defer connection.Close()

	bksvr, ok := getBackendSvr(connection)
	if !ok {
		return
	}
	remote, err := net.Dial("tcp", bksvr.svrStr)

	if err != nil {
		pLog.Error(err)
		bksvr.failTimes += 1
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
	remote.Close()
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

			from.SetReadDeadline(time.Now().Add(time.Duration(pConfig.Timeout) * time.Second))
			read, err = from.Read(bytes)
			if err != nil {
				complete <- true
				one_side <- true
				return
			}

			to.SetReadDeadline(time.Now().Add(time.Duration(pConfig.Timeout) * time.Second))
			_, err = to.Write(bytes[:read])
			if err != nil {
				complete <- true
				one_side <- true
				return
			}
		}
	}
}
