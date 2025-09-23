package etcd

import (
	"LogCollector/common"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

var (
	client *clientv3.Client
)

func Init(address []string) (err error) {
	// 1. 连接到etcd
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   address,         // etcd服务器地址
		DialTimeout: 5 * time.Second, // 连接超时时间
	})
	if err != nil {
		fmt.Printf("连接etcd失败: %v\n", err)
		return
	}
	return
}

// 拉取日志收集配置项的函数
func GetConf(key string) (collectEntryList []common.CollectEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	resp, err := client.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get config from etcd by key :%s failed,err:%v", key, err)
		return
	}

	if len(resp.Kvs) == 0 {
		logrus.Errorf("get len:0 config from etcd by key :%s ", key)
		return
	}

	ret := resp.Kvs[0]
	err = json.Unmarshal(ret.Value, &collectEntryList)
	if err != nil {
		logrus.Errorf("unmarshal config from etcd by key :%s failed,err:%v", key, err)
		return
	}
	return
}
