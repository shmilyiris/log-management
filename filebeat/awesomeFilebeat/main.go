package main

import (
	"awesomeFilebeat/conf"
	"awesomeFilebeat/etcd"
	"awesomeFilebeat/tools"
	"awesomeFilebeat/watchlog"
	"fmt"
	"gopkg.in/ini.v1"
	"sync"
	"time"
)

var config = new(conf.Config)

func main() {
	// 0. 加载配置文件
	err := ini.MapTo(config, "./conf/config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		return
	}

	// 1. 初始化etcd
	err = etcd.Init(config.Etcd.Address, time.Duration(config.Etcd.Timeout)*time.Second)
	if err != nil {
		fmt.Printf("init etcd failed, err: %v\n", err)
		return
	}
	fmt.Println("Init etcd success.")

	// 1.1 获取对外 IP 地址
	ip, err := tools.GetOurboundIP()
	if err != nil {
		panic(err)
	}
	etcdConfKey := fmt.Sprintf(config.Etcd.Key, ip)

	// 1.2 启动 watchlog 用于监控指定的 filebeat.yml 文件
	watchlog.Init("usr/share/filebeat/filebeat.yml")

	// 2. 监控 etcd 有关filebeat的日志收集配置
	var wg sync.WaitGroup
	wg.Add(1)
	go etcd.WatchConf(etcdConfKey, watchlog.NewConfChan()) // 哨兵发现最新的配置信息会通知上面的通道
	wg.Wait()
}
