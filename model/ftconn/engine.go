package ftconn

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xcthings.com/micro/svc"
)

var (
	Orm *xorm.Engine
)

// InitModelEngine init engine
func InitModelEngine(cfg svc.ValueDbconf) (err error) {
	if cfg.Type == "mysql" {
		err = MySQL(cfg)
	} else if cfg.Type == "sqlite3" {
		err = SQLite(cfg)
	} else {
		err = fmt.Errorf("InitEngine, not support type: %s", cfg.Type)
	}
	return
}

// MySQL  init mysql db
func MySQL(cfg svc.ValueDbconf) (err error) {

	var dbOption string

	if cfg.MaxIdle < 1 {
		cfg.MaxIdle = 5
	}
	if cfg.MaxConn < 1 {
		cfg.MaxConn = 5
	}

	if cfg.Host == "localhost" || cfg.Host == "127.0.0.1" {
		dbOption = fmt.Sprintf("%s:%s@unix(%s)/%s?charset=%s", cfg.User, cfg.Pass,
			cfg.Socket, cfg.Name, cfg.Charset)
	} else {
		dbOption = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", cfg.User, cfg.Pass,
			cfg.Host, cfg.Port, cfg.Name, cfg.Charset)
	}
	Orm, err = xorm.NewEngine("mysql", dbOption)
	if err != nil {
		err = fmt.Errorf("xorm.NewEngine(mysql, %s), error: %s", dbOption, err)
		return
	}
	if cfg.MaxIdle > 0 {
		Orm.SetMaxIdleConns(cfg.MaxIdle)
	}
	if cfg.MaxConn > 0 {
		Orm.SetMaxOpenConns(cfg.MaxConn)
	}
	if cfg.Debug == true {
		Orm.ShowSQL(cfg.Debug)
	}

	return
}

// SQLite sqlite engine
func SQLite(cfg svc.ValueDbconf) (err error) {
	Orm, err = xorm.NewEngine(cfg.Type, cfg.Name)
	if err != nil {
		err = fmt.Errorf("xorm.NewEngine(%s, %s), error: %s", cfg.Type, cfg.Name, err)
		return
	}
	if cfg.Debug == true {
		Orm.ShowSQL(cfg.Debug)
	}
	// init table
	// err = Orm.Sync2(new(Relayinfo))
	return
}
