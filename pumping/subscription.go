package pumping

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"sync"
)

// ResponseHandler 响应处理器
type ResponseHandler struct {
	RequestType REQUEST_TYPE
	Handler     func(response *TCPResponse) error
}

// TypedResponseHandler 类型化响应处理器
type TypedResponseHandler struct {
	RequestType REQUEST_TYPE
	PayloadType interface{} // payload的类型
	Handler     func(response *TCPResponse, payload interface{}) error
}

// SubscriptionManager 订阅管理器
type SubscriptionManager struct {
	mu             sync.RWMutex
	handlers       map[REQUEST_TYPE]ResponseHandler
	typedHandlers  map[REQUEST_TYPE]TypedResponseHandler
	defaultHandler func(response *TCPResponse) error
}

// NewSubscriptionManager 创建订阅管理器
func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		handlers:      make(map[REQUEST_TYPE]ResponseHandler),
		typedHandlers: make(map[REQUEST_TYPE]TypedResponseHandler),
	}
}

// RegisterHandler 注册基础处理器
func (sm *SubscriptionManager) RegisterHandler(
	requestType REQUEST_TYPE,
	handler func(response *TCPResponse) error,
) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.handlers[requestType] = ResponseHandler{
		RequestType: requestType,
		Handler:     handler,
	}
}

// RegisterTypedHandler 注册类型化处理器
func (sm *SubscriptionManager) RegisterTypedHandler(
	requestType REQUEST_TYPE,
	payloadType interface{},
	handler func(response *TCPResponse, payload interface{}) error,
) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.typedHandlers[requestType] = TypedResponseHandler{
		RequestType: requestType,
		PayloadType: payloadType,
		Handler:     handler,
	}
}

// SetDefaultHandler 设置默认处理器
func (sm *SubscriptionManager) SetDefaultHandler(handler func(response *TCPResponse) error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.defaultHandler = handler
}

// HandleMessage 处理消息
func (sm *SubscriptionManager) HandleMessage(data []byte) error {
	// 解析基础响应
	var response TCPResponse
	err := Decode[TCPResponse](data, &response)
	if err != nil {
		return err
	}

	// 检查状态
	if response.Status != "ok" {
		return fmt.Errorf("server returned error status: %s", response.Status)
	}

	sm.mu.RLock()
	defer sm.mu.RUnlock()

	// 查找类型化处理器
	if handler, exists := sm.typedHandlers[REQUEST_TYPE(response.Type)]; exists {
		return sm.handleTypedResponse(&response, handler)
	}

	// 查找基础处理器
	if handler, exists := sm.handlers[REQUEST_TYPE(response.Type)]; exists {
		return handler.Handler(&response)
	}

	// 使用默认处理器
	if sm.defaultHandler != nil {
		return sm.defaultHandler(&response)
	}

	return fmt.Errorf("no handler registered for request type: %s", response.Type)
}

// handleTypedResponse 处理类型化响应
func (sm *SubscriptionManager) handleTypedResponse(response *TCPResponse, handler TypedResponseHandler) error {
	// 将payload转换为JSON
	payloadJSON, err := jsoniter.Marshal(response.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// 创建payload类型的新实例
	var payload interface{}
	switch handler.PayloadType.(type) {
	case []MT5Tick:
		var tickPayload []MT5Tick
		if err := jsoniter.Unmarshal(payloadJSON, &tickPayload); err != nil {
			return fmt.Errorf("failed to unmarshal tick payload: %w", err)
		}
		payload = tickPayload
	case []MTOrderExtra:
		var orderPayload []MTOrderExtra
		if err := jsoniter.Unmarshal(payloadJSON, &orderPayload); err != nil {
			return fmt.Errorf("failed to unmarshal order payload: %w", err)
		}
		payload = orderPayload
	case []MTPositionExtra:
		var posPayload []MTPositionExtra
		if err := jsoniter.Unmarshal(payloadJSON, &posPayload); err != nil {
			return fmt.Errorf("failed to unmarshal pos payload: %w", err)
		}
		payload = posPayload
	case []Mt5Deal:
		var dealPayload []Mt5Deal
		if err := jsoniter.Unmarshal(payloadJSON, &dealPayload); err != nil {
			return fmt.Errorf("failed to unmarshal deal payload: %w", err)
		}
		payload = dealPayload
	case []MT5MarginCall:
		var margincallPayload []MT5MarginCall
		if err := jsoniter.Unmarshal(payloadJSON, &margincallPayload); err != nil {
			return fmt.Errorf("failed to unmarshal margincall payload: %w", err)
		}
		payload = margincallPayload
	case []MT5StopOut:
		var soPayload []MT5StopOut
		if err := jsoniter.Unmarshal(payloadJSON, &soPayload); err != nil {
			return fmt.Errorf("failed to unmarshal so payload: %w", err)
		}
		payload = soPayload
	default:
		// 通用处理
		newPayload := handler.PayloadType
		if err := jsoniter.Unmarshal(payloadJSON, &newPayload); err != nil {
			return fmt.Errorf("failed to unmarshal payload to type: %w", err)
		}
		payload = newPayload
	}

	return handler.Handler(response, payload)
}

// Unregister 取消注册消息处理器
func (sm *SubscriptionManager) Unregister(requestType REQUEST_TYPE) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.handlers, requestType)
	delete(sm.typedHandlers, requestType)
}

// GetRegisteredTypes 获取已注册的请求类型
func (sm *SubscriptionManager) GetRegisteredTypes() []REQUEST_TYPE {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	types := make([]REQUEST_TYPE, 0)
	for requestType := range sm.handlers {
		types = append(types, requestType)
	}
	for requestType := range sm.typedHandlers {
		types = append(types, requestType)
	}
	return types
}
