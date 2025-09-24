package tailfile

import (
	"LogCollector/Kafka"
	"context"
	"github.com/IBM/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type tailTask struct {
	path   string
	topic  string
	tObj   *tail.Tail
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
		select {
		case <-tt.ctx.Done():
			logrus.Infof("task collect path:%s need to stop", tt.path)
			return
		case line, ok := <-tt.tObj.Lines:
			//循环读取数据
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
}
