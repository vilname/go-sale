package repositoryPromoCode

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"telemetry-sale/internal/config/storage"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/util/constant"
	"telemetry-sale/internal/util/helper"
)

type PromoCodeExchangeRepository struct {
	db  *pgxpool.Pool
	ctx *gin.Context
	tr  *helper.Transaction
}

func NewPromoCodeExchangeRepository(ctx *gin.Context) *PromoCodeExchangeRepository {
	return &PromoCodeExchangeRepository{
		db:  storage.GetDB(),
		ctx: ctx,
	}
}

func (repository *PromoCodeExchangeRepository) CheckValid(
	code string, organizationId uint64,
) (model.PromoCodeCheckResult, error) {
	promoCodeResult := model.PromoCodeCheckResult{Success: true}

	query := `select pc.id, pc.discount_type, pc.discount_amount, pc.period_from, pc.period_to, pc.qty, pc.used,
       					pc.category_ids, pc.view_ids, pc.brand_ids, pc.ingredient_line_ids, pc.ingredient_ids, 
       					pc.machine_ids, schedules.weekday, schedules.period_from, schedules.period_to
				from promo_codes pc
				left join schedules on schedules.promo_code_id = pc.id
				where code=$1 and organization_id=$2`

	rows, err := repository.db.Query(repository.ctx, query, code, organizationId)

	if err != nil {
		return promoCodeResult, err
	}

	isFound := false
	for rows.Next() {
		isFound = true

		var schedule model.ScheduleCheckPromoCode

		err = rows.Scan(
			&promoCodeResult.Id, &promoCodeResult.Type, &promoCodeResult.Amount, &promoCodeResult.PeriodFrom,
			&promoCodeResult.PeriodTo, &promoCodeResult.Qty, &promoCodeResult.Used, &promoCodeResult.CategoryIds,
			&promoCodeResult.ViewIds, &promoCodeResult.BrandIds, &promoCodeResult.IngredientLineIds, &promoCodeResult.IngredientIds,
			&promoCodeResult.MachineIds, &schedule.Weekday, &schedule.From, &schedule.To,
		)
		if err != nil {
			return model.PromoCodeCheckResult{}, err
		}

		if schedule.Weekday != nil || schedule.From != nil || schedule.To != nil {
			promoCodeResult.Schedules = append(promoCodeResult.Schedules, &schedule)
		}

	}

	rows.Close()

	if !isFound {
		promoCodeResult.Success = false
		promoCodeResult.Reason = constant.NotFound
	}

	return promoCodeResult, nil
}

func (repository *PromoCodeExchangeRepository) IncrementUse(id uint64) error {
	query := `update promo_codes set used=(select used from promo_codes where id=$1)+1 where id=$2`

	_, err := repository.db.Exec(repository.ctx, query, id, id)
	if err != nil {
		return err
	}
	return nil
}

func (repository *PromoCodeExchangeRepository) SaveImport(importPromoCodes []model.ImportPromoCode) error {
	repository.tr = helper.NewTransaction()
	repository.tr.StartTransaction()

	for _, item := range importPromoCodes {

		description := ""
		if item.AllAutomats == 0 {
			description = "Была привязка к автоматам"
		}

		query := `insert into promo_codes (
					code,
					qty,
					discount_amount,
                   	period_from,
					period_to,
					organization_id,
					description,
					used,
				    discount_type,
					created
				) values ($1, $2, $3, $4, $5, $6, $7, $8, 'PERCENT', NOW())`

		_, err := repository.tr.Transaction.Exec(repository.ctx, query, item.PromoCode, item.UsageAmount, item.Discount, item.From,
			item.To, item.OrganizationId, description, item.Used)

		if err != nil {
			repository.tr.MustRollback()
			return err
		}
	}

	if err := repository.tr.Transaction.Commit(repository.tr.Ctx); err != nil {
		repository.tr.MustRollback()
	}

	return nil
}
