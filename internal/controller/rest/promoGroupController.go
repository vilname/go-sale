package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"telemetry-sale/internal/model"
	promoCodeService "telemetry-sale/internal/service/promoCode"
	"telemetry-sale/internal/util/helper"
)

// CreatePromoGroup godoc
// @Tags Группы промокодов
// @Param request body model.PromoGroupDto true "Группа промокода"
// @Summary Создание группы промокода
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200	{object} model.IdResult
// @Failure	400	{object} helper.ErrorValidate "Ошибка валидации"
// @Failure	500	{object} helper.ErrorResponse "Другие ошибки"
// @Router /promo-group/create [post]
func CreatePromoGroup(ctx *gin.Context) {
	var promoGroupDto model.PromoGroupDto

	body, err := io.ReadAll(ctx.Request.Body)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	if err := json.Unmarshal(body, &promoGroupDto); err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	service := promoCodeService.NewPromoGroupService(ctx)
	result, err := service.Create(promoGroupDto)
	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// ListPromoGroup godoc
// @Tags Группы промокодов
// @Summary Список групп промокодов
// @Param organizationId path uint true "ID организации"
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security OAuth2Implicit
//
// @Success 200	{object} model.PromoGroupResult
// @Failure	500	{object} helper.ErrorResponse "Другие ошибки"
// @Router /promo-group/list/{organizationId} [get]
func ListPromoGroup(ctx *gin.Context) {

	organizationId, err := strconv.ParseUint(ctx.Param("organizationId"), 10, 64)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	service := promoCodeService.NewPromoGroupService(ctx)
	promoGroupList, err := service.List(organizationId)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, promoGroupList)
}
