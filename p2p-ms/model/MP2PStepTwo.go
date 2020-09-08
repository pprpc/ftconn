package model

import (
	"xcthings.com/hjyz/common"
	"xcthings.com/pprpc"
	"xcthings.com/ftconn/model/ftconn"
	"xcthings.com/protoc/ftconnp2p/P2PStepTwo"
)

// MP2PStepTwo P2PStepTwo  
func MP2PStepTwo(c pprpc.RPCConn, req *P2PStepTwo.Req) (resp *P2PStepTwo.Resp, code uint64, err error) {
	p2p := new(ftconn.P2pinfo)
	p2p.DeviceOutsideIP, p2p.DeviceOutsidePort = common.GetRemoteIPPort(c.RemoteAddr())
	if req.Nat != nil {
		p2p.DeviceLocalIP = req.Nat.LocalIp
		p2p.DeviceLocalPort = req.Nat.LocalPort
	}
	p2p.SessionKey = req.SessionKey
	p2p.DeviceTime = common.GetTimeMs()

	code, err = p2p.UpdateBySessionKey()
	if err != nil {
		return
	}

	resp = new(P2PStepTwo.Resp)
	resp.OutsideIpaddr = p2p.DeviceOutsideIP
	resp.OutsidePort = p2p.DeviceOutsidePort
	return
}
