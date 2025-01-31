package promoCode

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/repository/repositoryPromoCode"
	"telemetry-sale/internal/service/webclient"
	"telemetry-sale/internal/util/constant"
	"time"
)

type Service struct {
	ctx             *gin.Context
	repository      *repositoryPromoCode.PromoCodeRepository
	repositoryGroup *repositoryPromoCode.PromoGroupRepository
	machineClient   *webclient.MachineClient
}

func NewService(ctx *gin.Context) *Service {
	return &Service{
		ctx:        ctx,
		repository: repositoryPromoCode.NewPromoCodeRepository(ctx),
	}
}

func (service *Service) Create(promoCodeDto model.PromoCodeCreateDto) error {

	if promoCodeDto.IsGenerateSeveral {
		if err := service.saveMultiple(promoCodeDto); err != nil {
			return err
		}
	} else {
		if err := service.saveSingle(promoCodeDto); err != nil {
			return err
		}
	}

	return nil
}

func (service *Service) Edit(id uint64, promoCodeDto model.PromoCodeEditDto) error {
	paramPromoCode := service.getParamIdsString(promoCodeDto.MachineIds, promoCodeDto.Product)

	err := service.repository.Edit(id, promoCodeDto, paramPromoCode)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) SwitchSelected(id uint64, switchSelected model.SwitchSelected) error {

	err := s.repository.SwitchSelected(id, switchSelected)
	if err != nil {
		return err
	}

	return nil
}

func (service *Service) List(promoFilter model.PromoCodeFilterList, page uint16, limit uint16) ([]model.PromoCodeListResult, error) {
	offset := page*limit - limit

	return service.repository.GetList(promoFilter, limit, offset)
}

func (s *Service) GetFilter(promoFilter model.PromoCodeFilterList) (model.PromoCodeFilter, error) {
	var promoCodeFilter model.PromoCodeFilter

	filterValue, err := s.repository.GetFilter(promoFilter)
	if err != nil {
		return promoCodeFilter, err
	}

	promoCodeFilter.Current = promoFilter
	promoCodeFilter.Value = filterValue

	return promoCodeFilter, nil
}

func (service *Service) Element(id uint64) (model.PromoCodeElementResult, error) {
	service.repositoryGroup = repositoryPromoCode.NewPromoGroupRepository(service.ctx)
	service.machineClient = webclient.NewMachineClient()

	promoCodeElement, paramIdsString, err := service.repository.GetElement(id)
	if err != nil {
		return promoCodeElement, err
	}

	promoCodeIdsStringService := NewPromoCodeIdsStringService(service.ctx)
	if err = promoCodeIdsStringService.AddMachine(&promoCodeElement, paramIdsString.MachineIds); err != nil {
		return promoCodeElement, err
	}

	if err = promoCodeIdsStringService.AddCategory(&promoCodeElement, paramIdsString.CategoryIds); err != nil {
		return promoCodeElement, err
	}

	if err = promoCodeIdsStringService.AddView(&promoCodeElement, paramIdsString.ViewIds); err != nil {
		return promoCodeElement, err
	}

	err = promoCodeIdsStringService.AddProduct(
		&promoCodeElement,
		paramIdsString.BrandIds,
		paramIdsString.IngredientLineIds,
		paramIdsString.IngredientIds,
	)

	group := promoCodeElement.Group

	if group != nil {
		qtyElementInGroup, err := service.repositoryGroup.GetQtyElementGroup(*group.Id)
		if err != nil {
			return model.PromoCodeElementResult{}, err
		}

		promoCodeElement.QtyInGroup = qtyElementInGroup
	}

	return promoCodeElement, nil
}

func (service *Service) Qty(organizationId uint64) (model.Qty, error) {
	return service.repository.GetQty(organizationId)
}

