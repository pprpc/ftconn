package app

import (
	"crypto/tls"
	"fmt"
	"net/url"
	"strconv"

	g "xcthings.com/ftconn/checknat-ms/common/global"
	"github.com/pprpc/util/logs"
	"xcthings.com/micro/svc"
	"github.com/pprpc/core"
	"github.com/pprpc/core/packets"
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
	if g.UDPSrv1 == nil || g.UdpSrv2 == nil {
		err = fmt.Errorf("Must start two UDP listening ports")
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
		usrv.SetReadTimeout(lis.ReadTimeout)
		logs.Logger.Infof("Listen UDPServer: %s.", lis.URI)
		go usrv.Serve()
		if g.UDPSrv1 == nil {
			g.UDPSrv1 = usrv
		} else if g.UdpSrv2 == nil {
			g.UdpSrv2 = usrv
		}
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
