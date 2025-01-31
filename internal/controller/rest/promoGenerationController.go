package rest

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/service/promoCode"
	"telemetry-sale/internal/util/helper"
)

// GenerationCode godoc
// @Tags Генерация промокода
// @Param request body model.PrefixPromoCodeDto true "Префикс промокода"
// @Summary Генерация одного промокода
// @Param organizationId path uint true "ID организации"
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200	{object} model.PromoCodeGenerationResult
// @Failure	400	{object} helper.ErrorValidate "Ошибка валидации"
// @Failure	500	{object} helper.ErrorResponse "Другие ошибки"
// @Router /promo-generation/single/{organizationId} [post]
func GenerationCode(ctx *gin.Context) {
	var prefixPromoCode model.PrefixPromoCodeDto

	body, err := io.ReadAll(ctx.Request.Body)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	if err := json.Unmarshal(body, &prefixPromoCode); err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	validateError := helper.RegisterValidate(prefixPromoCode)
	if validateError != nil {
		helper.ErrorValidateResponse(ctx, validateError)
		return
	}

	organizationId, err := strconv.ParseUint(ctx.Param("organizationId"), 10, 64)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	promoGenerationService := promoCode.NewPromoGenerationService(ctx)
	codeResult, err := promoGenerationService.GenerationPromoCode(organizationId, prefixPromoCode)

	if err != nil {
		helper.ErrorResponseMethod(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, codeResult)
}
