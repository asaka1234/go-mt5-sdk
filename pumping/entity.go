package pumping

// JSON 消息结构
type TCPRequest struct {
	Type   string    `json:"type"`   // 请求类型
	Params TCPParams `json:"params"` // 对应请求类型的附加信息
}

// 不同类型的负载结构
type TCPParams struct {
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

//----------------------------------------------------------------------

type MT5MarginCall struct {
	Login       uint64  `json:"login"  msgpack:"login"`
	Equity      float64 `json:"equity"  msgpack:"equity"`           //净值
	MarginLevel float64 `json:"marginLevel"  msgpack:"marginLevel"` //保证金率
}

type MT5StopOut struct {
	Login    uint64  `json:"login"  msgpack:"login"`
	SOLevel  float64 `json:"soLevel"  msgpack:"soLevel"`
	SOEquity float64 `json:"soEquity"  msgpack:"soEquity"`
	SOMargin float64 `json:"soMargin"  msgpack:"soMargin"`
}

//----------------------------------------------------------------------

type MT5Tick struct {
	Symbol string `json:"symbol"  msgpack:"symbol"`
	AskE8  int64  `json:"ask" msgpack:"ask"`   //都是扩大了 10e8倍
	BidE8  int64  `json:"bid" msgpack:"bid"`   //都是扩大了 10e8倍
	LastE8 int64  `json:"last" msgpack:"last"` //都是扩大了 10e8倍
	Volume uint64 `json:"volume" msgpack:"volume"`
	Time   int64  `json:"time" msgpack:"time"` //unix时间戳(ms毫秒)
}

//-----------------------------------------------------------------------

type Mt5Order struct {
	OrderId  uint64  `json:"orderId"  msgpack:"orderId"`
	Symbol   string  `json:"symbol"  msgpack:"symbol"`
	Login    uint64  `json:"login"  msgpack:"login"`
	Volume   float64 `json:"volume"  msgpack:"volume"`
	Direct   int     `json:"direct"  msgpack:"direct"`     //0-buy, 1-sell
	OpenTime int64   `json:"openTime"  msgpack:"openTime"` //unix时间戳（看是否平仓了）,决策是否在指定时间段内 (一个过滤条件)
	Price    float64 `json:"price"  msgpack:"price"`
	Comment  string  `json:"comment"  msgpack:"comment"`
}

//-----------------------------------------------------------------------

type Mt5Position struct {
	Operation int `json:"operation"` // 1-add, 2-remove, 3-modify

	PositionId uint64  `json:"positionId"  msgpack:"positionId"`
	Symbol     string  `json:"symbol"  msgpack:"symbol"`
	Login      uint64  `json:"login"  msgpack:"login"`
	Volume     float64 `json:"volume"  msgpack:"volume"`
	Direct     int     `json:"direct"  msgpack:"direct"`     //0-buy, 1-sell
	OpenTime   int64   `json:"openTime"  msgpack:"openTime"` //unix时间戳（看是否平仓了）,决策是否在指定时间段内 (一个过滤条件)
	OpenPrice  float64 `json:"openPrice"  msgpack:"openPrice"`
}

//-----------------------------------------------------------------------

type Mt5Deal struct {
	DealId     uint64  `json:"dealId"  msgpack:"dealId"`
	Symbol     string  `json:"symbol"  msgpack:"symbol"`
	Login      uint64  `json:"login"  msgpack:"login"`
	Volume     float64 `json:"volume"  msgpack:"volume"`
	Direct     int     `json:"direct"  msgpack:"direct"` //0-buy, 1-sell
	Entry      int     `json:"entry"  msgpack:"entry"`   //0-ENTRY_IN 开仓, 1-ENTRY_OUT 平仓
	Time       int64   `json:"time"  msgpack:"time"`
	Price      float64 `json:"price"  msgpack:"price"`           //执行价格
	CloseTime  int64   `json:"closeTime"  msgpack:"closeTime"`   //平仓的话,这里是平仓时间
	PositionId uint64  `json:"positionId"  msgpack:"positionId"` //平仓单的话,这里就是对应的positionId. 开仓单的话这里为0
}

//-----------------------------------------------------------------------
