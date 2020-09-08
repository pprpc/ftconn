package logic

import (
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"
	"xcthings.com/protoc/ftconnnat/NatProbe"

	m "xcthings.com/ftconn/checknat-ms/model"
)

// LNatProbe NatProbe Business logic
func LNatProbe(c pprpc.RPCConn, pkg *packets.CmdPacket, req *NatProbe.Req) (resp *NatProbe.Resp, code uint64, err error) {

	resp, code, err = m.MNatProbe(req)
	return
}
