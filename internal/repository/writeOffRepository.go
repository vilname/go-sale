package repository

import (
	"fmt"
	"strconv"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/util/helper"
)

func CreateWriteOff(saleId uint, writeOff []model.WriteOff, transaction *helper.Transaction) error {
	query := `insert into write_offs (
					cell_number,
                    volume,
					ingredient_id,
					sale_id,
                    unit,
                    cell_type,
					created 
			) values `
	args := []interface{}{}

	key := 1
	step := 6

	//init := nil

	for _, writeOff := range writeOff {

		query += fmt.Sprintf(
			"($%s, $%s, $%s, $%s, $%s, $%s, NOW()),",
			strconv.Itoa(key), strconv.Itoa(key+1), strconv.Itoa(key+2),
			strconv.Itoa(key+3), strconv.Itoa(key+4), strconv.Itoa(key+5),
		)

		var unit interface{} = nil

		if writeOff.Unit != "" {
			unit = writeOff.Unit
		}

		args = append(
			args, writeOff.CellNumber, writeOff.Volume, writeOff.IngredientId, saleId, unit, writeOff.CellType,
		)

		key += step
	}

	err := transaction.SaveAll(query, args, "writeOffSave"+string(rune(key)))
	if err != nil {
		return err
	}
	return nil
}
