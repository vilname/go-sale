package promoCode

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"telemetry-sale/internal/model"
	"telemetry-sale/internal/repository/repositoryPromoCode"
)

type UseListener struct {
	repository *repositoryPromoCode.PromoCodeExchangeRepository
}

func NewUseListener() *UseListener {
	return &UseListener{}
}

func (listener *UseListener) ReaderTopic(message []byte) {
	if err := listener.use(message); err != nil {
		slog.Error(err.Error())
	}
}

func (listener *UseListener) use(message []byte) error {
	var body model.Body
	var promoCodeUse model.PromoCodeUse

	if err := json.Unmarshal(message, &body); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(body.Body), &promoCodeUse); err != nil {
		return err
	}

	if promoCodeUse.Id == 0 {
		return errors.New("empty body")
	}

	listener.repository = repositoryPromoCode.NewPromoCodeExchangeRepository(&gin.Context{})

	if err := listener.repository.IncrementUse(promoCodeUse.Id); err != nil {
		return err
	}

	return nil
}
