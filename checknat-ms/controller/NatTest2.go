package controller

import (
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"xcthings.com/protoc/ftconnnat/NatTest2"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"

	l "github.com/pprpc/ftconn/checknat-ms/logic"
)

// NatTest2er .
type NatTest2er struct{}

// ReqHandle NatTest2 request handle
func (t *NatTest2er) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *NatTest2.Req) (err error) {
	var code uint64
	var resp *NatTest2.Resp
	resp, code, err = l.LNatTest2(c, pkg, req)
	if code != 0 {
		pkg.Code = code
		logs.Logger.Warnf("%s, l.LNatTest2, code: %d, err: %s.", c, code, err)
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

// RespHandle NatTest2 response handle
func (t *NatTest2er) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *NatTest2.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	} else {
		logs.Logger.Debugf("OutIP: %s, OutPort: %d .", resp.OutsideIpaddr, resp.OutsidePort)
	}
	return
}

// DestructHandle NatTest2.
func (t *NatTest2er) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, NatTest2, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}
