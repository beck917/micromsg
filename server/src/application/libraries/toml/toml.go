package toml

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type TomlConfig struct {
	Etcd   DBConfig
	Pillx  PillConfig
	Redis0 DBConfig
	Mysql  DBConfig
	Log    ValueConfig
	//Elasticsearch DBConfig
	//DB       database `toml:"database"`
}

var GlobalTomlConfig TomlConfig

func init() {
	var err error
	//etc/config.toml
	GlobalTomlConfig, err = LoadTomlConfig("etc/config.toml")
	if err != nil {
		panic(err)
	}
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBname   string
}

type PlatformConfig struct {
	Url       string
	AppId     string
	AppSecret string
}

type ValueConfig struct {
	Value string
}

type PillConfig struct {
	GatewayOuterHost string
	GatewayOuterPort int
	GatewayInnerHost string
	GatewayInnerPort int
	WorkerInnerHost  string
	WorkerInnerPort  int
	GatewayName      string
	WorkerName       string
}

func LoadTomlConfig(filename string) (tc TomlConfig, err error) {
	tomlData, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Read failed", err)
		return
	}

	if _, err1 := toml.Decode(string(tomlData), &tc); err1 != nil {
		err = err1
		fmt.Println("ReadToml failed", err)
		return
	}

	GlobalTomlConfig = tc
	return
}
