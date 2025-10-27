package pumping

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// TCPClient TCP客户端
type TCPClient struct {
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

	// 添加消息队列用于异步发送
	sendQueue     chan []byte
	sendQueueSize int
}

// NewTCPClient 创建新的TCP客户端
func NewTCPClient(config *Config, handler MessageHandler) *TCPClient {
	if config == nil {
		config = DefaultConfig()
	}
	if handler == nil {
		handler = &DefaultMessageHandler{}
	}

	return &TCPClient{
		config:        config,
		handler:       handler,
		sendQueueSize: 1000, // 默认发送队列大小
	}
}

// Connect 连接到服务器
func (c *TCPClient) Connect() error {
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

	// 创建发送队列
	c.sendQueue = make(chan []byte, c.sendQueueSize)

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

// Send 发送数据（同步方式）
func (c *TCPClient) Send(data []byte) error {
	if !c.isConnected.Load() {
		return fmt.Errorf("client is not connected")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.writer == nil {
		return fmt.Errorf("writer is not initialized")
	}

	// 写入数据
	if _, err := c.writer.Write(data); err != nil {
		c.handler.OnError(fmt.Errorf("failed to write data: %w", err))
		return err
	}

	// 刷新缓冲区
	if err := c.writer.Flush(); err != nil {
		c.handler.OnError(fmt.Errorf("failed to flush writer: %w", err))
		return err
	}

	return nil
}

// SendAsync 异步发送数据
func (c *TCPClient) SendAsync(data []byte) error {
	if !c.isConnected.Load() {
		return fmt.Errorf("client is not connected")
	}

	select {
	case c.sendQueue <- data:
		return nil
	default:
		return fmt.Errorf("send queue is full")
	}
}

// SendJSON 发送JSON数据
func (c *TCPClient) SendJSON(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// 添加换行符作为消息分隔符
	data = append(data, '\n')
	return c.Send(data)
}

// SendJSONAsync 异步发送JSON数据
func (c *TCPClient) SendJSONAsync(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	data = append(data, '\n')
	return c.SendAsync(data)
}

// SendWithDelimiter 发送带分隔符的数据
func (c *TCPClient) SendWithDelimiter(data []byte, delimiter byte) error {
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

// Subscribe 发送订阅请求
func (c *TCPClient) Subscribe(request *TCPRequest) error {
	return c.SendJSON(request)
}

// SubscribeTick 订阅tick数据
func (c *TCPClient) SubscribeTick(symbols string) error {
	req := &TCPRequest{
		Type: string(REQUEST_TYPE_TICK),
		Params: TCPParams{
			Symbols: symbols,
		},
	}
	return c.SendJSON(req)
}

func (c *TCPClient) SubscribePosition() error {
	req := &TCPRequest{
		Type: string(REQUEST_TYPE_POSITION),
	}
	return c.SendJSON(req)
}

func (c *TCPClient) SubscribeDeal() error {
	req := &TCPRequest{
		Type: string(REQUEST_TYPE_DEAL),
	}
	return c.SendJSON(req)
}

func (c *TCPClient) SubscribeOrder() error {
	req := &TCPRequest{
		Type: string(REQUEST_TYPE_ORDER),
	}
	return c.SendJSON(req)
}

// Unsubscribe 发送取消订阅请求
func (c *TCPClient) Unsubscribe(requestType REQUEST_TYPE) error {
	// 根据实际协议实现取消订阅逻辑
	req := &TCPRequest{
		Type:   string(requestType),
		Params: TCPParams{
			// 取消订阅可能需要特定的参数
		},
	}
	return c.SendJSON(req)
}

// writeLoop 写入循环（处理异步发送）
func (c *TCPClient) writeLoop(ctx context.Context) {
	defer c.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case data := <-c.sendQueue:
			if err := c.Send(data); err != nil {
				c.handler.OnError(fmt.Errorf("async send failed: %w", err))
			}
		}
	}
}

// 其他方法保持不变...
func (c *TCPClient) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected.Load() {
		return nil
	}

	// 取消所有goroutine
	if c.cancel != nil {
		c.cancel()
	}

	// 关闭发送队列
	if c.sendQueue != nil {
		close(c.sendQueue)
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

func (c *TCPClient) readLoop(ctx context.Context) {
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

func (c *TCPClient) heartbeatLoop(ctx context.Context) {
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
				// 发送心跳包
				heartbeat := map[string]string{
					"type": "heartbeat",
					"time": time.Now().Format(time.RFC3339),
				}
				if err := c.SendJSON(heartbeat); err != nil {
					c.handler.OnError(fmt.Errorf("failed to send heartbeat: %w", err))
				}
			}
		}
	}
}

func (c *TCPClient) handleDisconnect() {
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

// IsConnected 检查是否连接
func (c *TCPClient) IsConnected() bool {
	return c.isConnected.Load()
}

// SetSendQueueSize 设置发送队列大小
func (c *TCPClient) SetSendQueueSize(size int) {
	c.sendQueueSize = size
}
