package promoCode

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/repository/repositoryPromoCode"
	"telemetry-sale/internal/util/constant"
	"telemetry-sale/internal/util/helper"
)

type PromoGenerationService struct {
	repository *repositoryPromoCode.PromoGenerationRepository
}

func NewPromoGenerationService(ctx *gin.Context) *PromoGenerationService {
	return &PromoGenerationService{
		repository: repositoryPromoCode.NewPromoGenerationRepository(ctx),
	}
}

func (s *PromoGenerationService) GenerationPromoCode(
	organizationId uint64, prefixPromoCode model.PrefixPromoCodeDto,
) (model.PromoCodeGenerationResult, error) {
	var promoCodeGeneration model.PromoCodeGenerationResult
	var numberAttempt uint8
	numberAttempt = 7

	code, err := s.findUniquePromo(prefixPromoCode, organizationId, numberAttempt)
	if err != nil {
		return promoCodeGeneration, err
	}

	promoCodeGeneration.Code = code

	return promoCodeGeneration, nil
}

func (s *PromoGenerationService) findUniquePromo(
	prefixPromoCode model.PrefixPromoCodeDto, organizationId uint64, numberAttempt uint8,
) (string, error) {
	qtyLetterPrefix := len(prefixPromoCode.Prefix)

	randomString := helper.GeneratePromoCode(constant.MaxQtyLetterPromoCode - qtyLetterPrefix)

	code := strings.ToUpper(prefixPromoCode.Prefix) + randomString

	idPromoCode, err := s.repository.GetByCodeAndOrganizationId(code, organizationId)

	if err != nil {
		return "", err
	}

	if idPromoCode == 0 {
		return code, nil
	}

	if numberAttempt == 0 {
		err = errors.New(string(constant.MaxAttemptGenerateCode))
		return "", err
	}

	numberAttempt--

	return s.findUniquePromo(prefixPromoCode, organizationId, numberAttempt)
}
