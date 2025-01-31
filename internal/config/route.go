package config

import (
	"telemetry-sale/docs"
	"telemetry-sale/internal/controller/rest"
	"telemetry-sale/internal/util/middleware"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoute() *gin.Engine {
	router := gin.New()
	router.Use(middleware.AuthenticationMiddleware)

	router.GET("/sale-period/all/:serialNumber", rest.GetAllSalePeriod)
	router.GET("/sale-period/list/by-serial-numbers", rest.GetSalePeriodBySerialNumbers)
	router.GET("/sale-period/list/by-serial-numbers-period", rest.GetSalePeriodBySerialNumbersPeriod)

	promoGroup := router.Group("/promo-group")
	promoGroup.POST("/create", rest.CreatePromoGroup)
	promoGroup.GET("/list/:organizationId", rest.ListPromoGroup)

	// Промокоды
	promoCode := router.Group("/promo-code")
	promoCode.POST("/create", rest.CreatePromoCode)
	promoCode.POST("/edit/:id", rest.EditPromoCode)
	promoCode.POST("/switch-selected/:id", rest.SwitchSelected)
	promoCode.GET("/list/:organizationId", rest.ListPromoCode)
	promoCode.GET("/filter/:organizationId", rest.FilterPromoCode)
	promoCode.GET("/element/:id", rest.ElementPromoCode)
	promoCode.GET("/qty/:organizationId", rest.QtyPromoCode)

	promoCode.GET("/check", rest.CheckPromoCode)
	promoCode.POST("/import", rest.ImportPromoCode)

	// Продажи
	sale := router.Group("/sale")
	sale.GET("/list/:serialNumber", rest.ListSale)
	sale.GET("/qty/:serialNumber", rest.QtySale)
	sale.GET("/date-last-sale/:serialNumber", rest.DateLastSale)
	sale.POST("/get-serial-numbers-count", rest.GetSerialNumbersCount)

	// Генерация промокода
	promoGeneration := router.Group("/promo-generation")
	promoGeneration.POST("/single/:organizationId", rest.GenerationCode)

	docs.SwaggerInfo.BasePath = ""
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
