package controller

import (
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
	"xcthings.com/protoc/ftconnnat/NatProbe"

	l "xcthings.com/ftconn/checknat-ms/logic"
)

// NatProbeer .
type NatProbeer struct{}

// ReqHandle NatProbe request handle
func (t *NatProbeer) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *NatProbe.Req) (err error) {
	var code uint64
	var resp *NatProbe.Resp
	resp, code, err = l.LNatProbe(c, pkg, req)
	if code != 0 {
		pkg.Code = code
		logs.Logger.Warnf("%s, l.LNatProbe, code: %d, err: %s.", c, code, err)
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

// RespHandle NatProbe response handle
func (t *NatProbeer) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *NatProbe.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	return
}

// DestructHandle NatProbe.
func (t *NatProbeer) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, NatProbe, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}
