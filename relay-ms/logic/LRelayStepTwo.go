package logic

import (
	"fmt"

	errc "xcthings.com/ftconn/common/errorcode"
	"xcthings.com/ftconn/model/ftconn"
	"github.com/pprpc/util/common"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
	"xcthings.com/protoc/ftconnrelay/RelayStepTwo"

	lc "xcthings.com/ftconn/relay-ms/common"
	m "xcthings.com/ftconn/relay-ms/model"
)

// LRelayStepTwo RelayStepTwo Business logic
func LRelayStepTwo(c pprpc.RPCConn, pkg *packets.CmdPacket, req *RelayStepTwo.Req) (resp *RelayStepTwo.Resp, code uint64, err error) {
	if lc.AuthDevice(req.Did, req.DidSignkey) == false {
		code = errc.ACCOUNTNOTMATCH
		err = fmt.Errorf("Did: %s, SignKey: %s not match", req.Did, req.DidSignkey)
		return
	}
	q := new(ftconn.Relayinfo)
	q.SessionKey = req.SessionKey

	code, err = q.GetBySessionKey()
	if err != nil {
		return
	}
	if q.UserID != req.RemoteUid || q.Did != req.Did {
		code = errc.CONNSESSIONKEYNOTMATCH
		err = fmt.Errorf("SessionKey: %s not match(%d - %s)", req.SessionKey, req.RemoteUid, req.Did)
		return
	}
	curMS := common.GetTimeMs()
	if curMS-q.UserTime > 3000 {
		code = errc.SESSIONKEYTIMEOUT
		err = fmt.Errorf("SessionKey: %s timeout(%d - %s), %d - %d = %d",
			req.SessionKey, req.RemoteUid, req.Did, curMS, q.UserTime, curMS-q.UserTime)
		return
	}

	resp, code, err = m.MRelayStepTwo(req)
	return
}
