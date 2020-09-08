package errorcode

const (
	// CmdIDNotReg cmdid not register
	CmdIDNotReg uint64 = 1
	// INVOKEError invoke error
	INVOKEError uint64 = 2
	// UNKNOWNERROR unknown error
	UNKNOWNERROR uint64 = 99
	// NOTINIT not init obj
	NOTINIT uint64 = 100
	// DBERROR db error
	DBERROR uint64 = 101
	// ParameterError  .
	ParameterError uint64 = 102
	// NOTEXISTRECORD .
	NOTEXISTRECORD uint64 = 103
	// ParameterIllegal .
	ParameterIllegal uint64 = 104
	// MarshalError .
	MarshalError uint64 = 105
	// ACLDENY .
	ACLDENY uint64 = 106

	// ACCESSTOKENExpired access_token expired
	ACCESSTOKENExpired uint64 = 107
	// REFRESHTOKENExpired refresh_token expired
	REFRESHTOKENExpired uint64 = 108

	// ACCESSKEYNOTEXIST AccessKey not exist.
	ACCESSKEYNOTEXIST uint64 = 120
	// UIDNOTEXIST user id not exist.
	UIDNOTEXIST uint64 = 121

	// PUBLICAUTHDENY auth deny
	PUBLICAUTHDENY uint64 = 150

	// ACCOUNTNOTMATCH account not match
	ACCOUNTNOTMATCH uint64 = 201
	// NOTRUNPPMQCONNECT  NOTRUNPPMQCONNECT
	NOTRUNPPMQCONNECT uint64 = 202

	// ACCESSDEVICEDENY .
	ACCESSDEVICEDENY uint64 = 250
)

// authdevice-ms

// authuser-ms

// ftconnp2p
const (
	CONNSESSIONKEYNOTMATCH uint64 = 500
	SESSIONKEYTIMEOUT      uint64 = 501
)

// user-ms
const (
	// SENDVCODETOOOFTEN  .
	SENDVCODETOOOFTEN uint64 = 801
	// VCODENOTEXIST vcode .
	VCODENOTEXIST uint64 = 802
	// VCODENOTMATCH  vcode not match
	VCODENOTMATCH uint64 = 803
	// VCODEEXPIRED.
	VCODEEXPIRED uint64 = 804
	// USEREXIST .
	USEREXIST uint64 = 805
	//
	USERLOCK    uint64 = 806
	USERDISABLE uint64 = 807
	// USERPASSWORDError .
	USERPASSWORDError uint64 = 808
	//  USERMOBILEEXIST .
	USERMOBILEEXIST uint64 = 809
	// USERHEADPICSAVEError .
	USERHEADPICSAVEError uint64 = 810
	// REFRESHTOKENUSERIDNOTMATCH .
	REFRESHTOKENUSERIDNOTMATCH uint64 = 811
	// ACCESSTOKENNOTEXIST access_token not exist
	ACCESSTOKENNOTEXIST uint64 = 812
	// ACCESSTOKENLOGOUT access_token logout, not match
	ACCESSTOKENLOGOUT uint64 = 813
	// OLDMOBILENOTMATCH old mobile not match
	OLDMOBILENOTMATCH uint64 = 815
	// USERNOTEXIST .
	USERNOTEXIST uint64 = 816
)

// userdevice-ms (900-1000)
const (
	// DeviceAdded .
	DeviceAdded uint64 = 901
	// DeviceNotBelong .
	DeviceNotBelong uint64 = 902
	// SHAREEXIST .
	SHAREEXIST uint64 = 903
)

// clouds-ms/ cloudsindex-ms
const (
	// DeviceUnbind  .
	DeviceUnbind uint64 = 1001
	// UnsupportRecTyp  .
	UnsupportRecTyp uint64 = 1002
	// NoDCSConfig .
	NoDCSConfig uint64 = 1003
	// OssMaxFidxIll
	OssMaxFidxIll uint64 = 1004
	// AssumeRoleFault
	AssumeRoleFault uint64 = 1005
	// DCSConfigExpired dev clouds service has expired
	DCSConfigExpired uint64 = 1006
)
