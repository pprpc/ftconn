package global

import (
	"xcthings.com/micro/pprpcpool"
	"xcthings.com/micro/svc"
	"xcthings.com/ppmq/ppmqcli"
	"xcthings.com/ppmq/ppmqcli/msg"
	"github.com/pprpc/core"
	"github.com/pprpc/core/sess"
)

// ConnAttr
type ConnAttr struct {
	Connid        string // SessionKey-Class-Did-UserID
	Class         string // D Device; U User
	Did           string
	UserID        int64
	SessionKey    string
	CallPrehookCB bool
	InitTimeSec   int64
	RecvByte      int64
	SendByte      int64
}

// Conf .
var Conf svc.MSConfig
var PConf PrivateConf
var SvcAgent *svc.Agent

var SvcWatcher *svc.Watcher

// MicrosConn micro service connections.
var MicrosConn *pprpcpool.MicroClientConn

// Service global service
var Service, AuthService *pprpc.Service

var PCli *ppmqcli.PpmqCli
var PMsg *msg.PPMQMsg

var EtcdPoint, Region, Ethname, MSName string

// Sess all session
var Sess *sess.Sessions

func init() {
	Service = pprpc.NewService()
	AuthService = pprpc.NewService()
}
