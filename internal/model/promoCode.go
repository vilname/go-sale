package model

import (
	_ "github.com/go-playground/validator/v10"
	"telemetry-sale/internal/util/constant"
	"time"
)

// Date преобразование типа даты, по другому свагер не хочет нормально работать
type Date struct {
	Time  time.Time
	Valid bool
}

type PromoCodeCreateDto struct {
	Code              string            `json:"code" validate:"max=12,regexp=^[A-Z0-9]*$"` // Промокод
	IsGenerateSeveral bool              `json:"isGenerateSeveral"`                         // Сгенерировать несколько
	Qty               uint16            `json:"qty"`                                       // Количество промокодов
	OrganizationId    uint64            `json:"organizationId" validate:"required"`
	GenerationSetting GenerationSetting `json:"generationSetting"` // Настройки генерации
	Discount          Discount          `json:"discount"`
	PeriodUse         PeriodUseDto      `json:"periodUse"`
	GroupId           uint64            `json:"groupId"`
	Description       string            `json:"description"`
	Schedules         []NewSchedule     `json:"schedules"`
	MachineIds        []uint64          `json:"machineIds"`
	Product           *Product          `json:"product"`
}

type PromoCodeEditDto struct {
	Code        string       `json:"code" validate:"required,min=4,max=12,regexp=^[A-Z0-9]*$"` // Промокод
	Qty         uint16       `json:"qty"`                                                      // Количество промокодов
	IsActive    bool         `json:"isActive"`
	Discount    Discount     `json:"discount"`
	PeriodUse   PeriodUseDto `json:"periodUse"`
	GroupId     uint64       `json:"groupId"`
	Description string       `json:"description"`
	Schedules   []Schedule   `json:"schedules"`
	MachineIds  []uint64     `json:"machineIds"`
	Product     *Product     `json:"product"`
}

type PromoCodeFilter struct {
	Current PromoCodeFilterList        `json:"current"`
	Value   PromoCodeFilterValueResult `json:"value"`
}

type PromoCodeListResult struct {
	Id          uint64          `json:"id"`
	IsSelected  bool            `json:"isSelected"`
	Code        string          `json:"code"`
	Qty         *uint16         `json:"qty"`
	Used        uint32          `json:"used"`
	IsActive    bool            `json:"isActive"`
	Discount    Discount        `json:"discount"`
	PeriodUse   PeriodUseResult `json:"periodUse"`
	GroupName   *string         `json:"groupName"`
	Description *string         `json:"description"`
}

type PromoCodeFilterValueResult struct {
	QtyMax            *uint64 `json:"maxQty"`
	DiscountAmountMax *uint8  `json:"discountAmountMax"`
}

type PromoCodeFilterCurrentResult struct {
	IsSelected        interface{}               `json:"isSelected"`
	IsActive          interface{}               `json:"isActive"`
	PeriodFrom        string                    `json:"periodFrom"`
	PeriodTo          string                    `json:"periodTo"`
	UseMin            uint16                    `json:"useMin"`
	UseMax            uint16                    `json:"useMax"`
	DiscountType      constant.TypeDiscountEnum `json:"discountType"`
	DiscountAmountMin uint8                     `json:"discountAmountMin"`
	DiscountAmountMax uint8                     `json:"discountAmountMax"`
}

type PromoCodeElementResult struct {
	Id             uint64            `json:"id"`
	Code           string            `json:"code"`
	Qty            uint16            `json:"qty"`
	IsActive       bool              `json:"isActive"`
	Discount       Discount          `json:"discount"`
	PeriodUse      PeriodUseResult   `json:"periodUse"`
	Group          *PromoGroupResult `json:"group"`
	QtyInGroup     uint16            `json:"qtyInGroup"`
	Description    string            `json:"description"`
	Schedules      []*Schedule       `json:"schedules"`
	MachinesString string            `json:"-"`
	Machines       []Default         `json:"machines"`
	Product        ProductResult     `json:"product"`
}

type GenerationSetting struct {
	Prefix    string `json:"prefix" validate:"max=6,regexp=^[A-Z0-9]*$"`
	QtyLetter uint8  `json:"qtyLetter"`
	Qty       uint8  `json:"qty"`
}

type Discount struct {
	Type   constant.TypeDiscountEnum `json:"type"`
	Amount uint16                    `json:"amount"`
}

type PeriodUseDto struct {
	From *string `json:"from"`
	To   *string `json:"to"`
}

type PeriodUseResult struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type PaginationPromoCode struct {
	Data       []PromoCodeListResult `json:"data"`
	Pagination Pagination            `json:"pagination"`
}

func (pagination *PaginationPromoCode) GetPagination(list []PromoCodeListResult, page uint16, limit uint16) {
	pagination.Data = list
	pagination.Pagination = Pagination{
		Page:  page,
		Limit: limit,
	}
}

type PromoCodeFilterList struct {
	OrganizationId    *uint64                    `json:"-"`
	IsSelected        *bool                      `json:"isSelected"`
	IsActive          *bool                      `json:"isActive"`
	PeriodFrom        *time.Time                 `json:"periodFrom"`
	PeriodTo          *time.Time                 `json:"periodTo"`
	QtyMin            *uint16                    `json:"qtyMin"`
	QtyMax            *uint16                    `json:"qtyMax"`
	DiscountType      *constant.TypeDiscountEnum `json:"discountType"`
	DiscountAmountMin *uint8                     `json:"discountAmountMin"`
	DiscountAmountMax *uint8                     `json:"discountAmountMax"`
	Code              *string                    `json:"code"`
	CreatedSort       constant.SortDirection     `json:"createdSort"`
}

type SwitchSelected struct {
	IsSelected bool `json:"isSelected"`
}

func (d *Date) UnmarshalJSON(bytes []byte) error {
	var timeForm time.Time
	var valid bool

	dd, err := time.Parse(`"2006-01-02"`, string(bytes))
	if err != nil {
		timeForm = time.Time{}
		valid = false
	} else {
		timeForm = dd
		valid = true
	}

	*d = Date{
		Time:  timeForm,
		Valid: valid,
	}

	return nil
}
