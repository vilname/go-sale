package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"telemetry-sale/internal/service"
	"telemetry-sale/internal/util"
	"telemetry-sale/internal/util/constant"
)

type Test struct {
	Name string `json:"name"`
}

// GetAllSalePeriod godoc
// @Tags Продажи по преодам
// @Summary Получение данных о количестве продаж за день, неделю, месяц
// @Param serialNumber path string true "Серийный номер автомата"
// @Accept json
// @Produce json
// @Success 200	{object} model.SalePeriod
// @Failure	400	{object} helper.ErrorValidate "Ошибка валидации"
// @Failure	500	{object} helper.ErrorResponse "Другие ошибки"
// @Router /sale-period/all/{serialNumber} [get]
func GetAllSalePeriod(ctx *gin.Context) {

	serialNumber := ctx.Param("serialNumber")

	salePeriodService := service.NewSalePeriodService(ctx)
	salePeriodResult, err := salePeriodService.GetSalePeriod(serialNumber)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{util.ErrorKey: err.Error()})
	}

	ctx.JSON(http.StatusOK, salePeriodResult)
}

// GetSalePeriodBySerialNumbers godoc
// @Tags Продажи по преодам
// @Summary Получение данных о количестве продаж либо за день, либо за неделю, либо за месяц
// @Param serialNumbers query []string true "Серийный номер автомата"
// @Param period query constant.PeriodEnum true "Период выборки"
// @Accept json
// @Produce json
// @Success 200	{object} model.SaleQtyBySerialNumberResult
// @Failure	400	{object} helper.ErrorValidate "Ошибка валидации"
// @Failure	500	{object} helper.ErrorResponse "Другие ошибки"
// @Router /sale-period/list/by-serial-numbers [get]
func GetSalePeriodBySerialNumbers(ctx *gin.Context) {

	serialNumbers := ctx.QueryArray("serialNumbers")
	period := constant.PeriodEnum(ctx.Query("period"))

	salePeriodService := service.NewSalePeriodService(ctx)
	qtySales, err := salePeriodService.GetBySerialNumberAndPeriod(serialNumbers, period)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{util.ErrorKey: err.Error()})
	}

	ctx.JSON(http.StatusOK, qtySales)
}

// GetSalePeriodBySerialNumbersPeriod godoc
// @Tags Продажи по преодам
// @Summary Получение данных о количестве продаж либо за день, либо за неделю, либо за месяц
// @Param serialNumbers query []string true "Серийный номер автомата"
// @Param period query constant.PeriodEnum true "Период выборки"
// @Accept json
// @Produce json
// @Success 200	{object} model.SaleQtyBySerialNumberResult
// @Failure	400	{object} helper.ErrorValidate "Ошибка валидации"
// @Failure	500	{object} helper.ErrorResponse "Другие ошибки"
// @Router /sale-period/list/by-serial-numbers [get]
func GetSalePeriodBySerialNumbersPeriod(ctx *gin.Context) {

	serialNumbers := ctx.QueryArray("serialNumbers")
	period := constant.NewPeriodEnum(ctx.Query("period"))

	salePeriodService := service.NewSalePeriodService(ctx)
	qtySales, err := salePeriodService.GetBySerialNumberAndNewPeriod(serialNumbers, period)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{util.ErrorKey: err.Error()})
	}

	ctx.JSON(http.StatusOK, qtySales)
}
