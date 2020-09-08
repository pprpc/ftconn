package logic

import (
	"fmt"

	"xcthings.com/hjyz/common"
	errc "xcthings.com/ftconn/common/errorcode"
	"xcthings.com/protoc/ftconnnat/NatTest2"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"

	lc "xcthings.com/ftconn/checknat-ms/common"
	m "xcthings.com/ftconn/checknat-ms/model"
)

// LNatTest2 NatTest2 Business logic
func LNatTest2(c pprpc.RPCConn, pkg *packets.CmdPacket, req *NatTest2.Req) (resp *NatTest2.Resp, code uint64, err error) {
	if req.AuthType == 1 {
		if req.Did == "" || req.DidSignkey == "" {
			code = errc.ParameterIllegal
			err = fmt.Errorf("The parameter is invalid: Did, DidSignkey")
			return
		}
		if lc.AuthDevice(req.Did, req.DidSignkey) == false {
			code = errc.PUBLICAUTHDENY
			err = fmt.Errorf("Did: %s, DidSignkey: %s not match", req.Did, req.DidSignkey)
			return
		}
	} else if req.AuthType == 2 {
		if req.UserId == 0 || req.UserPass == "" {
			code = errc.ParameterIllegal
			err = fmt.Errorf("The parameter is invalid: UserId, UserPass")
			return
		}
		if lc.AuthUserID(req.UserId, req.UserPass) == false {
			code = errc.PUBLICAUTHDENY
			err = fmt.Errorf("UserID: %d, UserPass: %s not match", req.UserId, req.UserPass)
			return
		}
	} else {
		code = errc.ParameterError
		err = fmt.Errorf(" Parameter error, AuthType: %d not support", req.AuthType)
		return
	}

	resp, code, err = m.MNatTest2(req)
	if err != nil {
		return
	}
	resp.OutsideIpaddr, resp.OutsidePort = common.GetRemoteIPPort(c.RemoteAddr())

	return
}
