package order

type InitParams struct {
	OpenUrl    string `json:"openUrl" mapstructure:"openUrl" config:"openUrl" yaml:"openUrl"`             //开仓url
	CloseUrl   string `json:"closeUrl" mapstructure:"closeUrl" config:"closeUrl" yaml:"closeUrl"`         //平仓url
	PendingUrl string `json:"pendingUrl" mapstructure:"pendingUrl" config:"pendingUrl" yaml:"pendingUrl"` //挂单url
}

// -----------------------------------

// 普通开仓单
type OpenPositionRequest struct {
	//required
	Login  uint64        `json:"login"` //下单人
	Lots   float64       `json:"lots"`  // lots手数
	Symbol string        `json:"symbol"`
	Type   MtRequestType `json:"type"` // 只支持类型: 0-buy, 1-sell

	//option
	Comment string  `json:"comment,omitempty"`
	Sl      float64 `json:"sl,omitempty"`
	Tp      float64 `json:"tp,omitempty"`
}

// 平通平仓单
type ClosePositionRequest struct {
	Lots   float64 `json:"lots,omitempty"`
	Ticket int     `json:"ticket"` //是要平掉的order/position的id (通过这个可以拿到symbol和login)

	//option
	Comment string `json:"comment,omitempty"`
}

// 挂单
type PlacePendingOrderRequest struct {
	//required
	Login          uint64        `json:"login"` //下单人
	Symbol         string        `json:"symbol"`
	Lots           float64       `json:"lots"`             // lots手数
	Type           MtRequestType `json:"type"`             // 只支持如下6种类型: 2-OP_BUY_LIMIT, 3-OP_SELL_LIMIT, 4-OP_BUY_STOP, 5-OP_SELL_STOP，6-OP_BUY_STOP_LIMIT，7-OP_SELL_STOP_LIMIT
	Price          float64       `json:"price"`            // 挂单的价格(手动指定的)
	ExpireTimeType MtOrderTime   `json:"expire_time_type"` // 不传默认gtc

	//option
	ExpireTime   int64   `json:"expire_time"`   //到期时间,unix时间戳,传0则是不限制,
	TriggerPrice float64 `json:"trigger_price"` //只有6/7生效.  只有 Set the price, at which a Limit order is placed when the Stop Limit order triggers.

	//option
	Comment string  `json:"comment,omitempty"`
	Sl      float64 `json:"sl,omitempty"`
	Tp      float64 `json:"tp,omitempty"`
}

//------------------------------------------------------------------------

type CommonResp struct {
	Code    int         `json:"code"`           //错误码 0是成功
	Success bool        `json:"success"`        //是否成功
	Message string      `json:"message"`        //错误信息
	Data    interface{} `json:"data,omitempty"` //数据
}
