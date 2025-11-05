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
	AskE8  int64  `json:"ask" msgpack:"ask"`   //都是扩大了 10e8倍
	BidE8  int64  `json:"bid" msgpack:"bid"`   //都是扩大了 10e8倍
	LastE8 int64  `json:"last" msgpack:"last"` //都是扩大了 10e8倍
	Volume uint64 `json:"volume" msgpack:"volume"`
	Time   int64  `json:"time" msgpack:"time"` //unix时间戳(ms毫秒)
}

//-----------------------------------------------------------------------

type MTOrderExtra struct {
	Operation uint `json:"operation"` //1-add, 2-remove, 3-modify
	MTOrder   `json:",inline"`
}

type MTOrder struct {
	Login  uint64 `json:"login"`
	Ticket uint64 `json:"ticket"` //order_id
	Symbol string `json:"symbol"`
	//State        uint    //1是挂单  ORDER_STATE_PLACED
	TimeSetup    int64   `json:"time_setup"`    //下单时间
	Type         uint    `json:"type"`          //0-buy, 1-sell,2-buy limit ,3-sell limit, 4-buy stop, 5-sell stop, 6-buy stop limit, 7-sell stop limit,
	PriceOrder   float64 `json:"price_order"`   //下单价格 (stop/limit的价格)
	PriceTrigger float64 `json:"price_trigger"` //触发价格（stop limit 单）
	PriceSL      float64 `json:"price_sl"`
	PriceTP      float64 `json:"price_tp"`
	Volume       float64 `json:"volume"` //lots
}

//-----------------------------------------------------------------------

type MTPositionExtra struct {
	Operation  uint `json:"operation"` //1-add, 2-remove, 3-modify
	MTPosition `json:",inline"`
}

type MTPosition struct {
	Login          uint64  `json:"login"`
	Ticket         uint64  `json:"ticket"` //position_id
	Symbol         string  `json:"symbol"`
	Action         uint    `json:"action"`     // 0-buy, 1-sell
	PriceOpen      float64 `json:"price_open"` //开仓价
	PriceSL        float64 `json:"price_sl"`
	PriceTP        float64 `json:"price_tp"`
	Volume         float64 `json:"volume"` //lots
	Profit         float64 `json:"profit"`
	Storage        float64 `json:"storage"`
	ActivationMode uint    `json:"activation_mode"` //1-sl, 2-tp, 3-so
	ActivationTime int64   `json:"activation_time"` //unix时间戳(s)
	TimeCreate     int64   `json:"time_create"`     //unix时间戳(s)
}

//-----------------------------------------------------------------------

type Mt5DealExtra struct {
	Operation uint `json:"operation"` //1-add, 2-remove, 3-modify
	Mt5Deal   `json:",inline"`
}

type Mt5Deal struct {
	DealId     uint64  `json:"deal_id"`
	PositionId uint64  `json:"position_id"`
	Symbol     string  `json:"symbol"`
	Login      uint64  `json:"login"`
	Volume     float64 `json:"volume"`
	Entry      int     `json:"entry"`  //0-ENTRY_IN 开仓, 1-ENTRY_OUT 平仓
	Action     int     `json:"action"` //
	Time       int64   `json:"time"`
	Price      float64 `json:"price"` //执行价格
	PriceSL    float64 `json:"price_sl"`
	PriceTP    float64 `json:"price_tp"`
	Profit     float64 `json:"profit"`  //profit
	Storage    float64 `json:"storage"` //swap
}

//-----------------------------------------------------------------------
