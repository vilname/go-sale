package model

type SalePeriodResult struct {
	Day   uint16 `json:"day"`
	Week  uint16 `json:"week"`
	Month uint16 `json:"month"`
}

type SaleQtyBySerialNumberResult struct {
	SerialNumber string `json:"serialNumber"`
	Qty          uint16 `json:"qty"`
}
