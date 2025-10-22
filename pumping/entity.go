package pumping

// JSON 消息结构
type TCPRequest struct {
	Type   string    `json:"type"`   // 请求类型
	Params TCPParams `json:"params"` // 对应请求类型的附加信息
}

// 不同类型的负载结构
type TCPParams struct {
	//for login
	Username string `json:"username"`
	Password string `json:"password"`
	//for tick
	Symbols string `json:"symbols"` //是要订阅的symbol的列表,用comma拼接.  空则是所有
}
