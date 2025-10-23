package order

type MtRequestAction int
type MtRequestDirection int

const (
	MtRequestActionOpen  MtRequestAction = 0 //开仓 open
	MtRequestActionClose MtRequestAction = 1 //平仓 close

	//如下跟mt4的枚举定义一致 ( https://support.metaquotes.net/en/docs/mt5/api/reference_trading/trading_order/imtorder/imtorder_enum#enordertype )
	//如下也跟mt5的枚举定义一致:  https://support.metaquotes.net/en/docs/mt5/api/reference_trading/trading_order/imtorder/imtorder_enum#enordertype
	MtRequestDirectionBuy       MtRequestDirection = 0 //buy mt5也如此
	MtRequestDirectionSell      MtRequestDirection = 1 //sell mt5也如此
	MtRequestDirectionBuyLimit  MtRequestDirection = 2
	MtRequestDirectionSellLimit MtRequestDirection = 3
	MtRequestDirectionBuyStop   MtRequestDirection = 4
	MtRequestDirectionStop      MtRequestDirection = 5
)
