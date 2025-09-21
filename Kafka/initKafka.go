package Kafka

import (
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

var Client sarama.SyncProducer

func Init(addr []string) (err error) {
	//生产者配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	//连接Kafka
	Client, err = sarama.NewSyncProducer(addr, config)
	if err != nil {
		logrus.Error("kafka:producer closed, err:", err)
		return err
	}
	defer Client.Close()
	return nil
}
