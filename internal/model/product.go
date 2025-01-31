package model

type Product struct {
	CategoryIds []uint64 `json:"categoryIds"`
	ViewIds     []uint64 `json:"viewIds"`
	Brands      []Brand  `json:"brands"`
}

type Brand struct {
	Id              uint64           `json:"id"`
	IngredientLines []IngredientLine `json:"ingredientLines"`
}

type IngredientLine struct {
	Id            uint64   `json:"id"`
	IngredientIds []uint64 `json:"ingredientIds"`
}

type ParamPromoCodeIdsString struct {
	MachineIds        string
	CategoryIds       string
	ViewIds           string
	BrandIds          string
	IngredientLineIds string
	IngredientIds     string
}

type ProductResult struct {
	Categories []Default     `json:"categories"`
	Views      []Default     `json:"views"`
	Brands     []BrandResult `json:"brands"`
}

type BrandResult struct {
	Id              uint64                 `json:"id"`
	Name            string                 `json:"name"`
	IngredientLines []IngredientLineResult `json:"ingredientLines"`
}

type IngredientLineResult struct {
	Id          uint64             `json:"id"`
	Name        string             `json:"name"`
	Ingredients []IngredientResult `json:"ingredients"`
}

type IngredientLineDefaultResult struct {
	Id      uint64 `json:"id"`
	Name    string `json:"name"`
	BrandId uint64 `json:"brandId"`
}

type IngredientDefaultResult struct {
	Id               uint64 `json:"id"`
	Name             string `json:"name"`
	IngredientLineId uint64 `json:"ingredientLineId"`
}

type IngredientResult struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}
