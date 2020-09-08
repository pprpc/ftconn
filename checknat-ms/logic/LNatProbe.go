package logic

import (
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
	"xcthings.com/protoc/ftconnnat/NatProbe"

	m "github.com/pprpc/ftconn/checknat-ms/model"
)

// LNatProbe NatProbe Business logic
func LNatProbe(c pprpc.RPCConn, pkg *packets.CmdPacket, req *NatProbe.Req) (resp *NatProbe.Resp, code uint64, err error) {

	resp, code, err = m.MNatProbe(req)
	return
}
