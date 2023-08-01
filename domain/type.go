package domain

type Router struct {
	// in: query
	RouterSerial   string `json:"router-serial" form:"router-serial"`
	OperatorName   string `json:"operator-name" form:"operator-name"`
	IsoCountryCode string `json:"iso-country-code" form:"iso-country-code"`
	Mac            string `json:"mac" form:"mac"`
	RouterModel    string `json:"router-model" form:"router-model"`
	AccountID      string `json:"account-id" form:"account-id"`
}

type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}
