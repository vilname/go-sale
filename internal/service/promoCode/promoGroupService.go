package promoCode

import (
	"github.com/gin-gonic/gin"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/repository/repositoryPromoCode"
)

type PromoGroupService struct {
	repository *repositoryPromoCode.PromoGroupRepository
}

func NewPromoGroupService(ctx *gin.Context) *PromoGroupService {
	return &PromoGroupService{
		repository: repositoryPromoCode.NewPromoGroupRepository(ctx),
	}
}

func (service *PromoGroupService) Create(promoGroup model.PromoGroupDto) (model.IdResult, error) {

	result, err := service.repository.SavePromoGroup(promoGroup.Name, promoGroup.OrganizationId)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (service *PromoGroupService) List(organizationId uint64) ([]model.PromoGroupResult, error) {
	return service.repository.GetList(organizationId)
}
