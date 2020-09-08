package model

import (
	"fmt"

	"github.com/pprpc/util/common"
	"github.com/pprpc/util/crypto"
	g "xcthings.com/ftconn/relay-ms/common/global"
	"xcthings.com/ftconn/model/ftconn"
	"xcthings.com/protoc/ftconnrelay/RelayStepOne"
)

// MRelayStepOne RelayStepOne  
func MRelayStepOne(req *RelayStepOne.Req) (resp *RelayStepOne.Resp, code uint64, err error) {
	r := new(ftconn.Relayinfo)
	r.Did = req.RemoteDid
	r.UserID = req.UserId
	r.SessionKey = getSessionKey(req)
	r.UserTime = common.GetTimeMs()
	r.DeviceTime = 0
	r.RelayIpaddr = g.PConf.Relay.WanIP

	code, err = r.Set()
	if err != nil {
		return
	}
	resp = new(RelayStepOne.Resp)
	resp.SessionKey = r.SessionKey

	return
}

func getSessionKey(req *RelayStepOne.Req) string {
	return crypto.MD5([]byte(fmt.Sprintf("%s-%d-%d", req.RemoteDid, req.UserId, common.GetTimeNs())))
}
