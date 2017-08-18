package helpers

import (
	"application/libraries/toml"
	"fmt"
	"net/url"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
)

var dbMu sync.RWMutex
var _instanceDB map[string]*Database

const TablePrefix string = ""

func init() {
	_instanceDB = make(map[string]*Database)
}

type Database struct {
	XORM   *xorm.Engine
	Config toml.DBConfig
}

func InstanceDatabase(mysqlconfig toml.DBConfig) (*Database, error) {
	dbMu.Lock()
	defer dbMu.Unlock()

	name := mysqlconfig.DBname
	var err error
	if _instanceDB[name] == nil {
		_instanceDB[name] = new(Database)
		_instanceDB[name].Config = mysqlconfig
		_instanceDB[name].XORM, err = _instanceDB[name].MysqlDail()
	}

	return _instanceDB[name], err
}

func (database *Database) MysqlDail() (*xorm.Engine, error) {
	mysqlconfig := database.Config

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=", mysqlconfig.User, mysqlconfig.Password, mysqlconfig.Host, mysqlconfig.Port, mysqlconfig.DBname)

	dsn = dsn + url.QueryEscape("Asia/Shanghai")
	db, err := xorm.NewEngine("mysql", dsn)

	tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, TablePrefix)
	db.SetTableMapper(tbMapper)

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(10)

	//tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "_")
	//db.SetTableMapper(tbMapper)
	db.TZLocation = time.Local

	return db, err
}
