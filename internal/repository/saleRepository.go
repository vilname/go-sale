package repository

import (
	"telemetry-sale/internal/config/storage"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/util/helper"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SaleRepository struct {
	db  *pgxpool.Pool
	ctx *gin.Context
	tr  *helper.Transaction
}

func NewSaleRepository(ctx *gin.Context) *SaleRepository {
	return &SaleRepository{
		db:  storage.GetDB(),
		ctx: ctx,
	}
}

func (saleRepository *SaleRepository) CreateSaleFromMachine(sales []model.Sale) error {

	saleRepository.tr = helper.NewTransaction()
	saleRepository.tr.StartTransaction()

	for _, sale := range sales {

		var saleId uint

		promoCodeId := helper.Uint64ToNullInt64(sale.PromoCodeId)

		query := `insert into sales (
					serial_number,
					volume,
                   	name,
					price,
                   	unit,
					discount_id,
					promo_code_id,
					date_sale,
					date_sale_time_zone,
					created
				) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW()) returning id`

		row := saleRepository.tr.Transaction.QueryRow(
			saleRepository.ctx,
			query,
			sale.SerialNumber, sale.Volume, sale.Name, sale.Price, sale.Unit, sale.DiscountId,
			promoCodeId, sale.DateSale, sale.DateSale,
		)

		if err := row.Scan(&saleId); err != nil {
			saleRepository.tr.MustRollback()
			return err
		}

		if err := CreateWriteOff(saleId, sale.WriteOffs, saleRepository.tr); err != nil {
			saleRepository.tr.MustRollback()
			return err
		}

		if len(sale.Payments) == 0 {
			continue
		}

		if err := CreatePayment(saleId, sale.Payments, saleRepository.tr); err != nil {
			saleRepository.tr.MustRollback()
			return err
		}

	}

	if err := saleRepository.tr.Transaction.Commit(saleRepository.tr.Ctx); err != nil {
		saleRepository.tr.MustRollback()
	}

	return nil
}

func (saleRepository *SaleRepository) List(serialNumber string, offset uint16, limit uint16) ([]model.SaleListResult, error) {
	saleListsCheck := make(map[uint64]bool)
	saleLists := make([]model.SaleListResult, 0, 100)
	discountPricesMap := make(map[uint64][]model.DiscountPriceResult)

	query := `select sales.id, sales.name, sales.volume, sales.unit, sales.price, payments.price, payments.method,
       			promo_codes.code, promo_codes.discount_amount, promo_codes.discount_type, sales.date_sale
				from sales 
				    left join payments on payments.sale_id = sales.id
				    left join promo_codes on promo_codes.id = sales.promo_code_id
				where sales.serial_number=$1 order by sales.date_sale desc, sales.created desc offset $2 limit $3`

	rows, err := saleRepository.db.Query(saleRepository.ctx, query, serialNumber, offset, limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		saleList := model.NewSaleListResult()
		var discountPrice model.DiscountPriceResult
		var promoCode model.PromoCodeResult
		var dateSale time.Time

		err = rows.Scan(&saleList.Id, &saleList.Name, &saleList.Volume, &saleList.Unit, &saleList.Price, &discountPrice.Price,
			&discountPrice.Method, &promoCode.Code, &promoCode.Amount, &promoCode.Type, &dateSale)

		if err != nil {
			return nil, err
		}

		if discountPrice.Price != nil {
			discountPricesMap[saleList.Id] = append(discountPricesMap[saleList.Id], discountPrice)
		}

		saleList.DateSale = dateSale.Format("02.01.2006 15:04:05") //27.03.2024 16:03:22

		if promoCode.Code != nil {
			saleList.PromoCode = &promoCode
		}

		ok := saleListsCheck[saleList.Id]
		if ok {
			continue
		}

		saleListsCheck[saleList.Id] = true
		saleLists = append(saleLists, *saleList)
	}

	for idx, saleList := range saleLists {
		discountPrices, ok := discountPricesMap[saleList.Id]
		if !ok {
			continue
		}

		saleLists[idx].DiscountPrices = discountPrices
	}

	rows.Close()

	return saleLists, nil
}

func (saleRepository *SaleRepository) LastDate(serialNumber string) (model.DateLastSale, error) {
	var lastDate model.DateLastSale

	query := `select date_sale from sales where serial_number=$1 order by date_sale desc limit 1`

	err := saleRepository.db.QueryRow(saleRepository.ctx, query, serialNumber).Scan(&lastDate.Date)
	if err != nil {
		return model.DateLastSale{}, err
	}

	return lastDate, nil
}

func (saleRepository *SaleRepository) GetQty(serialNumber string) (model.Qty, error) {
	var qty model.Qty

	query := `select count(*) from sales where serial_number=$1`

	err := saleRepository.db.QueryRow(saleRepository.ctx, query, serialNumber).Scan(&qty.Qty)

	if err != nil {
		return qty, err
	}

	return qty, nil
}
