package repositoryPromoCode

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	"strings"
	"telemetry-sale/internal/config/storage"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/util/constant"
	"telemetry-sale/internal/util/helper"
)

type PromoCodeRepository struct {
	db  *pgxpool.Pool
	ctx *gin.Context
	tr  *helper.Transaction
}

func NewPromoCodeRepository(ctx *gin.Context) *PromoCodeRepository {
	return &PromoCodeRepository{
		db:  storage.GetDB(),
		ctx: ctx,
	}
}

func (repository *PromoCodeRepository) SavePromoCode(
	promoCode model.PromoCodeCreateDto, paramPromoCode model.ParamPromoCodeIdsString,
) (model.IdResult, error) {
	var result model.IdResult
	var groupId sql.NullInt64

	repository.tr = helper.NewTransaction()
	repository.tr.StartTransaction()

	userUuid, _ := repository.ctx.Get("userUuid")

	var promoCodeId uint64

	//periodFrom, periodTo := repository.preparePeriod(promoCode.PeriodUse)
	if promoCode.GroupId != 0 {
		groupId.Valid = true
		groupId.Int64 = int64(promoCode.GroupId)
	}

	query := `insert into promo_codes (
					code,
					qty,
					discount_type,
					discount_amount,
					period_from,
					period_to,
					organization_id,
					created_id,
					group_id,
					description,
				    machine_ids,
				    category_ids,
				    view_ids,
				    brand_ids,
                	ingredient_line_ids,
				    ingredient_ids,
					created
				) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, NOW()) returning id`

	row := repository.tr.Transaction.QueryRow(
		repository.tr.Ctx,
		query,
		promoCode.Code, promoCode.Qty, string(promoCode.Discount.Type), promoCode.Discount.Amount,
		promoCode.PeriodUse.From, promoCode.PeriodUse.To, promoCode.OrganizationId,
		userUuid, groupId, promoCode.Description, paramPromoCode.MachineIds, paramPromoCode.CategoryIds,
		paramPromoCode.ViewIds, paramPromoCode.BrandIds, paramPromoCode.IngredientLineIds, paramPromoCode.IngredientIds,
	)

	if err := row.Scan(&promoCodeId); err != nil {
		repository.tr.MustRollback()
		return result, err
	}

	if len(promoCode.Schedules) != 0 {

		scheduleRepository := NewScheduleRepository(repository.ctx)
		scheduleRepository.tr = repository.tr

		if err := scheduleRepository.save(promoCodeId, promoCode.Schedules); err != nil {
			repository.tr.MustRollback()
			return result, err
		}
	}

	if err := repository.tr.Transaction.Commit(repository.tr.Ctx); err != nil {
		repository.tr.MustRollback()
	}

	result.Id = promoCodeId

	return result, nil
}

func (repository *PromoCodeRepository) Edit(
	id uint64, promoCode model.PromoCodeEditDto, paramPromoCode model.ParamPromoCodeIdsString,
) error {
	var groupId sql.NullInt64

	repository.tr = helper.NewTransaction()
	repository.tr.StartTransaction()

	userUuid, _ := repository.ctx.Get("userUuid")

	//periodFrom, periodTo := repository.preparePeriod(promoCode.PeriodUse)

	if promoCode.GroupId != 0 {
		groupId.Valid = true
		groupId.Int64 = int64(promoCode.GroupId)
	}

	query := `update promo_codes set code=$1, qty=$2, is_active=$3, discount_type=$4, discount_amount=$5, 
                       period_from=$6, period_to=$7, updated_id=$8, group_id=$9, description=$10, machine_ids=$11, 
                       category_ids=$12, view_ids=$13, brand_ids=$14, ingredient_line_ids=$15, ingredient_ids=$16, updated=NOW()
                       where id=$17`

	_, err := repository.tr.Transaction.Exec(
		repository.tr.Ctx, query, promoCode.Code, promoCode.Qty, promoCode.IsActive, promoCode.Discount.Type,
		promoCode.Discount.Amount, promoCode.PeriodUse.From, promoCode.PeriodUse.To, userUuid, groupId,
		promoCode.Description, paramPromoCode.MachineIds, paramPromoCode.CategoryIds, paramPromoCode.ViewIds,
		paramPromoCode.BrandIds, paramPromoCode.IngredientLineIds, paramPromoCode.IngredientIds, id,
	)
	if err != nil {
		return err
	}

	if len(promoCode.Schedules) != 0 {
		scheduleRepository := NewScheduleRepository(repository.ctx)
		scheduleRepository.tr = repository.tr

		if err := scheduleRepository.changeSchedule(id, promoCode.Schedules); err != nil {
			repository.tr.MustRollback()
			return err
		}
	}

	if err := repository.tr.Transaction.Commit(repository.tr.Ctx); err != nil {
		repository.tr.MustRollback()
	}

	return nil
}

