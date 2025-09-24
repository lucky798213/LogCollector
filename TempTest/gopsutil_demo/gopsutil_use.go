package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"time"
)

func main() {
	// 获取CPU核心数
	count, _ := cpu.Counts(false) // false 表示物理核心，true 表示逻辑核心
	fmt.Printf("CPU核心数: %d\n", count)

	// 获取CPU使用率（间隔1秒，返回所有核心的使用率）
	percent, _ := cpu.Percent(1*time.Second, true)
	fmt.Println("CPU使用率（每个核心）:", percent)

	// 获取CPU型号
	info, _ := cpu.Info()
	fmt.Println("CPU型号:", info[0].ModelName)
}
