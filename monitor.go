// Author: zheng-ji.info

package main

import (
	"fmt"
	"net/http"
)

// 查询监控信息的接口
func statsHandler(w http.ResponseWriter, r *http.Request) {
	_str := ""
	for _, v := range pBackendSvrs {
		_str += fmt.Sprintf("Server:%s FailTimes:%d isUp:%t\n", v.svrStr, v.failTimes, v.isUp)
	}
	fmt.Fprintf(w, "%s", _str)
}

func initStats() {
	pLog.Infof("Start monitor on addr %s", pConfig.Stats)

	go func() {
		http.HandleFunc("/stats", statsHandler)
		http.ListenAndServe(pConfig.Stats, nil)
	}()
}
