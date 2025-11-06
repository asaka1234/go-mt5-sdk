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
	Status    string      `json:"status" msgpack:"status"`                       //ok是成功,其他是错误描述
	Type      string      `json:"type" msgpack:"type"`                           //映射:请求类型
	Payload   interface{} `json:"payload,omitempty" msgpack:"payload,omitempty"` // 响应负载数据
	Timestamp int64       `json:"timestamp" msgpack:"timestamp"`
}

//----------------------------------------------------------------------

type MT5MarginCall struct {
	Login       uint64  `json:"login"  msgpack:"login"`
	Equity      float64 `json:"equity"  msgpack:"equity"`             //净值
	MarginLevel float64 `json:"margin_level"  msgpack:"margin_level"` //保证金率
}

type MT5StopOut struct {
	Login    uint64  `json:"login"  msgpack:"login"`
	SOLevel  float64 `json:"so_level"  msgpack:"so_level"`
	SOEquity float64 `json:"so_equity"  msgpack:"so_equity"`
	SOMargin float64 `json:"so_margin"  msgpack:"so_margin"`
}

//----------------------------------------------------------------------

type MT5Tick struct {
	Symbol string `json:"symbol"  msgpack:"symbol"`
	AskE8  int64  `json:"ask_e8" msgpack:"ask_e8"`   //都是扩大了 10e8倍
	BidE8  int64  `json:"bid_e8" msgpack:"bid_e8"`   //都是扩大了 10e8倍
	LastE8 int64  `json:"last_e8" msgpack:"last_e8"` //都是扩大了 10e8倍
	Volume uint64 `json:"volume" msgpack:"volume"`
	Time   int64  `json:"time" msgpack:"time"` //unix时间戳(ms毫秒)
}

//-----------------------------------------------------------------------

type MTOrderExtra struct {
	Operation uint `json:"operation"  msgpack:"operation"` //1-add, 2-remove, 3-modify
	MTOrder   `json:",inline"  msgpack:",inline"`
}

type MTOrder struct {
	Login  uint64 `json:"login"  msgpack:"login"`
	Ticket uint64 `json:"ticket"  msgpack:"ticket"` //order_id
	Symbol string `json:"symbol"  msgpack:"symbol"`
	//State        uint    //1是挂单  ORDER_STATE_PLACED
	TimeSetup    int64   `json:"time_setup"  msgpack:"time_setup"`       //下单时间
	Type         uint    `json:"type"  msgpack:"type"`                   //0-buy, 1-sell,2-buy limit ,3-sell limit, 4-buy stop, 5-sell stop, 6-buy stop limit, 7-sell stop limit,
	PriceOrder   float64 `json:"price_order"  msgpack:"price_order"`     //下单价格 (stop/limit的价格)
	PriceTrigger float64 `json:"price_trigger"  msgpack:"price_trigger"` //触发价格（stop limit 单）
	PriceSL      float64 `json:"price_sl"  msgpack:"price_sl"`
	PriceTP      float64 `json:"price_tp"  msgpack:"price_tp"`
	Volume       float64 `json:"volume"  msgpack:"volume"` //lots
	RateMargin   float64 `json:"rate_margin"  msgpack:"rate_margin"`
}

//-----------------------------------------------------------------------

type MTPositionExtra struct {
	Operation  uint `json:"operation"  msgpack:"operation"` //1-add, 2-remove, 3-modify
	MTPosition `json:",inline"  msgpack:",inline"`
}

type MTPosition struct {
	Login          uint64  `json:"login"  msgpack:"login"`
	Ticket         uint64  `json:"ticket"  msgpack:"ticket"` //position_id
	Symbol         string  `json:"symbol"  msgpack:"symbol"`
	Action         uint    `json:"action"  msgpack:"action"`         // 0-buy, 1-sell
	PriceOpen      float64 `json:"price_open"  msgpack:"price_open"` //开仓价
	PriceSL        float64 `json:"price_sl"  msgpack:"price_sl"`
	PriceTP        float64 `json:"price_tp"  msgpack:"price_tp"`
	RateMargin     float64 `json:"rate_margin"  msgpack:"rate_margin"`
	RateProfit     float64 `json:"rate_profit"  msgpack:"rate_profit"`
	Volume         float64 `json:"volume"  msgpack:"volume"` //lots
	Profit         float64 `json:"profit"  msgpack:"profit"`
	Storage        float64 `json:"storage"  msgpack:"storage"`
	ActivationMode uint    `json:"activation_mode"  msgpack:"activation_mode"` //1-sl, 2-tp, 3-so
	ActivationTime int64   `json:"activation_time"  msgpack:"activation_time"` //unix时间戳(s)
	TimeCreate     int64   `json:"time_create"  msgpack:"time_create"`         //unix时间戳(s)
}

//-----------------------------------------------------------------------

type Mt5DealExtra struct {
	Operation uint `json:"operation"  msgpack:"operation"` //1-add, 2-remove, 3-modify
	Mt5Deal   `json:",inline"  msgpack:",inline"`
}

type Mt5Deal struct {
	DealId     uint64  `json:"deal_id"  msgpack:"deal_id"`
	PositionId uint64  `json:"position_id"  msgpack:"position_id"`
	Symbol     string  `json:"symbol"  msgpack:"symbol"`
	Login      uint64  `json:"login"  msgpack:"login"`
	Volume     float64 `json:"volume"  msgpack:"volume"`
	Entry      int     `json:"entry"  msgpack:"entry"`   //0-ENTRY_IN 开仓, 1-ENTRY_OUT 平仓
	Action     int     `json:"action"  msgpack:"action"` //
	Time       int64   `json:"time"  msgpack:"time"`
	Price      float64 `json:"price"  msgpack:"price"` //执行价格
	PriceSL    float64 `json:"price_sl"  msgpack:"price_sl"`
	PriceTP    float64 `json:"price_tp"  msgpack:"price_tp"`
	Profit     float64 `json:"profit"  msgpack:"profit"` //profit
	RateMargin float64 `json:"rate_margin"  msgpack:"rate_margin"`
	RateProfit float64 `json:"rate_profit"  msgpack:"rate_profit"`
	Storage    float64 `json:"storage"  msgpack:"storage"` //swap
}

//-----------------------------------------------------------------------
