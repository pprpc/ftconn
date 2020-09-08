package logic

import (
	"fmt"

	"github.com/pprpc/ftconn/model/ftconn"
	"github.com/pprpc/util/common"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
	"xcthings.com/protoc/ftconnp2p/P2PStepTwo"

	errc "github.com/pprpc/ftconn/common/errorcode"
	lc "github.com/pprpc/ftconn/p2p-ms/common"
	m "github.com/pprpc/ftconn/p2p-ms/model"
)

// LP2PStepTwo P2PStepTwo Business logic
func LP2PStepTwo(c pprpc.RPCConn, pkg *packets.CmdPacket, req *P2PStepTwo.Req) (resp *P2PStepTwo.Resp, code uint64, err error) {
	if lc.AuthDevice(req.Did, req.DidSignkey) == false {
		code = errc.ACCOUNTNOTMATCH
		err = fmt.Errorf("Did: %s, SignKey: %s not match", req.Did, req.DidSignkey)
		return
	}
	if req.Nat == nil {
		code = errc.ParameterIllegal
		err = fmt.Errorf("Natinfo no setting")
		return
	}
	// session_key
	q := new(ftconn.P2pinfo)
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

	resp, code, err = m.MP2PStepTwo(c, req)
	return
}
