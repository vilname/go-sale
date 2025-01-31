package rest

import (
	"net/http"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/service"
	"telemetry-sale/internal/util/constant"
	"telemetry-sale/internal/util/helper"

	"github.com/gin-gonic/gin"
)

// ListSale godoc
// @Tags Продажи
// @Summary Список продаж
// @Param serialNumber path string true "Серийный номер автомата"
// @Param page query integer false "Номер страницы"
// @Param limit query integer false "Количество элементов на странице"
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security OAuth2Implicit
//
// @Success 200	{object} model.PaginationSaleList
// @Failure	500	{object} helper.ErrorResponse "Другие ошибки"
// @Router /sale/list/{serialNumber} [get]
func ListSale(ctx *gin.Context) {

	serialNumber := ctx.Param("serialNumber")

	page := helper.GetDefaultPage(ctx.Query("page"))
	limit := helper.GetDefaultLimit(ctx.Query("limit"), constant.DefaultPageLimit)

	saleService := service.NewSaleService(ctx)
	saleListResult, err := saleService.List(serialNumber, page, limit)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	var saleList model.PaginationSaleList
	saleList.GetPagination(saleListResult, page, limit)

	ctx.JSON(http.StatusOK, saleList)
}

// DateLastSale godoc
// @Tags Продажи
// @Summary Дата последней продажи
// @Param serialNumber path string true "Серийный номер автомата"
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security OAuth2Implicit
//
// @Success 200	{object} model.DateLastSale
// @Failure	500	{object} helper.ErrorResponse "Другие ошибки"
// @Router /sale/date-last-sale/{serialNumber} [get]
func DateLastSale(ctx *gin.Context) {
	serialNumber := ctx.Param("serialNumber")

	saleService := service.NewSaleService(ctx)
	lastSale, err := saleService.LastDate(serialNumber)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, lastSale)
}

// QtySale godoc
// @Tags Продажи
// @Summary Всего элементов
// @Param serialNumber path string true "Серийный номер автомата"
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security OAuth2Implicit
//
// @Success 200	{object} model.Qty
// @Failure	500	{object} helper.ErrorResponse "Другие ошибки"
// @Router /sale/qty/{serialNumber} [get]
func QtySale(ctx *gin.Context) {

	serialNumber := ctx.Param("serialNumber")

	saleService := service.NewSaleService(ctx)
	qtyPromoCode, err := saleService.Qty(serialNumber)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, qtyPromoCode)
}
