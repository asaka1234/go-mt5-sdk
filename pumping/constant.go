package pumping

// RequestType 消息类型
type RequestType string

// 定义消息类型和结构
const (
	RequestTypeTick      RequestType = "tick"
	RequestTypeHeartBeat RequestType = "heartbeat"
)
