package controller

import (
	"fmt"
	"net"

	"github.com/pprpc/util/logs"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPublish"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
	"github.com/pprpc/core/ppudp"
	"xcthings.com/protoc/ftconnnat/NatProbe"
	"xcthings.com/protoc/ftconnnat/ProbeConfig"

	g "github.com/pprpc/ftconn/checknat-ms/common/global"
)

// SendNatProbe  .
func SendNatProbe(r *PPMQPublish.Req) {
	var err error
	conf := new(ProbeConfig.Req)
	if r.Format == int32(packets.TYPEPBBIN) {
		err = pprpc.Unmarshal(r.MsgBody, conf)
	} else if r.Format == int32(packets.TYPEPBJSON) {
		err = pprpc.UnmarshalJSON(r.MsgBody, conf)
	} else {
		logs.Logger.Errorf("not support format: %d.", r.Format)
		return
	}
	if err != nil {
		logs.Logger.Errorf("pprpc.Unmarshal/UnmarshalJSON(), %s.", err)
		return
	}
	var rAddr *net.UDPAddr
	rAddr, err = net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", conf.DstIp, conf.DstPort))
	if err != nil {
		logs.Logger.Errorf("net.ResolveUDPAddr(), %s.", err)
		return
	}
	var conn *ppudp.Connection
	conn, err = g.UdpSrv2.GetUDPConn(rAddr)
	if err != nil {
		logs.Logger.Errorf("g.UdpSrv2.GetUDPConn(rAddr), %s.", err)
		return
	}
	defer func() {
		conn = nil
	}()

	req := new(NatProbe.Req)
	req.Check = "NatProbe"

	seq := pprpc.GetSeqID()
	cmd := packets.NewCmdPacket(uint8(r.Format))
	cmd.FixHeader.SetProtocol(packets.PROTOUDP)
	cmd.CmdSeq = seq
	cmd.CmdID = NatProbe.CmdID
	cmd.EncType = g.PConf.Encrypt //packets.AESNONE //packets.AES256CBC FIXME
	cmd.RPCType = packets.RPCREQ

	if r.Format == int32(packets.TYPEPBBIN) {
		cmd.Payload, err = pprpc.Marshal(req)
	} else if r.Format == int32(packets.TYPEPBJSON) {
		cmd.Payload, err = pprpc.MarshalJSON(req)
	}
	if err != nil {
		logs.Logger.Errorf("pprpc.Marshal/MarshalJSON, %s.", err)
		return
	}

	_, err = cmd.Write(conn)
	if err != nil {
		logs.Logger.Errorf("cmd.Write(), %s.", err)
		return
	}
	logs.Logger.Debugf("MsgID: %s, NatProbe: %s, Write OK.", r.MsgId, rAddr)
	// write again FIXME
	cmd.Write(conn)
}
