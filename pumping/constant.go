package pumping

// RequestType 消息类型
type REQUEST_TYPE string

// 定义消息类型和结构
const (
	REQUEST_TYPE_DEAL     REQUEST_TYPE = "deal"
	REQUEST_TYPE_POSITION REQUEST_TYPE = "position"
	REQUEST_TYPE_ORDER    REQUEST_TYPE = "order"
	REQUEST_TYPE_TICK     REQUEST_TYPE = "tick"
	REQUEST_TYPE_UNKNOW   REQUEST_TYPE = "unknown" //未知
)
