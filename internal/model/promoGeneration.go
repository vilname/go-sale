package model

type PrefixPromoCodeDto struct {
	Prefix string `json:"prefix" validate:"max=6,regexp=^[A-Z0-9]*$"`
}

type PromoCodeGenerationResult struct {
	Code string `json:"code"`
}