func (repository *PromoCodeRepository) SwitchSelected(id uint64, switchSelected model.SwitchSelected) error {
	query := `update promo_codes set is_selected=$1 where id=$2`

	_, err := repository.db.Exec(
		repository.ctx, query, switchSelected.IsSelected, id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repository *PromoCodeRepository) GetList(
	promoFilter model.PromoCodeFilterList,
	limit uint16,
	offset uint16,
) ([]model.PromoCodeListResult, error) {
	var promoCodeList []model.PromoCodeListResult

	query, args := repository.getQueryForList(promoFilter, limit, offset)

	rows, err := repository.db.Query(repository.ctx, query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var promoCode model.PromoCodeListResult

		var periodFrom pq.NullTime
		var periodTo pq.NullTime

		err = rows.Scan(
			&promoCode.Id, &promoCode.IsSelected, &promoCode.IsActive, &promoCode.Code, &promoCode.Qty, &promoCode.Used,
			&promoCode.Discount.Type, &promoCode.Discount.Amount, &periodFrom, &periodTo,
			&promoCode.GroupName, &promoCode.Description,
		)

		if periodFrom.Valid {
			promoCode.PeriodUse.From = strings.Split(periodFrom.Time.String(), " ")[0]
		}

		if periodTo.Valid {
			promoCode.PeriodUse.To = strings.Split(periodTo.Time.String(), " ")[0]
		}

		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}

		promoCodeList = append(promoCodeList, promoCode)
	}

	rows.Close()

	return promoCodeList, nil
}

func (repository *PromoCodeRepository) GetFilter(promoFilter model.PromoCodeFilterList) (model.PromoCodeFilterValueResult, error) {
	var promoCode model.PromoCodeFilterValueResult

	query := `select max(promo_codes.qty) maxQty, max(promo_codes.discount_amount) discountAmountMax
				from promo_codes
				where promo_codes.organization_id=$1`

	query, args, _ := repository.getQueryFilter(promoFilter, query)

	query += " limit 1"

	err := repository.db.QueryRow(repository.ctx, query, args...).Scan(&promoCode.QtyMax, &promoCode.DiscountAmountMax)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return promoCode, err
	}

	return promoCode, nil
}

func (repository *PromoCodeRepository) GetElement(id uint64) (model.PromoCodeElementResult, model.ParamPromoCodeIdsString, error) {
	var promoCodeElement model.PromoCodeElementResult
	var paramIdsString model.ParamPromoCodeIdsString

	query := `select promo_codes.id, promo_codes.code, promo_codes.qty, promo_codes.is_active,
      				promo_codes.discount_type, promo_codes.discount_amount, promo_codes.period_from,
      				promo_codes.period_to, groups.id, groups.name, promo_codes.description, promo_codes.machine_ids,
      				promo_codes.category_ids, promo_codes.view_ids, promo_codes.brand_ids, promo_codes.ingredient_line_ids, 
      				promo_codes.ingredient_ids, schedules.id, schedules.name, schedules.period_from, schedules.period_to, 
      				schedules.weekday
      			from promo_codes left join groups on groups.id = promo_codes.group_id
      			left join schedules on schedules.promo_code_id = promo_codes.id where promo_codes.id=$1`

	rows, err := repository.db.Query(repository.ctx, query, id)
	if err != nil {
		return promoCodeElement, paramIdsString, err
	}

	var (
		qty         sql.NullInt16
		description sql.NullString
		periodFrom  pq.NullTime
		periodTo    pq.NullTime
	)

	for rows.Next() {
		scheduleResult := &model.ScheduleResult{}
		group := &model.PromoGroupResult{}

		var (
			weekDay sql.NullString
		)

		err = rows.Scan(
			&promoCodeElement.Id, &promoCodeElement.Code, &qty, &promoCodeElement.IsActive,
			&promoCodeElement.Discount.Type, &promoCodeElement.Discount.Amount, &periodFrom,
			&periodTo, &group.Id, &group.Name, &description, &paramIdsString.MachineIds, &paramIdsString.CategoryIds,
			&paramIdsString.ViewIds, &paramIdsString.BrandIds, &paramIdsString.IngredientLineIds, &paramIdsString.IngredientIds,
			&scheduleResult.Id, &scheduleResult.Name, &scheduleResult.From, &scheduleResult.To, &weekDay,
		)

		if err != nil {
			return promoCodeElement, paramIdsString, err
		}

		if qty.Valid {
			promoCodeElement.Qty = uint16(qty.Int16)
		}

		if description.Valid {
			promoCodeElement.Description = description.String
		}

		if periodFrom.Valid {
			promoCodeElement.PeriodUse.From = periodFrom.Time.Format("02.01.2006")
		}

		if periodTo.Valid {
			promoCodeElement.PeriodUse.To = periodTo.Time.Format("02.01.2006")
		}

		if group.Id != nil {
			promoCodeElement.Group = group
		}

		if scheduleResult.Id.Valid {
			schedule := &model.Schedule{
				Id:   uint64(scheduleResult.Id.Int64),
				Name: scheduleResult.Name.String,
				From: scheduleResult.From.String,
				To:   scheduleResult.To.String,
			}

			if weekDay.Valid {
				weekDaysArr := strings.Split(weekDay.String, ",")
				for _, wd := range weekDaysArr {
					schedule.Weekday = append(schedule.Weekday, constant.Weekday(wd))
				}
			}

			promoCodeElement.Schedules = append(promoCodeElement.Schedules, schedule)
		}

	}

	rows.Close()

	return promoCodeElement, paramIdsString, nil
}

