package controller

import (
	"fmt"

	"xcthings.com/hjyz/common"
	"xcthings.com/hjyz/logs"
	"xcthings.com/protoc/ftconnp2p/NotifyConn"
	"xcthings.com/protoc/ftconnrelay/RelayStepOne"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"

	pc "xcthings.com/protoc/common"

	g "xcthings.com/ftconn/relay-ms/common/global"
	l "xcthings.com/ftconn/relay-ms/logic"
)

// RelayStepOneer .
type RelayStepOneer struct{}

// ReqHandle RelayStepOne request handle
func (t *RelayStepOneer) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *RelayStepOne.Req) (err error) {
	var code uint64
	var resp *RelayStepOne.Resp
	resp, code, err = l.LRelayStepOne(c, pkg, req)
	if code != 0 {
		pkg.Code = code
		logs.Logger.Warnf("%s, l.LRelayStepOne, code: %d, err: %s.", c, code, err)
		logs.Logger.Debugf("%s, RelayStepOne, req: %v.", c, req)
	}

	_, err = pprpc.WriteResp(c, pkg, resp)
	if err != nil {
		logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
	}
	if err != nil || pkg.Code != 0 {
		return
	}

	// add connection
	ci := new(g.ConnAttr)
	ci.Did = req.RemoteDid
	ci.UserID = req.UserId
	ci.SessionKey = resp.SessionKey
	ci.Class = "U"
	ci.Connid = fmt.Sprintf("%s-%s-%s-%d", ci.SessionKey, ci.Class, ci.Did, ci.UserID)
	ci.CallPrehookCB = false
	ci.InitTimeSec = common.GetTimeSec()

	err = c.SetAttr(ci)
	if err != nil {
		logs.Logger.Errorf("c.SetAttr(ci), error: %s.", err)
		return
	}
	g.Sess.Push(ci.Connid, c)

	go notifyConn(req, resp)
	return
}

// RespHandle RelayStepOne response handle
func (t *RelayStepOneer) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *RelayStepOne.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	return
}

// DestructHandle RelayStepOne.
func (t *RelayStepOneer) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, RelayStepOne, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}

func notifyConn(req *RelayStepOne.Req, resp *RelayStepOne.Resp) {
	r := new(NotifyConn.Req)
	r.ConnType = 2
	r.SessionKey = resp.SessionKey
	r.UserId = req.UserId

	srv := new(pc.RelaySrv)
	srv.Ipaddr = g.PConf.Relay.WanIP
	srv.Ports = g.PConf.Relay.WanPort
	srv.ProxyIpaddr = req.ProxyIpaddr
	srv.ProxyPorts = append(srv.ProxyPorts, req.ProxyPort)

	r.Relay = srv
	msgid, err := g.PMsg.NotifyConn(r, req.RemoteDid)
	if err != nil {
		logs.Logger.Errorf("g.PMsg.NotifyConn(), UserID: %d, Did: %s , error: %s.", req.UserId, req.RemoteDid, err)
		return
	} else {
		logs.Logger.Debugf("g.PMsg.NotifyConn(), UserID: %d, Did: %s , MsgID: %s .",
			req.UserId, req.RemoteDid, msgid)
	}
}
