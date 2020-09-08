package common

import (
	"fmt"
	"io"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"

	g "xcthings.com/ftconn/checknat-ms/common/global"
)

// UDPInvoke udp invoke sync
func UDPInvoke(rAddr *net.UDPAddr, cmdid uint64, req interface{}, mt, crypt uint8) (pkg *packets.CmdPacket, resp interface{}, err error) {
	conn, e := g.UdpSrv2.GetUDPConn(rAddr)
	if e != nil {
		err = fmt.Errorf("g.UdpSrv2.GetUDPConn(), %s", e)
		return
	}

	seq := pprpc.GetSeqID()
	cmd := packets.NewCmdPacket(mt)
	cmd.FixHeader.SetProtocol(packets.PROTOUDP)
	cmd.CmdSeq = seq
	cmd.CmdID = cmdid
	cmd.EncType = crypt
	cmd.RPCType = packets.RPCREQ

	if mt == packets.TYPEPBBIN {
		cmd.Payload, err = proto.Marshal(req.(proto.Message))
	} else if mt == packets.TYPEPBJSON {
		cmd.Payload, err = proto.MarshalMessageSetJSON(req)
	}
	if err != nil {
		return
	}

	_, err = cmd.Write(conn)
	if err != nil {
		return
	}

	var p packets.PPPacket
	p, err = packets.ReadUDPPacket(conn)
	if err == io.EOF {
		return
	} else if err != nil {
		err = fmt.Errorf("packets.ReadUDPPacket(), %s", err)
		return
	}
	switch p.(type) {
	case *packets.CmdPacket:
		cmd := p.(*packets.CmdPacket)
		resp, err = pprpc.DecodePkg(cmd, g.Service)
	default:
		err = fmt.Errorf("not support PPPacket type")
	}
	return
}
