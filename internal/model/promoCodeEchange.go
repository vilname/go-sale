package model

import (
	"database/sql"
	"github.com/lib/pq"
	"telemetry-sale/internal/util/constant"
	"time"
)

type PromoCodeCheckResult struct {
	Id                uint64                    `json:"id"`
	Success           bool                      `json:"success"`
	Type              constant.TypeDiscountEnum `json:"type"`
	Amount            uint8                     `json:"amount"`
	Reason            constant.ReasonEnum       `json:"reason"`
	PeriodFrom        pq.NullTime               `json:"-"`
	PeriodTo          pq.NullTime               `json:"-"`
	Qty               sql.NullInt16             `json:"-"`
	Used              sql.NullInt16             `json:"-"`
	MachineIds        string                    `json:"-"`
	CategoryIds       string                    `json:"-"`
	ViewIds           string                    `json:"-"`
	BrandIds          string                    `json:"-"`
	IngredientLineIds string                    `json:"-"`
	IngredientIds     string                    `json:"-"`
	Schedules         []*ScheduleCheckPromoCode `json:"-"`
}

type PromoCodeUse struct {
	Id uint64 `json:"id"`
}

type PromoCodeCheckDto struct {
	Code         string    `form:"code"`
	IngredientId uint64    `form:"ingredientId"`
	SerialNumber string    `form:"serialNumber"`
	Time         time.Time `form:"time"`
}

type ImportPromoCode struct {
	Id             int
	UserId         int
	PromoCode      string
	From           time.Time
	To             time.Time
	UsageAmount    int
	RemainsUsage   int
	Discount       int
	Tastes         string
	Created        int
	Updated        int
	Used           int
	AllAutomats    int
	CompanyId      int
	IsDeleted      int
	OrganizationId int
}
