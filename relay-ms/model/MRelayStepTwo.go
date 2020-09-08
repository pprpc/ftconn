package model

import (
	"github.com/pprpc/util/common"
	"xcthings.com/ftconn/model/ftconn"
	"xcthings.com/protoc/ftconnrelay/RelayStepTwo"
)

// MRelayStepTwo RelayStepTwo  
func MRelayStepTwo(req *RelayStepTwo.Req) (resp *RelayStepTwo.Resp, code uint64, err error) {

	r := new(ftconn.Relayinfo)
	r.DeviceTime = common.GetTimeMs()
	r.SessionKey = req.SessionKey

	code, err = r.UpdateBySessionKey()
	if err != nil {
		return
	}
	resp = new(RelayStepTwo.Resp)
	resp.SessionKey = req.SessionKey

	return
}
