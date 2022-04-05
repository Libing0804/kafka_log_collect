package main

import (
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/tailfile"
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
	Address string `ini:"address"`
	Topic string `ini:"topic"`
	ChanSize int64 `ini:"chan_size"`
}
type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}
type EtcConfig struct {
	Address string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}
//正真的业务逻辑
//	tailObj-->log-->client-->kafka


func main(){
//1.读配置文件
	var configObj =new(Config)
	//cfg,err:= ini.Load("./conf/config.ini")
	//if err!=nil{
	//	logrus.Error("load config faild :%v",err)
	//	return
	//}
	//kafkaaddr:=cfg.Section("kafka").Key("address").String()
	//fmt.Println(kafkaaddr)
	err:=ini.MapTo(configObj,"./conf/config.ini")
	if err!=nil{
		logrus.Errorf("load config faild :%v",err)
		return
	}

//2.初始化
	err= kafka.Init([]string{configObj.KafkaConfig.Address},configObj.KafkaConfig.ChanSize)
	if err!=nil{
		logrus.Errorf("init kafka failed :%v",err)
		return
	}
	logrus.Info("init kafka sucess !")
	//从etcd拉数据  配置日志项
	err=etcd.Init([]string{configObj.EtcConfig.Address})
	if err!=nil{
		logrus.Errorf("init etcd failed :%v",err)
		return
	}

	allConf,err:=etcd.GetConf(configObj.EtcConfig.CollectKey)
	if err!=nil{

		logrus.Errorf("get conf from etcd  failed :%v",err)
		return
	}
	//派小弟去监视etcd中configObj.EtcConfig.CollectKey的变化
	go etcd.WatchConf(configObj.EtcConfig.CollectKey)
	//3.使用tail读日志
	err = tailfile.Init(allConf)
	if err!=nil{
		logrus.Errorf("init tail failed :%v",err)
		return
	}
	logrus.Info("init tail sucess !")
	//死循环等着那些后台的日志收集任务
	for{
		select {

		}
	}

}
