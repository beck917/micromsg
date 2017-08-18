package main

import (
	"flag"
	"fmt"

	"application/libraries/helpers"
	"application/libraries/toml"

	"github.com/beck917/pillX/pillx"
)

var (
	conf = flag.String("conf", "./etc/config.toml", "toml config")
)

func main() {
	flag.Parse()
	fmt.Println(*conf)
	tomlConfig, err := toml.LoadTomlConfig(*conf)
	if err != nil {
		panic(err)
	}

	etcdClient, etcdError := helpers.EtcdDail(tomlConfig.Etcd)
	if etcdError != nil {
		panic(etcdError)
	}

	gateway := &pillx.GatewayWebsocket{
		InnerAddr:   fmt.Sprintf("%s:%d", tomlConfig.Pillx.GatewayInnerHost, tomlConfig.Pillx.GatewayInnerPort),
		OuterAddr:   fmt.Sprintf("%s:%d", tomlConfig.Pillx.GatewayOuterHost, tomlConfig.Pillx.GatewayOuterPort),
		GatewayName: fmt.Sprintf("%s%d", tomlConfig.Pillx.GatewayName, 1),
		WatchName:   tomlConfig.Pillx.WorkerName,
	}

	gateway.OuterProtocol = &pillx.WebSocketProtocol{}
	gateway.EtcdClient = etcdClient
	gateway.Init()
	<-(chan int)(nil)
}
