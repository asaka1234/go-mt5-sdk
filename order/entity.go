package order

type InitParams struct {
	OpenUrl  string `json:"openUrl" mapstructure:"openUrl" config:"openUrl" yaml:"openUrl"`     //开仓url
	CloseUrl string `json:"closeUrl" mapstructure:"closeUrl" config:"closeUrl" yaml:"closeUrl"` //平仓url
}

// -----------------------------------

// 开仓请求
type OpenRequest struct {
	Amount float64 `json:"amount"` //must 金额(不需要做单位转换) 只是THB泰铢
	Ref1   string  `json:"ref1"`   //must 放业务自己的orderNo (只能是数字/字母) The ref1 use bank format can’t more than 18 digit
	//option
	Ref2 string `json:"ref2"`
	Ref3 string `json:"ref3"` //放customer name
	Ref4 string `json:"ref4"` // 这个用来做签名, 是amount/ref1/authkey的一个md5签名的截断值(18位)
}

type OpenResponse struct {
	Error    string `json:"error"`     //如果返回错误，则有该字段
	RespCode int    `json:"resp_code"` //201是正确
	RespMsg  string `json:"resp_msg"`
	Data     struct {
		Method            string  `json:"method"`
		Channel           string  `json:"channel"`
		Ref1              string  `json:"ref1"`
		Ref2              string  `json:"ref2"`
		Ref3              string  `json:"ref3"`
		Ref4              string  `json:"ref4"` // 这个用来做签名, 是amount/ref1/authkey的一个md5签名的截断值(18位)
		Amount            float64 `json:"amount"`
		Currency          string  `json:"currency"` //THB
		Location          string  `json:"location"`
		Device            string  `json:"device"`
		PartnerCode       string  `json:"partner_code"`
		CodeType          string  `json:"code_type"` //ThaiQRCode
		CodeImage         string  `json:"code_image"`
		CodeURL           string  `json:"code_url"`
		TransID           string  `json:"trans_id"`
		WalletCode        string  `json:"wallet_code"`
		Ref               string  `json:"_ref"`
		StoreID           string  `json:"store_id"`
		TerminalID        string  `json:"terminal_id"`
		MobileNo          string  `json:"mobile_no"`
		ExpiredDate       string  `json:"expired_date"`
		CreatedDate       string  `json:"created_date"`
		CodeImageTemplate string  `json:"code_image_template"`
	} `json:"data"`
}
