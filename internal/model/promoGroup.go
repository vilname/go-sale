package model

type PromoGroupDto struct {
	Name           string `json:"name"`
	OrganizationId uint64 `json:"organizationId"`
}

type PromoGroupResult struct {
	Id   *uint64 `json:"id"`
	Name *string `json:"name"`
}
