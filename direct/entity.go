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
	Symbol   string `json:"symbol"`
	Digit    uint   `json:"digit"` //symbol的精度
	Desc     string `json:"desc"`
	Category string `json:"category"` //分组(enum)
	//-----------currency-----------------------
	CurrencyBase        string `json:"currency_base"`
	CurrencyBaseDigit   uint   `json:"currency_base_digit"`
	CurrencyProfit      string `json:"currency_profit"`
	CurrencyProfitDigit uint   `json:"currency_profit_digit"`
	CurrencyMargin      string `json:"currency_margin"` //保证金货币
	CurrencyMarginDigit uint   `json:"currency_margin_digit"`
	//-----------交易------------------------
	ContractSize float64 `json:"contract_size"` //合约量
	CalcMode     uint    `json:"calc_mode"`     //利润/swap计算
	TradeMode    uint    `json:"trade_mode"`    //交易模式,比如:long_only https://support.metaquotes.net/en/docs/mt5/api/config_symbol/imtconsymbol/imtconsymbol_enum#entrademode
	GTCMode      uint    `json:"gtc_mode"`      //直到挂单取消
	VolumeMin    float64 `json:"volume_min"`    //最小下单手数
	VolumeMax    float64 `json:"volume_max"`    //最大下单手数
	VolumeStep   float64 `json:"volume_step"`   //下单步长

	//-----------保证金------------------------
	MarginInitial      float64 `json:"margin_initial"`       //初始保证金
	MarginHedged       float64 `json:"margin_hedged"`        //保证金对冲
	MarginRateCurrency float64 `json:"margin_rate_currency"` //保证金百分比

	MarginRateInitBuy  float64 `json:"margin_rate_init_buy"`  //percentage初始保证金比例针对buy方向
	MarginRateInitSell float64 `json:"margin_rate_init_sell"` //percentage针对sell方向
	MarginRateMainBuy  float64 `json:"margin_rate_main_buy"`  //percentage维持保证金比例针对buy方向
	MarginRateMainSell float64 `json:"margin_rate_main_sell"` //percentage针对sell方向

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

//-------------------------------------

type UserCreateReq struct {
	Uid      int64 `json:"uid,omitempty"`      //yubit的uid,用来设置备注
	Internal uint  `json:"internal,omitempty"` //看是否是内部测试账户(1是,2否)
	Leverage uint  `json:"leverage,omitempty"` //杠杆.默认500
}

type UserCreateResp struct {
	CommonResp `json:",inline"`
	Data       Mt5User `json:"data,omitempty"` //数据
}

type Mt5User struct {
	Login        uint64 `json:"login"` //account的login
	MasterPass   string `json:"master_pass"`
	InvestorPass string `json:"investor_pass"`
}

//-------------------------------------

type BalanceOperationReq struct {
	Login   uint64  `json:"login,omitempty"`   //mt5的login
	Balance float64 `json:"balance,omitempty"` //上账多少,支持浮点数和负数
}

type BalanceOperationResp struct {
	CommonResp `json:",inline"`
	Data       MtRecharge `json:"data,omitempty"` //数据
}

type MtRecharge struct {
	DealId uint64 `json:"deal_id"` //充提的deal id
}

//-----------------------------------------------

type UserAccountDetailResp struct {
	CommonResp `json:",inline"`
	Data       MTUserAccount `json:"data,omitempty"` //数据
}

type MTUserAccount struct {
	Login          uint64  `json:"login"`  //当前要操作的account的login
	Balance        float64 `json:"symbol"` //余额
	Margin         float64 `json:"margin"` //已用保证金
	MarginFree     float64 `json:"margin_free"`
	MarginLevel    float64 `json:"margin_level"`
	MarginLeverage uint    `json:"margin_leverage"` //杠杆
	Equity         float64 `json:"equity"`
	Storage        float64 `json:"storage"`
	Floating       float64 `json:"floating"`
}

//-----------------------------------------------

type ListPositionResp struct {
	CommonResp `json:",inline"`
	Data       []*MTPosition `json:"data,omitempty"` //数据
}
type MTPosition struct {
	Login          uint64  `json:"login"`
	Ticket         uint64  `json:"ticket"` //position_id
	Symbol         string  `json:"symbol"`
	Action         uint    `json:"action"`     // 0-buy, 1-sell
	PriceOpen      float64 `json:"price_open"` //开仓价
	PriceSL        float64 `json:"price_sl"`
	PriceTP        float64 `json:"price_tp"`
	RateMargin     float64 `json:"rate_margin"`
	RateProfit     float64 `json:"rate_profit"`
	Volume         float64 `json:"volume"` //lots
	Profit         float64 `json:"profit"`
	Storage        float64 `json:"storage"`
	ActivationMode uint    `json:"activation_mode"` //1-sl, 2-tp, 3-so
	ActivationTime int64   `json:"activation_time"` //unix时间戳(s)
	TimeCreate     int64   `json:"time_create"`     //unix时间戳(s)
}

//-----------------------------------------------

type ListPendingOrderResp struct {
	CommonResp `json:",inline"`
	Data       []*MTOrder `json:"data,omitempty"` //数据
}

type MTOrder struct {
	Login        uint64  `json:"login"`
	Ticket       uint64  `json:"ticket"` //order_id
	Symbol       string  `json:"symbol"`
	State        uint    `json:"state"`         //1是挂单  ORDER_STATE_PLACED
	TimeSetup    int64   `json:"time_setup"`    //下单时间
	Type         uint    `json:"type"`          //0-buy, 1-sell,2-buy limit ,3-sell limit, 4-buy stop, 5-sell stop, 6-buy stop limit, 7-sell stop limit,
	PriceOrder   float64 `json:"price_order"`   //下单价格 (stop/limit的价格)
	PriceTrigger float64 `json:"price_trigger"` //触发价格（stop limit 单）
	PriceSL      float64 `json:"price_sl"`
	PriceTP      float64 `json:"price_tp"`
	Volume       float64 `json:"volume"` //lots
	RateMargin   float64 `json:"rate_margin"`
}

//------------------------------------------------------

type GetOrderResp struct {
	CommonResp `json:",inline"`
	Data       MTOrder `json:"data,omitempty"` //数据
}

type GetPositionResp struct {
	CommonResp `json:",inline"`
	Data       MTPosition `json:"data,omitempty"` //数据
}
