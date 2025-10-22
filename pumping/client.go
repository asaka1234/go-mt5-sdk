package pumping

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// PumpingClient TCP客户端
type PumpingClient struct {
	config         *Config
	handler        MessageHandler
	conn           net.Conn
	reader         *bufio.Reader
	writer         *bufio.Writer
	mu             sync.RWMutex
	isConnected    atomic.Bool
	reconnectCount int
	cancel         context.CancelFunc
	wg             sync.WaitGroup
}

// NewPumpingClient 创建新的TCP客户端
func NewPumpingClient(config *Config, handler MessageHandler) *PumpingClient {
	if config == nil {
		config = DefaultConfig()
	}
	if handler == nil {
		handler = &DefaultMessageHandler{}
	}

	return &PumpingClient{
		config:  config,
		handler: handler,
	}
}

// Connect 连接到服务器
func (c *PumpingClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isConnected.Load() {
		return fmt.Errorf("client is already connected")
	}

	conn, err := net.DialTimeout("tcp", c.config.ServerAddr, c.config.Timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}

	c.conn = conn
	c.reader = bufio.NewReaderSize(conn, c.config.BufferSize)
	c.writer = bufio.NewWriterSize(conn, c.config.BufferSize)
	c.isConnected.Store(true)
	c.reconnectCount = 0

	// 创建上下文用于控制goroutine退出
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	// 启动读写goroutine
	c.wg.Add(3)
	go c.readLoop(ctx)
	go c.writeLoop(ctx)
	go c.heartbeatLoop(ctx)

	c.handler.OnConnected()
	return nil
}

// Disconnect 断开连接
func (c *PumpingClient) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected.Load() {
		return nil
	}

	// 取消所有goroutine
	if c.cancel != nil {
		c.cancel()
	}

	// 关闭连接
	if c.conn != nil {
		c.conn.Close()
	}

	// 等待所有goroutine退出
	c.wg.Wait()

	c.isConnected.Store(false)
	c.handler.OnDisconnected()
	return nil
}

// Send 发送数据
func (c *PumpingClient) Send(data []byte) error {
	if !c.isConnected.Load() {
		return fmt.Errorf("client is not connected")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if _, err := c.writer.Write(data); err != nil {
		c.handler.OnError(fmt.Errorf("failed to write data: %w", err))
		return err
	}

	return c.writer.Flush()
}

// SendWithDelimiter 发送带分隔符的数据
func (c *PumpingClient) SendWithDelimiter(data []byte, delimiter byte) error {
	if !c.isConnected.Load() {
		return fmt.Errorf("client is not connected")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if _, err := c.writer.Write(data); err != nil {
		c.handler.OnError(fmt.Errorf("failed to write data: %w", err))
		return err
	}

	if err := c.writer.WriteByte(delimiter); err != nil {
		c.handler.OnError(fmt.Errorf("failed to write delimiter: %w", err))
		return err
	}

	return c.writer.Flush()
}

// IsConnected 检查是否连接
func (c *PumpingClient) IsConnected() bool {
	return c.isConnected.Load()
}

// readLoop 读取循环
func (c *PumpingClient) readLoop(ctx context.Context) {
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, err := c.reader.ReadBytes('\n') // 以换行符作为分隔符
			if err != nil {
				if err == io.EOF {
					c.handleDisconnect()
					return
				}
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				c.handler.OnError(fmt.Errorf("read error: %w", err))
				c.handleDisconnect()
				return
			}

			// 处理接收到的消息
			if len(data) > 0 {
				c.handler.OnMessage(data)
			}
		}
	}
}

// writeLoop 写入循环（可用于处理写入队列）
func (c *PumpingClient) writeLoop(ctx context.Context) {
	defer c.wg.Done()

	// 这里可以实现消息队列等高级功能
	<-ctx.Done()
}

// heartbeatLoop 心跳循环
func (c *PumpingClient) heartbeatLoop(ctx context.Context) {
	defer c.wg.Done()

	if c.config.HeartbeatInterval <= 0 {
		return
	}

	ticker := time.NewTicker(c.config.HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if c.isConnected.Load() {
				// 发送心跳包，这里可以根据协议自定义
				heartbeat := []byte("PING\n")
				if err := c.Send(heartbeat); err != nil {
					c.handler.OnError(fmt.Errorf("failed to send heartbeat: %w", err))
				}
			}
		}
	}
}

// handleDisconnect 处理连接断开
func (c *PumpingClient) handleDisconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected.Load() {
		return
	}

	c.isConnected.Store(false)
	if c.conn != nil {
		c.conn.Close()
	}

	c.handler.OnDisconnected()

	// 自动重连逻辑
	if c.config.Reconnect && c.reconnectCount < c.config.MaxReconnects {
		c.reconnectCount++
		time.AfterFunc(c.config.ReconnectInterval, func() {
			if err := c.Connect(); err != nil {
				c.handler.OnError(fmt.Errorf("reconnect failed: %w", err))
			}
		})
	}
}
