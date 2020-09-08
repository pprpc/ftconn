package app

import (
	g "xcthings.com/ftconn/relay-ms/common/global"
	"xcthings.com/protoc/ftconnrelay/RelayStepOne"
	"xcthings.com/protoc/ftconnrelay/RelayStepThree"
	"xcthings.com/protoc/ftconnrelay/RelayStepTwo"

	ctrl "xcthings.com/ftconn/relay-ms/controller"
)

func regService() {
	RelayStepOne.RegisterService(g.Service, &ctrl.RelayStepOneer{})
	RelayStepTwo.RegisterService(g.Service, &ctrl.RelayStepTwoer{})
	RelayStepThree.RegisterService(g.Service, &ctrl.RelayStepThreeer{})
}
