package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	g "xcthings.com/ftconn/checknat-ms/common/global"
	"xcthings.com/hjyz/logs"
	xcpprof "xcthings.com/hjyz/pprof"
)

func pprofInit() {
	var cb = func(info xcpprof.ServiceInfo) {
		logs.Logger.Infof("HeapAlloc: %d KB, Gos: %d, StartLen: %d, tcp count: %d, udp count: %d.",
			info.HeapAlloc, info.Gos, info.StartLen, tcpCount(), udpCount())
	}
	xcpprof.ReportService(int(g.Conf.Public.ReportInterval), cb)
	if g.Conf.Public.AdminProf == false {
		return
	}
	addr := fmt.Sprintf("0.0.0.0:%d", g.Conf.Public.AdminPort)
	go xcpprof.StartPprof(addr)
}

func tcpCount() int32 {
	var n int32
	for _, v := range tcplis {
		n = n + v.Count()
	}
	return n
}

func udpCount() int32 {
	var n int32
	for _, v := range udplis {
		n = n + v.Count()
	}
	return n
}

func infoCB(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	info := xcpprof.ShowSysInfo()
	info.TCPCount = tcpCount()
	info.UDPCount = udpCount()
	b, _ := json.Marshal(&info)
	fmt.Fprintf(w, string(b))
}

func errCB(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	info := xcpprof.ShowSysInfo()
	info.TCPCount = tcpCount()
	info.UDPCount = udpCount()
	b, _ := json.Marshal(&info)
	fmt.Fprintf(w, string(b))
}
