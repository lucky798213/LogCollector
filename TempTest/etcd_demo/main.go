package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	// 示例：创建 etcd 客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"}, // etcd 服务地址（集群可填多个）
		DialTimeout: 5 * time.Second,            // 连接超时时间
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	//写
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = cli.Put(ctx, "name", "etcd-demo")
	cancel() // 及时释放 context，避免资源泄漏
	if err != nil {
		fmt.Printf("写入键值失败: %v\n", err)
	} else {
		fmt.Println("成功写入键值对 name=etcd-demo")
	}

	//取
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "name")
	cancel()
	if err != nil {
		fmt.Printf("读取键值失败: %v\n", err)
	} else {
		// 遍历结果（支持多键查询，此处仅查单个键）
		for _, kv := range resp.Kvs {
			fmt.Printf("键: %s, 值: %s\n", kv.Key, kv.Value)
		}
	}

	// 监听键 "name" 的变化
	rch := cli.Watch(context.Background(), "name")
	fmt.Println("开始监听键 name 的变化...")

	// 从通道循环读取变更事件
	for wresp := range rch {
		for _, ev := range wresp.Events {
			// 打印事件类型（PUT/DELETE）、键、值
			fmt.Printf("事件类型: %s, 键: %s, 值: %s\n",
				ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}
