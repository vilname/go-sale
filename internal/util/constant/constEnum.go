package constant

type ResponseMessage string
type KafkaTopic string
type Unit string
type PaymentMethod string
type SortDirection string

const DefaultPageLimit = 20

const (
	ValidationError ResponseMessage = "VALIDATION_ERROR"
)

const (
	SaleImportTopic KafkaTopic = "saleImportTopic"
	// PromoCodeUseImportTopic увеличивает счетчик использования промокода
	PromoCodeUseImportTopic KafkaTopic = "promoCodeUseImportTopic"
)

const (
	ML   Unit = "ML"   // миллилитры
	G    Unit = "G"    // грамм
	MG   Unit = "MG"   // милиграмм
	KCAL Unit = "KCAL" // ккал
	FLOZ Unit = "FLOZ" // жидкая унция
	KG   Unit = "KG"   // килограмм
	MCG  Unit = "MCG"  // микрограмм
	OZ   Unit = "OZ"   // унция
	KJ   Unit = "KJ"   // кДж
)

const (
	INGREDIENT PaymentMethod = "INGREDIENT"
	CUP        PaymentMethod = "CUP"
	WATER      PaymentMethod = "WATER"
	DISPOSABLE PaymentMethod = "DISPOSABLE"
)

const (
	Asc  SortDirection = "ASC"
	Desc SortDirection = "DESC"
)
