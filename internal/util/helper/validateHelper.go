package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"telemetry-sale/internal/util/constant"
)

type ErrorValidate struct {
	Message constant.ResponseMessage `json:"message"`
	Field   []map[string]Field       `json:"field"`
}

type Field struct {
	Messages []string `json:"messages"`
}

func RegisterValidate(structItem interface{}) error {
	var validate = validator.New()
	err := validate.RegisterValidation("regexp", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(fl.Param())

		return re.MatchString(fl.Field().Interface().(string))
	})

	if err != nil {
		slog.Error("error validate ", err.Error())
	}

	validateError := validate.Struct(structItem)

	return validateError
}

func ValidateHelper(validateError error) *ErrorValidate {
	var errorValidateResponse = &ErrorValidate{
		Message: constant.ValidationError,
	}

	for _, errField := range validateError.(validator.ValidationErrors) {
		switch errField.Tag() {
		case "required":
			createNullError(errorValidateResponse, errField)
		case "min", "max":
			createLengthError(errorValidateResponse, errField)
		case "regexp":
			createRegExpError(errorValidateResponse, errField)
		}

	}

	return errorValidateResponse
}

func ErrorValidateResponse(ctx *gin.Context, validateError error) {
	errorValidateResponse := ValidateHelper(validateError)
	ctx.JSON(http.StatusBadRequest, errorValidateResponse)
}

func createLengthError(errorValidateResponse *ErrorValidate, errField validator.FieldError) {
	errorValidateResponse.Field = append(errorValidateResponse.Field, map[string]Field{
		strings.ToLower(errField.StructField()): {
			Messages: []string{
				fmt.Sprintf("%s_%s", strings.ToUpper(errField.Tag()), errField.Param()),
			},
		},
	})
}

func createNullError(errorValidateResponse *ErrorValidate, errField validator.FieldError) {
	errorValidateResponse.Field = append(errorValidateResponse.Field, map[string]Field{
		strings.ToLower(errField.StructField()): {
			Messages: []string{"NOT_NULL"},
		},
	})
}

func createRegExpError(errorValidateResponse *ErrorValidate, errField validator.FieldError) {
	re := regexp.MustCompile(`\[(.*?)]`)
	matches := re.FindStringSubmatch(errField.Param())

	errorValidateResponse.Field = append(errorValidateResponse.Field, map[string]Field{
		strings.ToLower(errField.StructField()): {
			Messages: []string{
				fmt.Sprintf("ALLOWED_CHARACTERS: %s", matches[1]),
			},
		},
	})
}
