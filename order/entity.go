package order

type InitParams struct {
	Address string `json:"address" mapstructure:"address" config:"address" yaml:"address"` // http://ip:port这样的地址
}

// -----------------------------------

// 普通开仓单
type OpenPositionRequest struct {
	//required
	Login  uint64        `json:"login"` //下单人
	Lots   string        `json:"lots"`  // lots手数  float64
	Symbol string        `json:"symbol"`
	Type   MtRequestType `json:"type"` // 只支持类型: 0-buy, 1-sell

	//option
	Comment string `json:"comment,omitempty"`
	Sl      string `json:"sl,omitempty"` //float64 (为了避免精度损失)
	Tp      string `json:"tp,omitempty"` //float64 (为了避免精度损失)
}

// 平通平仓单
type ClosePositionRequest struct {
	Lots   string `json:"lots,omitempty"` // lots手数  float64
	Ticket int    `json:"ticket"`         //是要平掉的order/position的id (通过这个可以拿到symbol和login)

	//option
	Comment string `json:"comment,omitempty"`
}

// 一键平仓
type CloseAllPositionsRequest struct {
	Login uint64 `json:"login"`
	//option
	Comment string `json:"comment,omitempty"`
}

// 挂单
type PlacePendingOrderRequest struct {
	//required
	Login          uint64        `json:"login"` //下单人
	Symbol         string        `json:"symbol"`
	Lots           string        `json:"lots"`             // lots手数 float64
	Type           MtRequestType `json:"type"`             // 只支持如下6种类型: 2-OP_BUY_LIMIT, 3-OP_SELL_LIMIT, 4-OP_BUY_STOP, 5-OP_SELL_STOP，6-OP_BUY_STOP_LIMIT，7-OP_SELL_STOP_LIMIT
	Price          float64       `json:"price"`            // 挂单的价格(手动指定的)
	ExpireTimeType MtOrderTime   `json:"expire_time_type"` // 不传默认gtc

	//option
	ExpireTime   int64  `json:"expire_time"`   //到期时间,unix时间戳,传0则是不限制,
	TriggerPrice string `json:"trigger_price"` //float64 只有6/7生效.  只有 Set the price, at which a Limit order is placed when the Stop Limit order triggers.

	//option
	Comment string `json:"comment,omitempty"`
	Sl      string `json:"sl,omitempty"` //float64 (为了避免精度损失)
	Tp      string `json:"tp,omitempty"` //float64 (为了避免精度损失)
}

// 修改挂单
// type类型、volume 和 symbol等禁止修改. 只能修改price、time和comment
type ModifyPendingOrderRequest struct {
	//required
	Ticket int `json:"ticket"` //是要修改的 order 的id, 通过它可以拿到: symbol, login, type,

	//option
	Price        string `json:"price"`         //float64     // 新的价格 (不改就还是以前的价格)
	TriggerPrice string `json:"trigger_price"` //float64  new - 只有6/7生效.  只有 Set the price, at which a Limit order is placed when the Stop Limit order triggers.
	Sl           string `json:"sl,omitempty"`  //float64 (为了避免精度损失)
	Tp           string `json:"tp,omitempty"`  //float64 (为了避免精度损失)

	ExpireTimeType MtOrderTime `json:"expire_time_type"` // 不传默认gtc
	ExpireTime     int64       `json:"expire_time"`      //到期时间,unix时间戳,传0则是不限制,

	//option
	Comment string `json:"comment,omitempty"` //该modify操作的备注
}

// 关掉挂单
// type类型、volume 和 symbol等禁止修改. 只能修改price、time和comment
type RemovePendingOrderRequest struct {
	//required
	Ticket int `json:"ticket"` //是要删掉的 order 的id, 通过它可以拿到: symbol, login, type等信息

	//option
	Comment string `json:"comment,omitempty"` //该modify操作的备注
}

//------------------------------------------------------------------------

type CommonResp struct {
	Code    int         `json:"code"`           //错误码 0是成功
	Success bool        `json:"success"`        //是否成功
	Message string      `json:"message"`        //错误信息
	Data    interface{} `json:"data,omitempty"` //数据
}
