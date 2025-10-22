package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"safexapp.com/tradfi/go-mt5-sdk/pumping"
)

func main() {
	// 创建配置
	config := &pumping.Config{
		ServerAddr:        "localhost:8080",
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

	// 注册tick消息处理器
	handler.RegisterHandlerWithUnmarshal(pumping.RequestTypeTick, &pumping.TCPRequest{}, func(msg interface{}) error {
		//todo 处理函数实现
		return nil
	})

	// 设置默认处理器
	handler.SetDefaultHandler(pumping.SubscriptionHandlerFunc(func(data []byte) error {
		//todo 处理函数实现
		return nil
	}))

	// 创建客户端
	client := pumping.NewTCPClient(config, handler)

	// 连接服务器
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// 等待连接建立后发送订阅请求
	go func() {
		time.Sleep(1 * time.Second)

		// 订阅tick数据
		tickReq := pumping.TCPRequest{
			Type: string(pumping.RequestTypeTick),
			Params: pumping.TCPParams{
				Symbols: "XAUUSD",
			},
		}

		if err := client.Subscribe(tickReq); err != nil {
			log.Printf("Failed to subscribe to stock: %v", err)
		} else {
			fmt.Println("📈 Subscribed to stock data")
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
