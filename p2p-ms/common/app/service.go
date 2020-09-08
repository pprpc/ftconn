package app

import (
	g "github.com/pprpc/ftconn/p2p-ms/common/global"
	"xcthings.com/protoc/ftconnp2p/P2PStepOne"
	"xcthings.com/protoc/ftconnp2p/P2PStepThree"
	"xcthings.com/protoc/ftconnp2p/P2PStepTwo"

	ctrl "github.com/pprpc/ftconn/p2p-ms/controller"
)

func regService() {
	P2PStepOne.RegisterService(g.Service, &ctrl.P2PStepOneer{})
	P2PStepTwo.RegisterService(g.Service, &ctrl.P2PStepTwoer{})
	P2PStepThree.RegisterService(g.Service, &ctrl.P2PStepThreeer{})
}
