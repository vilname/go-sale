package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"telemetry-sale/internal/config"
	"telemetry-sale/internal/config/storage"
	"telemetry-sale/internal/kafka"
)

// @title Orders API
// @version 1.0
// @description This is a sample service for managing orders
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email soberkoder@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @description Адрес модуля <b>http://dev.ishaker.ru:8765/telemetry-sale</b>, перед всеми апи запросами нужно добавлять его
//
//	@securitydefinitions.oauth2.implicit	OAuth2Implicit
//	@authorizationUrl						http://localhost:8180/realms/shaker-realm/protocol/openid-connect/auth
//
// @host localhost:8325
// @BasePath /
func main() {

	fmt.Println("Start App v0.01")
	err := godotenv.Load()

	if err != nil {
		fmt.Println("env: ", err.Error())
	}

	// на время разработки, отключаем проверку ssl
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	//_, err = http.Get("https://golang.org/")
	if err != nil {
		fmt.Println(err)
	}

	mode := os.Getenv("MODE")
	if mode == "DEV" {
		config.EurekaConfig()
	}

	storage.InitDB()
	router := config.InitRoute()

	//слушатели kafka
	go kafka.SaleTopicConsumer()
	go kafka.PromoCodeUseConsumer()

	defer func(ctx context.Context) {
		db := storage.GetDB()
		db.Close()
	}(context.Background())

	fmt.Println("init")

	err = router.Run(":" + os.Getenv("PORT"))
	if err != nil {
		return
	}
}
