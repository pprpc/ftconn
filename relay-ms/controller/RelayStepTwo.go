package controller

import (
	"fmt"

	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
	"xcthings.com/protoc/ftconnrelay/RelayStepThree"
	"xcthings.com/protoc/ftconnrelay/RelayStepTwo"

	g "xcthings.com/ftconn/relay-ms/common/global"
	l "xcthings.com/ftconn/relay-ms/logic"
)

// RelayStepTwoer .
type RelayStepTwoer struct{}

// ReqHandle RelayStepTwo request handle
func (t *RelayStepTwoer) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *RelayStepTwo.Req) (err error) {
	var code uint64
	var resp *RelayStepTwo.Resp
	resp, code, err = l.LRelayStepTwo(c, pkg, req)
	if code != 0 {
		pkg.Code = code
		logs.Logger.Warnf("%s, l.LRelayStepTwo, code: %d, err: %s.", c, code, err)
		logs.Logger.Debugf("%s, RelayStepTwo, req: %v.", c, req)
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
	ci.Did = req.Did
	ci.UserID = req.RemoteUid
	ci.SessionKey = req.SessionKey
	ci.Class = "D"
	ci.Connid = fmt.Sprintf("%s-%s-%s-%d", ci.SessionKey, ci.Class, ci.Did, ci.UserID)
	ci.CallPrehookCB = false
	ci.InitTimeSec = common.GetTimeSec()

	err = c.SetAttr(ci)
	if err != nil {
		logs.Logger.Errorf("c.SetAttr(ci), error: %s.", err)
		return
	}
	g.Sess.Push(ci.Connid, c)

	go sendRelayStepThree(c, pkg, req)
	return
}

// RespHandle RelayStepTwo response handle
func (t *RelayStepTwoer) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *RelayStepTwo.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	return
}

// DestructHandle RelayStepTwo.
func (t *RelayStepTwoer) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, RelayStepTwo, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}

func sendRelayStepThree(c pprpc.RPCConn, pkg *packets.CmdPacket, req *RelayStepTwo.Req) {
	// find user conn，write rpc: RelayStepThree
	_t, e := c.GetAttr()
	if e != nil {
		logs.Logger.Errorf("c.GetAttr(), error: %s.", e)
		return
	}
	if _t == nil {
		logs.Logger.Errorf("conn.attr is nil.")
		return
	}
	ci := _t.(*g.ConnAttr)

	var uc pprpc.RPCConn

	uConnid := fmt.Sprintf("%s-U-%s-%d", ci.SessionKey, ci.Did, ci.UserID)
	v, e := g.Sess.Get(uConnid)
	if e != nil {
		logs.Logger.Errorf("Get connid: %s, %s.", uConnid, e)
		return
	}
	uc = v.(pprpc.RPCConn)

	r := new(RelayStepThree.Req)
	r.Did = req.Did
	r.SessionKey = req.SessionKey
	r.Code = req.Code

	err := pprpc.InvokeAsync(uc, RelayStepThree.CmdID, r, pkg.MessageType, pkg.EncType)
	if err != nil {
		logs.Logger.Errorf("pprpc.InvokeAsync(), %s.", err)
		return
	}
	logs.Logger.Infof("Notify, UserID: %d, Did: %s , SessionKey: %s OK.", ci.UserID, ci.Did, req.SessionKey)

	_u, e := uc.GetAttr()
	if e != nil {
		logs.Logger.Errorf("uc.GetAttr(), error: %s.", e)
		return
	}
	if _u == nil {
		logs.Logger.Errorf("conn.attr is nil.")
		return
	}
	uci := _u.(*g.ConnAttr)

	ci.CallPrehookCB = true
	uci.CallPrehookCB = true
	c.SetAutoCrypt(false)
	uc.SetAutoCrypt(false)
}
