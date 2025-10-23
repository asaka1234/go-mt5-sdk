package order

type MtRequestType int

// https://support.metaquotes.net/en/docs/mt5/api/reference_trading/trading_order/imtorder/imtorder_enum#enordertype
const (
	MtRequestTypeBuy       MtRequestType = 0 //buy mt5也如此
	MtRequestTypeSell      MtRequestType = 1 //sell mt5也如此
	MtRequestTypeBuyLimit  MtRequestType = 2
	MtRequestTypeSellLimit MtRequestType = 3
	MtRequestTypeBuyStop   MtRequestType = 4
	MtRequestTypeStop      MtRequestType = 5
)
