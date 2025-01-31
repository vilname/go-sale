package repositoryPromoCode

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"strings"
	"telemetry-sale/internal/config/storage"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/util/constant"
	"telemetry-sale/internal/util/helper"
)

type ScheduleRepository struct {
	db  *pgxpool.Pool
	ctx *gin.Context
	tr  *helper.Transaction
}

func NewScheduleRepository(ctx *gin.Context) *ScheduleRepository {
	return &ScheduleRepository{
		db:  storage.GetDB(),
		ctx: ctx,
	}
}

func (repository *ScheduleRepository) FindSchedulesByPromoCodeId(promoCodeId uint64) ([]model.Schedule, error) {
	var schedules []model.Schedule

	query := `select id, name, period_from, period_to, weekday from schedules where promo_code_id = $1`

	rows, err := repository.db.Query(repository.ctx, query, promoCodeId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var weekDays sql.NullString
		var schedule model.Schedule

		if err = rows.Scan(&schedule.Id, &schedule.Name, &schedule.From, &schedule.To, &weekDays); err != nil {
			return nil, err
		}

		weekDaysSlice := strings.Split(weekDays.String, ",")

		for _, wd := range weekDaysSlice {
			schedule.Weekday = append(schedule.Weekday, constant.Weekday(wd))
		}

		schedules = append(schedules, schedule)
	}

	rows.Close()

	return schedules, nil
}

func (repository *ScheduleRepository) changeSchedule(promoCodeId uint64, schedules []model.Schedule) error {
	var newSchedules []model.NewSchedule
	var oldSchedules []model.Schedule
	var editSchedules []model.Schedule
	var deleteScheduleIds []uint64

	schedulesBase, err := repository.FindSchedulesByPromoCodeId(promoCodeId)
	if err != nil {
		return err
	}

	for _, schedule := range schedules {
		if schedule.Id == 0 {
			scheduleDto := model.NewSchedule{
				Name:    schedule.Name,
				From:    schedule.From,
				To:      schedule.To,
				Weekday: schedule.Weekday,
			}

			newSchedules = append(newSchedules, scheduleDto)
		} else {
			oldSchedules = append(oldSchedules, schedule)
		}
	}

	for _, scheduleBase := range schedulesBase {
		isExist := false

		for _, schedule := range oldSchedules {
			if scheduleBase.Id == schedule.Id {
				editSchedules = append(editSchedules, schedule)
				isExist = true
				break
			}
		}

		if isExist == false {
			deleteScheduleIds = append(deleteScheduleIds, scheduleBase.Id)
		}
	}

	if len(newSchedules) != 0 {
		if err := repository.save(promoCodeId, newSchedules); err != nil {
			return err
		}
	}

	if len(editSchedules) != 0 {
		if err := repository.edit(editSchedules); err != nil {
			return err
		}
	}

	if len(deleteScheduleIds) != 0 {
		if err := repository.deleteAll(deleteScheduleIds); err != nil {
			return err
		}
	}

	return nil
}

func (repository *ScheduleRepository) edit(schedules []model.Schedule) error {
	query := `update schedules set name=$1, period_from=$2, period_to=$3, weekday=$4, updated=NOW() where id=$5`

	for _, schedule := range schedules {
		weekDay := helper.WeekDayEnumToString(schedule.Weekday)

		_, err := repository.tr.Transaction.Exec(
			repository.tr.Ctx, query, schedule.Name, schedule.From, schedule.To, weekDay, schedule.Id,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repository *ScheduleRepository) deleteAll(scheduleIds []uint64) error {

	args := make([]interface{}, len(scheduleIds))
	for idx, id := range scheduleIds {
		args[idx] = id
	}

	placeholderString := helper.CreateArgsForIn(len(scheduleIds))

	query := fmt.Sprintf(`delete from schedules where id in (%s)`, placeholderString)
	_, err := repository.tr.Transaction.Exec(repository.tr.Ctx, query, args...)

	if err != nil {
		return err
	}

	return nil
}

func (repository *ScheduleRepository) save(promoCodeId uint64, schedules []model.NewSchedule) error {
	key := 1
	step := 5

	args := []interface{}{}

	query := `insert into schedules (
					name,
					period_from,
					period_to,
					weekday,
					promo_code_id,
					created
				) values `

	for _, schedule := range schedules {
		query += fmt.Sprintf(
			"($%s, $%s, $%s, $%s, $%s, NOW()),",
			strconv.Itoa(key), strconv.Itoa(key+1), strconv.Itoa(key+2),
			strconv.Itoa(key+3), strconv.Itoa(key+4),
		)

		var weekDayBuild strings.Builder

		for _, day := range schedule.Weekday {
			weekDayBuild.WriteString(string(day + ","))
		}

		weekDay := strings.TrimRight(weekDayBuild.String(), ",")

		args = append(
			args, schedule.Name, schedule.From, schedule.To, weekDay, promoCodeId,
		)

		key += step
	}

	if err := repository.tr.SaveAll(query, args, "scheduleSave"); err != nil {
		return err
	}

	return nil
}
