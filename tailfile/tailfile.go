package tailfile

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"logagent/kafka"
	"strings"
	"time"
)

type tailTask struct {
	path  string
	topic string
	//收集日志的实例
	tobj *tail.Tail
	//	创建一个cancel函数  用于停止goroutine
	ctx    context.Context
	cancel context.CancelFunc
}

func newTailTask(path, topic string) *tailTask {
	ctx, cancel := context.WithCancel(context.Background())
	tt := &tailTask{
		path:   path,
		topic:  topic,
		ctx:    ctx,
		cancel: cancel,
	}

	return tt
}
func (t *tailTask) Init() (err error) {
	//默认项
	cfg := tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	}
	t.tobj, err = tail.TailFile(t.path, cfg)
	return
}
func (t *tailTask) run() {
	//	读取日志发送到kafka
	logrus.Infof("create a tail task for path: %s success", t.path)

	//循环读数据
	for {
		select {
		case <-t.ctx.Done(): //调用cancel（）就会触发
			logrus.Infof("path : %s is stopping... ", t.path)
			return
		case line, ok := <-t.tobj.Lines:

			if !ok {
				logrus.Warn("tail file close reopen,path:%s\n", t.path)
				time.Sleep(time.Second * 1)
				continue
			}
			//如果是空行直接略过 但是在windows是有换行符
			if len(strings.Trim(line.Text, "\r")) == 0 {
				continue
			}
			//	利用通道  改为异步并发
			//读出来line改为msg信息
			msg := &sarama.ProducerMessage{}
			msg.Topic = t.topic //自己的topic
			msg.Value = sarama.StringEncoder(line.Text)
			kafka.ToMsgChan(msg)

		}
	}
}
