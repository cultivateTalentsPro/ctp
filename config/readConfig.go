package config

import (
	"ctp/databases"
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
)

type YycConfig struct {
	MP databases.MysqlParam `ini:mysql`
}
func ReadConfig() (YycConfig, error){
	var configDir string
	var yycConfig YycConfig
	flag.StringVar(&configDir, "c", "", "配置文件路径")
	flag.Parse()
	fmt.Println(configDir)
	configDir = "E:/goyyc/src/ctp/config/yyc.ini"
	cfgs, err := ini.Load(configDir)
	if err != nil {
		fmt.Println(err)
		return yycConfig, err
	}
	//tt := cfgs.Section("mysql").Key("driver_name").Value()
	//fmt.Println(tt)
	err = cfgs.Section("mysql").MapTo(&yycConfig.MP)
	if err != nil {
		fmt.Println(err)
		return yycConfig, err
	}
	//fmt.Println(yycConfig.MP)

	return yycConfig, nil
}