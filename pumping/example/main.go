package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"safexapp.com/tradfi/go-mt5-sdk/pumping"
	"syscall"
	"time"
)

// CustomHandler 自定义消息处理器
type CustomHandler struct {
	pumping.DefaultMessageHandler
}

func (h *CustomHandler) OnMessage(data []byte) {
	fmt.Printf("Received: %s", string(data))
}

func (h *CustomHandler) OnConnected() {
	fmt.Println("Connected to server!")
}

func (h *CustomHandler) OnDisconnected() {
	fmt.Println("Disconnected from server!")
}

func (h *CustomHandler) OnError(err error) {
	fmt.Printf("Error: %v\n", err)
}

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

	// 创建客户端
	handler := &CustomHandler{}
	client := pumping.NewPumpingClient(config, handler)

	// 连接服务器
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// 发送测试消息
	go func() {
		time.Sleep(2 * time.Second)
		for i := 0; i < 5; i++ {
			msg := fmt.Sprintf("Hello, Server! %d\n", i)
			if err := client.Send([]byte(msg)); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			time.Sleep(1 * time.Second)
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
