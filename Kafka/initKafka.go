package Kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

var client sarama.SyncProducer
var msgChan chan *sarama.ProducerMessage

func Init(addr []string, chanSize int64) (err error) {
	//生产者配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	//连接Kafka
	client, err = sarama.NewSyncProducer(addr, config)
	if err != nil {
		logrus.Error("kafka:producer closed, err:", err)
		return err
	}

	msgChan = make(chan *sarama.ProducerMessage, chanSize)

	//起一个后台的goroutine 从msgchan中读数据
	go sendMsg()
	return nil
}

//从MsgChan中读取msg，发送给Kafka

func sendMsg() {
	for {
		select {
		case msg := <-msgChan:
			pid, offset, err := client.SendMessage(msg)
			if err != nil {
				logrus.Warning("send message err:", err)
				return
			}
			logrus.Info(fmt.Sprintf("pid:%v offset:%v\n", pid, offset))
		}
	}
}

func ToMsgChan(msg *sarama.ProducerMessage) {
	msgChan <- msg
}
