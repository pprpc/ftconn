package model

import (
	"fmt"

	"github.com/pprpc/ftconn/model/ftconn"
	"xcthings.com/micro/svc"
)

//InitEngine  .
func InitEngine(cfg svc.MSConfig) (err error) {
	for _, row := range cfg.Dbs {
		switch row.ConfName {
		case "ftconn":
			err = ftconn.InitModelEngine(row)
			if err != nil {
				return
			}
		default:
			err = fmt.Errorf("InitEngine, not support conf_name: %s", row.ConfName)
			return
		}

	}
	return
}
