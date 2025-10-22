package pumping

// MessageHandler 消息处理器接口
type MessageHandler interface {
	OnMessage(data []byte) // 收到消息时的回调
	OnConnected()          // 连接建立时的回调
	OnDisconnected()       // 连接断开时的回调
	OnError(err error)     // 发生错误时的回调
}

// DefaultMessageHandler 默认消息处理器
type DefaultMessageHandler struct {
	OnMessageFunc      func(data []byte)
	OnConnectedFunc    func()
	OnDisconnectedFunc func()
	OnErrorFunc        func(err error)
}

// 实现 MessageHandler 接口
func (h *DefaultMessageHandler) OnMessage(data []byte) {
	if h.OnMessageFunc != nil {
		h.OnMessageFunc(data)
	}
}

func (h *DefaultMessageHandler) OnConnected() {
	if h.OnConnectedFunc != nil {
		h.OnConnectedFunc()
	}
}

func (h *DefaultMessageHandler) OnDisconnected() {
	if h.OnDisconnectedFunc != nil {
		h.OnDisconnectedFunc()
	}
}

func (h *DefaultMessageHandler) OnError(err error) {
	if h.OnErrorFunc != nil {
		h.OnErrorFunc(err)
	}
}

// SubscriptionMessageHandler 支持订阅的消息处理器
type SubscriptionMessageHandler struct {
	DefaultMessageHandler
	subscriptionManager *SubscriptionManager
}

// NewSubscriptionMessageHandler 创建订阅消息处理器
func NewSubscriptionMessageHandler() *SubscriptionMessageHandler {
	return &SubscriptionMessageHandler{
		subscriptionManager: NewSubscriptionManager(),
	}
}

// OnMessage 处理接收到的消息
func (h *SubscriptionMessageHandler) OnMessage(data []byte) {
	if err := h.subscriptionManager.HandleMessage(data); err != nil {
		h.OnError(err)
	}
}

// RegisterHandler 注册消息处理器
func (h *SubscriptionMessageHandler) RegisterHandler(msgType RequestType, handler SubscriptionHandler, messagePtr interface{}) {
	h.subscriptionManager.Register(msgType, handler, messagePtr)
}

// RegisterHandlerFunc 使用函数注册消息处理器
func (h *SubscriptionMessageHandler) RegisterHandlerFunc(msgType RequestType, handlerFunc func(data []byte) error, messagePtr interface{}) {
	h.subscriptionManager.RegisterFunc(msgType, handlerFunc, messagePtr)
}

// RegisterHandlerWithUnmarshal 注册支持自动反序列化的处理器
func (h *SubscriptionMessageHandler) RegisterHandlerWithUnmarshal(msgType RequestType, messagePtr interface{}, handlerFunc func(interface{}) error) {
	h.subscriptionManager.RegisterWithUnmarshal(msgType, messagePtr, handlerFunc)
}

// SetDefaultHandler 设置默认处理器
func (h *SubscriptionMessageHandler) SetDefaultHandler(handler SubscriptionHandler) {
	h.subscriptionManager.SetDefaultHandler(handler)
}

// GetSubscriptionManager 获取订阅管理器（用于直接操作）
func (h *SubscriptionMessageHandler) GetSubscriptionManager() *SubscriptionManager {
	return h.subscriptionManager
}
