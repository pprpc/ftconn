package app

import (
	g "xcthings.com/ftconn/checknat-ms/common/global"
	"xcthings.com/protoc/ftconnnat/NatProbe"
	"xcthings.com/protoc/ftconnnat/NatTest1"
	"xcthings.com/protoc/ftconnnat/NatTest2"
	"xcthings.com/protoc/ftconnnat/ReportNat"

	ctrl "xcthings.com/ftconn/checknat-ms/controller"
)

func regService() {
	NatTest1.RegisterService(g.Service, &ctrl.NatTest1er{})
	NatTest2.RegisterService(g.Service, &ctrl.NatTest2er{})
	NatProbe.RegisterService(g.Service, &ctrl.NatProbeer{})
	ReportNat.RegisterService(g.Service, &ctrl.ReportNater{})
}
