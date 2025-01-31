package repository

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"telemetry-sale/internal/config/storage"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/util/helper"
	"time"
)

type SalePeriodRepository struct {
	db  *pgxpool.Pool
	ctx *gin.Context
	tr  *helper.Transaction
}

func NewSalePeriodRepository(ctx *gin.Context) *SalePeriodRepository {
	return &SalePeriodRepository{
		db:  storage.GetDB(),
		ctx: ctx,
	}
}

func (repository *SalePeriodRepository) GetCountSaleByDate(serialNumber string, date time.Time) (uint16, error) {
	var countSale uint16

	query := "select count(sales.id) from sales " +
		"where sales.serial_number = $1 and sales.date_sale > $2"

	err := repository.db.QueryRow(repository.ctx, query, serialNumber, date).Scan(&countSale)
	if err != nil {
		return countSale, err
	}

	return countSale, nil
}

func (repository *SalePeriodRepository) GetCountBySerialNumbersAndTimeAgo(
	serialNumbers []string, timeAgo time.Time,
) ([]model.SaleQtyBySerialNumberResult, error) {
	qtyPeriods := make([]model.SaleQtyBySerialNumberResult, 0, 100)
	args := make([]interface{}, len(serialNumbers), 100)

	for idx, serialNumber := range serialNumbers {
		args[idx] = serialNumber
	}

	args = append(args, timeAgo)

	placeholderString := helper.CreateArgsForIn(len(serialNumbers))

	query := fmt.Sprintf("select sales.serial_number, count(sales.id) from sales "+
		"where sales.serial_number in (%s) and sales.date_sale > $%d group by sales.serial_number", placeholderString, len(serialNumbers)+1)

	row, err := repository.db.Query(repository.ctx, query, args...)

	if err != nil {
		return nil, err
	}

	for row.Next() {
		var qtyPeriod model.SaleQtyBySerialNumberResult

		err := row.Scan(&qtyPeriod.SerialNumber, &qtyPeriod.Qty)
		if err != nil {
			return nil, err
		}

		qtyPeriods = append(qtyPeriods, qtyPeriod)
	}

	row.Close()

	return qtyPeriods, nil
}
