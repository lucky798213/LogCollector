package main

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"time"
)

func main() {
	// 1. 配置 InfluxDB 连接信息（替换为你的实际信息）
	influxURL := "http://localhost:8086" // 本地默认地址
	token := "your-auth-token"           // 你的 InfluxDB 认证 token
	org := "your-organization"           // 组织名
	bucket := "your-bucket"              // 桶名（数据存储的容器）

	// 2. 创建客户端
	client := influxdb2.NewClient(influxURL, token)
	defer client.Close() // 程序结束时关闭客户端

	// 3. 获取写入 API
	writeAPI := client.WriteAPI(org, bucket)

	// 4. 创建数据点（示例：写入一条 CPU 使用率数据）
	point := write.NewPoint(
		"cpu_usage", // 测量名（measurement）
		map[string]string{ // 标签（tags）
			"host": "localhost",
			"cpu":  "cpu0",
		},
		map[string]interface{}{ // 字段（fields，实际测量值）
			"usage": 25.5, // CPU 使用率
		},
		time.Now(), // 时间戳（可选，默认当前时间）
	)

	// 5. 写入数据
	writeAPI.WritePoint(point)
	// 等待写入完成
	writeAPI.Flush()

	// 6. 验证写入（可选，查询数据）
	queryAPI := client.QueryAPI(org)
	query := fmt.Sprintf(`from(bucket: "%s") 
		|> range(start: -1m) 
		|> filter(fn: (r) => r._measurement == "cpu_usage")`, bucket)

	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("查询失败: %v\n", err)
		return
	}

	// 打印查询结果
	for result.Next() {
		record := result.Record()
		fmt.Printf("时间: %v, 测量名: %s, CPU使用率: %v\n",
			record.Time(), record.Measurement(), record.Value())
	}
	if result.Err() != nil {
		fmt.Printf("查询结果错误: %v\n", result.Err())
	}
}
