package app

import (
	"fmt"

	g "xcthings.com/ftconn/checknat-ms/common/global"
	ctrl "xcthings.com/ftconn/checknat-ms/controller"
	"xcthings.com/hjyz/logs"
	mqcli "xcthings.com/ppmq/ppmqcli"
	"xcthings.com/ppmq/ppmqcli/msg"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPublish"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQSubscribe"
	"xcthings.com/pprpc/packets"
	"xcthings.com/protoc/ftconnnat/ProbeConfig"
)

func initPPMQCli() (err error) {
	g.PCli, err = mqcli.NewPpmqcli(g.Conf.Ppmqclis[0].URL, g.Conf.Ppmqclis[0].Account,
		g.Conf.Ppmqclis[0].Password, g.Conf.Ppmqclis[0].HWFeature)
	if err != nil {
		err = fmt.Errorf("ppmqcli.NewPpmqcli(), %s", err)
		return
	}
	g.PCli.SetRecivePublishCB(recivrPublish)
	err = g.PCli.Dail()
	if err != nil {
		err = fmt.Errorf("ppmqcli.Dail(), %s", err)
		return
	}
	subscribe()
	g.PMsg = msg.New(g.PCli)

	return
}

func subscribe() {
	req := new(PPMQSubscribe.Req)

	s := new(PPMQSubscribe.TopicInfo)
	//s.Topic = g.Conf.Ppmqclis[0].TopicPrefix + g.Conf.Public.ServerID
	s.Topic = g.Conf.Ppmqclis[0].TopicPrefix + g.Wanip
	s.Qos = 1
	s.Cluster = 2
	s.ClusterSubid = ""

	req.Topics = append(req.Topics, s)

	_, err := g.PCli.Subscribe(req)
	if err != nil {
		logs.Logger.Errorf("ppmqcli.Subscribe(), %s.", err)
		return
	}

}

func recivrPublish(pkg *packets.CmdPacket, req *PPMQPublish.Req) {
	logs.Logger.Debugf("MsgID: %s, Topic: %s , Cmdid: %d, CmdType: %d, Format: %d.",
		req.MsgId, req.Topic, req.Cmdid, req.CmdType, req.Format)

	if req.Cmdid != ProbeConfig.CmdID {
		logs.Logger.Debugf("not support cmdid: %d.", req.Cmdid)
		return
	}
	if req.CmdType != int32(packets.RPCREQ) {
		logs.Logger.Debugf("not support cmdtype: %d.", req.CmdType)
		return
	}
	if req.Format != int32(packets.TYPEPBBIN) && req.Format != int32(packets.TYPEPBJSON) {
		logs.Logger.Debugf("not support format: %d.", req.Format)
		return
	}
	ctrl.SendNatProbe(req)
}
