package direct

type InitParams struct {
	Address string `json:"address" mapstructure:"address" config:"address" yaml:"address"` // http://ip:port这样的地址
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

//------------------------------------------------------------------------

type CommonResp struct {
	Code    int    `json:"code"`    //错误码 0是成功
	Success bool   `json:"success"` //是否成功
	Message string `json:"message"` //错误信息
	//Data    interface{} `json:"data,omitempty"` //数据
}

type ListSymbolResp struct {
	CommonResp `json:",inline"`
	Data       []MT5SymbolBase `json:"data,omitempty"` //数据
}

type MT5SymbolBase struct {
	//-----------base----------------------------
	Symbol    string `json:"symbol"`
	Desc      string `json:"desc"`
	TradeMode uint   `json:"trade_mode"` //交易模式,比如:long_only https://support.metaquotes.net/en/docs/mt5/api/config_symbol/imtconsymbol/imtconsymbol_enum#entrademode
	Category  string `json:"category"`   //分组(enum)

	//-----------交易------------------------
	Digit        uint    `json:"digit"`         //symbol的精度
	GTCMode      uint    `json:"gtc_mode"`      //直到挂单取消
	ContractSize float64 `json:"contract_size"` //合约量
	VolumeMin    float64 `json:"volume_min"`    //最小下单手数
	VolumeMax    float64 `json:"volume_max"`    //最大下单手数
	CalcMode     uint    `json:"calc_mode"`     //利润计算

	//-----------保证金------------------------
	MarginInitial      float64 `json:"margin_initial"`       //初始保证金
	MarginHedged       float64 `json:"margin_hedged"`        //保证金对冲
	MarginRateCurrency float64 `json:"margin_rate_currency"` //保证金百分比
	CurrencyMargin     string  `json:"currency_margin"`      //保证金货币

	//-------------利息----------------------
	SwapMode  uint    `json:"swap_mode"`  //库存费类型
	SwapLong  float64 `json:"swap_long"`  //买入库存费
	SwapShort float64 `json:"swap_short"` //卖出库存费
	Swap3Day  uint    `json:"swap_3_day"` //3日库存费

	//-------------交易时间----------------------
	SessionTrade []SessionTrade `json:"session_trade"` //每周七天的一个列表
}

type SessionTrade struct {
	Wday     uint     `json:"wday"` //Day of the week. The day is specified by a value 0 (Sunday) to 6 (Saturday).
	Sessions []string `json:"sessions"`
}

//-------------------------------------

type TickReviewResp struct {
	CommonResp `json:",inline"`
	Data       []MT5Tick `json:"data,omitempty"` //数据
}

type MT5Tick struct {
	Symbol string `json:"symbol"`
	AskE8  int64  `json:"ask"`  //都是扩大了 10e8倍
	BidE8  int64  `json:"bid"`  //都是扩大了 10e8倍
	LastE8 int64  `json:"last"` //都是扩大了 10e8倍
	Volume uint64 `json:"volume"`
	Time   int64  `json:"time"` //unix时间戳(ms毫秒)
}
