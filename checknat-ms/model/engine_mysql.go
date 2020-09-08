package model

import (
	"fmt"

	"xcthings.com/ftconn/model/ftconn"
	"xcthings.com/micro/svc"
)

//InitEngine  db engine
func InitEngine(cfg svc.MSConfig) (err error) {
	for _, row := range cfg.Dbs {
		if row.Type != "mysql" {
			err = fmt.Errorf("InitEngine, not support type: %s", row.Type)
			return
		}
		switch row.ConfName {
		case "ftconn":
			err = ftconn.MySQL(row)
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
