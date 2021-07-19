package main

import (
	"fmt"
	"github.com/wsxzhmzwdfxeds/kafka2es/conf"
	"github.com/wsxzhmzwdfxeds/kafka2es/utils"
	"gopkg.in/ini.v1"
	"os"
)

func main() {
	// 0. load config file
	var cfg conf.LogTransferCfg
	err := ini.MapTo(&cfg, "cfg.ini")
	if err != nil {
		panic(err)
	}
	fmt.Printf("cfg:%v\n", cfg)

	// 0. init etcd
	err = utils.InitEtcd(cfg.EtcdCfg.Address, cfg.EtcdCfg.Timeout)
	if err != nil {
		fmt.Printf("init etcd faild, err : %v\n", err)
		os.Exit(1)
	}

	// 1. initialize es
	err = utils.InitEs(cfg.EsCfg.Address)
	if err != nil {
		fmt.Printf("init es failed, err: %v\n", err)
		os.Exit(1)
	}

	// 2. initialize kafka
	utils.InitKfk(cfg.KafkaCfg)
	if err != nil {
		fmt.Printf("init kafka consumer failed, err: %v\n", err)
		return
	}
}
