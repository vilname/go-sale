package webclient

import (
	"encoding/json"
	"net/url"
	"os"
	"strconv"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/util/helper"
)

type CategoryClient struct{}

func NewCategoryClient() *CategoryClient {
	return &CategoryClient{}
}

func (c *CategoryClient) GetByIds(ids []uint64) ([]model.Default, error) {
	categoryResult := make([]model.Default, 0, 100)

	urlPath := &url.URL{
		Scheme: os.Getenv("SCHEME"),
		Host:   os.Getenv("URL_PRODUCT_BASE"),
		Path:   "/cell-category/list/by-ids",
	}

	query := urlPath.Query()

	for _, id := range ids {
		query.Add("ids", strconv.Itoa(int(id)))
	}

	urlPath.RawQuery = query.Encode()
	body := helper.GetWebClient(urlPath.String())

	err := json.Unmarshal(body, &categoryResult)
	if err != nil {
		return categoryResult, err
	}

	return categoryResult, nil
}
