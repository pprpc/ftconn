package ftconn

import (
	"fmt"

	errc "xcthings.com/ftconn/common/errorcode"
)

// Relayinfo table name relayinfo struct define
type Relayinfo struct {
	ID          uint32 `json:"_" xorm:"id"`
	RelayIpaddr string `json:"relay_ipaddr" xorm:"relay_ipaddr"`
	Did         string `json:"did" xorm:"did"`
	UserID      int64  `json:"user_id" xorm:"user_id"`
	SessionKey  string `json:"session_key"`
	UserTime    int64  `json:"user_time"`
	DeviceTime  int64  `json:"device_time"`
}

// Add add record
func (obj *Relayinfo) Add() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Relayinfo")
		return
	}
	_, err = Orm.Insert(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Get get record
func (obj *Relayinfo) Get() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Relayinfo")
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
func (obj *Relayinfo) Update() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Relayinfo")
		return
	}
	if obj.UserID == 0 || obj.Did == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:UserID/Did Can not be empty")
		return
	}

	_, err = Orm.Where("user_id = ?", obj.UserID).And("did = ?", obj.Did).AllCols().NoAutoCondition().Update(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Delete delete record
func (obj *Relayinfo) Delete() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Relayinfo")
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
func (obj *Relayinfo) Reset() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Relayinfo")
		return
	}
	*obj = Relayinfo{}
	return
}

// GetAdv get record
func (obj *Relayinfo) GetAdv() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Relayinfo")
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
func (obj *Relayinfo) Set() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Relayinfo")
		return
	}
	if obj.UserID == 0 || obj.Did == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:UserID/Did Can not be empty")
		return
	}

	_t := new(Relayinfo)
	_t.UserID = obj.UserID
	_t.Did = obj.Did
	code, err = _t.GetAdv()
	if code == errc.NOTEXISTRECORD {
		code, err = obj.Add()
		if err != nil {
			err = fmt.Errorf("Add, %s", err)
		}
	} else if code == 0 {
		obj.ID = _t.ID
		code, err = obj.Update()
		if err != nil {
			err = fmt.Errorf("Update, %s", err)
		}
	}
	return
}

// GetBySessionKey get record
func (obj *Relayinfo) GetBySessionKey() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Relayinfo")
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
func (obj *Relayinfo) UpdateBySessionKey() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Relayinfo")
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
