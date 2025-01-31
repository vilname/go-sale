package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"telemetry-sale/internal/model"
	promoCodeService "telemetry-sale/internal/service/promoCode"
	"telemetry-sale/internal/util/constant"
	"telemetry-sale/internal/util/helper"
)

// CreatePromoCode godoc
// @Tags Промокоды
// @Param request body model.PromoCodeCreateDto true "Промокод"
// @Summary Создание промокода
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Failure	400	{object} helper.ErrorValidate "Ошибка валидации"
// @Failure	500	{object} helper.ErrorResponse "Другие ошибки"
// @Router /promo-code/create [post]
func CreatePromoCode(ctx *gin.Context) {
	var promoCodeDto model.PromoCodeCreateDto

	body, err := io.ReadAll(ctx.Request.Body)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	if err := json.Unmarshal(body, &promoCodeDto); err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	validateError := helper.RegisterValidate(promoCodeDto)
	if validateError != nil {
		helper.ErrorValidateResponse(ctx, validateError)
		return
	}

	service := promoCodeService.NewService(ctx)
	err = service.Create(promoCodeDto)
	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// EditPromoCode godoc
// @Tags Промокоды
// @Param request body model.PromoCodeEditDto true "Промокод"
// @Summary Редактирование промокода
// @Param id path uint true "ID промокода"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @param Authorization header string true "Authorization"
// @Failure	400	{object} helper.ErrorValidate "Ошибка валидации"
// @Failure	500	{object} helper.ErrorResponse "Другие ошибки"
// @Router /promo-code/edit/{id} [post]
func EditPromoCode(ctx *gin.Context) {
	var promoCodeDto model.PromoCodeEditDto

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	body, err := io.ReadAll(ctx.Request.Body)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	if err := json.Unmarshal(body, &promoCodeDto); err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	validateError := helper.RegisterValidate(promoCodeDto)
	if validateError != nil {
		helper.ErrorValidateResponse(ctx, validateError)
		return
	}

	service := promoCodeService.NewService(ctx)
	err = service.Edit(id, promoCodeDto)
	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// SwitchSelected godoc
// @Tags Промокоды
// @Param request body model.SwitchSelected true "Избранные"
// @Summary Переключатель избранных
// @Param id path uint true "ID промокода"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Failure	400	{object} helper.ErrorValidate "Ошибка валидации"
// @Failure	500	{object} helper.ErrorResponse "Другие ошибки"
// @Router /promo-code/switch-selected/{id} [post]
func SwitchSelected(ctx *gin.Context) {
	var switchSelected model.SwitchSelected

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	body, err := io.ReadAll(ctx.Request.Body)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	if err := json.Unmarshal(body, &switchSelected); err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	service := promoCodeService.NewService(ctx)
	err = service.SwitchSelected(id, switchSelected)
	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// ListPromoCode godoc
// @Tags Промокоды
// @Summary Список промокодов
// @Param organizationId path uint true "ID организации"
// @Param page query integer false "Номер страницы"
// @Param limit query integer false "Количество элементов на странице"
// @Param code query string false "Код промокода"
// @Param isSelected query bool false "Избранное"
// @Param isActive query bool false "Активные"
// @Param periodFrom query string false "Дата начала периода" example(2006-01-02)
// @Param periodTo query string false "Дата окончания периода" example(2006-01-02)
// @Param qtyMin query uint false "Количество минимальных использований"
// @Param qtyMax query uint false "Количество максимальных использований"
// @Param discountType query constant.TypeDiscountEnum false "Тип скидки"
// @Param discountAmountMin query uint false "Размер скидки минимум"
// @Param discountAmountMax query uint false "Размер скидки максимум"
// @Param createdSort query constant.SortDirection false "Сортировка по дате создания"
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security OAuth2Implicit
//
// @Success 200	{object} model.PaginationPromoCode
// @Failure	500	{object} helper.ErrorResponse "Другие ошибки"
// @Router /promo-code/list/{organizationId} [get]
func ListPromoCode(ctx *gin.Context) {

	page := helper.GetDefaultPage(ctx.Query("page"))
	limit := helper.GetDefaultLimit(ctx.Query("limit"), constant.DefaultPageLimit)

	service := promoCodeService.NewService(ctx)
	listFilter, err := service.CreateFilterForList(ctx)
	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}
	promoCodeList, err := service.List(listFilter, page, limit)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	var pagePagination model.PaginationPromoCode
	pagePagination.GetPagination(promoCodeList, page, limit)

	ctx.JSON(http.StatusOK, pagePagination)
}

// FilterPromoCode godoc
// @Tags Промокоды
// @Summary Фильтр списка промокодов
// @Param organizationId path uint true "ID организации"
// @Param isSelected query bool false "Избранное"
// @Param isActive query bool false "Активные"
// @Param periodFrom query string false "Дата начала периода" example(2006-01-02)
// @Param periodTo query string false "Дата окончания периода" example(2006-01-02)
// @Param qtyMin query uint false "Количество минимальных использований"
// @Param qtyMax query uint false "Количество максимальных использований"
// @Param discountType query constant.TypeDiscountEnum false "Тип скидки"
// @Param discountAmountMin query uint false "Размер скидки минимум"
// @Param discountAmountMax query uint false "Размер скидки максимум"
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security OAuth2Implicit
//
// @Success 200	{object} model.PromoCodeFilter
// @Failure	500	{object} helper.ErrorResponse "Другие ошибки"
// @Router /promo-code/filter/{organizationId} [get]
func FilterPromoCode(ctx *gin.Context) {

	service := promoCodeService.NewService(ctx)
	listFilter, err := service.CreateFilterForList(ctx)
	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}
	promoCodeFilter, err := service.GetFilter(listFilter)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, promoCodeFilter)
}

// ElementPromoCode godoc
// @Tags Промокоды
// @Summary Детальная страница
// @Param id path uint true "ID промокода"
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security OAuth2Implicit
//
// @Success 200	{object} model.PromoCodeElementResult
// @Failure	500	{object} helper.ErrorResponse "Другие ошибки"
// @Router /promo-code/element/{id} [get]
func ElementPromoCode(ctx *gin.Context) {

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	service := promoCodeService.NewService(ctx)
	promoCodeElement, err := service.Element(id)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, promoCodeElement)
}

// QtyPromoCode godoc
// @Tags Промокоды
// @Summary Всего элементов
// @Param organizationId path uint true "ID организации"
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security OAuth2Implicit
//
// @Success 200	{object} model.Qty
// @Failure	500	{object} helper.ErrorResponse "Другие ошибки"
// @Router /promo-code/qty/{organizationId} [get]
func QtyPromoCode(ctx *gin.Context) {

	organizationId, err := strconv.ParseUint(ctx.Param("organizationId"), 10, 64)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	service := promoCodeService.NewService(ctx)
	qtyPromoCode, err := service.Qty(organizationId)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, qtyPromoCode)
}
