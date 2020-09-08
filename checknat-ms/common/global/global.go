package global

import (
	"xcthings.com/micro/pprpcpool"
	"xcthings.com/micro/svc"
	"xcthings.com/ppmq/ppmqcli"
	"xcthings.com/ppmq/ppmqcli/msg"
	"xcthings.com/pprpc"
)

// Conf .
var Conf svc.MSConfig
var PConf PrivateConf
var EtcdPoint, Region, Ethname, MSName, Wanip string
var SvcAgent *svc.Agent

var SvcWatcher *svc.Watcher

// MicrosConn micro service connections.
var MicrosConn *pprpcpool.MicroClientConn

var PCli *ppmqcli.PpmqCli
var PMsg *msg.PPMQMsg

var UDPSrv1, UdpSrv2 *pprpc.RPCUDPServer

// Service global service
var Service, AuthService *pprpc.Service

func init() {
	Service = pprpc.NewService()
	AuthService = pprpc.NewService()
}
