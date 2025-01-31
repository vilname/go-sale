package promoCode

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/url"
	"os"
	"strconv"
	"strings"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/repository/repositoryPromoCode"
	"telemetry-sale/internal/service/webclient"
	"telemetry-sale/internal/util/constant"
	"telemetry-sale/internal/util/helper"
	"time"
)

type ExchangeService struct {
	repository *repositoryPromoCode.PromoCodeExchangeRepository
}

func NewExchangeService(ctx *gin.Context) *ExchangeService {
	return &ExchangeService{
		repository: repositoryPromoCode.NewPromoCodeExchangeRepository(ctx),
	}
}

type Machine struct {
	Id int `json:"id"`
}

func (service *ExchangeService) Check(promoCodeCheckDto model.PromoCodeCheckDto, organizationId uint64) (model.PromoCodeCheckResult, error) {

	promoCode, err := service.repository.CheckValid(promoCodeCheckDto.Code, organizationId)
	if err != nil {
		return model.PromoCodeCheckResult{}, err
	}

	service.checkPeriodExpired(&promoCode)

	if promoCode.Success {
		service.checkTimeAndDay(&promoCode, promoCodeCheckDto.Time)
	}

	if promoCode.Success {
		service.checkUseExceed(&promoCode)
	}

	if promoCode.Success {
		service.checkIngredient(&promoCode, promoCodeCheckDto.IngredientId)
	}

	if promoCode.Success {
		service.checkMachine(&promoCode, promoCodeCheckDto.SerialNumber)
	}

	return promoCode, nil
}

func (service *ExchangeService) checkPeriodExpired(promoCode *model.PromoCodeCheckResult) {

	now := time.Now()
	dateTo := now.AddDate(0, 0, -1)

	if promoCode.PeriodFrom.Valid && promoCode.PeriodTo.Valid {
		if !(promoCode.PeriodFrom.Time.Before(now) && promoCode.PeriodTo.Time.After(dateTo)) {
			promoCode.Success = false
			promoCode.Reason = constant.Expired

			return
		}
	}

	// если нет даты начала но есть дата конца
	if promoCode.PeriodFrom.Valid && !promoCode.PeriodTo.Valid && !promoCode.PeriodFrom.Time.Before(now) {
		promoCode.Success = false
		promoCode.Reason = constant.Expired
	}

	// если есть дата начала но нет даты конца
	if promoCode.PeriodTo.Valid && !promoCode.PeriodFrom.Valid && !promoCode.PeriodTo.Time.After(dateTo) {
		promoCode.Success = false
		promoCode.Reason = constant.Expired
	}
}

func (service *ExchangeService) checkTimeAndDay(promoCode *model.PromoCodeCheckResult, time time.Time) {
	weekDay := time.Weekday().String()
	hour := uint8(time.Hour())
	minute := uint8(time.Minute())

	shortWeekday := helper.GetWeedDay(weekDay)

	if len(promoCode.Schedules) == 0 {
		return
	}

	isDay := false

	for _, schedule := range promoCode.Schedules {
		var hourFrom uint8
		var minuteFrom uint8
		var hourTo uint8
		var minuteTo uint8
		var err error

		isFindDay := true
		if schedule.Weekday != nil {
			isFindDay = strings.Contains(*schedule.Weekday, shortWeekday)
		}

		if schedule.From != nil {
			hourFrom, minuteFrom, err = service.parseStrTime(*schedule.From)
			if err != nil {
				promoCode.Success = false
				promoCode.Reason = constant.WrongTime
			}
		}

		if schedule.To != nil {
			hourTo, minuteTo, err = service.parseStrTime(*schedule.To)
			if err != nil {
				promoCode.Success = false
				promoCode.Reason = constant.WrongTime
			}
		}

		// если совпадает день недели и временные отрезки в течении дня если они указаны
		if isFindDay &&
			(schedule.From == nil || hour > hourFrom || (hour == hourFrom && minute >= minuteFrom)) &&
			(schedule.To == nil || hour < hourTo || (hour == hourTo && minute <= minuteTo)) {

			isDay = true
		}
	}

	if !isDay {
		promoCode.Success = false
		promoCode.Reason = constant.WrongTime
	}
}

func (service *ExchangeService) checkUseExceed(promoCode *model.PromoCodeCheckResult) {
	if promoCode.Qty.Valid && promoCode.Qty.Int16 > 0 && promoCode.Used.Int16 >= promoCode.Qty.Int16 {
		promoCode.Success = false
		promoCode.Reason = constant.UseExceed
	}
}

func (service *ExchangeService) checkIngredient(promoCode *model.PromoCodeCheckResult, ingredientId uint64) {
	if promoCode.IngredientIds != "" {
		ingredientIds := strings.Split(promoCode.IngredientIds, ",")

		for _, id := range ingredientIds {
			parseUint, _ := strconv.ParseUint(id, 10, 64)

			if ingredientId == parseUint {
				return
			}
		}

		promoCode.Success = false
		promoCode.Reason = constant.InvalidProduct

		return
	}

	// если ограничения по продукты в промокоде не заданы
	if promoCode.IngredientLineIds == "" && promoCode.BrandIds == "" && promoCode.ViewIds == "" && promoCode.CategoryIds == "" {
		return
	}

	productClient := webclient.NewProductClient()
	isPromoCodeIngredient := productClient.IsPromoCodeIngredient(
		ingredientId,
		promoCode.CategoryIds,
		promoCode.ViewIds,
		promoCode.BrandIds,
		promoCode.IngredientLineIds,
	)

	if !isPromoCodeIngredient {
		promoCode.Success = false
		promoCode.Reason = constant.InvalidProduct
	}
}

func (service *ExchangeService) checkMachine(promoCode *model.PromoCodeCheckResult, serialNumber string) {
	var machine Machine

	if promoCode.MachineIds == "" {
		return
	}

	urlPath := &url.URL{
		Scheme: os.Getenv("SCHEME"),
		Host:   os.Getenv("URL_MACHINE_CONTROLEE"),
		Path:   "/machine-exchange/element/" + serialNumber,
	}

	body := helper.GetWebClient(urlPath.String())

	err := json.Unmarshal(body, &machine)
	if err != nil {
		slog.Error("checkMachine", err.Error())
		promoCode.Success = false
		promoCode.Reason = constant.WrongMachine

		return
	}

	machineIds := strings.Split(promoCode.MachineIds, ",")

	for _, machineIdStr := range machineIds {
		machineId, _ := strconv.Atoi(machineIdStr)

		if machineId == machine.Id {
			return
		}
	}

	promoCode.Success = false
	promoCode.Reason = constant.WrongMachine
}

func (service *ExchangeService) parseStrTime(timeStr string) (uint8, uint8, error) {
	timeSlice := strings.Split(timeStr, ":")
	hour, err := strconv.ParseUint(timeSlice[0], 0, 64)
	if err != nil {
		return 0, 0, err
	}

	minute, err := strconv.ParseUint(timeSlice[1], 0, 64)
	if err != nil {
		return 0, 0, err
	}

	return uint8(hour), uint8(minute), nil
}
