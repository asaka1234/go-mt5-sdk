package pumping

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

// 基础消息结构，用于解析type字段
type BaseMessage struct {
	Type RequestType `json:"type"`
}

// SubscriptionHandler 订阅处理器
type SubscriptionHandler interface {
	HandleMessage(data []byte) error
}

// SubscriptionHandlerFunc 订阅处理器函数类型
type SubscriptionHandlerFunc func(data []byte) error

func (f SubscriptionHandlerFunc) HandleMessage(data []byte) error {
	return f(data)
}

// TypedHandler 带类型的消息处理器
type TypedHandler struct {
	MessageType RequestType
	Handler     SubscriptionHandler
	MessagePtr  interface{} // 用于JSON反序列化的消息指针
}

// SubscriptionManager 订阅管理器
type SubscriptionManager struct {
	mu             sync.RWMutex
	handlers       map[RequestType]*TypedHandler
	defaultHandler SubscriptionHandler
}

// NewSubscriptionManager 创建订阅管理器
func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		handlers: make(map[RequestType]*TypedHandler),
	}
}

// Register 注册消息处理器
func (sm *SubscriptionManager) Register(msgType RequestType, handler SubscriptionHandler, messagePtr interface{}) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.handlers[msgType] = &TypedHandler{
		MessageType: msgType,
		Handler:     handler,
		MessagePtr:  messagePtr,
	}
}

// RegisterFunc 使用函数注册消息处理器
func (sm *SubscriptionManager) RegisterFunc(msgType RequestType, handlerFunc func(data []byte) error, messagePtr interface{}) {
	sm.Register(msgType, SubscriptionHandlerFunc(handlerFunc), messagePtr)
}

// RegisterWithUnmarshal 注册支持自动反序列化的处理器
func (sm *SubscriptionManager) RegisterWithUnmarshal(msgType RequestType, messagePtr interface{}, handlerFunc func(interface{}) error) {
	sm.Register(msgType, SubscriptionHandlerFunc(func(data []byte) error {
		// 创建新的消息实例
		msg := reflect.New(reflect.TypeOf(messagePtr).Elem()).Interface()

		if err := json.Unmarshal(data, &msg); err != nil {
			return fmt.Errorf("failed to unmarshal message type %s: %w", msgType, err)
		}

		return handlerFunc(msg)
	}), messagePtr)
}

// SetDefaultHandler 设置默认处理器
func (sm *SubscriptionManager) SetDefaultHandler(handler SubscriptionHandler) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.defaultHandler = handler
}

// HandleMessage 处理消息
func (sm *SubscriptionManager) HandleMessage(data []byte) error {
	// 先解析基础消息获取type字段
	var baseMsg BaseMessage
	if err := json.Unmarshal(data, &baseMsg); err != nil {
		return fmt.Errorf("failed to unmarshal base message: %w", err)
	}

	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// 查找对应的处理器
	if handler, exists := sm.handlers[baseMsg.Type]; exists {
		return handler.Handler.HandleMessage(data)
	}

	// 使用默认处理器
	if sm.defaultHandler != nil {
		return sm.defaultHandler.HandleMessage(data)
	}

	return fmt.Errorf("no handler registered for message type: %s", baseMsg.Type)
}

// Unregister 取消注册消息处理器
func (sm *SubscriptionManager) Unregister(msgType RequestType) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.handlers, msgType)
}

// GetRegisteredTypes 获取已注册的消息类型
func (sm *SubscriptionManager) GetRegisteredTypes() []RequestType {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	types := make([]RequestType, 0, len(sm.handlers))
	for msgType := range sm.handlers {
		types = append(types, msgType)
	}
	return types
}
