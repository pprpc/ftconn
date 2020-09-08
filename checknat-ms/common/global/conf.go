package global

import (
	"encoding/json"
	"fmt"

	"xcthings.com/hjyz/common"
	"xcthings.com/hjyz/logs"
	"xcthings.com/micro/svc"
)

// MicroClient micrl service client
type MicroClient struct {
	Name string   `json:"name,omitempty"`
	URIS []string `json:"uris,omitempty"`
}

// NatConf .
type NatConf struct {
	Port1   int32  `json:"port1"`
	Port2   int32  `json:"port2"`
	Ipaddr2 string `json:"ipaddr2"`
}

type PrivateConf struct {
	MaxSession int32         `json:"max_session"`
	Encrypt    uint8         `json:"encrypt,omitempty"`
	Nat        NatConf       `json:"nat"`
	Micros     []MicroClient `json:"micros,omitempty"` // authuser, authdevice
}

// LoadConf load conf
func LoadConf(filePath string) (conf svc.MSConfig, err error) {
	conf, err = loadConfFromETCD()
	if err != nil {
		logs.Logger.Warnf("LoadConf, loadConfFromETCD(), %s.", err)
		conf, err = loadConfFromFile(filePath)
	}
	return
}

func loadConfFromETCD() (conf svc.MSConfig, err error) {
	var ep []string
	ep = append(ep, EtcdPoint)
	var ag *svc.Agent
	var cfg *svc.Config

	ag, err = svc.NewAgent(svc.ValueRegService{}, 5, ep)
	if err != nil {
		err = fmt.Errorf("svc.NewAgent(), %s", err)
		return
	}
	defer ag.Close()

	ips, e := common.GetIPAddrByName(Ethname)
	if e != nil {
		err = fmt.Errorf("common.GetIPAddrByName(%s), %s", Ethname, e)
		return
	}

	cfg, err = svc.NewConfig(ag, Region, ips[0], MSName, []string{"ftconn"}, true)
	if err != nil {
		err = fmt.Errorf("svc.NewConfig(), %s", err)
		return
	}

	err = cfg.GetAll()
	if err != nil {
		err = fmt.Errorf("cfg.GetAll(), %s", err)
		return
	}
	conf = *cfg.Conf

	err = json.Unmarshal(conf.PrivateConfig, &PConf)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal(), %s", err)
		return
	}

	Wanip, err = cfg.GetWanIP(ips[0])
	//
	return
}

func loadConfFromFile(filePath string) (conf svc.MSConfig, err error) {
	if common.PathIsExist(filePath) != true {
		err = fmt.Errorf("conf file not exist")
		return
	}
	var buf []byte
	buf, err = common.LoadFileToByte(filePath)
	if err != nil {
		return
	}
	err = json.Unmarshal(buf, &conf)
	err = json.Unmarshal(conf.PrivateConfig, &PConf)
	return
}
