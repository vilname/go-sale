package helper

import (
	"strconv"
)

func GetDefaultPage(pageStr string) uint16 {
	page, _ := strconv.ParseUint(pageStr, 10, 16)

	if page == 0 {
		page = 1
	}

	return uint16(page)
}

func GetDefaultLimit(limitStr string, defaultLimit uint16) uint16 {
	limit, _ := strconv.ParseUint(limitStr, 10, 16)

	limit16 := uint16(limit)

	if limit16 == 0 {
		limit16 = defaultLimit
	}

	return limit16
}

func GetOffset(page uint16, limit uint16) uint16 {
	return (page - 1) * limit
}
