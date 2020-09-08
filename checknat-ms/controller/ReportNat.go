package controller

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
	"xcthings.com/protoc/ftconnnat/NatProbe"
	"xcthings.com/protoc/ftconnnat/ReportNat"

	lc "xcthings.com/ftconn/checknat-ms/common"
	g "xcthings.com/ftconn/checknat-ms/common/global"
	l "xcthings.com/ftconn/checknat-ms/logic"
)

// ReportNater .
type ReportNater struct{}

// ReqHandle ReportNat request handle
func (t *ReportNater) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *ReportNat.Req) (err error) {
	ip := common.GetIPAddr(c.RemoteAddr())
	upnp := detectUPNP(req, ip)
	if upnp {
		req.Nat.Upnp = lc.UPNPENABLE
	} else {
		req.Nat.Upnp = lc.UPNPDISABLE
	}

	var code uint64
	var resp *ReportNat.Resp
	resp, code, err = l.LReportNat(c, pkg, req)
	if code != 0 {
		pkg.Code = code
		logs.Logger.Warnf("%s, l.LReportNat, code: %d, err: %s.", c, code, err)
	}

	_, err = pprpc.WriteResp(c, pkg, resp)
	if err != nil {
		logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
	}
	//if err != nil || pkg.Code != 0 {
	// 	return
	//}
	return
}

// RespHandle ReportNat response handle
func (t *ReportNater) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *ReportNat.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	return
}

// DestructHandle ReportNat.
func (t *ReportNater) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, ReportNat, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}

func detectUPNP(req *ReportNat.Req, ipAddr string) bool {
	nat := req.Nat
	if nat.Upnp == lc.UPNPDISABLE {
		return false
	}
	if nat.UpnpProtocol != lc.UPNPTCP && nat.UpnpProtocol != lc.UPNPUDP {
		return false
	}
	if nat.UpnpPort < 1 {
		return false
	}
	if nat.UpnpProtocol == lc.UPNPTCP {
		return detectTCP(ipAddr, nat.UpnpPort)
	} else {
		return detectUDP(ipAddr, nat.UpnpPort)
	}

}

func detectTCP(ip string, port int32) bool {
	u, e := url.ParseRequestURI(fmt.Sprintf("tcp://%s:%d", ip, port))
	if e != nil {
		logs.Logger.Errorf("detectTCP, url.ParseRequestURI(\"tcp://%s:%d\"), error: %s.", ip, port, e)
		return false
	}
	conn, err := pprpc.Dail(u, nil, g.Service, 2*time.Second, nil)
	if err != nil {
		logs.Logger.Errorf("pprpc.Dail(), error: %s.", err)
		conn.Close()
		return false
	}
	defer conn.Close()

	req := new(NatProbe.Req)
	req.Check = "upnp"
	rpc := NatProbe.NewNatProbeClient(conn)
	_, _, err = rpc.PPCall(context.Background(), req)
	if err != nil {
		logs.Logger.Errorf("NatProbe.PPCall(), error: %s.", err)
		return false
	}

	return true
}

func detectUDP(ip string, port int32) bool {
	conn, err := pprpc.DailUDP(fmt.Sprintf("%s:%d", ip, port), g.Service, 30, nil)
	if err != nil {
		logs.Logger.Errorf("pprpc.DailUDP(%s:%d), error: %s.", ip, port, err)
		return false
	}
	defer conn.Close()
	conn.SyncWriteTimeoutMs = 2000

	req := new(NatProbe.Req)
	req.Check = "upnp"
	rpc := NatProbe.NewNatProbeClient(conn)
	_, _, err = rpc.PPCall(context.Background(), req)
	if err != nil {
		logs.Logger.Errorf("NatProbe.PPCall(), error: %s.", err)
		return false
	}

	return true
}
