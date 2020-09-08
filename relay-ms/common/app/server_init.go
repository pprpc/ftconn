package app

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"strconv"

	g "xcthings.com/ftconn/relay-ms/common/global"
	"xcthings.com/hjyz/common"
	"xcthings.com/hjyz/logs"
	"xcthings.com/micro/svc"
	"xcthings.com/pprpc"
	"xcthings.com/pprpc/packets"
	"xcthings.com/pprpc/pptcp"
	"xcthings.com/pprpc/ppudp"
)

var tcplis []*pprpc.RPCTCPServer
var udplis []*pprpc.RPCUDPServer

func serverInit() (err error) {
	for _, lis := range g.Conf.Listen {
		err = runServer(lis)
		if err != nil {
			logs.Logger.Errorf("runServer(lis), error: %s.", err)
			return
		}
	}
	return
}

func runServer(lis svc.LisConf) error {
	u, e := url.ParseRequestURI(lis.URI)
	if e != nil {
		return e
	}
	switch u.Scheme {
	case "udp":
		p := u.Port()
		_t, e := strconv.Atoi(p)
		if e != nil {
			return e
		}
		usrv, err := pprpc.NewRPCUDPServer(u.Hostname(), int(_t), int(g.PConf.MaxSession))
		if err != nil {
			return fmt.Errorf("pprpc.NewRPCUDPServer(), error: %s", err)
		}
		usrv.HBCB = cb
		usrv.Service = g.Service
		usrv.DisconnectCB = closeUDP
		usrv.PreHookCB = preHookCB
		usrv.SetReadTimeout(lis.ReadTimeout)
		logs.Logger.Infof("Listen UDPServer: %s.", lis.URI)
		go usrv.Serve()
		udplis = append(udplis, usrv)
	default:
		var tlsc *tls.Config
		if lis.TLSCrt != "" && lis.TLSKey != "" {
			tlsc, e = pprpc.GetTLSConfig(lis.TLSCrt, lis.TLSKey)
			if e != nil {
				return fmt.Errorf("pprpc.GetTLSConfig(), %s", e)
			}
		} else {
			tlsc = nil
		}

		srv, err := pprpc.NewRPCTCPServer(u, tlsc)
		if err != nil {
			return fmt.Errorf("pprpc.NewRPCTCPServer(), error: %s", err)
		}
		srv.HBCB = cb
		srv.Service = g.Service
		srv.DisconnectCB = closeTCP
		srv.PreHookCB = preHookCB // prehookCallBack func(packets.PPPacket, RPCConn) continue bool
		srv.SetReadTimeout(lis.ReadTimeout)

		logs.Logger.Infof("Listen TCPServer: %s.", lis.URI)
		go srv.Serve()
		tcplis = append(tcplis, srv)
	}
	return nil
}

func cb(pkg *packets.HBPacket, c pprpc.RPCConn) error {
	logs.Logger.Debugf("%s, HBPacket, MessageType: %d.", c, pkg.MessageType)
	_, err := pkg.Write(c)
	return err
}

