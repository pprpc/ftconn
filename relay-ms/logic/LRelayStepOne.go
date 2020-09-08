package logic

import (
	"fmt"

	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
	"xcthings.com/protoc/ftconnrelay/RelayStepOne"

	errc "github.com/pprpc/ftconn/common/errorcode"
	lc "github.com/pprpc/ftconn/relay-ms/common"
	m "github.com/pprpc/ftconn/relay-ms/model"
)

// LRelayStepOne RelayStepOne Business logic
func LRelayStepOne(c pprpc.RPCConn, pkg *packets.CmdPacket, req *RelayStepOne.Req) (resp *RelayStepOne.Resp, code uint64, err error) {
	// 1 auth
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

	resp, code, err = m.MRelayStepOne(req)
	return
}
