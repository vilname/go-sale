package webclient

import (
	"encoding/json"
	"net/url"
	"os"
	"strconv"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/util/helper"
)

type MachineClient struct{}

func NewMachineClient() *MachineClient {
	return &MachineClient{}
}

func (c *MachineClient) GetMachineOrganization(serialNumber string) (model.MachineOrganization, error) {
	var machineOrganization model.MachineOrganization

	urlPath := &url.URL{
		Scheme: os.Getenv("SCHEME"),
		Host:   os.Getenv("URL_MACHINE_CONTROLEE"),
		Path:   "/machine-exchange/element/" + serialNumber,
	}

	body := helper.GetWebClient(urlPath.String())

	err := json.Unmarshal(body, &machineOrganization)
	if err != nil {
		return machineOrganization, err
	}

	return machineOrganization, nil
}

func (c *MachineClient) GetByIds(ids []uint64) ([]model.Default, error) {
	machinesResult := make([]model.Default, 0, 100)

	urlPath := &url.URL{
		Scheme: os.Getenv("SCHEME"),
		Host:   os.Getenv("URL_MACHINE_CONTROLEE"),
		Path:   "/machine-exchange/list/by-ids",
	}

	query := urlPath.Query()

	for _, id := range ids {
		query.Add("ids", strconv.Itoa(int(id)))
	}

	urlPath.RawQuery = query.Encode()
	body := helper.GetWebClient(urlPath.String())

	err := json.Unmarshal(body, &machinesResult)
	if err != nil {
		return machinesResult, err
	}

	return machinesResult, nil
}
