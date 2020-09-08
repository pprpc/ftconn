package model

import (
	"fmt"

	"xcthings.com/ftconn/model/ftconn"
	"xcthings.com/micro/svc"
)

//InitEngine  db engine
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
