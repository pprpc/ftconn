package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	g "xcthings.com/ftconn/relay-ms/common/global"
	"github.com/pprpc/util/logs"
	"xcthings.com/micro/pprpcpool"
	"xcthings.com/micro/svc"
	"xcthings.com/protoc/authdevice/CheckDevice"
	"xcthings.com/protoc/authuser/CheckUser"
	"xcthings.com/protoc/authuser/CheckUserDeviceAcl"
)

func initAuthConns() (err error) {
	authRegService()
	err = initMicroClientConn()
	if err != nil {
		return
	}
	err = etcdWatcher()
	if err != nil {
		return
	}
	// load
	loadRegister()
	return
}

func authRegService() {
	CheckDevice.RegisterService(g.AuthService, nil)
	CheckUser.RegisterService(g.AuthService, nil)
	CheckUserDeviceAcl.RegisterService(g.AuthService, nil)
}

func initMicroClientConn() (err error) {
	g.MicrosConn = pprpcpool.NewMicroClientConn(g.AuthService)
	for _, v := range g.PConf.Micros {
		err = g.MicrosConn.AddMicro(v.Name)
		if err != nil {
			err = fmt.Errorf("g.MicrosConn.AddMicro(%s), %s", v.Name, err)
			return
		}
	}
	return
}

func loadRegister() {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	key := fmt.Sprintf("/register/%s/", g.Region)
	kvs, err := g.SvcWatcher.GetValues(ctx, key)
	if err != nil {
		logs.Logger.Errorf("g.SvcAgent.GetValues(%s), %s.", key, err)
		return
	}

	for _, row := range kvs {
		_t := new(svc.ValueRegService)
		err := json.Unmarshal([]byte(row.Value), _t)
		if err != nil {
			logs.Logger.Errorf("json.Unmarshal(), %s.", err)
			continue
		}
		err = g.MicrosConn.AddHost(key, *_t)
		if err != nil {
			logs.Logger.Errorf("g.MicroConn.AddHost(%s), %s.", key, err)
		}
	}
}

func etcdWatcher() (err error) {
	var ep []string
	ep = append(ep, g.EtcdPoint)
	g.SvcWatcher, err = svc.NewWatcher("/register", ep, regcb)
	if err != nil {
		err = fmt.Errorf("svc.NewWatcher(/register), %s", err)
		return
	}
	go g.SvcWatcher.Start()

	return
}

func regcb(action, key, value string) {
	if action == "PUT" {
		_t := new(svc.ValueRegService)
		err := json.Unmarshal([]byte(value), _t)
		if err != nil {
			logs.Logger.Errorf("json.Unmarshal(), %s.", err)
			return
		}
		err = g.MicrosConn.AddHost(key, *_t)
		if err != nil {
			logs.Logger.Errorf("g.MicroConn.AddHost(%s), %s.", key, err)
		}
	} else {
		err := g.MicrosConn.DelHost(key)
		if err != nil {
			logs.Logger.Errorf("g.MicroConn.DelHost(%s), %s.", key, err)
		}
	}
}
