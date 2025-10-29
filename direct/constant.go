package direct

type MtRequestType int
type MtOrderTime uint

// https://support.metaquotes.net/en/docs/mt5/api/reference_trading/trading_order/imtorder/imtorder_enum#enordertype
const (
	MtRequestTypeBuy       MtRequestType = 0 //buy mt5也如此
	MtRequestTypeSell      MtRequestType = 1 //sell mt5也如此
	MtRequestTypeBuyLimit  MtRequestType = 2
	MtRequestTypeSellLimit MtRequestType = 3
	MtRequestTypeBuyStop   MtRequestType = 4
	MtRequestTypeStop      MtRequestType = 5
)

// -----------------------------
// https://support.metaquotes.net/en/docs/mt5/api/reference_trading/trading_order/imtorder/imtorder_enum#enordertime
// 挂单挂到什么时候?
const (
	MtOrderTimeGTC          MtOrderTime = 0 //gtc
	MtOrderTimeDay          MtOrderTime = 1 //
	MtOrderTimeSpecified    MtOrderTime = 0 //
	MtOrderTimeSpecifiedDay MtOrderTime = 1 //
)
