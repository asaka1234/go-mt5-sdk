package order

import "github.com/shopspring/decimal"

type InitParams struct {
	OpenUrl  string `json:"openUrl" mapstructure:"openUrl" config:"openUrl" yaml:"openUrl"`     //开仓url
	CloseUrl string `json:"closeUrl" mapstructure:"closeUrl" config:"closeUrl" yaml:"closeUrl"` //平仓url
}

// -----------------------------------

// 开仓请求
type OpenOrderReq struct {
	Login   uint64             `json:"login"` //下单人
	Lots    float64            `json:"lots"`  // lots手数
	Symbol  string             `json:"symbol"`
	Type    MtRequestDirection `json:"type"`  // 方向: 0-buy, 1-sell,  方向//2-OP_BUY_LIMIT, 3-OP_SELL_LIMIT, 4-OP_BUY_STOP, 5-OP_SELL_STOP
	Price   decimal.Decimal    `json:"price"` // 如果type是0/1 则不需要传该参数, 程序自动计算
	Comment string             `json:"comment,omitempty"`
	Sl      float64            `json:"sl,omitempty"`
	Tp      float64            `json:"tp,omitempty"`
}

type CloseOrderReq struct {
	Lots   float64 `json:"lots,omitempty"`
	Ticket int     `json:"ticket"` //是要平掉的order/position的id (通过这个可以拿到symbol和login)

	//-----以下非入参----------------------
	Symbol string             `json:"symbol"`
	Login  uint64             `json:"login"`
	Type   MtRequestDirection `json:"type"` //这个是自己填的,并不是参数传递的!!!!    方向: 0-buy, 1-sell, 3-OP_BUY_LIMIT, 4-OP_SELL_LIMIT, 5-OP_BUY_STOP, 6-OP_SELL_STOP
}

//------------------------------------------------------------------------

type CommonResp struct {
	Code    int         `json:"code"`           //错误码 0是成功
	Success bool        `json:"success"`        //是否成功
	Message string      `json:"message"`        //错误信息
	Data    interface{} `json:"data,omitempty"` //数据
}
