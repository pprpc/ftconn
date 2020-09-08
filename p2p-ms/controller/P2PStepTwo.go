package controller

import (
	"fmt"
	"net"

	"xcthings.com/hjyz/common"
	"xcthings.com/hjyz/logs"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"
	"xcthings.com/protoc/ftconnp2p/P2PStepThree"
	"xcthings.com/protoc/ftconnp2p/P2PStepTwo"

	g "xcthings.com/ftconn/p2p-ms/common/global"
	l "xcthings.com/ftconn/p2p-ms/logic"
	m "xcthings.com/ftconn/p2p-ms/model"
)

// P2PStepTwoer .
type P2PStepTwoer struct{}

// ReqHandle P2PStepTwo request handle
func (t *P2PStepTwoer) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *P2PStepTwo.Req) (err error) {
	var code uint64
	var resp *P2PStepTwo.Resp
	resp, code, err = l.LP2PStepTwo(c, pkg, req)
	if code != 0 {
		pkg.Code = code
		logs.Logger.Warnf("%s, l.LP2PStepTwo, code: %d, err: %s.", c, code, err)
	}

	_, err = pprpc.WriteResp(c, pkg, resp)
	if err != nil {
		logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
	}
	if err != nil || pkg.Code != 0 {
		return
	}
	go p2pStepThree(req, resp, pkg)
	return
}

// RespHandle P2PStepTwo response handle
func (t *P2PStepTwoer) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *P2PStepTwo.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	return
}

// DestructHandle P2PStepTwo.
func (t *P2PStepTwoer) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, P2PStepTwo, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}

func p2pStepThree(req *P2PStepTwo.Req, resp *P2PStepTwo.Resp, pkg *packets.CmdPacket) {
	sPort, uPort, uIP, err := m.GetP2pinfoBySessionKey(req.SessionKey)
	if err != nil {
		logs.Logger.Errorf("m.GetP2pinfoBySessionKey(%s), %s.", req.SessionKey, err)
		return
	}
	uAddr := fmt.Sprintf("%s:%d", uIP, uPort)
	_, lp := g.UDPSrv1.GetListenInfo()
	if lp == sPort {
		// 发送 P2PStepThree
		sendP2PStepThree(req, resp, pkg, g.UDPSrv1, uAddr)
		return
	}
	if g.UDPSrv2 != nil {
		sendP2PStepThree(req, resp, pkg, g.UDPSrv2, uAddr)
	} else {
		logs.Logger.Warnf("not find P2pServer Port: %d.", sPort)
	}
}

func sendP2PStepThree(req *P2PStepTwo.Req, resp *P2PStepTwo.Resp, pkg *packets.CmdPacket, udps *pprpc.RPCUDPServer, uAddr string) {
	var err error
	var rAddr *net.UDPAddr
	rAddr, err = net.ResolveUDPAddr("udp4", uAddr)
	if err != nil {
		logs.Logger.Errorf("net.ResolveUDPAddr(), %s.", err)
		return
	}

	conn, err := udps.GetUDPConn(rAddr)
	if err != nil {
		logs.Logger.Errorf("RPCUDPServer.GetUDPConn(%s), %s.", uAddr, err)
		return
	}

	r := new(P2PStepThree.Req)
	r.Did = req.Did
	r.SessionKey = req.SessionKey
	r.OutsideIpaddr = resp.OutsideIpaddr
	r.OutsidePort = resp.OutsidePort
	r.Code = req.Code
	r.Nat = req.Nat

	// FIXME:
	err = pprpc.InvokeAsync(conn, P2PStepThree.CmdID, r, pkg.MessageType, pkg.EncType)
	if err != nil {
		logs.Logger.Errorf("pprpc.InvokeAsync(), %s.", err)
	}
	// write again
	common.SleepMs(100)
	pprpc.InvokeAsync(conn, P2PStepThree.CmdID, r, pkg.MessageType, pkg.EncType)
}
