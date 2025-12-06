package pumping

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asaka1234/go-mt5-sdk/utils"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrNotConnected   = errors.New("client is not connected")
	ErrNotInitialized = errors.New("not initialized")
)

// TCPClient TCP客户端
type TCPClient struct {
	config         *Config
	handler        MessageHandler
	conn           net.Conn
	reader         *bufio.Reader
	writer         *bufio.Writer
	mu             sync.RWMutex
	isSubscribed   atomic.Bool
	isConnected    atomic.Bool
	reconnectCount atomic.Int32
	cancel         context.CancelFunc
	wg             sync.WaitGroup
	stopChan       chan struct{}
	reconnectCh    chan struct{}
	initialBackoff time.Duration
	alert          AlertHandler
	subRequests    []*TCPRequest

	// 添加消息队列用于异步发送
	sendQueue     chan []byte
	sendQueueSize int
}

// NewTCPClient 创建新的TCP客户端
func NewTCPClient(config *Config, handler MessageHandler, alert AlertHandler) *TCPClient {
	if config == nil {
		config = DefaultConfig()
	}
	if handler == nil {
		handler = &DefaultMessageHandler{}
	}

	return &TCPClient{
		config:         config,
		handler:        handler,
		reconnectCh:    make(chan struct{}, 1),
		stopChan:       make(chan struct{}),
		initialBackoff: time.Second,
		sendQueueSize:  1000, // 默认发送队列大小
		alert:          alert,
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
		if c.alert != nil {
			go c.alert.SendAlert(AlertLevelError, AlertActionConnect, "failed to connect server, server: "+c.config.ServerAddr)
		}

		return fmt.Errorf("failed to connect to server: %w", err)
	}

	if c.alert != nil {
		go c.alert.SendAlert(AlertLevelWarning, AlertActionConnect, "connect to server successfully, server: "+c.config.ServerAddr)
	}

	c.conn = conn
	c.reader = bufio.NewReaderSize(conn, c.config.BufferSize)
	c.writer = bufio.NewWriterSize(conn, c.config.BufferSize)
	c.isConnected.Store(true)
	c.isSubscribed.Store(false)
	c.reconnectCount.Store(0)

	// 创建发送队列
	c.sendQueue = make(chan []byte, c.sendQueueSize)

	// 创建上下文用于控制goroutine退出
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel

	// monitor goroutine
	go c.monitorLoop()

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
		return ErrNotConnected
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conn == nil || c.writer == nil {
		return ErrNotInitialized
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
		return ErrNotConnected
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
		return ErrNotConnected
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

func (c *TCPClient) SetSubscribes(requests []*TCPRequest) {
	c.subRequests = requests
}

// SetSubscribeTick 订阅tick数据
func (c *TCPClient) SetSubscribeTick(symbols string) {
	req := &TCPRequest{
		Type: string(REQUEST_TYPE_TICK),
		Params: TCPParams{
			Symbols: symbols,
		},
	}
	c.subRequests = append(c.subRequests, req)
}

func (c *TCPClient) SetSubscribePosition() {
	req := &TCPRequest{
		Type: string(REQUEST_TYPE_POSITION),
	}
	c.subRequests = append(c.subRequests, req)
}

func (c *TCPClient) SetSubscribeDeal() {
	req := &TCPRequest{
		Type: string(REQUEST_TYPE_DEAL),
	}
	c.subRequests = append(c.subRequests, req)
}

func (c *TCPClient) SetSubscribeUserAdd() {
	req := &TCPRequest{
		Type: string(REQUEST_TYPE_USER_ADD),
	}
	c.subRequests = append(c.subRequests, req)
}

func (c *TCPClient) SetSubscribeOrder() {
	req := &TCPRequest{
		Type: string(REQUEST_TYPE_ORDER),
	}
	c.subRequests = append(c.subRequests, req)
}

func (c *TCPClient) SetSubscribeMarginCall() {
	req := &TCPRequest{
		Type: string(REQUEST_TYPE_MARGINCAL),
	}
	c.subRequests = append(c.subRequests, req)
}

func (c *TCPClient) SetSubscribeStopOut() {
	req := &TCPRequest{
		Type: string(REQUEST_TYPE_STOPOUT),
	}
	c.subRequests = append(c.subRequests, req)
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

func (c *TCPClient) SubscribeUserAdd() error {
	req := &TCPRequest{
		Type: string(REQUEST_TYPE_USER_ADD),
	}
	return c.SendJSON(req)
}

func (c *TCPClient) SubscribeOrder() error {
	req := &TCPRequest{
		Type: string(REQUEST_TYPE_ORDER),
	}
	return c.SendJSON(req)
}

func (c *TCPClient) SubscribeMarginCall() error {
	req := &TCPRequest{
		Type: string(REQUEST_TYPE_MARGINCAL),
	}
	return c.SendJSON(req)
}

func (c *TCPClient) SubscribeStopOut() error {
	req := &TCPRequest{
		Type: string(REQUEST_TYPE_STOPOUT),
	}
	return c.SendJSON(req)
}

// Unsubscribe 发送取消订阅请求
func (c *TCPClient) Unsubscribe(requestType REQUEST_TYPE) error {
	// update requests
	if len(c.subRequests) > 0 {
		requests := make([]*TCPRequest, 0, len(c.subRequests))
		for _, req := range c.subRequests {
			if req.Type == string(requestType) {
				continue
			}
			obj := req
			requests = append(requests, obj)
		}
		c.subRequests = requests
	}

	// 根据实际协议实现取消订阅逻辑
	req := &TCPRequest{
		Type: string(requestType),
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
		default:
			if !c.isConnected.Load() {
				continue
			}

			if c.conn == nil || c.writer == nil {
				continue
			}

			if data, ok := <-c.sendQueue; ok {
				if err := c.Send(data); err != nil {
					// handle error
					if !errors.Is(err, ErrNotConnected) && !errors.Is(err, ErrNotInitialized) {
						c.isConnected.Store(false)
						select {
						case c.reconnectCh <- struct{}{}:
						default:
						}

						c.handler.OnError(fmt.Errorf("async send failed: %w", err))

						if c.alert != nil {
							go c.alert.SendAlert(AlertLevelError, AlertActionClose, "connection is closed, server: "+c.config.ServerAddr)
						}
					}
				}
			}
		}
	}
}

// 其他方法保持不变...
func (c *TCPClient) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	//if !c.isConnected.Load() {
	//	return nil
	//}

	// close monitor and others goroutine
	if c.stopChan != nil {
		close(c.stopChan)
	}

	// 取消所有goroutine
	if c.cancel != nil {
		c.cancel()
		c.cancel = nil
	}

	// 关闭发送队列
	if c.sendQueue != nil {
		close(c.sendQueue)
	}

	// 关闭连接
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}

	// 等待所有goroutine退出
	c.wg.Wait()

	c.isSubscribed.Store(false)
	c.isConnected.Store(false)
	c.handler.OnDisconnected()
	return nil
}

