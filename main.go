package main

import (
	"LogCollector/Kafka"
	"LogCollector/tailfile"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

type Config struct {
	KafkaConfig   `ini:"kafka"`
	CollectConfig `ini:"collect"`
}

type KafkaConfig struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}

type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
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
	err = Kafka.Init([]string{configObj.Address})
	if err != nil {
		logrus.Error("init kafka failed,err:%v", err)
		return
	}
	logrus.Info("init kafka success")

	//初始化tail
	tailfile.Init(configObj.LogFilePath)
	logrus.Info("tailfile init success")

}
