package domain

type Router struct {
	// in: query
	RouterSerial string `json:"router-serial" form:"router-serial"`
}

type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}
