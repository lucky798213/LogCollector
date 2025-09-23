package tailfile

import (
	"LogCollector/Kafka"
	"LogCollector/common"
	"github.com/IBM/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type tailTask struct {
	path  string
	topic string
	tObj  *tail.Tail
}

func newTailTask(path, topic string) *tailTask {
	tt := &tailTask{
		path:  path,
		topic: topic,
	}

	return tt
}

func (tt *tailTask) Init() (err error) {
	config := tail.Config{
		Follow:    true,
		ReOpen:    true,
		MustExist: false,
		Poll:      true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
	}
	tt.tObj, err = tail.TailFile(tt.path, config)
	return
}

func (tt *tailTask) run() {
	for {
		//循环读取数据
		line, ok := <-tt.tObj.Lines
		//是空行就略过
		if len(strings.Trim(line.Text, "\r")) == 0 {
			logrus.Info("出现空行，直接跳过")
			continue
		}
		if !ok {
			logrus.Warn("tail file close reopen, filename:%s\n", tt.path)
			time.Sleep(time.Second)
			continue
		}
		//包装成Kafka中的msg类型
		msg := &sarama.ProducerMessage{}
		msg.Topic = tt.topic
		msg.Value = sarama.StringEncoder(line.Text)

		//丢到通道中
		Kafka.ToMsgChan(msg)

	}
}

func Init(allConf []common.CollectEntry) (err error) {
	//allConf中存了若干个日志的收集项
	//针对每一个日志收集项创建一个对应的tailObj
	for _, conf := range allConf {
		tt := newTailTask(conf.Path, conf.Topic)
		err := tt.Init()
		if err != nil {
			logrus.Errorf("create tailObj for path:%s failed,err:%v", conf.Path, err)
			continue
		}
		go tt.run()
	}

	return
}
