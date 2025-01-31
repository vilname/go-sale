package model

import (
	"telemetry-sale/internal/util"
	"telemetry-sale/internal/util/constant"
	"time"
)

type Sale struct {
	SerialNumber     string        `json:"serialNumber"`
	DateSale         string        `json:"dateSale"`
	MachineProductId uint          `json:"machineProductId"`
	PromoCodeId      uint64        `json:"promoCodeId"`
	DiscountId       uint          `json:"discountId"`
	Volume           uint          `json:"volume"`
	Name             string        `json:"name"`
	Price            float64       `json:"price"`
	Unit             constant.Unit `json:"unit"`
	WriteOffs        []WriteOff    `json:"writeOffs"`
	Payments         []Payment     `json:"payments"`
}

type WriteOff struct {
	CellNumber   uint          `json:"cellNumber"`
	IngredientId uint          `json:"ingredientId"`
	Volume       uint          `json:"volume"`
	CellType     util.CellType `json:"cellType"`
	Unit         string        `json:"unit"`
}

type Payment struct {
	Price  float64 `json:"price"`
	Method string  `json:"method"`
}

type SalePeriod struct {
	Day   int `json:"day"`
	Week  int `json:"week"`
	Month int `json:"month"`
}

type SaleListResult struct {
	Id             uint64                `json:"id"`
	Name           string                `json:"name"`
	Volume         uint16                `json:"volume"`
	Unit           constant.Unit         `json:"unit"`
	Price          float32               `json:"price"`
	DiscountPrices []DiscountPriceResult `json:"discountPrices"`
	PromoCode      *PromoCodeResult      `json:"promoCode"`
	DateSale       string                `json:"dateSale" example:"27.03.2024 16:03:22"`
}

func NewSaleListResult() *SaleListResult {
	return &SaleListResult{
		DiscountPrices: make([]DiscountPriceResult, 0, 100),
	}
}

type PaginationSaleList struct {
	Data       []SaleListResult `json:"data"`
	Pagination Pagination       `json:"pagination"`
}

func (pagination *PaginationSaleList) GetPagination(list []SaleListResult, page uint16, limit uint16) {
	pagination.Data = list
	pagination.Pagination = Pagination{
		Page:  page,
		Limit: limit,
	}
}

type PromoCodeResult struct {
	Code   *string                    `json:"code"`
	Amount *uint16                    `json:"amount"`
	Type   *constant.TypeDiscountEnum `json:"type"`
}

type DiscountPriceResult struct {
	Price  *float32                `json:"price"`
	Method *constant.PaymentMethod `json:"method"`
}

type DateLastSale struct {
	Date time.Time `json:"date"`
}
