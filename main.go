package main

import (
	"LogCollector/Kafka"
	"LogCollector/etcd"
	"LogCollector/tailfile"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"strings"
	"time"
)

type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
	EtcdConfig    `ini:"etcd"`
}

type KafkaConfig struct {
	Address  string `ini:"address"`
	ChanSize int64  `ini:"chan_size"`
}

type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}

type EtcdConfig struct {
	Address    string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

// 真正的业务逻辑
func run() (err error) {
	for {
		//循环读取数据
		line, ok := <-tailfile.TailObj.Lines
		//是空行就略过
		if len(strings.Trim(line.Text, "\r")) == 0 {
			logrus.Info("出现空行，直接跳过")
			continue
		}
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
		Kafka.ToMsgChan(msg)

		//将消息放进Kafka
		fmt.Println(line.Text)
	}
}

func main() {
	var configObj = new(Config)
	//读取配置
	err := ini.MapTo(configObj, "./conf/config.ini")
	if err != nil {
		logrus.Errorf("load config file failed,err:%v", err)
		return
	}

	//初始化Kafka
	err = Kafka.Init([]string{configObj.KafkaConfig.Address}, configObj.ChanSize)
	if err != nil {
		logrus.Errorf("init kafka failed,err:%v", err)
		return
	}
	logrus.Info("init kafka success")

	//初始化etcd连接
	err = etcd.Init([]string{configObj.EtcdConfig.Address})
	if err != nil {
		logrus.Errorf("init etcd failed,err:%v", err)
		return
	}
	//从etcd中拉取要收集日志的配置项
	allConf, err := etcd.GetConf(configObj.EtcdConfig.CollectKey)
	if err != nil {
		logrus.Errorf("get conf from etcd failed,err:%v", err)
		return
	}
	fmt.Println(allConf)

	//初始化tail
	tailfile.Init(configObj.LogFilePath)
	logrus.Info("tailfile init success")

	err = run()
	if err != nil {
		logrus.Errorf("run failed,err:%v", err)
	}
}
