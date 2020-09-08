package common

import (
	"context"

	g "github.com/pprpc/ftconn/p2p-ms/common/global"
	"github.com/pprpc/util/logs"
	"xcthings.com/protoc/authdevice/CheckDevice"
	"xcthings.com/protoc/authuser/CheckUser"
	"xcthings.com/protoc/authuser/CheckUserDeviceAcl"
)

// AuthDevice .
func AuthDevice(did, signkey string) bool {
	req := new(CheckDevice.Req)
	req.Did = did
	req.DidSignkey = signkey
	pkg, resp, err := g.MicrosConn.Invoke(context.Background(), CheckDevice.Module, CheckDevice.CmdID, req)
	if err != nil {
		logs.Logger.Errorf("g.MicrosConn.Invoke(CheckDevice), %s.", err)
		return false
	}
	if pkg.Code != 0 {
		logs.Logger.Errorf("AuthDevice, pkg.Code: %d.", pkg.Code)
		return false
	}
	if resp.(*CheckDevice.Resp).Code != 0 {
		logs.Logger.Errorf("AuthDevice, resp.Code: %d.", resp.(*CheckDevice.Resp).Code)
		return false
	}
	logs.Logger.Debugf("AuthDevice, Did: %s, SignKey: %s OK.", did, signkey)
	return true
}

// AuthUserID .
func AuthUserID(uid int64, pass string) bool {
	req := new(CheckUser.Req)
	req.UserId = uid
	req.Password = pass
	req.CountryCode = ""
	req.AccessKey = ""
	pkg, resp, err := g.MicrosConn.Invoke(context.Background(), CheckUser.Module, CheckUser.CmdID, req)
	if err != nil {
		logs.Logger.Errorf("AuthUserID, g.MicrosConn.Invoke(CheckUser), %s.", err)
		return false
	}
	if pkg.Code != 0 {
		logs.Logger.Errorf("AuthUserID, pkg.Code: %d.", pkg.Code)
		return false
	}
	if resp.(*CheckUser.Resp).Code != 0 {
		logs.Logger.Errorf("AuthUserID, resp.Code: %d.", resp.(*CheckUser.Resp).Code)
		return false
	}
	logs.Logger.Debugf("AuthUserID, Uid: %d, Pass: %s OK.", uid, pass)
	return true
}

// ACLUserIDAndDID Whether the user has access to the device
func ACLUserIDAndDID(uid int64, did string) bool {
	req := new(CheckUserDeviceAcl.Req)
	req.UserId = uid
	req.Did = did

	pkg, resp, err := g.MicrosConn.Invoke(context.Background(), CheckUserDeviceAcl.Module, CheckUserDeviceAcl.CmdID, req)
	if err != nil {
		logs.Logger.Errorf("ACLUserIDAndDID, g.MicrosConn.Invoke(CheckUserDeviceAcl), %s.", err)
		return false
	}
	if pkg.Code != 0 {
		logs.Logger.Errorf("ACLUserIDAndDID, pkg.Code: %d.", pkg.Code)
		return false
	}
	if resp.(*CheckUserDeviceAcl.Resp).Code == 0 {
		logs.Logger.Errorf("ACLUserIDAndDID, resp.Code: %d, not acl.", resp.(*CheckUserDeviceAcl.Resp).Code)
		return false
	}
	logs.Logger.Debugf("ACLUserIDAndDID, Uid: %d, Did: %s ,Code: %d, OK.", uid, did, resp.(*CheckUserDeviceAcl.Resp).Code)
	return true
}
