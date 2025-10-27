package main

import (
	"fmt"
	"github.com/asaka1234/go-mt5-sdk/pumping"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 创建配置
	config := &pumping.Config{
		ServerAddr:        "localhost:8355",
		Timeout:           5 * time.Second,
		Reconnect:         true,
		MaxReconnects:     10,
		ReconnectInterval: 3 * time.Second,
		HeartbeatInterval: 20 * time.Second,
		BufferSize:        1024,
	}

	// 创建支持订阅的消息处理器
	handler := pumping.NewSubscriptionMessageHandler()
	// 设置连接事件回调
	handler.OnConnectedFunc = func() {
		fmt.Println("✅ Connected to server!")
	}
	handler.OnDisconnectedFunc = func() {
		fmt.Println("❌ Disconnected from server!")
	}
	handler.OnErrorFunc = func(err error) {
		fmt.Printf("⚠️ Error: %v\n", err)
	}

	// 注册tick消息处理器 - 使用类型化处理器
	handler.RegisterTypedHandler(pumping.REQUEST_TYPE_TICK, []pumping.MT5Tick{},
		func(response *pumping.TCPResponse, payload interface{}) error {
			tickItems := payload.([]pumping.MT5Tick)

			for i, item := range tickItems {
				tickTime := time.Unix(item.Time, 0)

				fmt.Printf("  [%d] %s - Ask: %s, Bid: %s, Time: %s\n",
					i+1,
					item.Symbol,
					item.AskE8,
					item.BidE8,
					tickTime.Format("15:04:05"))
			}
			return nil
		})

	// 设置默认处理器
	handler.SetDefaultHandler(func(response *pumping.TCPResponse) error {
		//todo 实现默认处理器
		return nil
	})

	//-------------------------------------------------------------------

	// 创建客户端
	client := pumping.NewTCPClient(config, handler)

	// 连接服务器
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// 等待连接建立后发送请求
	go func() {
		time.Sleep(1 * time.Second)

		// 订阅tick数据
		if err := client.SubscribeTick("XAUUSD,XAGUSD,EURUSD"); err != nil {
			log.Printf("Failed to subscribe to tick: %v", err)
		} else {
			fmt.Println("📈 Subscribed to tick data")
		}
	}()

	// 等待中断信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// 断开连接
	if err := client.Disconnect(); err != nil {
		log.Printf("Error during disconnect: %v", err)
	}

	fmt.Println("Client stopped")
}
