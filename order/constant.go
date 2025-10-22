package order

// 返回值的code枚举

type One2payDepositStatus int

const (
	One2payDepositStatusSuccess One2payDepositStatus = iota + 200
	One2payDepositStatusCreated
)

func (s One2payDepositStatus) Code() int {
	return int(s)
}

func (s One2payDepositStatus) Name() string {
	switch s {
	case One2payDepositStatusSuccess:
		return "Success"
	case One2payDepositStatusCreated:
		return "Created"
	default:
		return "Unknown"
	}
}

func (s One2payDepositStatus) Desc() string {
	switch s {
	case One2payDepositStatusSuccess:
		return "Success"
	case One2payDepositStatusCreated:
		return "Created"
	default:
		return "Unknown"
	}
}

func (s One2payDepositStatus) Eq(value interface{}) bool {
	switch v := value.(type) {
	case int:
		return s.Code() == v
	case One2payDepositStatus:
		return s == v
	default:
		return false
	}
}

//-------------------------------------------------------

type One2PayBankCode struct {
	Code string
	Name string
	Desc string
}

// BankCodeOne2PayEnum represents the collection of bank codes
var One2PayBankCodeEnum = struct {
	BBL   One2PayBankCode
	KBANK One2PayBankCode
	KTB   One2PayBankCode
	TTB   One2PayBankCode
	SCB   One2PayBankCode
	BAY   One2PayBankCode
	KKP   One2PayBankCode
	CIMBT One2PayBankCode
	TISCO One2PayBankCode
	UOBT  One2PayBankCode
	TCD   One2PayBankCode
	LHFG  One2PayBankCode
	ICBCT One2PayBankCode
	SME   One2PayBankCode
	BAAC  One2PayBankCode
	EXIM  One2PayBankCode
	GSB   One2PayBankCode
	GHB   One2PayBankCode
	ISBT  One2PayBankCode
}{
	BBL:   One2PayBankCode{"002", "BBL", "BBL"},
	KBANK: One2PayBankCode{"004", "KBANK", "KBANK"},
	KTB:   One2PayBankCode{"006", "KTB", "KTB"},
	TTB:   One2PayBankCode{"011", "TTB", "TTB"},
	SCB:   One2PayBankCode{"014", "SCB", "SCB"},
	BAY:   One2PayBankCode{"025", "BAY", "BAY"},
	KKP:   One2PayBankCode{"069", "KKP", "KKP"},
	CIMBT: One2PayBankCode{"022", "CIMBT", "CIMBT"},
	TISCO: One2PayBankCode{"067", "TISCO", "TISCO"},
	UOBT:  One2PayBankCode{"024", "UOBT", "UOBT"},
	TCD:   One2PayBankCode{"071", "TCD", "TCD"},
	LHFG:  One2PayBankCode{"073", "LHFG", "LHFG"},
	ICBCT: One2PayBankCode{"070", "ICBCT", "ICBCT"},
	SME:   One2PayBankCode{"098", "SME", "SME"},
	BAAC:  One2PayBankCode{"034", "BAAC", "BAAC"},
	EXIM:  One2PayBankCode{"035", "EXIM", "EXIM"},
	GSB:   One2PayBankCode{"030", "GSB", "GSB"},
	GHB:   One2PayBankCode{"033", "GHB", "GHB"},
	ISBT:  One2PayBankCode{"066", "ISBT", "ISBT"},
}

// GetAcronymCode returns the bank code for the given acronym
func GetAcronymCode(acronym string) string {
	codes := map[string]string{
		"BBL":   "002",
		"KBANK": "004",
		"KTB":   "006",
		"TTB":   "011",
		"SCB":   "014",
		"BAY":   "025",
		"KKP":   "069",
		"CIMBT": "022",
		"TISCO": "067",
		"UOBT":  "024",
		"TCD":   "071",
		"LHFG":  "073",
		"ICBCT": "070",
		"SME":   "098",
		"BAAC":  "034",
		"EXIM":  "035",
		"GSB":   "030",
		"GHB":   "033",
		"ISBT":  "066",
	}

	if _, ok := codes[acronym]; ok {
		return codes[acronym]
	}
	return ""
}

// Eq checks if the bank code matches the given value
func (b One2PayBankCode) Eq(value string) bool {
	return b.Code == value
}
