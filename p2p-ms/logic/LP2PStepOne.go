package logic

import (
	"fmt"

	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"
	"xcthings.com/protoc/ftconnp2p/P2PStepOne"

	errc "xcthings.com/ftconn/common/errorcode"
	lc "xcthings.com/ftconn/p2p-ms/common"
	m "xcthings.com/ftconn/p2p-ms/model"
)

// LP2PStepOne P2PStepOne Business logic
func LP2PStepOne(c pprpc.RPCConn, pkg *packets.CmdPacket, req *P2PStepOne.Req) (resp *P2PStepOne.Resp, code uint64, err error) {
	if lc.AuthUserID(req.UserId, req.UserPass) == false {
		code = errc.ACCOUNTNOTMATCH
		err = fmt.Errorf("userid, userpass not match")
		return
	}
	if lc.ACLUserIDAndDID(req.UserId, req.RemoteDid) == false {
		code = errc.ACCESSDEVICEDENY
		err = fmt.Errorf("UserID: %d access device: %s, deny", req.UserId, req.RemoteDid)
		return
	}
	// if req.Nat == nil {
	// 	code = errc.ParameterIllegal
	// 	err = fmt.Errorf("Natinfo no setting")
	// 	return
	// }

	resp, code, err = m.MP2PStepOne(c, req)
	return
}
