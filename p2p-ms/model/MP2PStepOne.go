package model

import (
	"fmt"

	"github.com/pprpc/ftconn/model/ftconn"
	g "github.com/pprpc/ftconn/p2p-ms/common/global"
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/crypto"
	"github.com/pprpc/util/logs"
	"github.com/pprpc/core"
	"xcthings.com/protoc/ftconnp2p/P2PStepOne"
)

// MP2PStepOne P2PStepOne  
func MP2PStepOne(c pprpc.RPCConn, req *P2PStepOne.Req) (resp *P2PStepOne.Resp, code uint64, err error) {
	p2p := new(ftconn.P2pinfo)
	p2p.Did = req.RemoteDid
	p2p.UserID = req.UserId
	p2p.DeviceOutsideIP = ""
	p2p.DeviceOutsidePort = 0
	p2p.DeviceLocalIP = ""
	p2p.DeviceLocalPort = 0
	p2p.UserOutsideIP, p2p.UserOutsidePort = common.GetRemoteIPPort(c.RemoteAddr())
	if req.Nat != nil {
		p2p.UserLocalIP = req.Nat.LocalIp
		p2p.UserLocalPort = req.Nat.LocalPort
	} else {
		p2p.UserLocalIP = ""
		p2p.UserLocalPort = 0
	}
	if req.IsRetry == 0 {
		p2p.SessionKey = getSessionKey(req)
	} else {
		r, e := ftconn.GetP2pInfo(req.RemoteDid, req.UserId)
		if e != nil {
			logs.Logger.Warnf("ftconn.GetP2pInfo(%s,%d), %s.", req.RemoteDid, req.UserId, e)
			p2p.SessionKey = getSessionKey(req)
		} else {
			p2p.SessionKey = r.SessionKey
		}
	}
	p2p.P2psrvIP = g.PConf.P2p.WanIP
	_, p2p.P2psrvPort = common.GetRemoteIPPort(c.LocalAddr()) //g.PConf.P2p.WanPort
	p2p.UserTime = common.GetTimeMs()
	p2p.DeviceTime = 0

	code, err = p2p.Set()
	if err != nil {
		return
	}

	resp = new(P2PStepOne.Resp)
	resp.OutsideIpaddr = p2p.UserOutsideIP
	resp.OutsidePort = p2p.UserOutsidePort
	resp.SessionKey = p2p.SessionKey
	return
}

func getSessionKey(req *P2PStepOne.Req) string {
	return crypto.MD5([]byte(fmt.Sprintf("%s-%d-%d", req.RemoteDid, req.UserId, common.GetTimeNs())))
}
