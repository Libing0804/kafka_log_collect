package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"logagent/common"
	"logagent/tailfile"
	"time"
)

var cli *clientv3.Client

func Init(addr []string) (err error) {
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	return
}

//拉取日志收集配置项的函数
func GetConf(key string) (collectEntryList []common.CollectEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	resp, err := cli.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get conf from etcd error by %s failed err :%s", key, err)
		return
	}
	if len(resp.Kvs) == 0 {
		logrus.Warningf("get len=0 conf from etcd error by %s", key)
		return
	}
	ret := resp.Kvs[0]

	err = json.Unmarshal(ret.Value, &collectEntryList)
	if err != nil {
		logrus.Errorf("json unmarshal failed err :%s", err)
		return
	}
	return
}

// 监视etcd中日志收集项的配置变化
func WatchConf(key string) {
	for {
		wCh := cli.Watch(context.Background(), key)

		for wresp := range wCh {
			logrus.Info("new conf from etc ")
			for _, evt := range wresp.Events {
				var newConf []common.CollectEntry
				fmt.Printf("type :%s ,key: %s , value :%s", evt.Type, evt.Kv.Key, evt.Kv.Value)

				if evt.Type == clientv3.EventTypeDelete {
					//	删除操作 发送空的
					logrus.Warning("etcd delete the key !!!")
					tailfile.SendNewConf(newConf) //没有人接收 阻塞住了
					continue
				}

				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					logrus.Errorf("json newconf failed err: %v", err)
					continue
				}
				//告诉tailfile 模块启用新的配置！
				tailfile.SendNewConf(newConf) //没有人接收 阻塞住了
			}
		}
	}
}
