package logic

import (
	"xcthings.com/protoc/ftconnrelay/RelayStepThree"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"

	m "xcthings.com/ftconn/relay-ms/model"
)

// LRelayStepThree RelayStepThree Business logic
func LRelayStepThree(c pprpc.RPCConn, pkg *packets.CmdPacket, req *RelayStepThree.Req) (resp *RelayStepThree.Resp, code uint64, err error) {

	resp, code, err = m.MRelayStepThree(req)
	return
}
