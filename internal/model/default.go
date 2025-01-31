package model

type IdResult struct {
	Id uint64 `json:"id"`
}

type Pagination struct {
	Page  uint16 `json:"page"`
	Limit uint16 `json:"limit"`
}

type Qty struct {
	Qty uint16 `json:"qty"`
}

type Default struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}
