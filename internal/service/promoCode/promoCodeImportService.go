package promoCode

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/repository/repositoryPromoCode"
	"time"
)

type ImportService struct {
	repository *repositoryPromoCode.PromoCodeExchangeRepository
}

func NewPromoCodeImportService(ctx *gin.Context) *ImportService {
	return &ImportService{
		repository: repositoryPromoCode.NewPromoCodeExchangeRepository(ctx),
	}
}

func (s *ImportService) Import(ctx *gin.Context) error {
	reader, file := s.openBufferFile(ctx)
	defer file.Close()

	importPromoCodes := make([]model.ImportPromoCode, 0, 1000)

	for {
		row, err := s.readRow(reader)

		if err != nil {
			break
		}

		importPromoCode := s.createStructImport(row)

		if importPromoCode.IsDeleted != 0 {
			continue
		}

		if importPromoCode.To.Before(time.Now()) {
			continue
		}

		year := importPromoCode.From.Year()
		if year == 3000 {
			continue
		}

		if importPromoCode.UsageAmount > 0 && importPromoCode.RemainsUsage <= 0 {
			continue
		}

		if importPromoCode.CompanyId != 35 {
			continue
		}

		importPromoCodes = append(importPromoCodes, importPromoCode)

	}

	err := s.repository.SaveImport(importPromoCodes)
	if err != nil {
		return err
	}
	return nil
}

func (s *ImportService) openBufferFile(ctx *gin.Context) (*bufio.Reader, multipart.File) {
	fileForm, err := ctx.FormFile("file")

	file, err := fileForm.Open()
	if err != nil {
		fmt.Println("Unable to open file:", err)
		return nil, nil
	}

	var reader *bufio.Reader
	reader = bufio.NewReader(file)

	return reader, file
}

func (s *ImportService) readRow(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	line = strings.TrimRight(line, "\n")

	if err != nil {
		return "", err
	}

	return line, nil
}

func (s *ImportService) createStructImport(row string) model.ImportPromoCode {
	itemRow := strings.Split(row, ",")
	id, _ := strconv.Atoi(itemRow[0])
	userId, _ := strconv.Atoi(itemRow[1])

	layout := "2006-01-02 15:04:05"
	from, _ := time.Parse(layout, itemRow[3])
	to, _ := time.Parse(layout, itemRow[4])

	usageAmount, _ := strconv.Atoi(itemRow[5])
	remainsUsage, _ := strconv.Atoi(itemRow[6])
	discount, _ := strconv.Atoi(itemRow[7])
	created, _ := strconv.Atoi(itemRow[9])
	updated, _ := strconv.Atoi(itemRow[10])
	used, _ := strconv.Atoi(itemRow[11])
	allAutomats, _ := strconv.Atoi(itemRow[12])
	companyId, _ := strconv.Atoi(itemRow[13])
	isDeleted, _ := strconv.Atoi(itemRow[14])

	organizationId, _ := strconv.Atoi(os.Getenv("FITVEND_ORGANIZATION_ID"))

	return model.ImportPromoCode{
		Id:             id,
		UserId:         userId,
		PromoCode:      itemRow[2],
		From:           from,
		To:             to,
		UsageAmount:    usageAmount,
		RemainsUsage:   remainsUsage,
		Discount:       discount,
		Tastes:         itemRow[8],
		Created:        created,
		Updated:        updated,
		Used:           used,
		AllAutomats:    allAutomats,
		CompanyId:      companyId,
		IsDeleted:      isDeleted,
		OrganizationId: organizationId,
	}
}