func (s *Service) saveSingle(promoCodeDto model.PromoCodeCreateDto) error {
	paramPromoCode := s.getParamIdsString(promoCodeDto.MachineIds, promoCodeDto.Product)

	_, err := s.repository.SavePromoCode(promoCodeDto, paramPromoCode)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateFilterForList(ctx *gin.Context) (model.PromoCodeFilterList, error) {
	var promoCodeListFilter model.PromoCodeFilterList

	organizationId, err := strconv.ParseUint(ctx.Param("organizationId"), 10, 64)

	if err != nil {
		return promoCodeListFilter, err
	}

	isSelectedString := ctx.Query("isSelected")
	if isSelectedString != "" {
		isSelected, err := strconv.ParseBool(isSelectedString)
		if err != nil {
			return model.PromoCodeFilterList{}, err
		}

		promoCodeListFilter.IsSelected = &isSelected
	}

	isActiveString := ctx.Query("isActive")
	if isActiveString != "" {
		isActive, err := strconv.ParseBool(isActiveString)
		if err != nil {
			return model.PromoCodeFilterList{}, err
		}

		promoCodeListFilter.IsActive = &isActive
	}

	periodFromString := ctx.Query("periodFrom")
	if periodFromString != "" {
		periodFrom, err := time.Parse(constant.DateFormat, periodFromString)
		if err != nil {
			return model.PromoCodeFilterList{}, err
		}

		promoCodeListFilter.PeriodFrom = &periodFrom
	}

	periodToString := ctx.Query("periodTo")
	if periodToString != "" {
		periodTo, err := time.Parse(constant.DateFormat, periodToString)
		if err != nil {
			return model.PromoCodeFilterList{}, err
		}

		promoCodeListFilter.PeriodTo = &periodTo
	}

	qtyMinString := ctx.Query("qtyMin")
	if qtyMinString != "" {
		qtyMin, err := strconv.ParseUint(qtyMinString, 10, 16)
		if err != nil {
			return model.PromoCodeFilterList{}, err
		}

		qtyMinUint16 := uint16(qtyMin)
		promoCodeListFilter.QtyMin = &qtyMinUint16
	}

	qtyMaxString := ctx.Query("qtyMax")
	if qtyMaxString != "" {
		qtyMax, err := strconv.ParseUint(qtyMaxString, 10, 16)
		qtyMaxUint16 := uint16(qtyMax)

		if err != nil {
			return model.PromoCodeFilterList{}, err
		}

		promoCodeListFilter.QtyMax = &qtyMaxUint16
	}

	discountType := constant.TypeDiscountEnum(ctx.Query("discountType"))
	if discountType != "" {
		promoCodeListFilter.DiscountType = &discountType
	}

	discountAmountMinString := ctx.Query("discountAmountMin")
	if discountAmountMinString != "" {
		discountAmountMin, err := strconv.ParseUint(discountAmountMinString, 10, 8)
		if err != nil {
			return model.PromoCodeFilterList{}, err
		}

		discountAmountMinUint8 := uint8(discountAmountMin)
		promoCodeListFilter.DiscountAmountMin = &discountAmountMinUint8
	}

	discountAmountMaxString := ctx.Query("discountAmountMax")
	if discountAmountMinString != "" {
		discountAmountMax, err := strconv.ParseUint(discountAmountMaxString, 10, 8)
		if err != nil {
			return model.PromoCodeFilterList{}, err
		}

		discountAmountMaxUint8 := uint8(discountAmountMax)
		promoCodeListFilter.DiscountAmountMax = &discountAmountMaxUint8
	}

	code := ctx.Query("code")
	if code != "" {
		promoCodeListFilter.Code = &code
	}

	promoCodeListFilter.CreatedSort = constant.Asc
	createdSort := ctx.Query("createdSort")
	if createdSort != "" {
		promoCodeListFilter.CreatedSort = constant.SortDirection(createdSort)
	}

	promoCodeListFilter.OrganizationId = &organizationId

	return promoCodeListFilter, nil
}

func (s *Service) saveMultiple(promoCodeDto model.PromoCodeCreateDto) error {
	var prefixPromoCode model.PrefixPromoCodeDto
	generationSetting := promoCodeDto.GenerationSetting

	prefixPromoCode.Prefix = generationSetting.Prefix
	promoGenerationService := NewPromoGenerationService(s.ctx)

	for i := 0; i < int(generationSetting.Qty); i++ {
		code, err := promoGenerationService.GenerationPromoCode(promoCodeDto.OrganizationId, prefixPromoCode)
		if err != nil {
			return err
		}

		promoCodeDto.Code = code.Code

		if err := s.saveSingle(promoCodeDto); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) sliceIdsToString(ids []uint64) string {
	var idsString string

	if len(ids) == 0 {
		return idsString
	}

	for _, id := range ids {
		idsString += strconv.FormatUint(id, 10) + ","
	}

	if idsString != "" {
		idsString = strings.Trim(idsString, ",")
	}

	return idsString
}

func (s *Service) getParamIdsString(machineIds []uint64, product *model.Product) model.ParamPromoCodeIdsString {
	var paramPromoCode model.ParamPromoCodeIdsString
	brandIds := make([]uint64, 0, 100)
	ingredientLineIds := make([]uint64, 0, 100)
	ingredientIds := make([]uint64, 0, 100)

	paramPromoCode.MachineIds = s.sliceIdsToString(machineIds)
	paramPromoCode.CategoryIds = s.sliceIdsToString(product.CategoryIds)
	paramPromoCode.ViewIds = s.sliceIdsToString(product.ViewIds)

	brands := product.Brands
	for _, brand := range brands {
		brandIds = append(brandIds, brand.Id)

		ingredientLines := brand.IngredientLines
		for _, ingredientLine := range ingredientLines {
			ingredientLineIds = append(ingredientLineIds, ingredientLine.Id)

			ingredients := ingredientLine.IngredientIds
			for _, ingredient := range ingredients {
				ingredientIds = append(ingredientIds, ingredient)
			}
		}
	}

	paramPromoCode.BrandIds = s.sliceIdsToString(brandIds)
	paramPromoCode.IngredientLineIds = s.sliceIdsToString(ingredientLineIds)
	paramPromoCode.IngredientIds = s.sliceIdsToString(ingredientIds)

	return paramPromoCode
}
