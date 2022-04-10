package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"logagent/common"
	"logagent/etcd"
	"logagent/getSexAble"
	"logagent/influxDB"
	"logagent/kafka"
	"logagent/sendToinflux"
	"logagent/tailfile"
	"time"
)

//日志收集客户端
//类似GitHub中的filebeat
//往kafka发数据
//使用tail读取日志
type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
	EtcConfig     `ini:"etc"`
}
type KafkaConfig struct {
	Address  string `ini:"address"`
	Topic    string `ini:"topic"`
	ChanSize int64  `ini:"chan_size"`
}
type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}
type EtcConfig struct {
	Address    string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

//真正的业务逻辑
//	tailObj-->log-->client-->kafka

func main() {
	//	0： 获取本机IP为后续去etcd获取配置文件使用
	iP, err := common.GetOutIp()
	if err != nil {
		logrus.Errorf("local ip get failed err:%s", err)
		return
	}
	//1.读配置文件
	var configObj = new(Config)
	//cfg,err:= ini.Load("./conf/config.ini")
	//if err!=nil{
	//	logrus.Error("load config faild :%v",err)
	//	return
	//}
	//kafkaaddr:=cfg.Section("kafka").Key("address").String()
	//fmt.Println(kafkaaddr)
	err = ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		logrus.Errorf("load config faild :%v", err)
		return
	}

	//2.初始化  并且去一个管道中等着读取日志文件到kafka中
	err = kafka.Init([]string{configObj.KafkaConfig.Address}, configObj.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Errorf("init kafka failed :%v", err)
		return
	}
	logrus.Info("init kafka sucess !")
	//************************ 监视服务器的性能
	//1.初始化influx
	err = influxDB.InitInfluxDB()
	if err != nil {
		logrus.Errorf("influxDB init failed err:%s", err)
		return
	}
	//去做性能监控的管道获取数据并且发送到数据库
	go sendToinflux.SendMsgToDB()
	//2、获取当前系统的各种性能状态 发送到性能管道中
	go getSexAble.Run(time.Second)
	//*************************

	//初始化etcd链接
	err = etcd.Init([]string{configObj.EtcConfig.Address})
	if err != nil {
		logrus.Errorf("init etcd failed :%v", err)
		return
	}
	//从etcd拉数据  配置日志项
	collectKey := fmt.Sprintf(configObj.EtcConfig.CollectKey, iP)
	allConf, err := etcd.GetConf(collectKey)
	if err != nil {

		logrus.Errorf("get conf from etcd  failed :%v", err)
		return
	}
	//派小弟去监视etcd中configObj.EtcConfig.CollectKey的变化
	go etcd.WatchConf(collectKey)
	//3.使用tail读日志
	err = tailfile.Init(allConf)
	if err != nil {
		logrus.Errorf("init tail failed :%v", err)
		return
	}
	logrus.Info("init tail sucess !")
	//死循环等着那些后台的日志收集任务
	for {
		select {}
	}

}
