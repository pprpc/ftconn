package ftconn

import (
	"fmt"

	errc "xcthings.com/ftconn/common/errorcode"
)

// P2pinfo table name p2pinfo struct define
type P2pinfo struct {
	ID                uint32 `json:"_" xorm:"id"`
	Did               string `json:"did" xorm:"did"`
	UserID            int64  `json:"user_id" xorm:"user_id"`
	DeviceOutsideIP   string `json:"device_outside_ip" xorm:"device_outside_ip"`
	DeviceOutsidePort int32  `json:"device_outside_port" xorm:"device_outside_port"`
	DeviceLocalIP     string `json:"device_local_ip" xorm:"device_local_ip"`
	DeviceLocalPort   int32  `json:"device_local_port" xorm:"device_local_port"`
	UserOutsideIP     string `json:"user_outside_ip" xorm:"user_outside_ip"`
	UserOutsidePort   int32  `json:"user_outside_port" xorm:"user_outside_port"`
	UserLocalIP       string `json:"user_local_ip" xorm:"user_local_ip"`
	UserLocalPort     int32  `json:"user_local_port" xorm:"user_local_port"`
	SessionKey        string `json:"session_key"`
	P2psrvIP          string `json:"p2psrv_ip" xorm:"p2psrv_ip"`
	P2psrvPort        int32  `json:"p2psrv_port"`
	UserTime          int64  `json:"user_time"`
	DeviceTime        int64  `json:"device_time"`
}

// Add add record
func (obj *P2pinfo) Add() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： P2pinfo")
		return
	}
	_, err = Orm.Insert(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Get get record
func (obj *P2pinfo) Get() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： P2pinfo")
		return
	}
	if obj.UserID == 0 {
		code = errc.ParameterError
		err = fmt.Errorf("Incomplete parameters: UserID Can not be empty")
		return
	}
	var has bool
	has, err = Orm.Where("user_id = ?", obj.UserID).NoAutoCondition().Get(obj)
	if err != nil {
		code = errc.DBERROR
		return
	}
	if has != true {
		code = errc.NOTEXISTRECORD
		err = fmt.Errorf("The record of the search,ClientID: %vdoes not exist", obj.UserID)
	}
	return
}

// Update update record
func (obj *P2pinfo) Update() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： P2pinfo")
		return
	}
	if obj.UserID == 0 {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:UserID Can not be empty")
		return
	}

	_, err = Orm.Where("user_id = ?", obj.UserID).And("did = ?", obj.Did).AllCols().NoAutoCondition().Update(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Delete delete record
func (obj *P2pinfo) Delete() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： P2pinfo")
	}
	if obj.UserID == 0 {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:UserID Can not be empty")
		return
	}
	/*
		up := new(MsgConnect)
		up.IsDel = 2
		_, err = Orm.Where("code = ?", obj.Code).NoAutoCondition().Update(up)
		if err != nil {
			code = errc.CodeDB
		}
	*/
	_, err = Orm.Where("user_id = ?", obj.UserID).NoAutoCondition().Delete(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Reset reset struct
func (obj *P2pinfo) Reset() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： P2pinfo")
		return
	}
	*obj = P2pinfo{}
	return
}

// GetAdv get record
func (obj *P2pinfo) GetAdv() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： P2pinfo")
		return
	}
	if obj.UserID == 0 {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:UserID Can not be empty")
		return
	}
	var has bool
	has, err = Orm.Get(obj)
	if err != nil {
		code = errc.DBERROR
		return
	}
	if has != true {
		code = errc.NOTEXISTRECORD
		err = fmt.Errorf("The record of the search,UserID: %vdoes not exist", obj.UserID)
	}
	return
}

// Set add or update record
func (obj *P2pinfo) Set() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： P2pinfo")
		return
	}
	if obj.UserID == 0 {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:UserID Can not be empty")
		return
	}

	_t := new(P2pinfo)
	_t.UserID = obj.UserID
	_t.Did = obj.Did
	//_t.P2psrvIP = obj.P2psrvIP
	code, err = _t.GetAdv()
	if code == errc.NOTEXISTRECORD {
		code, err = obj.Add()
	} else if code == 0 {
		obj.ID = _t.ID
		code, err = obj.Update()
	}
	return
}

// GetBySessionKey get record
func (obj *P2pinfo) GetBySessionKey() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： P2pinfo")
		return
	}
	if obj.SessionKey == "" {
		code = errc.ParameterError
		err = fmt.Errorf("Incomplete parameters: SessionKey Can not be empty")
		return
	}
	var has bool
	has, err = Orm.Where("session_key = ?", obj.SessionKey).NoAutoCondition().Get(obj)
	if err != nil {
		code = errc.DBERROR
		return
	}
	if has != true {
		code = errc.NOTEXISTRECORD
		err = fmt.Errorf("The record of the search,SessionKey: %vdoes not exist", obj.SessionKey)
	}
	return
}

// UpdateBySessionKey update record
func (obj *P2pinfo) UpdateBySessionKey() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： P2pinfo")
		return
	}
	if obj.SessionKey == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:SessionKey Can not be empty")
		return
	}

	_, err = Orm.Where("session_key = ?", obj.SessionKey).NoAutoCondition().Update(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

//GetP2pInfo get p2p info by did and uid.
func GetP2pInfo(did string, uid int64) (row *P2pinfo, err error) {
	var has bool
	row = new(P2pinfo)
	has, err = Orm.Where("did = ?", did).And("user_id = ?", uid).NoAutoCondition().Get(row)
	if err != nil {
		return
	}
	if has != true {
		err = fmt.Errorf("The record of the search,did: %v, uid: %vdoes not exist", did, uid)
	}
	return
	// if common.GetTimeMs() - row.UserTime > 3000  {
	// 	err = fmt.Errorf("Timeout")
	// }
}
