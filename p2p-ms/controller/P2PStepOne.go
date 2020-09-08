package controller

import (
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
	pc "xcthings.com/protoc/common"
	"xcthings.com/protoc/ftconnp2p/NotifyConn"
	"xcthings.com/protoc/ftconnp2p/P2PStepOne"

	g "xcthings.com/ftconn/p2p-ms/common/global"
	l "xcthings.com/ftconn/p2p-ms/logic"
)

// P2PStepOneer .
type P2PStepOneer struct{}

// ReqHandle P2PStepOne request handle
func (t *P2PStepOneer) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *P2PStepOne.Req) (err error) {
	var code uint64
	var resp *P2PStepOne.Resp
	resp, code, err = l.LP2PStepOne(c, pkg, req)
	if code != 0 {
		pkg.Code = code
		logs.Logger.Warnf("%s, l.LP2PStepOne, code: %d, err: %s.", c, code, err)
	}

	_, err = pprpc.WriteResp(c, pkg, resp)
	if err != nil {
		logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
	}
	if err != nil || pkg.Code != 0 {
		return
	}
	go notifyConn(req, resp)
	return
}

// RespHandle P2PStepOne response handle
func (t *P2PStepOneer) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *P2PStepOne.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	return
}

// DestructHandle P2PStepOne.
func (t *P2PStepOneer) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, P2PStepOne, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}

func notifyConn(req *P2PStepOne.Req, resp *P2PStepOne.Resp) {
	r := new(NotifyConn.Req)
	r.ConnType = 1
	r.SessionKey = resp.SessionKey
	r.UserId = req.UserId

	srv := new(pc.P2PSrv)
	srv.Ipaddr = g.PConf.P2p.WanIP
	srv.Ports = append(srv.Ports, g.PConf.P2p.WanPort)

	p2p := new(pc.P2PInfo)
	p2p.Nat = req.Nat
	p2p.P2PServer = srv
	p2p.OutsideIp = resp.OutsideIpaddr
	p2p.OutsidePort = resp.OutsidePort

	r.P2P = p2p

	msgid, err := g.PMsg.NotifyConn(r, req.RemoteDid)
	if err != nil {
		logs.Logger.Errorf("g.PMsg.NotifyConn(), UserID: %d, Did: %s , error: %s.", req.UserId, req.RemoteDid, err)
		return
	} else {
		logs.Logger.Debugf("g.PMsg.NotifyConn(), UserID: %d, Did: %s , MsgID: %s .",
			req.UserId, req.RemoteDid, msgid)
	}
}
