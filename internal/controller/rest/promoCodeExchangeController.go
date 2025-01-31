package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"telemetry-sale/internal/model"
	promoCodeService "telemetry-sale/internal/service/promoCode"
	"telemetry-sale/internal/service/webclient"
	"telemetry-sale/internal/util/helper"
)

// CheckPromoCode godoc
// @Tags Промокод обмены
// @Summary Проверка валидности промокода
// @Param code query string true "Код промокода"
// @Param ingredientId query integer true "ID ингредиента"
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security OAuth2Implicit
//
// @Success 200	{object} model.PromoCodeElementResult
// @Failure	500	{object} helper.ErrorResponse "Другие ошибки"
// @Router /promo-code/check [get]
func CheckPromoCode(ctx *gin.Context) {
	var promoCodeCheckDto model.PromoCodeCheckDto

	if err := ctx.ShouldBindQuery(&promoCodeCheckDto); err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	machineClient := webclient.NewMachineClient()
	organization, err := machineClient.GetMachineOrganization(promoCodeCheckDto.SerialNumber)
	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	service := promoCodeService.NewExchangeService(ctx)
	promoCodeElement, err := service.Check(promoCodeCheckDto, organization.OrganizationId)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, promoCodeElement)
}

// ImportPromoCode godoc
// @Tags Промокод обмены
// @Summary Импорт промокодов
// @Param file formData file true "File to upload"
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Security OAuth2Implicit
//
// @Success 200
// @Failure	500	{object} helper.ErrorResponse "Другие ошибки"
// @Router /promo-code/import [post]
func ImportPromoCode(ctx *gin.Context) {
	service := promoCodeService.NewPromoCodeImportService(ctx)
	err := service.Import(ctx)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}
