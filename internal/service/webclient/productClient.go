package webclient

import (
	"encoding/json"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/util/helper"
)

type ProductClient struct{}

func NewProductClient() *ProductClient {
	return &ProductClient{}
}

func (c *ProductClient) GetAllBrandByIds(ids []uint64) ([]model.Default, error) {
	brandResult := make([]model.Default, 0, 100)

	urlPath := &url.URL{
		Scheme: os.Getenv("SCHEME"),
		Host:   os.Getenv("URL_PRODUCT_BASE"),
		Path:   "/brand/list/by-ids",
	}

	query := urlPath.Query()

	for _, id := range ids {
		query.Add("ids", strconv.Itoa(int(id)))
	}

	urlPath.RawQuery = query.Encode()
	body := helper.GetWebClient(urlPath.String())

	err := json.Unmarshal(body, &brandResult)
	if err != nil {
		return brandResult, err
	}

	return brandResult, nil
}

func (c *ProductClient) GetAllIngredientLineByIds(ids []uint64) ([]model.IngredientLineDefaultResult, error) {
	ingredientLineResult := make([]model.IngredientLineDefaultResult, 0, 100)

	urlPath := &url.URL{
		Scheme: os.Getenv("SCHEME"),
		Host:   os.Getenv("URL_PRODUCT_BASE"),
		Path:   "/ingredient-line/list/by-ids",
	}

	query := urlPath.Query()

	for _, id := range ids {
		query.Add("ids", strconv.Itoa(int(id)))
	}

	urlPath.RawQuery = query.Encode()
	body := helper.GetWebClient(urlPath.String())

	err := json.Unmarshal(body, &ingredientLineResult)
	if err != nil {
		return ingredientLineResult, err
	}

	return ingredientLineResult, nil
}

func (c *ProductClient) GetAllIngredientByIds(ids []uint64) ([]model.IngredientDefaultResult, error) {
	ingredientResult := make([]model.IngredientDefaultResult, 0, 100)

	urlPath := &url.URL{
		Scheme: os.Getenv("SCHEME"),
		Host:   os.Getenv("URL_PRODUCT_BASE"),
		Path:   "/ingredient/list-default/by-ids",
	}

	query := urlPath.Query()

	for _, id := range ids {
		query.Add("ids", strconv.Itoa(int(id)))
	}

	urlPath.RawQuery = query.Encode()
	body := helper.GetWebClient(urlPath.String())

	err := json.Unmarshal(body, &ingredientResult)
	if err != nil {
		return ingredientResult, err
	}

	return ingredientResult, nil
}

func (c *ProductClient) IsPromoCodeIngredient(
	ingredientId uint64,
	categoryIdsString string,
	viewIdsString string,
	brandIdsString string,
	ingredientLineIdsString string,
) bool {
	urlPath := &url.URL{
		Scheme: os.Getenv("SCHEME"),
		Host:   os.Getenv("URL_PRODUCT_BASE"),
		Path:   "/ingredient/is-promocode-ingredient/" + strconv.FormatUint(ingredientId, 10),
	}

	query := urlPath.Query()

	if categoryIdsString != "" {
		ids := helper.GetIdsSlice(categoryIdsString)
		helper.AddQueryParamSlice(&query, "categoryIds", ids)
	}

	if viewIdsString != "" {
		ids := helper.GetIdsSlice(viewIdsString)
		helper.AddQueryParamSlice(&query, "viewIds", ids)
	}

	if brandIdsString != "" {
		ids := helper.GetIdsSlice(brandIdsString)
		helper.AddQueryParamSlice(&query, "brandIds", ids)
	}

	if ingredientLineIdsString != "" {
		ids := helper.GetIdsSlice(ingredientLineIdsString)
		helper.AddQueryParamSlice(&query, "ingredientLineIds", ids)
	}

	urlPath.RawQuery = query.Encode()
	body := helper.GetWebClient(urlPath.String())

	var isPromoCodeIngredient bool

	err := json.Unmarshal(body, &isPromoCodeIngredient)
	if err != nil {
		slog.Error("IsPromoCodeIngredient", err.Error())
	}

	return isPromoCodeIngredient
}
