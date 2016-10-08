// Author: zheng-ji.info

package main

import (
	//"fmt"
	"math/rand"
	"net"
	"stathat.com/c/consistent"
	"time"
)

// BackendSvr Type
type BackendSvr struct {
	svrStr    string
	isUp      bool // is Up or Down
	failTimes int
}

var (
	pConsisthash *consistent.Consistent
	pBackendSvrs map[string]*BackendSvr
)

func initBackendSvrs(svrs []string) {
	pConsisthash = consistent.New()
	pBackendSvrs = make(map[string]*BackendSvr)

	for _, svr := range svrs {
		pConsisthash.Add(svr)
		pBackendSvrs[svr] = &BackendSvr{
			svrStr:    svr,
			isUp:      true,
			failTimes: 0,
		}
	}
	go checkBackendSvrs()
}

func getBackendSvr(conn net.Conn) (*BackendSvr, bool) {
	remoteAddr := conn.RemoteAddr().String()
	svr, _ := pConsisthash.Get(remoteAddr)

	bksvr, ok := pBackendSvrs[svr]
	return bksvr, ok
}

func checkBackendSvrs() {
	// scheduler every 10 seconds
	rand.Seed(time.Now().UnixNano())
	t := time.Tick(time.Duration(10)*time.Second + time.Duration(rand.Intn(100))*time.Millisecond*100)

	for _ = range t {
		for _, v := range pBackendSvrs {
			if v.failTimes >= pConfig.FailOver && v.isUp == true {
				v.isUp = false
				pConsisthash.Remove(v.svrStr)
			}
		}

	}
}
