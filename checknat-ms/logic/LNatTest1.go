package logic

import (
	"fmt"

	"xcthings.com/hjyz/common"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"
	pc "xcthings.com/protoc/common"
	"xcthings.com/protoc/ftconnnat/NatTest1"

	lc "xcthings.com/ftconn/checknat-ms/common"
	g "xcthings.com/ftconn/checknat-ms/common/global"
	m "xcthings.com/ftconn/checknat-ms/model"
	errc "xcthings.com/ftconn/common/errorcode"
)

// LNatTest1 NatTest1 Business logic
func LNatTest1(c pprpc.RPCConn, pkg *packets.CmdPacket, req *NatTest1.Req) (resp *NatTest1.Resp, code uint64, err error) {
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
	resp, code, err = m.MNatTest1(req)
	if err != nil {
		return
	}
	resp.OutsideIpaddr, resp.OutsidePort = common.GetRemoteIPPort(c.RemoteAddr())
	resp.NatPort2 = g.PConf.Nat.Port2

	srv := new(pc.SrvAddr)
	srv.Ipaddr = g.PConf.Nat.Ipaddr2
	srv.Ports = append(srv.Ports, g.PConf.Nat.Port2)

	resp.Srv2 = srv

	return
}
