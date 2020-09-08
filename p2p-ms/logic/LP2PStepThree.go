package logic

import (
	"xcthings.com/protoc/ftconnp2p/P2PStepThree"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"

	m "xcthings.com/ftconn/p2p-ms/model"
)

// LP2PStepThree P2PStepThree Business logic
func LP2PStepThree(c pprpc.RPCConn, pkg *packets.CmdPacket, req *P2PStepThree.Req) (resp *P2PStepThree.Resp, code uint64, err error) {

	resp, code, err = m.MP2PStepThree(req)
	return
}
