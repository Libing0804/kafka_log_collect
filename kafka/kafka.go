package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	Client sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
	)
func Init(address []string,chanSize int64)(err error){
	//	初始化生产者
	//生产者配置
	config:=sarama.NewConfig()
	//配置kafka的ack回传级别
	config.Producer.RequiredAcks=sarama.WaitForAll
	//发送到那个分区
	config.Producer.Partitioner=sarama.NewRandomPartitioner
	//	成功交付信息
	config.Producer.Return.Successes=true



	//	连接kafka
	Client,err=sarama.NewSyncProducer(address,config)
	if err!=nil {
		logrus.Error("kafka:producer closer err:",err)
		return
	}
	//初始化卡夫卡msgchan
	msgChan=make(chan *sarama.ProducerMessage,chanSize)
	//起一个后台的goroutine 读数据
	go sendMsg()
	return
}

//从通道中读取消息到msg发送到kafka
func sendMsg(){
	for {
		select {
		case msg:= <-msgChan:
			pid,offset,err:=Client.SendMessage(msg)
			if err!=nil{
				logrus.Warning("send msg to kafka failed ,err:",err)
				return
			}
			logrus.Info("send msg to kafka success,pid:%v,offset :%v",pid,offset)
		}
	}
}

func ToMsgChan(msg *sarama.ProducerMessage){
	msgChan <-msg
}