package model

import (
	"xcthings.com/hjyz/common"
	lc "xcthings.com/ftconn/checknat-ms/common"
	"xcthings.com/ftconn/model/ftconn"
	"xcthings.com/protoc/ftconnnat/ReportNat"
	"xcthings.com/pprpc"
)

// MReportNat ReportNat  
func MReportNat(c pprpc.RPCConn, req *ReportNat.Req) (resp *ReportNat.Resp, code uint64, err error) {
	t := new(ftconn.Natinfo)
	if req.AuthType == 1 {
		t.Did = req.Did
		t.UserID = 0
	} else {
		t.Did = ""
		t.UserID = req.UserId
	}
	if req.Nat.Upnp == lc.UPNPENABLE {
		t.Upnp = req.Nat.Upnp
		t.UpnpProtocol = req.Nat.UpnpProtocol
		t.UpnpPort = req.Nat.UpnpPort
	} else {
		t.Upnp = req.Nat.Upnp
		t.UpnpProtocol = 0
		t.UpnpPort = 0
	}
	t.NatType = req.Nat.NatType
	t.LocalIP = req.Nat.LocalIp
	t.LocalPort = req.Nat.LocalPort
	t.OutsideIpaddr = common.GetIPAddr(c.RemoteAddr())
	t.LocalSsid = req.Nat.LocalSsid
	t.LocalGateway = req.Nat.LocalGateway
	t.GatewayMac = req.Nat.GatewayMac
	t.LastTime = common.GetTimeMs()

	code, err = t.Set()

	resp = new(ReportNat.Resp)
	resp.Upnp = req.Nat.Upnp

	return
}

/*
	Did           string `json:"did" xorm:"did"`
	UserID        int64  `json:"user_id" xorm:"user_id"`
	Upnp          int32  `json:"upnp" xorm:"upnp"`
	UpnpProtocol  int32  `json:"upnp_protocol"`
	UpnpPort      int32  `json:"upnp_pprt"`
	NatType       int32  `json:"nat_type"`
	LocalIP       string `json:"local_ip" xorm:"local_ip"`
	LocalPort     int32  `json:"local_port"`
	OutsideIpaddr string `json:"outside_ipaddr" xorm:"outside_ipaddr"`
	LocalSsid     string `json:"local_ssid"`
	LocalGateway  string `json:"local_gateway"`
	GatewayMac    string `json:"gateway_mac"`
	LastTime      int64  `json:"last_time"`
*/
