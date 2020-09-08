package main

import (
	"flag"
	"fmt"

	"xcthings.com/ftconn/relay-ms/common/app"
	g "xcthings.com/ftconn/relay-ms/common/global"
	"github.com/pprpc/util/common"
	"github.com/pprpc/util/logs"
)

var (
	bDate       string //
	ciHash      string //
	mainVersion string = "0.0.9"
)

var (
	etcdipaddr = flag.String("ipaddr", "127.0.0.1:2379", "etcd server ipaddr")
	region     = flag.String("region", "cn-shenzhen", "region name")
	ethname    = flag.String("i", "eth0", "network device name")
	ver        = flag.Bool("v", false, "show version")
	confFile   = flag.String("conf", "../conf/server.json", "Specify configuration files")
)

func main() {
	flag.Parse()
	// show version
	if *ver {
		version := mainVersion
		if len(bDate) > 0 {
			version += ("+" + bDate)
		}
		fmt.Println("version:", version)

		if len(ciHash) > 0 {
			fmt.Println("git commit hash:", ciHash)
		}
		return
	}
	defer logs.Logger.Flush()

	g.EtcdPoint = *etcdipaddr
	g.Region = *region
	g.Ethname = *ethname
	g.MSName = "relay"

	err := app.Run(*confFile)
	if err != nil {
		logs.Logger.Errorf("app.Run(), error: %s.", err)
	}
	common.WaitCtrlC()
}
