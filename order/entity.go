package order

type InitParams struct {
	OpenUrl  string `json:"openUrl" mapstructure:"openUrl" config:"openUrl" yaml:"openUrl"`     //开仓url
	CloseUrl string `json:"closeUrl" mapstructure:"closeUrl" config:"closeUrl" yaml:"closeUrl"` //平仓url
}

// -----------------------------------

// 普通开仓单
type OpenPositionRequest struct {
	//required
	Login  uint64        `json:"login"` //下单人
	Lots   float64       `json:"lots"`  // lots手数
	Symbol string        `json:"symbol"`
	Type   MtRequestType `json:"type"` // 类型: 0-buy, 1-sell, 2-OP_BUY_LIMIT, 3-OP_SELL_LIMIT, 4-OP_BUY_STOP, 5-OP_SELL_STOP

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

//------------------------------------------------------------------------

type CommonResp struct {
	Code    int         `json:"code"`           //错误码 0是成功
	Success bool        `json:"success"`        //是否成功
	Message string      `json:"message"`        //错误信息
	Data    interface{} `json:"data,omitempty"` //数据
}