func preHookCB(pkg packets.PPPacket, c pprpc.RPCConn) bool {
	_t, e := c.GetAttr()
	if e != nil {
		logs.Logger.Debugf("c.GetAttr(), error: %s.", e)
		return true
	}
	if _t == nil {
		logs.Logger.Debugf("conn.attr is nil.")
		return true
	}
	ci := _t.(*g.ConnAttr)

	runTime := common.GetTimeSec() - ci.InitTimeSec
	if ci.CallPrehookCB == false {
		return true
	}
	var peerConnid string
	if ci.Class == "D" {
		peerConnid = fmt.Sprintf("%s-U-%s-%d", ci.SessionKey, ci.Did, ci.UserID)
	} else {
		peerConnid = fmt.Sprintf("%s-D-%s-%d", ci.SessionKey, ci.Did, ci.UserID)
	}
	v, e := g.Sess.Get(peerConnid)
	if e == nil {
		switch pkg.(type) {
		case *packets.AVPacket:
			_t := pkg.(*packets.AVPacket)
			b := append(_t.RawHeader, _t.VarHeader...)
			b = append(b, _t.RAWPayload...)
			ci.RecvByte = ci.RecvByte + int64(len(b)) + 54
			v.(pprpc.RPCConn).Write(b)
			// logs.Logger.Debugf("%s -> %s, AVPacket, Length: %d , RunTime: %d, RecvByte: %d .",
			// 	c, v.(pprpc.RPCConn), _t.FixHeader.Length, runTime, ci.RecvByte)
			// logs.Logger.Debugf("Msgtype: %d, Length: %d.", _t.FixHeader.MessageType, _t.FixHeader.Length)
			// logs.Logger.Debugf("AVIFrame: %d, AVFormat: %d, EncType: %d, AVChannel: %d, AVSeq: %d, Timestamp: %d, EncLength: %d.",
			// 	_t.AVIFrame, _t.AVFormat, _t.EncType, _t.AVChannel, _t.AVSeq, _t.Timestamp, _t.EncLength)
			// logs.Logger.Debugf("Header: [%s].", common.ByteConvertString(append(_t.FixHeader.RawHeader, _t.VarHeader...)))
			// logs.Logger.Debugf("Frame: [%s]", common.ByteConvertString(_t.RAWPayload))
		case *packets.FilePacket:
			_t := pkg.(*packets.FilePacket)
			b := append(_t.RawHeader, _t.VarHeader...)
			b = append(b, _t.Payload...)
			ci.RecvByte = ci.RecvByte + int64(len(b)) + 54
			v.(pprpc.RPCConn).Write(b)
			logs.Logger.Debugf("%s -> %s, FilePacket, fid: %d, length: %d.",
				c, v.(pprpc.RPCConn), _t.FileID, _t.FixHeader.Length)

		case *packets.CmdPacket:
			_t := pkg.(*packets.CmdPacket)
			b := append(_t.RawHeader, _t.VarHeader...)
			b = append(b, _t.RAWPayload...)
			ci.RecvByte = ci.RecvByte + int64(len(b)) + 54
			v.(pprpc.RPCConn).Write(b)
			logs.Logger.Debugf("%s -> %s, CmdPacket, Length: %d, CmdID: %d, RunTime: %d,RecvByte: %d.",
				c, v.(pprpc.RPCConn), _t.FixHeader.Length, _t.CmdID, runTime, ci.RecvByte)
		case *packets.CustomerPacket:
			_t := pkg.(*packets.CustomerPacket)
			b := append(_t.RawHeader, _t.Payload...)
			ci.RecvByte = ci.RecvByte + int64(len(b)) + 54
			v.(pprpc.RPCConn).Write(b)

			logs.Logger.Debugf("%s -> %s, CustomerPacket, Length: %d , RunTime: %d, RecvByte: %d .",
				c, v.(pprpc.RPCConn), _t.FixHeader.Length, runTime, ci.RecvByte)
		}
	}
	return false
}

func closeUDP(conn *ppudp.Connection) {
	disconnectCB(conn)
}

func closeTCP(conn *pptcp.Connection) {
	disconnectCB(conn)
}

func disconnectCB(c pprpc.RPCConn) {
	// 1, offlineEvent
	// 2, Session Remove
	_t, e := c.GetAttr()
	if e != nil {
		logs.Logger.Debugf("c.GetAttr(), error: %s.", e)
		return
	}
	if _t == nil {
		logs.Logger.Debugf("conn.attr is nil.")
		return
	}
	ci := _t.(*g.ConnAttr)

	runTime := common.GetTimeSec() - ci.InitTimeSec
	if runTime == 0 {
		runTime = 1
	}
	var OtherConnid string
	v, e := g.Sess.Get(ci.Connid)
	if e == nil {
		if v.(pprpc.RPCConn).RemoteAddr() == c.RemoteAddr() && v.(pprpc.RPCConn).Type() == c.Type() {
			logs.Logger.Infof("%s, Sess.Remove(%s), recv speed: (%d/%d) = %d(Byte/Sec).", c, ci.Connid, ci.RecvByte, runTime, ci.RecvByte/runTime)
			g.Sess.Remove(ci.Connid)
			// Close Other
			if ci.Class == "D" {
				// SessionKey-Class-Did-UserID
				OtherConnid = fmt.Sprintf("%s-U-%s-%d", ci.SessionKey, ci.Did, ci.UserID)
			} else {
				OtherConnid = fmt.Sprintf("%s-D-%s-%d", ci.SessionKey, ci.Did, ci.UserID)
			}
			closePeerConn(OtherConnid)
		} else {
			logs.Logger.Infof("%s, Connid: %s , conn not match, not remove.", c, ci.Connid)
		}
	}
}

func closePeerConn(connid string) {
	v, e := g.Sess.Get(connid)
	if e == nil {
		v.(pprpc.RPCConn).Close()
		logs.Logger.Infof("Close Peer ConnID: %s .", connid)
	}
}
