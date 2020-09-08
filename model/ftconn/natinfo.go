package ftconn

import (
	"fmt"

	errc "github.com/pprpc/ftconn/common/errorcode"
)

// Natinfo table name natinfo struct define
type Natinfo struct {
	ID            uint32 `json:"_" xorm:"id"`
	Did           string `json:"did" xorm:"did"`
	UserID        int64  `json:"user_id" xorm:"user_id"`
	Upnp          int32  `json:"upnp" xorm:"upnp"`
	UpnpProtocol  int32  `json:"upnp_protocol"`
	UpnpPort      int32  `json:"upnp_pprt"`
	NatType       int32  `json:"nat_type"`
	LocalIP       string `json:"local_ip" xorm:"local_ip"`
	LocalPort     int32  `json:"local_port"`
	OutsideIpaddr string `json:"outside_ipaddr" xorm:"outside_ipaddr"`
	LocalSsid     string `json:"local_ssid"`
	LocalGateway  string `json:"local_gateway"`
	GatewayMac    string `json:"gateway_mac"`
	LastTime      int64  `json:"last_time"`
}

// Add add record
func (obj *Natinfo) Add() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Natinfo")
		return
	}
	_, err = Orm.Insert(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Get get record
func (obj *Natinfo) Get() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Natinfo")
		return
	}
	if obj.UserID == 0 && obj.Did == "" {
		code = errc.ParameterError
		err = fmt.Errorf("Incomplete parameters: UserID/Did Cannot be empty at the same time")
		return
	}
	var has bool
	has, err = Orm.Where("user_id = ?", obj.UserID).And("did = ?", obj.Did).NoAutoCondition().Get(obj)
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
func (obj *Natinfo) Update() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Natinfo")
		return
	}
	if obj.UserID == 0 && obj.Did == "" {
		code = errc.ParameterError
		err = fmt.Errorf("Incomplete parameters: UserID/Did Cannot be empty at the same time")
		return
	}

	_, err = Orm.Where("user_id = ?", obj.UserID).And("did = ?", obj.Did).NoAutoCondition().Update(obj)
	if err != nil {
		code = errc.DBERROR
	}
	return
}

// Delete delete record
func (obj *Natinfo) Delete() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Natinfo")
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
func (obj *Natinfo) Reset() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Natinfo")
		return
	}
	*obj = Natinfo{}
	return
}

// GetAdv get record
func (obj *Natinfo) GetAdv() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Natinfo")
		return
	}
	if obj.UserID == 0 && obj.Did == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:UserID / Did Cannot be empty at the same time")
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
		err = fmt.Errorf("The record of the search,UserID: %v, Device: %vdoes not exist", obj.UserID, obj.Did)
	}
	return
}

// Set add or update record
func (obj *Natinfo) Set() (code uint64, err error) {
	if obj == nil {
		code = errc.NOTINIT
		err = fmt.Errorf("No initialization structure： Natinfo")
		return
	}
	if obj.UserID == 0 && obj.Did == "" {
		code = errc.ParameterError
		err = fmt.Errorf("The parameter is incorrect:UserID/Did Cannot be empty at the same time")
		return
	}

	_t := new(Natinfo)
	_t.UserID = obj.UserID
	_t.Did = obj.Did
	code, err = _t.GetAdv()
	if code == errc.NOTEXISTRECORD {
		code, err = obj.Add()
	} else if code == 0 {
		code, err = obj.Update()
	}
	return
}
