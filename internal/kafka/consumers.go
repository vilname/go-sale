package kafka

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"os"
	"telemetry-sale/internal/service"
	"telemetry-sale/internal/service/promoCode"
	"telemetry-sale/internal/util/constant"
)

func StartConsumer(topic constant.KafkaTopic) *kafka.Reader {
	kafkaUrl := os.Getenv("KAFKA_URL")

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{kafkaUrl},
		Topic:       string(topic),
		Partition:   0,
		GroupID:     "groupId",
		StartOffset: kafka.LastOffset,
	})
}

func SaleTopicConsumer() {

	ctx := &gin.Context{}

	r := StartConsumer(constant.SaleImportTopic)

	saleService := service.NewSaleService(ctx)

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			panic(fmt.Sprintf("%s: %s", constant.SaleImportTopic, err.Error()))
		}

		fmt.Println("message: ", m.Value)

		err = saleService.CreateSaleFromMachine(m.Value)
		if err != nil {
			fmt.Printf("Сообщение: %s\n", string(m.Value))
		}

	}
}

func PromoCodeUseConsumer() {
	ctx := context.Background()

	r := StartConsumer(constant.PromoCodeUseImportTopic)
	listener := promoCode.NewUseListener()

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			panic(fmt.Sprintf("%s: %s", constant.PromoCodeUseImportTopic, err.Error()))
		}

		fmt.Println("message: ", m.Value)

		listener.ReaderTopic(m.Value)
	}
}
