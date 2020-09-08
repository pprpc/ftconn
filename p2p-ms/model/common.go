package model

import "xcthings.com/ftconn/model/ftconn"

// GetP2pinfoBySessionKey .
func GetP2pinfoBySessionKey(sk string) (srvPort, uPort int32, uIP string, err error) {
	q := new(ftconn.P2pinfo)
	q.SessionKey = sk
	_, err = q.GetBySessionKey()
	if err != nil {
		return
	}
	srvPort = q.P2psrvPort
	uPort = q.UserOutsidePort
	uIP = q.UserOutsideIP

	return
}
