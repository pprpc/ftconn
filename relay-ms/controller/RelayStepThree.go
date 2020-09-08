package controller

import (
	"xcthings.com/hjyz/common"
	"xcthings.com/hjyz/logs"
	"xcthings.com/protoc/ftconnrelay/RelayStepThree"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"
)

// RelayStepThreeer .
type RelayStepThreeer struct{}

// ReqHandle RelayStepThree request handle
func (t *RelayStepThreeer) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *RelayStepThree.Req) (err error) {
	// logs.Logger.Debugf("Req: %v.", req)
	// resp := new(RelayStepThree.Resp)
	// resp.SessionKey = req.SessionKey

	// _, err = pprpc.WriteResp(c, pkg, resp)
	// if err != nil {
	// 	logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
	// }
	// // Write SyncConn
	// sc := new(SyncConn.Req)
	// sc.SyncFlag = 1
	// err = pprpc.InvokeAsync(c, SyncConn.CmdID, sc, pkg.MessageType, pkg.EncType)
	// if err != nil {
	// 	logs.Logger.Errorf("pprpc.InvokeAsync(), Write SyncConn, %s.", err)
	// 	return
	// }
	return
}

// RespHandle RelayStepThree response handle
func (t *RelayStepThreeer) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *RelayStepThree.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	}
	return
}

// DestructHandle RelayStepThree.
func (t *RelayStepThreeer) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, RelayStepThree, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}
