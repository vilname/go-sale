package promoCode

import (
	"github.com/gin-gonic/gin"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/service/webclient"
	"telemetry-sale/internal/util/helper"
)

type PromoCodeIdsStringService struct {
	ctx            *gin.Context
	machineClient  *webclient.MachineClient
	categoryClient *webclient.CategoryClient
	viewClient     *webclient.ViewClient
	productClient  *webclient.ProductClient
}

func NewPromoCodeIdsStringService(ctx *gin.Context) *PromoCodeIdsStringService {
	return &PromoCodeIdsStringService{
		ctx: ctx,
	}
}

func (s *PromoCodeIdsStringService) AddMachine(promoCodeElement *model.PromoCodeElementResult, machineIdsString string) error {
	if machineIdsString == "" {
		return nil
	}

	machineIds := helper.GetIdsSlice(machineIdsString)

	machines, err := s.machineClient.GetByIds(machineIds)
	if err != nil {
		return err
	}

	promoCodeElement.Machines = machines

	return nil
}

func (s *PromoCodeIdsStringService) AddCategory(promoCodeElement *model.PromoCodeElementResult, categoryIdsString string) error {
	promoCodeElement.Product.Categories = make([]model.Default, 0)

	if categoryIdsString == "" {
		return nil
	}

	categoryIds := helper.GetIdsSlice(categoryIdsString)

	categories, err := s.categoryClient.GetByIds(categoryIds)

	if err != nil {
		return err
	}

	promoCodeElement.Product.Categories = categories

	return nil
}

func (s *PromoCodeIdsStringService) AddView(promoCodeElement *model.PromoCodeElementResult, viewIdsString string) error {
	promoCodeElement.Product.Views = make([]model.Default, 0)

	if viewIdsString == "" {
		return nil
	}

	viewIds := helper.GetIdsSlice(viewIdsString)

	views, err := s.viewClient.GetByIds(viewIds)

	if err != nil {
		return err
	}

	promoCodeElement.Product.Views = views

	return nil
}

func (s *PromoCodeIdsStringService) AddProduct(
	promoCodeElement *model.PromoCodeElementResult,
	brandIdsString string,
	ingredientLineIdsString string,
	ingredientIdsString string,
) error {
	brandsResult, err := s.getBrands(brandIdsString)
	if err != nil {
		return err
	}

	ingredientLinesResult, err := s.getIngredientLines(ingredientLineIdsString)
	if err != nil {
		return err
	}

	ingredientsResult, err := s.getIngredients(ingredientIdsString)
	if err != nil {
		return err
	}

	promoCodeElement.Product.Brands = s.buildProduct(brandsResult, ingredientLinesResult, ingredientsResult)

	return nil
}

func (s *PromoCodeIdsStringService) getBrands(brandIdsString string) ([]model.Default, error) {
	if brandIdsString == "" {
		return nil, nil
	}

	brandIds := helper.GetIdsSlice(brandIdsString)

	brandResult, err := s.productClient.GetAllBrandByIds(brandIds)
	if err != nil {
		return nil, err
	}

	return brandResult, nil
}

func (s *PromoCodeIdsStringService) getIngredientLines(ingredientLineIdsString string) ([]model.IngredientLineDefaultResult, error) {
	if ingredientLineIdsString == "" {
		return nil, nil
	}

	ingredientLineIds := helper.GetIdsSlice(ingredientLineIdsString)

	ingredientLineResult, err := s.productClient.GetAllIngredientLineByIds(ingredientLineIds)
	if err != nil {
		return nil, err
	}

	return ingredientLineResult, nil
}

func (s *PromoCodeIdsStringService) getIngredients(ingredientIdsString string) ([]model.IngredientDefaultResult, error) {
	if ingredientIdsString == "" {
		return nil, nil
	}

	ingredientLineIds := helper.GetIdsSlice(ingredientIdsString)

	ingredientResult, err := s.productClient.GetAllIngredientByIds(ingredientLineIds)
	if err != nil {
		return nil, err
	}

	return ingredientResult, nil
}

func (s *PromoCodeIdsStringService) buildProduct(
	brandsResult []model.Default,
	ingredientLinesResult []model.IngredientLineDefaultResult,
	ingredientsResult []model.IngredientDefaultResult,
) []model.BrandResult {
	brandsPromoCodeResult := make([]model.BrandResult, 0, 100)

	for _, brandResult := range brandsResult {
		brandPromoCode := model.BrandResult{
			Id:              brandResult.Id,
			Name:            brandResult.Name,
			IngredientLines: make([]model.IngredientLineResult, 0),
		}

		for _, ingredientLineResult := range ingredientLinesResult {
			if brandResult.Id == ingredientLineResult.BrandId {
				ingredientLinePromoCode := model.IngredientLineResult{
					Id:          ingredientLineResult.Id,
					Name:        ingredientLineResult.Name,
					Ingredients: make([]model.IngredientResult, 0),
				}

				for _, ingredientResult := range ingredientsResult {
					if ingredientLinePromoCode.Id == ingredientResult.IngredientLineId {
						ingredientPromoCode := model.IngredientResult{
							Id:   ingredientResult.Id,
							Name: ingredientResult.Name,
						}

						ingredientLinePromoCode.Ingredients = append(ingredientLinePromoCode.Ingredients, ingredientPromoCode)
					}
				}

				brandPromoCode.IngredientLines = append(brandPromoCode.IngredientLines, ingredientLinePromoCode)
			}

		}

		brandsPromoCodeResult = append(brandsPromoCodeResult, brandPromoCode)
	}

	return brandsPromoCodeResult
}
