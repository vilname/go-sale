package service

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/repository"
	"telemetry-sale/internal/util/helper"
)

type SaleService struct {
	repository *repository.SaleRepository
}

func NewSaleService(ctx *gin.Context) *SaleService {
	return &SaleService{
		repository: repository.NewSaleRepository(ctx),
	}
}

func (saleService *SaleService) CreateSaleFromMachine(message []byte) error {
	var body model.Body
	var sales []model.Sale

	if err := json.Unmarshal(message, &body); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(body.Body), &sales); err != nil {
		return err
	}

	if len(sales) == 0 {
		return errors.New("empty body")
	}

	for i, _ := range sales {
		sales[i].SerialNumber = body.ClientId
	}

	if err := saleService.repository.CreateSaleFromMachine(sales); err != nil {
		return err
	}

	return nil
}

func (saleService *SaleService) List(serialNumber string, page uint16, limit uint16) ([]model.SaleListResult, error) {
	offset := helper.GetOffset(page, limit)

	return saleService.repository.List(serialNumber, offset, limit)
}

func (saleService *SaleService) LastDate(serialNumber string) (model.DateLastSale, error) {
	return saleService.repository.LastDate(serialNumber)
}

func (saleService *SaleService) Qty(serialNumber string) (model.Qty, error) {
	return saleService.repository.GetQty(serialNumber)
}
