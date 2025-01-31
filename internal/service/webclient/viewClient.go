package webclient

import (
	"encoding/json"
	"net/url"
	"os"
	"strconv"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/util/helper"
)

type ViewClient struct{}

func NewViewClient() *ViewClient {
	return &ViewClient{}
}

func (c *ViewClient) GetByIds(ids []uint64) ([]model.Default, error) {
	viewResult := make([]model.Default, 0, 100)

	urlPath := &url.URL{
		Scheme: os.Getenv("SCHEME"),
		Host:   os.Getenv("URL_PRODUCT_BASE"),
		Path:   "/view/list/by-ids",
	}

	query := urlPath.Query()

	for _, id := range ids {
		query.Add("ids", strconv.Itoa(int(id)))
	}

	urlPath.RawQuery = query.Encode()
	body := helper.GetWebClient(urlPath.String())

	err := json.Unmarshal(body, &viewResult)
	if err != nil {
		return viewResult, err
	}

	return viewResult, nil
}
