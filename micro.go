package main

import (
	"ctp/config"
	"ctp/databases"
	"ctp/routers"
	"fmt"
	"github.com/go-micro/plugins/v4/registry/etcd"
	_ "github.com/go-sql-driver/mysql"
	"ctp/logger"
	//"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/web"
)

var (
	//etcdCaPath       string = "F:/gotest/src/ca.pem"
	etcdReg registry.Registry
)

func RegistryInit() {
	//tlsInfo := transport.TLSInfo{
	//	TrustedCAFile: etcdCaPath,
	//}
	//tlsConfig, _ := tlsInfo.ClientConfig()
	etcdReg = etcd.NewRegistry(
		registry.Addrs("192.168.72.135:2379"),
		//registry.TLSConfig(tlsConfig),
		//etcd.Auth("root", "123456"),
	)
}

func main() {
	config, err := config.ReadConfig()
	if err != nil{
		fmt.Println(err)
		return
	}
	logger.LoggerInit()
	log := logger.GetLogger()
	log.Println(config.MP)
	RegistryInit()
	msqlDb, err := databases.InitDbMysql(config.MP)
	if err != nil{
		log.Error(err)
		return
	}
	router, _ := routers.InitRouter(msqlDb)
	//micro.NewService(
	//	micro.RegisterHandler()
	//	)
	go logger.TestLog(log)

	service := web.NewService(

		web.Name("api.miki1.com.userserver"),
		web.Address(":8002"),
		web.Handler(router),
		web.Registry(etcdReg),
		web.Version("v1.0.1"),
	)
	// initialise flags
	service.Init()
	// start the service
	service.Run()

}
