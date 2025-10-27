package pumping

// JSON 消息结构
type TCPRequest struct {
	Type   string    `json:"type"`   // 请求类型
	Params TCPParams `json:"params"` // 对应请求类型的附加信息
}

// 不同类型的负载结构
type TCPParams struct {
	//for login
	//for tick
	Symbols string `json:"symbols"` //是要订阅的symbol的列表,用comma拼接.  空则是所有
}

//----------------------------------

// TCPResponse 响应结构 - 改进版
type TCPResponse struct {
	Status    string      `json:"status"`            //ok是成功,其他是错误描述
	Type      string      `json:"type"`              //映射:请求类型
	Payload   interface{} `json:"payload,omitempty"` // 响应负载数据
	Timestamp int64       `json:"timestamp"`
}

// todo 暂时未用
type LoginPayloadItem struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Token    string `json:"token"`
	Expires  int64  `json:"expires"`
}

// todo 暂时未用
type HeartbeatPayloadItem struct {
	ServerTime    string `json:"server_time"`
	Timestamp     int64  `json:"timestamp"`
	Uptime        int64  `json:"uptime"`
	ActiveClients int    `json:"active_clients"`
}

type TickPayloadItem struct {
	Symbol string `json:"symbol"  msgpack:"symbol"`
	Ask    string `json:"ask" msgpack:"ask"`
	Bid    string `json:"bid" msgpack:"bid"`
	Last   string `json:"last" msgpack:"last"`
	Volume uint64 `json:"volume" msgpack:"volume"`
	Time   int64  `json:"time" msgpack:"time"` //unix时间戳(ms毫秒)
}
