package service

import (
	"github.com/gin-gonic/gin"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/repository"
	"telemetry-sale/internal/util/constant"
	"telemetry-sale/internal/util/helper"
)

type SalePeriodService struct {
	repository *repository.SalePeriodRepository
}

func NewSalePeriodService(ctx *gin.Context) *SalePeriodService {
	return &SalePeriodService{
		repository: repository.NewSalePeriodRepository(ctx),
	}
}

func (service *SalePeriodService) GetSalePeriod(serialNumber string) (model.SalePeriodResult, error) {
	var salePeriodResult model.SalePeriodResult

	yesterday := helper.GetDateByPeriod(constant.DAY)
	countSale, err := service.repository.GetCountSaleByDate(serialNumber, yesterday)
	if err != nil {
		return salePeriodResult, err
	}

	salePeriodResult.Day = countSale

	week := helper.GetDateByPeriod(constant.WEEK)
	countSale, err = service.repository.GetCountSaleByDate(serialNumber, week)
	if err != nil {
		return salePeriodResult, err
	}

	salePeriodResult.Week = countSale

	month := helper.GetDateByPeriod(constant.MONTH)
	countSale, err = service.repository.GetCountSaleByDate(serialNumber, month)
	if err != nil {
		return salePeriodResult, err
	}

	salePeriodResult.Month = countSale

	return salePeriodResult, nil
}

func (service *SalePeriodService) GetBySerialNumberAndPeriod(serialNumbers []string, period constant.PeriodEnum) ([]model.SaleQtyBySerialNumberResult, error) {
	timeAgo := helper.GetDateByPeriod(period)

	return service.repository.GetCountBySerialNumbersAndTimeAgo(serialNumbers, timeAgo)
}

func (service *SalePeriodService) GetBySerialNumberAndNewPeriod(
	serialNumbers []string, period constant.NewPeriodEnum,
) ([]model.SaleQtyBySerialNumberResult, error) {
	timeAgo := helper.GetDateByNewPeriod(period)

	return service.repository.GetCountBySerialNumbersAndTimeAgo(serialNumbers, timeAgo)
}
