package app

import (
	"fmt"

	g "github.com/pprpc/ftconn/relay-ms/common/global"
	"github.com/pprpc/util/logs"
	mqcli "xcthings.com/ppmq/ppmqcli"
	"xcthings.com/ppmq/ppmqcli/msg"
	"xcthings.com/ppmq/protoc/ppmqd/PPMQPublish"
	"github.com/pprpc/core/packets"
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
	g.PMsg = msg.New(g.PCli)

	return
}

func recivrPublish(pkg *packets.CmdPacket, req *PPMQPublish.Req) {
	logs.Logger.Warnf("MsgID: %s, Topic: %s , Cmdid: %d.", req.MsgId, req.Topic, req.Cmdid)
}
