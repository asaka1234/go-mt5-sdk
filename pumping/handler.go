package pumping

// MessageHandler 消息处理器接口
type MessageHandler interface {
	OnMessage(data []byte) // 收到消息时的回调
	OnConnected()          // 连接建立时的回调
	OnDisconnected()       // 连接断开时的回调
	OnError(err error)     // 发生错误时的回调
}

// DefaultMessageHandler 默认消息处理器
type DefaultMessageHandler struct{}

func (h *DefaultMessageHandler) OnMessage(data []byte) {
	// 默认实现，什么都不做
}

func (h *DefaultMessageHandler) OnConnected() {
	// 默认实现，什么都不做
}

func (h *DefaultMessageHandler) OnDisconnected() {
	// 默认实现，什么都不做
}

func (h *DefaultMessageHandler) OnError(err error) {
	// 默认实现，什么都不做
}
