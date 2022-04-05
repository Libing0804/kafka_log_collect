package tailfile

import (
	"github.com/sirupsen/logrus"
	"logagent/common"
)

// 用来管理 tailfile
type tailTaskMgr struct {
	tailTaskMap      map[string]*tailTask       //所有task任务
	CollectEntryList []common.CollectEntry      //所有配置项
	confChan         chan []common.CollectEntry //等待新配置的通道
}

//单例模式创建一个全局的变量
var (
	ttMgr *tailTaskMgr
)

func Init(allConf []common.CollectEntry) (err error) {
	//allConf 存了好多个日志收集项
	//每一个收集项 创建一个tailobj
	ttMgr = &tailTaskMgr{
		tailTaskMap:      make(map[string]*tailTask, 20),
		CollectEntryList: allConf,
		confChan:         make(chan []common.CollectEntry),
	}
	for _, conf := range allConf {
		tt := newTailTask(conf.Path, conf.Topic) //创建一个日志收集任务
		err := tt.Init()                         //初始化一个日志收集任务
		if err != nil {
			logrus.Errorf("create tailobj for path :%s failed err:%v", conf.Path, err)
			continue
		}
		logrus.Infof("create a tail task for path: %s success", conf.Path)
		ttMgr.tailTaskMap[tt.path] = tt //把创建的tailTask保存起来，方便后续管理
		//去干活 收集日志
		go tt.run()
	}
	go ttMgr.watch() //等新配置

	return
}

func (t *tailTaskMgr) watch() {
	for {
		//派一个小弟监视新配置来没来
		newConf := <-t.confChan //取到值说明新配置来了
		//来了新配置管理一下我之前的tailTask
		logrus.Infof("get new conf from etcd conf:%v,start manage tailTsak...", newConf)
		for _, conf := range newConf {
			//1、原来已经存在的就别动
			if t.isExist(conf) {
				continue
			}
			//2、原来没有的  新创建一个
			tt := newTailTask(conf.Path, conf.Topic) //创建一个日志收集任务
			err := tt.Init()                         //初始化一个日志收集任务
			if err != nil {
				logrus.Errorf("create tailobj for path :%s failed err:%v", conf.Path, err)
				continue
			}
			t.tailTaskMap[tt.path] = tt
			//去干活 收集日志
			go tt.run()

		}
		//3、原来有现在没有的  停止
		//	找出tailTaskMap中存在，但是新的newConf中不存在的
		for key, task := range t.tailTaskMap {
			var found bool
			for _, conf := range newConf {
				if key == conf.Path {
					found = true
					break
				}
			}
			if !found {
				//	没找到 这个任务应该停下来
				//	这里使用的是context上下文管理的停止，在日志收集结构体中加上ctx
				//和cancel（）
				logrus.Infof("the task collect path :%s need to stop ", task.path)
				//别忘了在管理的map中删除
				delete(t.tailTaskMap, key)
				task.cancel()
			}
		}
	}
}
func (t *tailTaskMgr) isExist(conf common.CollectEntry) bool {
	_, ok := t.tailTaskMap[conf.Path]
	return ok
}
func SendNewConf(newConf []common.CollectEntry) {
	ttMgr.confChan <- newConf
}
