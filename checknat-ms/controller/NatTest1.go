package controller

import (
	"fmt"
	"net"

	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
	"xcthings.com/protoc/ftconnnat/NatProbe"
	"xcthings.com/protoc/ftconnnat/NatTest1"
	"xcthings.com/protoc/ftconnnat/ProbeConfig"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"

	g "github.com/pprpc/ftconn/checknat-ms/common/global"
	l "github.com/pprpc/ftconn/checknat-ms/logic"
)

// NatTest1er .
type NatTest1er struct{}

// ReqHandle NatTest1 request handle
func (t *NatTest1er) ReqHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, req *NatTest1.Req) (err error) {
	var code uint64
	var resp *NatTest1.Resp
	resp, code, err = l.LNatTest1(c, pkg, req)
	if code != 0 {
		pkg.Code = code
		logs.Logger.Warnf("%s, l.LNatTest1, code: %d, err: %s.", c, code, err)
	}

	_, err = pprpc.WriteResp(c, pkg, resp)
	if err != nil {
		logs.Logger.Errorf("%s, %s, write response error:  %s.", c, pkg.CmdName, err)
	}
	if err != nil || pkg.Code != 0 {
		return
	}
	go notifyPort2(c, pkg, resp)
	go notifyNat2(req, resp)
	return
}

// RespHandle NatTest1 response handle
func (t *NatTest1er) RespHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *NatTest1.Resp) (err error) {
	if pkg.Code != 0 {
		logs.Logger.Errorf("%s, %s, Seq: %d, Response recv error code: %d.",
			c, pkg.CmdName, pkg.CmdSeq, pkg.Code)
		return
	} else {
		// resp.Srv2.Ipaddr
		logs.Logger.Debugf("Srv2: %s:%d, Port2: %d, OutIP: %s, OutPort: %d .",
			resp.Srv2.Ipaddr, resp.Srv2.Ports, resp.NatPort2, resp.OutsideIpaddr, resp.OutsidePort)
	}
	return
}

// DestructHandle NatTest1.
func (t *NatTest1er) DestructHandle(c pprpc.RPCConn, pkg *packets.CmdPacket, startMs int64) {
	logs.Logger.Infof("%s, NatTest1, DestructHandle, useMs: %d.", c, common.GetTimeMs()-startMs)
}

// notifyNat2 通知另外一个NAT Server.
func notifyNat2(req *NatTest1.Req, resp *NatTest1.Resp) {
	conf := new(ProbeConfig.Req)
	conf.SrcIp = resp.Srv2.Ipaddr
	conf.SrcPort = resp.Srv2.Ports[0]

	conf.DstIp = resp.OutsideIpaddr
	conf.DstPort = resp.OutsidePort

	msgid, err := g.PMsg.ProbeConfig(conf, resp.Srv2.Ipaddr)
	if err != nil {
		logs.Logger.Errorf("g.PMsg.ProbeConfig(),Did: %s , error: %s.", req.Did, err)
		return
	} else {
		logs.Logger.Debugf("g.PMsg.ProbeConfig(), Did: %s , MsgID: %s .",
			req.Did, msgid)
	}
}

func notifyPort2(c pprpc.RPCConn, pkg *packets.CmdPacket, resp *NatTest1.Resp) {
	req := new(NatProbe.Req)
	req.Check = "NatProbe"

	var err error
	var rAddr *net.UDPAddr
	rAddr, err = net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", resp.OutsideIpaddr, resp.OutsidePort))
	if err != nil {
		logs.Logger.Errorf("net.ResolveUDPAddr(), %s.", err)
		return
	}

	conn, e := g.UdpSrv2.GetUDPConn(rAddr)
	if e != nil {
		logs.Logger.Errorf("g.UdpSrv2.GetUDPConn(), %s.", e)
		return
	}
	defer func() {
		conn = nil
	}()
	// FIXME: 后续此处修改为超时重发更合理,当前策略时发送两次
	err = pprpc.InvokeAsync(conn, NatProbe.CmdID, req, pkg.MessageType, pkg.EncType)
	if err != nil {
		logs.Logger.Errorf("pprpc.InvokeAsync(), %s.", err)
	}
	// write again FIXME
	common.SleepMs(100)
	pprpc.InvokeAsync(conn, NatProbe.CmdID, req, pkg.MessageType, pkg.EncType)
}
