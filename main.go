package main

import (
	"LogCollector/Kafka"
	"LogCollector/tailfile"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"time"
)

type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
}

type KafkaConfig struct {
	Address  string `ini:"address"`
	Topic    string `ini:"topic"`
	ChanSize int64  `ini:"chan_size"`
}

type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}

// 真正的业务逻辑
func run() (err error) {
	for {
		//循环读取数据
		line, ok := <-tailfile.TailObj.Lines
		if !ok {
			logrus.Warn("tail file close reopen, filename:%s\n", tailfile.TailObj.Filename)
			time.Sleep(time.Second)
			continue
		}
		//包装成Kafka中的msg类型
		msg := &sarama.ProducerMessage{}
		msg.Topic = "web_log"
		msg.Value = sarama.StringEncoder("this is a test")

		//丢到通道中
		Kafka.MsgChan <- msg
 
		//将消息放进Kafka
		fmt.Println(line.Text)
	}
}

func main() {
	var configObj = new(Config)
	//读取配置
	err := ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		logrus.Error("load config file failed,err:%v", err)
		return
	}

	//初始化Kafka
	err = Kafka.Init([]string{configObj.Address}, configObj.ChanSize)
	if err != nil {
		logrus.Error("init kafka failed,err:%v", err)
		return
	}
	logrus.Info("init kafka success")

	//初始化tail
	tailfile.Init(configObj.LogFilePath)
	logrus.Info("tailfile init success")

	err = run()
	if err != nil {
		logrus.Error("run failed,err:%v", err)
	}
}