func (c *TCPClient) readLoop(ctx context.Context) {
	defer c.wg.Done()

	for {
		select {
		case <-c.stopChan:
			return
		case <-ctx.Done():
			return
		default:
			if !c.isConnected.Load() {
				continue
			}

			if c.conn == nil || c.reader == nil {
				continue
			}

			data, err := utils.DecodeTCPMsgWithLePrefix(c.reader)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				// setting reconnection
				c.isConnected.Store(false)
				select {
				case c.reconnectCh <- struct{}{}:
				default:
				}

				if err == io.EOF {
					return
				}

				c.handler.OnError(fmt.Errorf("read error: %w", err))

				if c.alert != nil {
					go c.alert.SendAlert(AlertLevelError, AlertActionClose, "connection is closed, server: "+c.config.ServerAddr)
				}

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
		case <-c.stopChan:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			if !c.isConnected.Load() {
				continue
			}

			if c.conn == nil || c.writer == nil {
				continue
			}

			// 发送心跳包
			heartbeat := map[string]string{
				"type": "heartbeat",
				"time": time.Now().Format(time.RFC3339),
			}

			if err := c.SendJSON(heartbeat); err != nil {
				// handle error
				if !errors.Is(err, ErrNotConnected) && !errors.Is(err, ErrNotInitialized) {
					c.isConnected.Store(false)
					select {
					case c.reconnectCh <- struct{}{}:
					default:
					}

					c.handler.OnError(fmt.Errorf("heartbeat error: %w", err))

					if c.alert != nil {
						go c.alert.SendAlert(AlertLevelError, AlertActionClose, "connection is closed, server: "+c.config.ServerAddr)
					}
				}
			}
		}
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

// monitorLoop 监控循环（处理重连）
func (c *TCPClient) monitorLoop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopChan:
			return

		case <-ticker.C:
			// send subscriptions
			if c.isSubscribed.Load() {
				continue
			}

			isReSend := false
			for _, req := range c.subRequests {
				if err := c.SendJSON(req); err != nil {
					if !errors.Is(err, ErrNotConnected) && !errors.Is(err, ErrNotInitialized) {
						c.handler.OnError(fmt.Errorf("subscribe error: %w", err))
					}
					isReSend = true
					break
				}
			}

			if isReSend {
				c.isSubscribed.Store(false)
			} else {
				c.isSubscribed.Store(true)
			}

		case <-c.reconnectCh:
			if c.isConnected.Load() {
				continue
			}

			// send subscriptions
			c.isSubscribed.Store(false)

			c.mu.Lock()
			// clean all worker group
			if c.cancel != nil {
				c.cancel()
				c.cancel = nil
			}

			if c.sendQueue != nil {
				// 确保只关闭一次，因为 handleDisconnect 和 Disconnect 都可能调用
				close(c.sendQueue)
				c.sendQueue = nil
			}

			if c.conn != nil {
				c.conn.Close()
				c.conn = nil
			}
			c.mu.Unlock()

			c.wg.Wait()

			attempt := int(c.reconnectCount.Add(1))

			// 判断是否达到最大重试次数
			if c.config.MaxReconnects > 0 && attempt > c.config.MaxReconnects {
				c.reconnectCount.Store(0)
				if c.alert != nil {
					go c.alert.SendAlert(AlertLevelError, AlertActionDisConnect, "failed to connect server, max reconnects exceeded, server: "+c.config.ServerAddr)
				}
			}

			// 计算指数退避时间：1s → 2s → 4s → 8s → ... 但不超过 ReconnectInterval
			backoff := c.initialBackoff * time.Duration(1<<(attempt-1))
			if backoff > c.config.ReconnectInterval {
				backoff = c.config.ReconnectInterval
			}

			// reconnect
			if err := c.reconnect(); err != nil {
				c.reconnectCh <- struct{}{}
				c.handler.OnError(fmt.Errorf("reconnect failed: %w", err))
				if c.alert != nil {
					go c.alert.SendAlert(AlertLevelError, AlertActionDisConnect, "failed to connect server, server: "+c.config.ServerAddr)
				}
				time.Sleep(backoff)
			}

			if c.alert != nil {
				go c.alert.SendAlert(AlertLevelWarning, AlertActionDisConnect, "connect to server successfully, server: "+c.config.ServerAddr)
			}
		}
	}
}

// performReconnect 执行一次重连（含指数退避逻辑）
func (c *TCPClient) reconnect() error {
	// 尝试重新建立连接
	conn, err := net.DialTimeout("tcp", c.config.ServerAddr, c.config.Timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}

	c.mu.Lock()
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	c.conn = conn
	c.reader = bufio.NewReaderSize(conn, c.config.BufferSize)
	c.writer = bufio.NewWriterSize(conn, c.config.BufferSize)
	c.isConnected.Store(true)
	c.isSubscribed.Store(false)
	c.reconnectCount.Store(0)
	// 创建发送队列
	c.sendQueue = make(chan []byte, c.sendQueueSize)
	c.mu.Unlock()

	// 启动读写goroutine
	c.wg.Add(3)
	go c.readLoop(ctx)
	go c.writeLoop(ctx)
	go c.heartbeatLoop(ctx)

	c.handler.OnConnected()

	return nil
}
