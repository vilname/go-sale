package util

type CellType string

const ErrorKey = "message"

const (
	Ingredient CellType = "INGREDIENT"
	Cup        CellType = "CUP"
	Water      CellType = "WATER"
	Disposable CellType = "DISPOSABLE"
)
