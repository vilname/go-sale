package repository

import (
	"fmt"
	"strconv"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/util/helper"
)

func CreatePayment(saleId uint, payments []model.Payment, transaction *helper.Transaction) error {
	query := `insert into payments (
					price,
                    sale_id,
					method,
					created`
	args := []interface{}{}

	key := 1
	step := 3

	for _, payment := range payments {

		query += fmt.Sprintf(
			") values ($%s, $%s, $%s, NOW()),",
			strconv.Itoa(key), strconv.Itoa(key+1), strconv.Itoa(key+2),
		)
		args = append(
			args, payment.Price, saleId, payment.Method,
		)

		key += step
	}

	err := transaction.SaveAll(query, args, "paymentSave")
	if err != nil {
		return err
	}
	return nil
}
