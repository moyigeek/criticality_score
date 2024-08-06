package main

import (
	"log"
	"os"
	"time"

	"github.com/HUSTSecLab/criticality_score/pkg/gitmetricsync"
)

func main() {
	// 检查命令行参数
	if len(os.Args) < 2 {
		log.Fatal("Usage: ./program <configPath>")
	}
	configPath := os.Args[1] // 第一个命令行参数是配置文件路径

	// 创建一个定时器，每30分钟触发一次
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	// 首次运行不等待定时器
	runSync(configPath)

	// 无限循环，直到程序被外部中断
	for {
		select {
		case <-ticker.C:
			runSync(configPath)
		}
	}
}

// runSync 封装了同步调用逻辑，以便可以在首次运行时直接调用，而不需要等待第一个ticker周期
func runSync(configPath string) {
	log.Println("Starting synchronization...")
	gitmetricsync.Run(configPath)
	log.Println("Synchronization complete.")
}
