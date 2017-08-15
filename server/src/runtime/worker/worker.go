package main

import (
	"flag"
	"fmt"

	"application/controllers"
	"application/libraries/helpers"
	"application/libraries/logger"
	"application/libraries/opcodes"
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

	etcdClient := helpers.EtcdDail(tomlConfig.Etcd)

	worker := &pillx.Worker{
		InnerAddr:  fmt.Sprintf("%s:%d", tomlConfig.Pillx.WorkerInnerHost, tomlConfig.Pillx.WorkerInnerPort),
		WorkerName: fmt.Sprintf("%s_%s:%d", tomlConfig.Pillx.WorkerName, tomlConfig.Pillx.WorkerInnerHost, tomlConfig.Pillx.WorkerInnerPort),
		WatchName:  tomlConfig.Pillx.GatewayName,
	}
	helpers.GlobalWorker = worker
	worker.EtcdClient = etcdClient

	logger.InitLog("worker")

	worker.Init()
	worker.InnerServer.HandleFunc(opcodes.APP_LOGIN, controllers.LoginHandler)
	worker.InnerServer.HandleFunc(opcodes.APP_SEND, controllers.SendHandler)
	worker.InnerServer.HandleFunc(opcodes.APP_ADD, controllers.AddHandler)
	worker.InnerServer.HandleFunc(opcodes.APP_OPEN, controllers.OpenHandler)
	worker.InnerServer.HandleFunc(opcodes.APP_DELETE, controllers.DeleteHandler)
	worker.InnerServer.HandleFunc(opcodes.APP_DELETE_MSG, controllers.DeleteMsgHandler)
	worker.InnerServer.HandleFunc(opcodes.APP_REGISTER, controllers.RegisterHandler)
	worker.InnerServer.HandleFunc(pillx.SYS_CLIENT_DISCONNECT, controllers.OnClientCloseHandler)
	worker.InnerServer.HandleFunc(pillx.SYS_CLIENT_DISCONNECT_WORKER, controllers.OnClientCloseWorkerHandler)
	worker.InnerServer.HandleFunc(pillx.SYS_ON_MESSAGE, controllers.OnMessageHandler)
	<-(chan int)(nil)
	//worker.Watch()
}
