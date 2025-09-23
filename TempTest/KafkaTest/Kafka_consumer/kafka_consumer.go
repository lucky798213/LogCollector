package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"sync"
)

func main() {
	//创建新的消费者(最小消费单元，并连接到Broker节点)
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, nil)
	if err != nil {
		fmt.Println("Error creating consumer:", err)
		return
	}
	//Kafka集群中的元数据是相通的，所以只通过一个Broker都可以获取到指定Topic的所有partition
	//拿到指定Topic下面的所有分区列表
	patitionList, err := consumer.Partitions("web_blog") //根据Topic取到所有分区
	if err != nil {
		fmt.Println("Error getting list of partitions:", err)
		return
	}
	fmt.Println("Partitions:", patitionList)
	var wg sync.WaitGroup
	for _, partition := range patitionList { //遍历所有分区
		//分区消费者” 是客户端为了实现 “按分区粒度消费” 而设计的组件，它不是 Kafka 服务器端的实体
		//针对每个分区创建一个对应的分区消费者
		pc, err := consumer.ConsumePartition("web_blog", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Println("Error creating partition consumer:", err)
			return
		}
		defer pc.AsyncClose()
		//异步从每个分区消费信息
		wg.Add(1)
		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d offset:%d Key:%s Value:%s ", msg.Partition, msg.Offset, msg.Key, msg.Value)
			}
		}(pc)
	}
	wg.Wait()
}