func (repository *PromoCodeRepository) GetQty(organizationId uint64) (model.Qty, error) {
	var qty model.Qty

	query := `select count(*) from promo_codes where organization_id=$1`

	err := repository.db.QueryRow(repository.ctx, query, organizationId).Scan(&qty.Qty)

	if err != nil {
		return qty, err
	}

	return qty, nil
}

func (repository *PromoCodeRepository) getQueryForList(
	promoFilter model.PromoCodeFilterList, limit uint16, offset uint16,
) (string, []interface{}) {

	query := `select promo_codes.id, promo_codes.is_selected, promo_codes.is_active, promo_codes.code, promo_codes.qty,
	  				promo_codes.used, promo_codes.discount_type, promo_codes.discount_amount,
	  				promo_codes.period_from, promo_codes.period_to, groups.name, promo_codes.description
				from promo_codes
				left join groups on groups.id = promo_codes.group_id
				where promo_codes.organization_id=$1`

	query, args, argIdx := repository.getQueryFilter(promoFilter, query)

	query += fmt.Sprintf(
		" order by promo_codes.is_active desc, promo_codes.is_selected desc, promo_codes.created %s limit $%d offset $%d",
		promoFilter.CreatedSort, argIdx+1, argIdx+2,
	)
	args = append(args, limit, offset)

	return query, args
}

func (repository *PromoCodeRepository) getQueryFilter(promoFilter model.PromoCodeFilterList, query string) (string, []interface{}, int) {
	args := make([]interface{}, 0, 100)
	args = append(args, promoFilter.OrganizationId)

	argIdx := 1
	if promoFilter.IsSelected != nil {
		argIdx++
		args = append(args, promoFilter.IsSelected)
		query += fmt.Sprintf(" and promo_codes.is_selected=$%d", argIdx)
	}

	if promoFilter.IsActive != nil {
		argIdx++
		args = append(args, promoFilter.IsActive)
		query += fmt.Sprintf(" and promo_codes.is_active=$%d", argIdx)
	}

	if promoFilter.PeriodFrom != nil {
		argIdx++
		args = append(args, promoFilter.PeriodFrom)
		query += fmt.Sprintf(" and promo_codes.period_from>=$%d", argIdx)
	}

	if promoFilter.PeriodTo != nil {
		argIdx++
		dateTo := promoFilter.PeriodTo.AddDate(0, 0, 1)

		args = append(args, dateTo)
		query += fmt.Sprintf(" and promo_codes.period_to<=$%d", argIdx)
	}

	if promoFilter.QtyMin != nil {
		argIdx++
		args = append(args, promoFilter.QtyMin)
		query += fmt.Sprintf(" and promo_codes.qty>=$%d", argIdx)
	}

	if promoFilter.QtyMax != nil {
		argIdx++
		args = append(args, promoFilter.QtyMax)
		query += fmt.Sprintf(" and promo_codes.qty<=$%d", argIdx)
	}

	if promoFilter.DiscountType != nil {
		argIdx++
		args = append(args, promoFilter.DiscountType)
		query += fmt.Sprintf(" and promo_codes.discount_type=$%d", argIdx)
	}

	if promoFilter.DiscountAmountMin != nil {
		argIdx++
		args = append(args, promoFilter.DiscountAmountMin)
		query += fmt.Sprintf(" and promo_codes.discount_amount>=$%d", argIdx)
	}

	if promoFilter.DiscountAmountMax != nil {
		argIdx++
		args = append(args, promoFilter.DiscountAmountMax)
		query += fmt.Sprintf(" and promo_codes.discount_amount<=$%d", argIdx)
	}

	if promoFilter.Code != nil {
		argIdx++
		codeUpper := strings.ToUpper(*promoFilter.Code)
		codeUpper += "%"

		args = append(args, &codeUpper)
		query += fmt.Sprintf(" and promo_codes.code like $%d", argIdx)
	}

	return query, args, argIdx
}
