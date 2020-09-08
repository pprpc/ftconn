package controller

import (
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"xcthings.com/protoc/ftconnp2p/P2PStepThree"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
)

// P2PStepThreeer .
type P2PStepThreeer struct{}

// ReqHandle P2PStepThree request handle
func (t *P2PStepThreeer) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *P2PStepThree.Req) (err error) {
	// 只有User端才会用到

	logs.Logger.Debugf("Req: %v.", req)
	resp := new(P2PStepThree.Resp)
	resp.SessionKey = req.SessionKey

	_, err = pprpc.WriteResp(c, pkg, resp)
	if err != nil {
		logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
	}
	//if err != nil || pkg.Code != 0 {
	// 	return
	//}
	return
}

// RespHandle P2PStepThree response handle
func (t *P2PStepThreeer) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *P2PStepThree.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	return
}

// DestructHandle P2PStepThree.
func (t *P2PStepThreeer) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, P2PStepThree, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}
